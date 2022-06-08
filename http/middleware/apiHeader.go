package middleware

import (
	"github.com/gin-gonic/gin"
)

func ApiHeader(c *gin.Context) {
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Expires", "0")
	c.Next()
}
