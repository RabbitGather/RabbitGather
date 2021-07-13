package web_server

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	//"log"
	"net/http"
	"net/url"
	"rabbit_gather/src/logger"
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
var log = logger.NewLogger("WebServer")

func init() {
	type Config struct {
		ServePath string
	}
	var config Config
	err := util.ParseJsonConfic(&config, "config/web_server.config.json")
	if err != nil {
		panic(err.Error())
	}
	ServePath, err = url.Parse(config.ServePath)
	log.DEBUG.Println("WebServer - ServePath : ", ServePath)
	if err != nil {
		panic(err.Error())
	}
}

func (w *WebServer) Startup(ctx context.Context, shutdownCallback util.ShutdownCallback) error {
	shutdownCallback(w.shutdown)
	w.ginEngine = gin.Default()
	fmt.Println("WebServer - ServePath.String() : ", ServePath.String())
	w.serverInst = &http.Server{
		Addr:    ":" + ServePath.Port(),
		Handler: w.ginEngine,
		TLSConfig: &tls.Config{
			ClientAuth: tls.NoClientCert,
		},
	}
	w.ginEngine.Use(func(c *gin.Context) {
		req := c.Request
		if !util.CheckIDENTIFICATION_SYMBOL(req) {
			c.AbortWithStatus(http.StatusForbidden)
			log.DEBUG.Printf("reject direct connection from : %s", req.RemoteAddr)
			return
		}
	})
	w.MountService(ctx)
	go func() {
		if err := w.serverInst.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.ERROR.Println(err.Error())
		}
	}()
	log.DEBUG.Println("WebServer Started .")
	return nil
}

func (w *WebServer) shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := w.serverInst.Shutdown(ctx); err != nil {
		log.ERROR.Println("WebServer fail to shutdown:", err)
	} else {
		log.DEBUG.Println("WebServer closed.")
	}
}

func (w *WebServer) MountService(ctx context.Context) {
	//w.ginEngine.Use(w.appendPageLoadBitmask)
	w.ginEngine.Use(static.Serve("/", static.LocalFile("public/web/", false)))
	w.ginEngine.LoadHTMLFiles("public/web/index.html")

	w.ginEngine.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Main website",
		})
	})
}

//
//func (w *WebServer) appendPageLoadBitmask(c *gin.Context) {
//	//	Add the PageLoad JWT in the head
//	tokenInRequest := c.GetHeader("token")
//	var token *auth.JWTToken
//	if tokenInRequest == "" {
//		//	new token
//		var err error
//		token, err = auth.NewSignedToken(auth.PermissionClaims{
//			StandardClaims:       *auth.NewStandardClaims(),
//			APIPermissionBitmask: auth.WaitVerificationCode,
//		})
//		if err != nil {
//			c.AbortWithStatus(http.StatusInternalServerError)
//			return
//		}
//	}else{
//		//	append Bitmask
//		var claims *auth.PermissionClaims
//		var err error
//		token, err = auth.ParseToken(tokenInRequest, claims)
//		if err != nil {
//			c.AbortWithStatus(http.StatusForbidden)
//			return
//		}
//		claims = token.Claims.(*auth.PermissionClaims)
//		if auth.BitMaskCheck(claims.APIPermissionBitmask, auth.WaitVerificationCode) {
//			c.AbortWithStatus(http.StatusConflict)
//			log.Println("SentVerificationCodeHandler - GetToken error")
//			return
//		}else{
//			claims.APIPermissionBitmask = claims.APIPermissionBitmask|auth.WaitVerificationCode
//		}
//	}
//
//	c.Header(auth.TokenHeaderKey, token.GetSignedString())
//}

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
