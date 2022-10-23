package api_server

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kr/pretty"
	"rabbit_gather/src/auth/token/claims"
	"rabbit_gather/src/server"

	//cors "github.com/rs/cors/wrapper/gin"
	"net/http"
	"net/url"
	"rabbit_gather/src/auth/bitmask"
	"rabbit_gather/src/logger"
	"rabbit_gather/src/service/account_management"
	"rabbit_gather/src/service/article_management"
	"rabbit_gather/util"
)

var ServePath *url.URL
var log = logger.NewLoggerWrapper("APIServer")

func init() {
	type Config struct {
		ServePath string
	}
	var config Config
	err := util.ParseFileJsonConfig(&config, "config/api_server.config.json")
	if err != nil {
		panic(err.Error())
	}
	ServePath, err = url.Parse(config.ServePath)
	if err != nil {
		panic(err.Error())
	}
}

var corsHandler gin.HandlerFunc

func init() {
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:8080", "https://meowalien.com:443"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "token"}
	corsHandler = cors.New(config)
}

// The APIServer provide all restful API, Websocket APIs.
type APIServer struct {
	serverInst              *http.Server
	ginEngine               *gin.Engine
	shutdownCallbackMethods []func() error
}

func (w *APIServer) Startup(ctx context.Context) error {
	log.DEBUG.Println("APIServer listen on : ", ServePath.String())
	w.shutdownCallbackMethods = []func() error{}
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
			log.ERROR.Println(err)
		}
	}()
	log.DEBUG.Println("APIServer Started .")
	return nil
}

func (w *APIServer) Shutdown() error {
	for i := len(w.shutdownCallbackMethods) - 1; i >= 0; i-- {
		err := w.shutdownCallbackMethods[i]()
		if err != nil {
			log.ERROR.Println(err.Error())
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), server.ShutdownWaitTime)
	defer cancel()

	if err := w.serverInst.Shutdown(ctx); err != nil {
		return fmt.Errorf("APIServer fail to Shutdown: %s", err.Error())
	} else {
		log.DEBUG.Println("APIServer closed.")
		return nil
	}
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w *APIServer) debugLogger(c *gin.Context) {
	var log = log.TempLog()
	log.Println("Request: ", pretty.Sprint(c.Request.Body))
	log.Println("Method: ", pretty.Sprint(c.Request.Method))
	log.Println("ContentType: ", pretty.Sprint(c.ContentType()))
	log.Println("Header: ", pretty.Sprint(c.Request.Header))

	blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
	c.Writer = blw
	c.Next()
	log.Println("--- AfterHandle ---")
	log.Println("Response Body: ", pretty.Sprint(blw.body.String()))
	log.Println("Response Header: ", pretty.Sprint(c.Writer.Header()))
	//log.Println("ContentType: ",pretty.Sprint(c.ContentType()))
	//log.Println("Header: ",pretty.Sprint(c.Request.Header))
}
func (w *APIServer) MountService(ctx context.Context) {

	w.ginEngine.Use(gin.Recovery())
	//w.ginEngine.Use(w.debugLogger)
	//w.useMiddleware(w.permissionCheck)
	w.ginEngine.Use(corsHandler)

	userAccount := account_management.AccountManagement{}
	w.ginEngine.POST("/account/login", w.permissionCheckHandler(bitmask.NoStatus), userAccount.LoginHandler)
	w.ginEngine.POST("/account/logout", w.permissionCheckHandler(bitmask.NoStatus), userAccount.LogoutHandler)
	w.ginEngine.POST("/account/signup", w.permissionCheckHandler(bitmask.WaitVerificationCode), userAccount.SignupHandler)
	w.ginEngine.POST("/account/sent_verification_code", w.permissionCheckHandler(bitmask.NoStatus), userAccount.SentVerificationCodeHandler)
	w.appendShutdownCallback(func() error {
		return userAccount.Close()
	})

	articleManagement := article_management.ArticleManagement{}
	// 修改文章設定
	w.ginEngine.POST("/article/settings", w.permissionCheckHandler(bitmask.NoStatus), articleManagement.UpdateAuthorityHandler)
	// 查詢文章設定
	w.ginEngine.GET("/article/settings", w.permissionCheckHandler(bitmask.NoStatus), articleManagement.AskAuthorityHandler)
	// 新增文章
	w.ginEngine.POST("/article/new", w.permissionCheckHandler(bitmask.NoStatus), articleManagement.PostArticleHandler)
	// 搜尋文章
	w.ginEngine.GET("/article/search", w.permissionCheckHandler(bitmask.NoStatus), articleManagement.SearchArticleHandler)
	// 取得指定文章
	w.ginEngine.GET("/article/:id", w.permissionCheckHandler(bitmask.NoStatus), articleManagement.GetArticleHandler)
	// 監聽文章變更 - 連線後接收請求監聽某文章狀態（不可更詳細，因為不信任請求內容），設定最大監聽量
	w.ginEngine.GET("/article/listen", w.permissionCheckHandler(bitmask.NoStatus), articleManagement.ListenArticleChangeHandler)
	// 刪除文章
	w.ginEngine.DELETE("/article/:id", w.permissionCheckHandler(bitmask.NoStatus), articleManagement.DeleteArticleHandler)
	// 更新文章
	w.ginEngine.PATCH("/article/:id", w.permissionCheckHandler(bitmask.NoStatus), articleManagement.UpdateArticleHandler)
	w.appendShutdownCallback(func() error {
		return articleManagement.Close()
	})
}

