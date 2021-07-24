package api_server

import (
	"context"
	"crypto/tls"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kr/pretty"
	"rabbit_gather/src/auth/claims"
	"rabbit_gather/src/auth/status_bitmask"
	"rabbit_gather/src/auth/token"
	"rabbit_gather/src/service/account_management"
	"rabbit_gather/src/service/article_management"
	"rabbit_gather/src/service/peer"

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
	//if util.DebugMode{
	//	config.AllowAllOrigins = true // debug
	//}
	config.AllowOrigins = []string{"http://meowalien.com:443"}
	config.AllowMethods = []string{"GET", "POST"}
	corsHandler = cors.New(config)
}

func (w *APIServer) MountService(ctx context.Context) {

	w.useMiddleware(gin.Recovery())
	w.useMiddleware(w.permissionCheck)
	w.useMiddleware(corsHandler)

	userAccount := account_management.AccountManagement{}

	w.HandlePost("/login", userAccount.LoginHandler, status_bitmask.NoStatus)
	w.HandlePost("/signup", userAccount.SignupHandler, status_bitmask.WaitVerificationCode)
	w.HandlePost("/sent_verification_code", userAccount.SentVerificationCodeHandler, status_bitmask.NoStatus)

	articleManagement := article_management.ArticleManagement{}
	w.HandlePost("/post_article", articleManagement.PostArticleHandler, status_bitmask.Login)
	w.HandlePost("/search_article", articleManagement.SearchArticleHandler, status_bitmask.Login)

	peerService := peer.PeerService{}
	w.HandleGet("/peerjs/id", peerService.GetPeerIDHandler, status_bitmask.Login)
	w.HandleGet("/peerjs", peerService.PeerWebsocketHandler, status_bitmask.Login)
}

func (w *APIServer) useMiddleware(middleware func(c *gin.Context)) {
	w.ginEngine.Use(middleware)
}

func (w *APIServer) HandlePost(path string, handler gin.HandlerFunc, permissionCode status_bitmask.StatusBitmask) {
	pathStatusRequirementMap[POST][path] = permissionCode
	w.ginEngine.POST(path, handler)
}

func (w *APIServer) HandleGet(path string, handler func(c *gin.Context), permissionCode status_bitmask.StatusBitmask) {
	pathStatusRequirementMap[GET][path] = permissionCode
	w.ginEngine.GET(path, handler)
}

// check if client can access the API
func (w *APIServer) permissionCheck(c *gin.Context) {
	fullPath := c.FullPath()
	statusBitmask, exist := w.getStatusRequirement(HttpMethod(c.Request.Method), fullPath)
	if !exist {
		panic("Error when getStatusRequirement")
	}
	// this should not happen
	if statusBitmask == status_bitmask.Reject {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"err": "something got wrong."})
		log.ERROR.Printf("Got Reject when getStatusRequirement")
		return
	}
	// if the path doesn't need status to access
	if statusBitmask == status_bitmask.NoStatus {
		c.Next()
		return
	}

	tokenRawString := c.GetHeader(util.TokenHeaderKey)
	if tokenRawString == "" {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"err": "no token"})
		log.DEBUG.Println("no token: ", pretty.Sprint(c.Request.Header))
		return
	}

	//var utilityClaims = claims.UtilityClaims{}
	utilityClaims, err := token.ParseToken(tokenRawString)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"err": "token not valid"})
		log.DEBUG.Printf("ParseToken error : %s", err.Error())
		return
	}
	statusClaims, exist := utilityClaims.GetSubClaims(claims.StatusClaimsName)
	if !exist {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"err": "status claims not exist"})
		log.DEBUG.Println("status claims not exist")
		return
	}

	statusClaimsAfterCast, ok := statusClaims.(claims.StatusClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"err": "Server error",
		})
		log.ERROR.Println("error when cast StatusClaims")
		return
	}
	if !status_bitmask.BitMaskCheck(statusBitmask, statusClaimsAfterCast.StatusBitmask) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"err": "you do not have permission to access this api"})
		log.DEBUG.Printf("reject access to: %s\nClaims: %s", fullPath, pretty.Sprint(utilityClaims))
		return
	}
	c.Next()
	return
}

type HttpMethod string

const (
	GET    HttpMethod = "GET"
	POST   HttpMethod = "POST"
	PATCH  HttpMethod = "PATCH"
	DELETE HttpMethod = "DELETE"
	PUT    HttpMethod = "PUT"
)

var pathStatusRequirementMap = map[HttpMethod]map[string]status_bitmask.StatusBitmask{
	GET:    {},
	POST:   {},
	PUT:    {},
	PATCH:  {},
	DELETE: {},
}

//var pathStatusRequirementMap_RWMutex = sync.RWMutex{}

func (w *APIServer) getStatusRequirement(method HttpMethod, path string) (status_bitmask.StatusBitmask, bool) {
	switch method {
	case GET:
	case POST:
	case PUT:
	case PATCH:
	case DELETE:
	default:
		log.ERROR.Println("Not supported http method")
		return status_bitmask.Reject, true
	}
	p, exist := pathStatusRequirementMap[method][path]
	return p, exist
}

//func getAPIStatusRequiredWithPath(path string) auth.StatusBitmask {
//	return pathPermissionMapPost[path]
//}

// Check if this path need status to access
//func needPermissionCheck(method, path string) bool {
//	exist := false
//	var status auth.StatusBitmask
//	switch method {
//	//an empty string means GET.
//	case "":
//		permissionBitMask, exist = pathPermissionMapGet[path]
//	case "GET":
//		permissionBitMask, exist = pathPermissionMapGet[path]
//	case "POST":
//		permissionBitMask, exist = pathPermissionMapPost[path]
//	default:
//		log.ERROR.Println("unexpected method: " + method)
//		return false
//	}
//	if !exist {
//		//log.ERROR.Println("")
//		return false
//		//log.TempLog("method, path: %s,%s",method, path)
//		//log.TempLog("permissionBitMask: %d ",permissionBitMask)
//		//return permissionBitMask != auth.AllStatus
//	} else {
//		return false
//	}
//}