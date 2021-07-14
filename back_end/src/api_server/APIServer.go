package api_server

import (
	"context"
	"crypto/tls"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kr/pretty"
	//"log"
	"net/http"
	"net/url"
	"rabbit_gather/src/auth"
	"rabbit_gather/src/logger"

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
var log = logger.NewLoggerWrapper("APIServer")

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
	if err != nil {
		panic(err.Error())
	}
}

func (w *APIServer) Startup(ctx context.Context, shutdownCallback util.ShutdownCallback) error {
	log.DEBUG.Println("APIServer listen on : ", ServePath)

	shutdownCallback(w.shutdown)
	w.ginEngine = gin.Default()
	//w.ginEngine.TrustedProxies = append(w.ginEngine.TrustedProxies,"127.0.0.1/0" )
	//log.DEBUG.Println("APIServer - ServePath.String() : ", ServePath.String())
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
			log.DEBUG.Printf("reject direct connection from : %s", req.RemoteAddr)
			return
		}
	})
	w.MountService(ctx)
	go func() {
		if err := w.serverInst.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.ERROR.Println(err)
		}
	}()
	log.DEBUG.Println("APIServer Started .")
	return nil
}

func (w *APIServer) shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := w.serverInst.Shutdown(ctx); err != nil {
		log.ERROR.Println("APIServer fail to shutdown:", err)
	} else {
		log.DEBUG.Println("APIServer closed.")
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
	w.HandlePost("/signup", userAccount.SignupHandler, auth.WaitVerificationCode)
	w.HandlePost("/sent_verification_code", userAccount.SentVerificationCodeHandler, auth.Public)

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

func (w *APIServer) HandlePost(path string, handler gin.HandlerFunc, permissionCode auth.APIPermissionBitmask) {
	pathPermissionMapPost[path] = permissionCode
	w.ginEngine.POST(path, handler)
}

func (w *APIServer) HandleGet(path string, handler func(c *gin.Context), permissionCode auth.APIPermissionBitmask) {
	pathPermissionMapGet[path] = permissionCode
	w.ginEngine.GET(path, handler)
}

func (w *APIServer) permissionCheck(c *gin.Context) {
	fillPath := c.FullPath()
	if !needPermissionCheck(c.Request.Method, fillPath) {
		c.Next()
		return
	}

	tokenRawString := c.GetHeader(auth.TokenHeaderKey)
	if tokenRawString == "" {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"err": "no token"})
		log.DEBUG.Println("no token: ", pretty.Sprint(c.Request.Header))
		return
	}
	var uc auth.PermissionClaims
	token, err := auth.ParseToken(tokenRawString, &uc)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"err": "token not valid"})
		log.DEBUG.Printf("ParseToken error : %s", err.Error())
		return
	}
	if !token.Valid {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"err": "token not valid"})
		log.DEBUG.Printf("token not valid: %s", pretty.Sprint(*token))
		return
	}

	userClaims, ok := token.Claims.(*auth.PermissionClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": "token not valid"})
		log.ERROR.Printf("token Claims not type of PermissionClaims: %s", pretty.Sprint(*userClaims))
		return
	}

	APIPermissionCode := getAPIPermissionCodeWithServePath(fillPath)
	if !auth.BitMaskCheck(APIPermissionCode, userClaims.PermissionBitmask) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"err": "this user don't have permission to access this api"})
		log.DEBUG.Printf("this user don't have permission to access this api: %s", pretty.Sprint(*userClaims))
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
	permissionBitMask := auth.Admin
	switch method {
	//an empty string means GET.
	case "":
		permissionBitMask, exist = pathPermissionMapGet[path]
	case "GET":
		permissionBitMask, exist = pathPermissionMapGet[path]
	case "POST":
		permissionBitMask, exist = pathPermissionMapPost[path]
	default:
		log.ERROR.Println("unexpected method: " + method)
	}
	if exist {
		//log.TempLog("method, path: %s,%s",method, path)
		//log.TempLog("permissionBitMask: %d ",permissionBitMask)
		return permissionBitMask != auth.Public
	} else {
		return false
	}
}
