package middleware

import "github.com/gin-gonic/gin"

func DefaultHeader(c *gin.Context) {
	c.Writer.Header().Set("Server", "CloudGoBrrr")
	c.Next()
}

func ApiHeader(c *gin.Context) {
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Expires", "0")
	c.Next()
}

func PublicHeader(c *gin.Context) {
	c.Writer.Header().Set("Content-Security-Policy", "default-src 'self'; img-src 'self' data:; object-src 'none'; script-src 'self'; style-src 'self' 'unsafe-inline'; frame-ancestors 'self'; base-uri 'self'; form-action 'self';")
	c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
	c.Writer.Header().Set("X-Frame-Options", "DENY")
	c.Writer.Header().Set("X-XSS-Protection", "1; mode=block")
	c.Writer.Header().Set("Referrer-Policy", "no-referrer")
	c.Next()
}
