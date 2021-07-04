package web_server

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/gin-contrib/static"
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
type WebServer struct {
	serverInst *http.Server
	ginEngine  *gin.Engine
}

var ServePath *url.URL

func init() {
	//ServePath, _= url.Parse( "http://127.0.0.1:2004")

	type Config struct {
		ServePath string
	}
	var config Config
	err := util.ParseJsonConfic(&config, "config/web_server.config.json")
	if err != nil {
		panic(err.Error())
	}
	ServePath, err = url.Parse(config.ServePath)
	log.Println("APIServer - ServePath : ",ServePath)
	if err != nil {
		panic(err.Error())
	}
}

func (w *WebServer) Startup(ctx context.Context, shutdownCallback util.ShutdownCallback) error {
	shutdownCallback(w.shutdown)
	w.ginEngine = gin.Default()
	fmt.Println("WebServer - ServePath.String() : ",ServePath.String())
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
		fmt.Println("WebServer Error : ", err)
		fmt.Println("WebServer - w.serverInst : ", w.serverInst)
		panic("WebServer - ListenAndServe Error")
	}
	}()
	//r.AfterServerStartup()
	fmt.Println("WebServer Started .")
	return nil
}

func (w *WebServer) shutdown() {
	//fmt.Println("APIServer - shutdown")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := w.serverInst.Shutdown(ctx); err != nil {
		log.Println("WebServer fail to shutdown:", err)
	} else {
		log.Println("WebServer closed.")
	}

}

func (w *WebServer) MountService(ctx context.Context) {
	//w.ginEngine.Static("/","public/web/")
	w.ginEngine.Use(static.Serve("/", static.LocalFile("public/web/", false)))
	w.ginEngine.LoadHTMLFiles("public/web/index.html")

	w.ginEngine.NoRoute(func(c *gin.Context) {
		//t, _ := template.ParseFiles("public/web/index.html")
		c.HTML(http.StatusOK,"index.html" ,gin.H{
			"title": "Main website",
		})
	})
}

//func (w *WebServer)  indexHandler(c *gin.Context)() {
//	fmt.Println("WebServer - c : ",c.Request)
//	//
//	//body := json.NewDecoder(c.Request.Body)
//	//body.DisallowUnknownFields()
//	//
//	//userinfo := struct {
//	//	Title string `json:"title"`
//	//	Content string `json:"content"`
//	//	Position string `json:"position"`
//	//}{}
//	//
//	//err := body.Decode(&userinfo)
//	//if err != nil {
//	//	c.AbortWithStatus(http.StatusForbidden)
//	//	log.Printf("postArticleHandler - body.Decode error : %s", err.Error())
//	//	return
//	//}
//	//fmt.Println("Title : ",userinfo.Title)
//	//fmt.Println("Content : ",userinfo.Content)
//	//fmt.Println("Position : ",userinfo.Position)
//}
