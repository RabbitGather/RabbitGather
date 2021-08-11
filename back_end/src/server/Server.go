package server

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"rabbit_gather/src/auth/token"
	"rabbit_gather/src/auth/token/claims"
	"rabbit_gather/src/logger"
	"time"
)

type analyzer struct {
	*gin.Context
}

const AnalyzerKey = "AnalyzerKey"
const UtilityClaimKey = "UtilityClaimKey"

// ContextAnalyzer will wrap gin.Context in analyzer and cache it in context for next time access.
func ContextAnalyzer(c *gin.Context) *analyzer {
	v, exist := c.Get(AnalyzerKey)
	if exist {
		return v.(*analyzer)
	}
	a := &analyzer{Context: c}
	c.Set(AnalyzerKey, a)
	return a
}

const (
	ShutdownWaitTime = 30 * time.Second
)

// GetClientIP get the Client's IP that ReverseProxyServer got
func (c *analyzer) GetClientIP() string {
	ip := c.Request.Header.Get(ClientIP_KEY)
	if ip == "" {
		panic("the ClientIP in \"\"")
	}
	return ip
}

var ErrTokenIsEmpty = errors.New("token is empty")

// GetUtilityClaim will parse the claims.UtilityClaim in the header token, and cache it in Context for next time access.
func (c *analyzer) GetUtilityClaim() (*claims.UtilityClaim, error) {
	if v, exist := c.Get(UtilityClaimKey); exist {
		return v.(*claims.UtilityClaim), nil
	}

	tokenRawString := c.GetHeader(token.TokenKey)
	if tokenRawString == "" {
		return nil, ErrTokenIsEmpty
	}

	utilityClaims, err := token.ParseToken(tokenRawString)
	if err != nil {
		return nil, err
	}
	c.Set(UtilityClaimKey, utilityClaims)
	return utilityClaims, nil
}

func (c *analyzer) ParseUserClaim() (*claims.UserClaim, error) {
	utilityClaim, err := c.GetUtilityClaim()
	if err != nil {
		return nil, err
	}
	var userClaim claims.UserClaim
	err = utilityClaim.GetSubClaims(&userClaim)
	if err != nil {
		return nil, err
	}
	return &userClaim, nil
}

type Server interface {
	// The Startup will start up the server.
	Startup(ctx context.Context) error
	// The Shutdown method will the server.
	Shutdown() error
}

type SubServer interface {
	Server
	MountService(context.Context)
}

var log = logger.NewLoggerWrapper("server")

func CheckIdentificationSymbol(c *gin.Context) {
	req := c.Request
	if !(req.Header.Get(IDENTIFICATION_SYMBOL_KEY) == IDENTIFICATION_SYMBOL) {
		c.AbortWithStatus(http.StatusForbidden)
		log.DEBUG.Printf("reject direct connection from : %s", req.RemoteAddr)
		return
	}
	c.Next()
}

// IDENTIFICATION_SYMBOL_KEY is a key to verify the request is come from ReverseProxyServer
const IDENTIFICATION_SYMBOL_KEY = "IDENTIFICATION_SYMBOL"

// IDENTIFICATION_SYMBOL is a value to verify the request is come from ReverseProxyServer
const IDENTIFICATION_SYMBOL = "fvqejfopj3/5<>?>9rm2ur#$TW 0924#$@T$#T$#^"

// ClientIP_KEY is a key to store the real client ip the ReverseProxyServer got.
const ClientIP_KEY = "fue8asodxn8fewj8snxfpei"
