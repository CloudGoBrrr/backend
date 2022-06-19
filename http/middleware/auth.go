package middleware

import (
	"cloudgobrrr/backend/database/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthenticateToken(c *gin.Context) {
	header := c.Request.Header.Get("Authorization")

	if header == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		c.Abort()
		return
	}

	user, authToken, err := model.GetUserFromAuthToken(header)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		c.Abort()
		return
	}

	c.Set("userID", user.ID)
	c.Set("userName", user.Username)
	c.Set("userEmail", user.Email)
	c.Set("tokenID", authToken.ID)

	c.Next()
}

func AuthenticateBasic(c *gin.Context) {
	if !ValidateBasic(c) {
		return
	}

	c.Next()
}

func ValidateBasic(c *gin.Context) bool {
	username, password, hasAuth := c.Request.BasicAuth()
	if !hasAuth {
		c.Writer.Header().Set("WWW-Authenticate", "Basic realm=\"User Visible Realm\", charset=\"UTF-8\"")
		c.AbortWithStatus(401)
		return false
	}

	user, _, err := model.GetUserFromBasicAuth(username, password)
	if err != nil {
		c.AbortWithStatus(401)
		return false
	}

	c.Set("userID", user.ID)
	c.Set("userName", user.Username)
	c.Set("userEmail", user.Email)

	return true
}
