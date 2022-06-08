package middleware

import "github.com/gin-gonic/gin"

func DefaultHeader(c *gin.Context) {
	c.Writer.Header().Set("Server", "CloudGoBrrr")
	c.Next()
}