// check if client can access the API
func (w *APIServer) permissionCheckHandler(statusBitmask bitmask.StatusBitmask) func(c *gin.Context) {
	return func(c *gin.Context) {
		fullPath := c.FullPath()

		// this should not happen
		if statusBitmask == bitmask.Reject {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"err": "this API is temporarily closed"})
			log.ERROR.Printf("Got Reject on: %s", fullPath)
			return
		}
		// if the path doesn't need status to access
		if statusBitmask == bitmask.NoStatus {
			c.Next()
			return
		}

		utilityClaims, err := server.ContextAnalyzer(c).GetUtilityClaim()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": "server error"})
			log.DEBUG.Println("error when GetUtilityClaim: ", err.Error())
			return
		}

		var statusClaims claims.StatusClaim
		err = utilityClaims.GetSubClaims(&statusClaims)
		if err != nil {
			if err == claims.NoSuchClaimsError {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"err": "status claims not exist"})
				log.DEBUG.Println("status claims not exist")
				return
			} else if err == claims.UnknownClaimsError {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": "server error"})
				log.ERROR.Println("UnknownClaimsError: ", err.Error())
				return
			} else {
				panic(err.Error())
			}

		}

		if !bitmask.MaskCheck(statusBitmask, statusClaims.StatusBitmask) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"err": "you do not have permission to access this api"})
			log.DEBUG.Printf("reject access to: %s\nClaims: %s", fullPath, pretty.Sprint(utilityClaims))
			return
		}
		c.Next()
	}
}

type HttpMethod string

const (
	GET     HttpMethod = "GET"
	POST    HttpMethod = "POST"
	PATCH   HttpMethod = "PATCH"
	DELETE  HttpMethod = "DELETE"
	PUT     HttpMethod = "PUT"
	OPTIONS HttpMethod = "OPTIONS"
)

var pathStatusRequirementMap = map[HttpMethod]map[string]bitmask.StatusBitmask{
	GET:     {},
	POST:    {},
	PUT:     {},
	PATCH:   {},
	DELETE:  {},
	OPTIONS: {},
}

func (w *APIServer) getStatusRequirement(method HttpMethod, path string) (bitmask.StatusBitmask, bool) {
	switch method {
	case GET:
	case POST:
	case PUT:
	case PATCH:
	case DELETE:
	case OPTIONS:
	default:
		log.ERROR.Println("Not supported http method: ", method)
		return bitmask.Reject, true
	}
	p, exist := pathStatusRequirementMap[method][path]
	return p, exist
}

func (w *APIServer) appendShutdownCallback(f func() error) {
	w.shutdownCallbackMethods = append(w.shutdownCallbackMethods, f)
}
