package server

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"rabbit_gather/src/auth/token"
	"rabbit_gather/src/auth/token/claims"
	"rabbit_gather/src/logger"
	"rabbit_gather/util"
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

// GetUtilityClaims will parse the claims.UtilityClaim in the header token, and cache it in Context for next time access.
func (c *analyzer) GetUtilityClaims() (*claims.UtilityClaim, error) {
	if v, exist := c.Get(UtilityClaimKey); exist {
		return v.(*claims.UtilityClaim), nil
	}

	tokenRawString := c.GetHeader(token.TokenKey)
	if tokenRawString == "" {
		return nil, errors.New("token is empty")
	}

	utilityClaims, err := token.ParseToken(tokenRawString)
	if err != nil {
		return nil, err
	}
	c.Set(UtilityClaimKey, utilityClaims)
	return utilityClaims, nil
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
	if !util.CheckIDENTIFICATION_SYMBOL(req) {
		c.AbortWithStatus(http.StatusForbidden)
		log.DEBUG.Printf("reject direct connection from : %s", req.RemoteAddr)
		return
	}
}
