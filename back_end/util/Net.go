package util

import "github.com/gin-gonic/gin"

// To get the GetClientIP that ReverseProxyServer got
func GetClientIP(c *gin.Context) string {
	return c.Request.Header.Get(ClientIP_KEY)
}
