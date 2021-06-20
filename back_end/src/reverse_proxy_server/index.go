package reverse_proxy_server

import (
	"bufio"
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"rabbit_gather/util"
	"strings"
	"time"
)

type ReverseProxyServer struct {
	serverInst *http.Server
	ginEngine                             *gin.Engine
}



var CERT_FILE string
var KEY_FILE string
var ServePath *url.URL
var SSLCertificationsCrts []string
var redirectAddrMap map[string]*url.URL

func init() {
	type Config struct {
		CERT_FILE string
		KEY_FILE string
		ServePath string
		MEOWALIEN_COM_CRT string
		SectigoRSADomainValidationSecureServerCA string
		USERTrustRSAAddTrustCA string
		AddTrustExternalCARoot string
		RedirectAddrMap map[string]string

	}
	var config Config
	err := util.ParseJsonConfic(&config,"config/reverse_proxy_server.config.json")
	if err != nil {
		panic(err.Error())
	}

	SSLCertificationsCrts = []string{
		config.MEOWALIEN_COM_CRT,
		config.SectigoRSADomainValidationSecureServerCA,
		config.USERTrustRSAAddTrustCA,
		config.AddTrustExternalCARoot,
	}
	ServePath , err = url.Parse(config.ServePath)
	if err!=nil{panic(err.Error())}
	CERT_FILE = config.CERT_FILE
	KEY_FILE = config.KEY_FILE
	redirectAddrMap = map[string]*url.URL{}
	for s, s2 := range config.RedirectAddrMap {
		u ,e := url.Parse(s2)
		if e != nil {
			panic(e.Error())
		}
		redirectAddrMap [s] =  u
	}

}

func (r *ReverseProxyServer) Startup(ctx context.Context , shutdownCallback util.ShutdownCallback)error {
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
	// 分配器掛載在根路由，轉發任何種類的請求
	r.ginEngine.Use(r.distributor)

	r.serverInst = &http.Server{
		Addr:    ":" +  ServePath.Port(),
		Handler: r.ginEngine,
		//ErrorLog: s.logger,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
			ClientCAs:  crtPool,
			ClientAuth: tls.VerifyClientCertIfGiven,
		},
	}

	go func() {
		if err := r.serverInst.ListenAndServeTLS(CERT_FILE, KEY_FILE); err != nil && err != http.ErrServerClosed {
			log.Printf("ReverseProxyServer Error: %s\n", err.Error())
		}
	}()

	//r.AfterServerStartup()
	log.Println("ReverseProxyServer Started .")
	return nil
}

// 分配請求的主要邏輯
func (s *ReverseProxyServer) distributor (c *gin.Context) {
	fmt.Println("ReverseProxyServer - Request Host : ",c.Request.Host)
	req := c.Request
	if req.Host == ""{
		c.AbortWithStatus(http.StatusForbidden)
		log.Println("Hostname is empty")
		return
	}
	sph := strings.Split(req.Host,".")
	subHost:=""
	if len(sph) > 1 {
		subHost = sph[0]
	}

	realAddrURL , exist := redirectAddrMap[subHost]
	if !exist{
		c.AbortWithStatus(http.StatusNotFound)
		log.Println("SubHost not exist : ", subHost)
		return
	}
	//fmt.Println("realAddrURL : ",realAddrURL)

	// 用以認證是由此轉發的請求
	req.Header.Add(util.IDENTIFICATION_SYMBOL_KEY, util.IDENTIFICATION_SYMBOL)
	// 改變請求對象為目標伺服器位置
	req.URL.Scheme = realAddrURL.Scheme
	req.URL.Host = realAddrURL.Host
	transport := http.DefaultTransport
	resp, err := transport.RoundTrip(req)
	if err != nil {
		log.Printf("error in roundtrip: %v", err)
		c.String(500, "error")
		return
	}
	//resp.StatusCode
	//fmt.Println("resp.StatusCode : ",resp.StatusCode)
	c.Status(resp.StatusCode)
	resp.Header.Del(util.IDENTIFICATION_SYMBOL_KEY)

	for k, vv := range resp.Header {
		for _, v := range vv {
			c.Header(k, v)
		}
	}
	if resp.StatusCode >= 400 {
		c.AbortWithStatus(resp.StatusCode)
		return
	}
	defer resp.Body.Close()
	i, err := bufio.NewReader(resp.Body).WriteTo(c.Writer)
	if err != nil {
		if err == http.ErrBodyNotAllowed{
			c.Writer.WriteHeader(resp.StatusCode)
		}else{
			panic(fmt.Sprintf("i : %d Error : %s", i, err.Error()))
		}
	}
	return
}
//func (s *ReverseProxyServer) AfterServerStartup() {
//	log.Println("ReverseProxyServer Start up .")
//	//s.runFuncsNeedToExecutedAfterServerStartup()
//}

func (r *ReverseProxyServer) shutdown() {
	fmt.Println("ReverseProxyServer - shutdown")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := r.serverInst.Shutdown(ctx); err != nil {
		log.Println("ReverseProxyServer fail to shutdown:", err)
	}else{
		log.Println("ReverseProxyServer closed.")
	}

}
