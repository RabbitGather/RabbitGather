package reverse_proxy_server

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	//"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"rabbit_gather/src/logger"
	"rabbit_gather/util"
	"strings"
	"time"
)

type ReverseProxyServer struct {
	serverInst *http.Server
	ginEngine  *gin.Engine
}

var CERT_FILE string
var KEY_FILE string
var ServePath *url.URL
var SSLCertificationsCrts []string
var redirectAddrMap map[string]*url.URL
var log = logger.NewLoggerWrapper("ReverseProxyServer")

func init() {
	type Config struct {
		CERT_FILE         string
		KEY_FILE          string
		ServePath         string
		MEOWALIEN_COM_CRT string
		RedirectAddrMap   map[string]string
	}
	var config Config
	err := util.ParseJsonConfic(&config, "config/reverse_proxy_server.config.json")
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
	CERT_FILE = config.CERT_FILE
	KEY_FILE = config.KEY_FILE
	redirectAddrMap = map[string]*url.URL{}
	for s, s2 := range config.RedirectAddrMap {
		u, e := url.Parse(s2)
		if e != nil {
			panic(e.Error())
		}
		redirectAddrMap[s] = u
	}
}

func (r *ReverseProxyServer) Startup(ctx context.Context, shutdownCallback util.ShutdownCallback) error {
	log.DEBUG.Println("ReverseProxyServer listen on : ", ServePath)

	shutdownCallback(r.shutdown)
	r.ginEngine = gin.Default()
	crtPool := x509.NewCertPool()
	for _, crt := range SSLCertificationsCrts {
		crtFile, err := ioutil.ReadFile(crt)
		if err != nil {
			panic(err.Error())
		}
		crtPool.AppendCertsFromPEM(crtFile)
	}

	//r.ginEngine.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
	//	return fmt.Sprintf("%s |%s %d %s| %s |%s %s %s %s | %s | %s | %s\n",
	//		param.TimeStamp.Format(time.RFC1123),
	//		param.StatusCodeColor(),
	//		param.StatusCode,
	//		param.ResetColor(),
	//		param.ClientIP,
	//		param.MethodColor(),
	//		param.Method,
	//		param.ResetColor(),
	//		param.Path,
	//		param.Latency,
	//		param.Request.UserAgent(),
	//		param.ErrorMessage,
	//	)
	//}))

	// 分配器掛載在根路由，轉發任何種類的請求
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
		if err := r.serverInst.ListenAndServeTLS(CERT_FILE, KEY_FILE); err != nil && err != http.ErrServerClosed {
			panic("ReverseProxyServer Error: %s\n" + err.Error())
		}
	}()

	//r.AfterServerStartup()
	log.DEBUG.Println("ReverseProxyServer Started .")
	return nil
}

// 分配請求的主要邏輯
func (s *ReverseProxyServer) distributor(c *gin.Context) {
	log.DEBUG.Println("ReverseProxyServer - Request Host : ", c.Request.Host)
	log.DEBUG.Println("ReverseProxyServer GetClientIP: ", c.ClientIP())

	//log.TempLog().Println("ReverseProxyServer Header: ",pretty.Sprint(c.Request.Header))
	req := c.Request
	if req.Host == "" {
		c.AbortWithStatusJSON(
			http.StatusForbidden,
			gin.H{
				"err": "Hostname is empty",
			})
		log.DEBUG.Println("Hostname is empty")
		return
	}
	sph := strings.Split(req.Host, ".")
	subHost := ""
	if len(sph) > 1 {
		subHost = sph[0]
	}

	realAddrURL, exist := redirectAddrMap[subHost]
	if !exist {
		c.AbortWithStatus(http.StatusNotFound)
		log.DEBUG.Println("SubHost not exist : ", subHost)
		return
	}
	fmt.Println("redirect to realAddr : ", realAddrURL)
	proxy := httputil.NewSingleHostReverseProxy(realAddrURL)
	proxy.Director = func(req *http.Request) {
		req.URL.Scheme = realAddrURL.Scheme
		req.URL.Host = realAddrURL.Host

		req.Header.Add(util.IDENTIFICATION_SYMBOL_KEY, util.IDENTIFICATION_SYMBOL)
		req.Header.Set(util.ClientIP_KEY, c.ClientIP())
	}
	proxy.ModifyResponse = func(response *http.Response) error {
		response.Header.Del(util.IDENTIFICATION_SYMBOL_KEY)
		response.Header.Del(util.ClientIP_KEY)
		return nil
	}
	proxy.ServeHTTP(c.Writer, c.Request)
	return
}

//func (s *ReverseProxyServer) AfterServerStartup() {
//	log.Println("ReverseProxyServer Start up .")
//	//s.runFuncsNeedToExecutedAfterServerStartup()
//}

func (r *ReverseProxyServer) shutdown() {
	log.DEBUG.Println("ReverseProxyServer - shutdown")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := r.serverInst.Shutdown(ctx); err != nil {
		log.ERROR.Println("ReverseProxyServer fail to shutdown:", err)
	} else {
		log.DEBUG.Println("ReverseProxyServer closed.")
	}

}
