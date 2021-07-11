package service

import (
	"github.com/gin-gonic/gin"
	"rabbit_gather/src/handler"
)

type Service interface {
	GetHandler(handlerName handler.HandlerNames) gin.HandlerFunc
}
