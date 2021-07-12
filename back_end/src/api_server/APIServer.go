package api_server

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kr/pretty"
	"log"
	"net/http"
	"net/url"
	"rabbit_gather/src/auth"
	//"rabbit_gather/src/handler"
	"rabbit_gather/src/service"
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

var corsHandler gin.HandlerFunc

func init() {
	// - No origin allowed by default
	// - GET,POST, PUT, HEAD methods
	// - Credentials share disabled
	// - Preflight requests cached for 12 hours
	config := cors.DefaultConfig()
	//config.AllowOrigins = []string{"http://localhost:8080"}
	//config.AllowMethods = []string{"POST"}
	config.AllowAllOrigins = true
	//config.AllowOrigins = []string{"http://localhost:8081","https://localhost:8081"}
	config.AllowMethods = []string{"GET", "POST"}
	corsHandler = cors.New(config)
}

func (w *APIServer) MountService(ctx context.Context) {

	w.useMiddleware(gin.Recovery())
	w.useMiddleware(w.permissionCheck)
	w.useMiddleware(corsHandler)

	userAccount := service.AccountManagement{}
	w.HandlePost("/login", userAccount.LoginHandler, auth.Public)
	w.HandlePost("/signup", userAccount.SignupHandler, auth.Public)

	articleManagement := service.ArticleManagement{}
	w.HandlePost("/post_article", articleManagement.PostArticleHandler, auth.Login)
	w.HandlePost("/search_article", articleManagement.SearchArticleHandler, auth.Login)

	peerService := service.PeerService{}
	w.HandleGet("/peerjs/id", peerService.GetPeerIDHandler, auth.Login)
	w.HandleGet("/peerjs", peerService.PeerWebsocketHandler, auth.Login)
}

func (w *APIServer) useMiddleware(middleware func(c *gin.Context)) {
	w.ginEngine.Use(middleware)
}

var pathPermissionMapPost = map[string]auth.APIPermissionBitmask{}
var pathPermissionMapGet = map[string]auth.APIPermissionBitmask{}

func (w *APIServer) HandlePost(path string, handler func(c *gin.Context), permissionCode auth.APIPermissionBitmask) {
	pathPermissionMapPost[path] = permissionCode
	w.ginEngine.POST(path, handler)
}

func (w *APIServer) HandleGet(path string, handler func(c *gin.Context), permissionCode auth.APIPermissionBitmask) {
	pathPermissionMapGet[path] = permissionCode
	w.ginEngine.GET(path, handler)
}

func (w *APIServer) permissionCheck(c *gin.Context) {
	fillPath := c.FullPath()
	if fillPath == "" {
		log.Println("permissionCheck - c.FullPath() is empty: ")
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	if !needPermissionCheck(c.Request.Method, fillPath) {
		c.Next()
		return
	}
	tokenRawString := c.GetHeader("token")
	if tokenRawString == "" {
		c.AbortWithStatus(http.StatusForbidden)
		log.Println("permissionCheck - Permission needed request with no token: ", pretty.Sprint(c.Request))
		return
	}
	var uc auth.PermissionClaims
	token, err := auth.ParseToken(tokenRawString, &uc)
	if err != nil {
		c.AbortWithStatus(http.StatusForbidden)
		log.Printf("searchArticleHandler - ParseToken error : %s", err.Error())
		return
	}
	if !token.Valid {
		c.AbortWithStatus(http.StatusForbidden)
		log.Printf("searchArticleHandler - token not Valid: %s", pretty.Sprint(*token))
		return
	}

	userClaims, ok := token.Claims.(*auth.PermissionClaims)
	if !ok {
		c.AbortWithStatus(http.StatusForbidden)
		log.Printf("searchArticleHandler - token Claims not type of PermissionClaims: %s", pretty.Sprint(*userClaims))
		return
	}
	APIPermissionCode := getAPIPermissionCodeWithServePath(fillPath)
	if !auth.APIAuthorizationCheck(APIPermissionCode, userClaims.APIPermissionBitmask) {
		c.AbortWithStatus(http.StatusForbidden)
		log.Printf("searchArticleHandler - this user don't have permission to access this api: %s", pretty.Sprint(*userClaims))
		return
	}
	c.Next()
	return
}

func getAPIPermissionCodeWithServePath(path string) auth.APIPermissionBitmask {
	return pathPermissionMapPost[path]
}

func needPermissionCheck(method, path string) bool {
	exist := false
	switch method {
	//an empty string means GET.
	case "":
		_, exist = pathPermissionMapGet[path]
	case "GET":
		_, exist = pathPermissionMapGet[path]
	case "POST":
		_, exist = pathPermissionMapPost[path]
	default:
		panic("unexpected method: " + method)
	}
	return exist
}
