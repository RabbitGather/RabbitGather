package websocket_server

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/gin-contrib/cors"
	//"github.com/gin-contrib/cors"

	//"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"net/url"
	"rabbit_gather/util"
	"time"
)

/*
WebsocketServer
*/
type WebsocketServer struct {
	serverInst *http.Server
	ginEngine  *gin.Engine
}

var ServePath *url.URL

func init() {
	type Config struct {
		ServePath string
	}
	var config Config
	err := util.ParseJsonConfic(&config, "config/websocket_server.config.json")
	if err != nil {
		panic(err.Error())
	}
	ServePath, err = url.Parse(config.ServePath)
	log.Println("WebsocketServer - ServePath : ", ServePath)
	if err != nil {
		panic(err.Error())
	}
}

func (w *WebsocketServer) Startup(ctx context.Context, shutdownCallback util.ShutdownCallback) error {
	shutdownCallback(w.shutdown)
	w.ginEngine = gin.Default()
	fmt.Println("WebsocketServer - ServePath.String() : ", ServePath.String())
	w.serverInst = &http.Server{
		Addr:    ":" + ServePath.Port(),
		Handler: w.ginEngine,
		//ErrorLog: w.logger,
		TLSConfig: &tls.Config{
			//MinVersion: tls.VersionTLS12,
			ClientAuth: tls.NoClientCert,
		},
	}
	w.ginEngine.Use(func(c *gin.Context) {
		req := c.Request
		if !util.CheckIDENTIFICATION_SYMBOL(req) {
			c.AbortWithStatus(http.StatusForbidden)
			log.Printf("reject direct connection from : %s", req.RemoteAddr)
			return
		}
	})
	w.MountService(ctx)
	go func() {
		if err := w.serverInst.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("WebsocketServer Error : ", err)
			fmt.Println("WebsocketServer - w.serverInst : ", w.serverInst)
			panic("WebsocketServer - ListenAndServe Error")
		}
	}()
	//r.AfterServerStartup()
	fmt.Println("WebsocketServer Started .")
	return nil
}

func (w *WebsocketServer) shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := w.serverInst.Shutdown(ctx); err != nil {
		log.Println("WebsocketServer fail to shutdown:", err)
	} else {
		log.Println("WebsocketServer closed.")
	}

}

func (w *WebsocketServer) MountService(ctx context.Context) {
	// - No origin allowed by default
	// - GET,POST, PUT, HEAD methods
	// - Credentials share disabled
	// - Preflight requests cached for 12 hours

	//
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	//config.AllowOrigins = []string{"http://localhost:8081","https://localhost:8081"}
	config.AllowMethods = []string{"GET"}
	w.ginEngine.Use(cors.New(config))
	//w.ginEngine.Use(w.IDGetter)
	w.ginEngine.GET("/peerjs/id", w.IDGetter)
	w.ginEngine.GET("/peerjs", w.peerHandler)
}

func (w *WebsocketServer) IDGetter(c *gin.Context) {
	//c.Writer.Write([]byte("b454f8be-1685-4ef2-9ee2-4cbb6f8a50ed"))
	fmt.Println("origin", c.Request.Header.Get("Origin"))
	//pretty.Println("Access-Control-Allow-Origin : ",c.GetHeader("Access-Control-Allow-Origin"))

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte("b454f8be-1685-4ef2-9ee2-4cbb6f8a50ed"))
	//c.Next()
}
