package api_server

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"net/url"
	"rabbit_gather/util"
	"time"
)

/*
靜態頁面服務器
*/
type APIServer struct {
	serverInst *http.Server
	ginEngine  *gin.Engine
}

var ServePath *url.URL

func init() {
	//ServePath, _= url.Parse( "http://127.0.0.1:2002")
	//
	type Config struct {
		ServePath string
	}
	var config Config
	err := util.ParseJsonConfic(&config, "config/api_server.config.json")
	if err != nil {
		panic(err.Error())
	}
	ServePath, err = url.Parse(config.ServePath)
	log.Println("APIServer - ServePath : ", ServePath)
	if err != nil {
		panic(err.Error())
	}
}

func (w *APIServer) Startup(ctx context.Context, shutdownCallback util.ShutdownCallback) error {
	shutdownCallback(w.shutdown)
	w.ginEngine = gin.Default()
	fmt.Println("APIServer - ServePath.String() : ", ServePath.String())
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
			fmt.Println("APIServer Error : ", err)
			fmt.Println("APIServer - w.serverInst : ", w.serverInst)
			panic("APIServer - ListenAndServe Error")
		}
	}()
	//r.AfterServerStartup()
	fmt.Println("APIServer Started .")
	return nil
}

func (w *APIServer) shutdown() {
	//fmt.Println("APIServer - shutdown")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := w.serverInst.Shutdown(ctx); err != nil {
		log.Println("APIServer fail to shutdown:", err)
	} else {
		log.Println("APIServer closed.")
	}

}

func (w *APIServer) MountService(ctx context.Context) {
	// - No origin allowed by default
	// - GET,POST, PUT, HEAD methods
	// - Credentials share disabled
	// - Preflight requests cached for 12 hours
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:8080"}
	config.AllowMethods = []string{"POST"}
	w.ginEngine.Use(cors.New(config))
	w.ginEngine.POST("/post_article", w.postArticleHandler)

}

func (w *APIServer) postArticleHandler(c *gin.Context) () {


	body := json.NewDecoder(c.Request.Body)
	body.DisallowUnknownFields()
	type PosistionStruct struct {
		Latitude  float32 `json:"latitude"`
		Longitude float32 `json:"longitude"`
	}
	articleReceived := struct {
		Title    string          `json:"title"`
		Content  string          `json:"content"`
		Position PosistionStruct `json:"position"`
	}{}

	err := body.Decode(&articleReceived)
	if err != nil {
		c.AbortWithStatus(http.StatusForbidden)
		log.Printf("postArticleHandler - body.Decode error : %s", err.Error())
		return
	}
	fmt.Println("Title : ", articleReceived.Title)
	fmt.Println("Content : ", articleReceived.Content)
	fmt.Println("Position : ", articleReceived.Position)

	//rsString , err := json.Marshal(articleReceived)
	c.JSON(200, gin.H{
		"result":articleReceived,
	})
}
