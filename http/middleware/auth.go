package middleware

import (
	"cloudgobrrr/backend/database/model"
	"cloudgobrrr/backend/http/binding"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthenticateToken(c *gin.Context) {
	header := c.Request.Header.Get("Authorization")

	if header == "" {
		c.JSON(http.StatusUnauthorized, binding.ResErrorUnauthorized)
		c.Abort()
		return
	}

	user, session, err := model.SessionGetUserToken(header)
	if err != nil {
		c.JSON(http.StatusUnauthorized, binding.ResErrorUnauthorized)
		c.Abort()
		return
	}

	c.Set("userID", user.ID)
	c.Set("userName", user.Username)
	c.Set("userEmail", user.Email)
	c.Set("sessionId", session.ID)
	c.Set("isAdmin", user.IsAdmin)

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

	user, session, err := model.SessionGetUserBasic(username, password)
	if err != nil {
		c.AbortWithStatus(401)
		return false
	}

	c.Set("userID", user.ID)
	c.Set("userName", user.Username)
	c.Set("userEmail", user.Email)
	c.Set("sessionId", session.ID)
	c.Set("isAdmin", user.IsAdmin)

	return true
}
