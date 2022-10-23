package reverse_proxy_server

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"rabbit_gather/src/server"

	//"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"rabbit_gather/src/logger"
	"rabbit_gather/util"
	"strings"
)

var CertFile string
var KeyFile string
var ServePath *url.URL
var SSLCertificationsCrts []string

// A map that recorde the SubHostName - real server url correlation
var redirectAddrMap map[string]*url.URL
var log = logger.NewLoggerWrapper("ReverseProxyServer")

func init() {
	type Config struct {
		CERT_FILE         string            `json:"cert_file"`
		KEY_FILE          string            `json:"key_file"`
		ServePath         string            `json:"serve_path"`
		MEOWALIEN_COM_CRT string            `json:"meowalien_com_crt"`
		RedirectAddrMap   map[string]string `json:"redirect_addr_map"`
	}
	var config Config
	err := util.ParseFileJsonConfig(&config, "config/reverse_proxy_server.config.json")
	if err != nil {
		panic(err.Error())
	}

	SSLCertificationsCrts = []string{
		config.MEOWALIEN_COM_CRT,
	}
	ServePath, err = url.Parse(config.ServePath)
	if err != nil {
		panic(err.Error())
	}
	CertFile = config.CERT_FILE
	KeyFile = config.KEY_FILE
	redirectAddrMap = map[string]*url.URL{}
	for s, s2 := range config.RedirectAddrMap {
		u, e := url.Parse(s2)
		if e != nil {
			panic(e.Error())
		}
		redirectAddrMap[s] = u
	}
}

// The ReverseProxyServer will deliver the requests and responses between the open port and the local port
type ReverseProxyServer struct {
	serverInst *http.Server
	ginEngine  *gin.Engine
}

func (r *ReverseProxyServer) Startup(ctx context.Context) error {
	log.DEBUG.Println("ReverseProxyServer listen on : ", ServePath.String())
	r.ginEngine = gin.Default()
	crtPool := x509.NewCertPool()
	for _, crt := range SSLCertificationsCrts {
		crtFile, err := ioutil.ReadFile(crt)
		if err != nil {
			panic(err.Error())
		}
		crtPool.AppendCertsFromPEM(crtFile)
	}
	// distribute all request
	r.ginEngine.Use(r.distributor)

	r.serverInst = &http.Server{
		Addr:    ":" + ServePath.Port(),
		Handler: r.ginEngine,
		//ErrorLog: s.logger,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
			ClientCAs:  crtPool,
			ClientAuth: tls.VerifyClientCertIfGiven,
		},
		ReadTimeout:       0,
		ReadHeaderTimeout: 0,
		WriteTimeout:      0,
		IdleTimeout:       0,
		MaxHeaderBytes:    0,
		TLSNextProto:      nil,
		ConnState:         nil,
		ErrorLog:          &log.ERROR.Logger,
		BaseContext:       nil,
		ConnContext:       nil,
	}

	go func() {
		if err := r.serverInst.ListenAndServeTLS(CertFile, KeyFile); err != nil && err != http.ErrServerClosed {
			panic("ReverseProxyServer Error: %s\n" + err.Error())
		}
	}()

	log.DEBUG.Println("ReverseProxyServer Started .")
	return nil
}

// distribute all request according to the sub host name
func (s *ReverseProxyServer) distributor(c *gin.Context) {
	log.DEBUG.Println("Request Host: ", c.Request.Host)
	req := c.Request
	if req.Host == "" {
		c.AbortWithStatusJSON(
			http.StatusNotFound,
			gin.H{
				"err": "Hostname is empty",
			})
		log.DEBUG.Println("Hostname is empty")
		return
	}

	sph := strings.Split(req.Host, ".")
	if len(sph) < 1 {
		c.AbortWithStatusJSON(
			http.StatusNotFound,
			gin.H{
				"err": "sub hostname is empty",
			})
		log.DEBUG.Println("sub hostname is empty")
		return
	}

	subHost := sph[0]
	realAddrURL, exist := redirectAddrMap[subHost]
	if !exist {
		mes := fmt.Sprint("SubHost not exist: ", subHost)
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"err": mes,
		})
		log.DEBUG.Println(mes)
		return
	}
	log.DEBUG.Printf("redirect to : %s\n", realAddrURL.Host)

	proxy := httputil.NewSingleHostReverseProxy(realAddrURL)
	proxy.Director = func(req *http.Request) {
		req.URL.Scheme = realAddrURL.Scheme
		req.URL.Host = realAddrURL.Host

		// to let the target server make sure the request come from here
		req.Header.Add(server.IDENTIFICATION_SYMBOL_KEY, server.IDENTIFICATION_SYMBOL)
		req.Header.Set(server.ClientIP_KEY, c.ClientIP())
	}
	proxy.ModifyResponse = func(response *http.Response) error {
		response.Header.Del(server.IDENTIFICATION_SYMBOL_KEY)
		response.Header.Del(server.ClientIP_KEY)
		return nil
	}
	proxy.ServeHTTP(c.Writer, c.Request)
	return
}

func (r *ReverseProxyServer) Shutdown() error {
	log.DEBUG.Println("ReverseProxyServer start to Shutdown")

	ctx, cancel := context.WithTimeout(context.Background(), server.ShutdownWaitTime)
	defer cancel()

	if err := r.serverInst.Shutdown(ctx); err != nil {
		return fmt.Errorf("ReverseProxyServer fail to Shutdown: %s", err.Error())
	} else {
		log.DEBUG.Println("ReverseProxyServer closed.")
		return nil
	}
}
