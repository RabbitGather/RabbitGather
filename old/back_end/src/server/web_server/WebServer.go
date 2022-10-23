package web_server

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/gin-gonic/gin"
	"rabbit_gather/src/server"

	//"log"
	"net/http"
	"net/url"
	"rabbit_gather/src/logger"
	"rabbit_gather/util"
)

/*
靜態頁面服務器
*/
type WebServer struct {
	serverInst *http.Server
	ginEngine  *gin.Engine
}

var ServePath *url.URL
var log = logger.NewLoggerWrapper("WebServer")

func init() {
	type Config struct {
		ServePath string
	}
	var config Config
	err := util.ParseFileJsonConfig(&config, "config/web_server.config.json")
	if err != nil {
		panic(err.Error())
	}
	ServePath, err = url.Parse(config.ServePath)
	//log.DEBUG.Println("WebServer - ServePath : ", ServePath)
	if err != nil {
		panic(err.Error())
	}
}

func (w *WebServer) Startup(ctx context.Context) error {
	log.DEBUG.Println("WebServer listen on : ", ServePath.String())

	w.ginEngine = gin.Default()
	w.serverInst = &http.Server{
		Addr:    ":" + ServePath.Port(),
		Handler: w.ginEngine,
		TLSConfig: &tls.Config{
			ClientAuth: tls.NoClientCert,
		},
	}
	w.ginEngine.Use(server.CheckIdentificationSymbol)

	w.MountService(ctx)
	go func() {
		if err := w.serverInst.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.ERROR.Println(err.Error())
		}
	}()
	log.DEBUG.Println("WebServer Started .")
	return nil
}

func (w *WebServer) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), server.ShutdownWaitTime)
	defer cancel()

	if err := w.serverInst.Shutdown(ctx); err != nil {
		return fmt.Errorf("WebServer fail to Shutdown: %s", err.Error())
	} else {
		log.DEBUG.Println("WebServer closed.")
		return nil
	}
}

func (w *WebServer) MountService(ctx context.Context) {
	w.ginEngine.NoRoute(func(c *gin.Context) {
		c.File("public/web/index.html")
	})
}
