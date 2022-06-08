package middleware

import (
	"cloudgobrrr/backend/database/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Authenticate(c *gin.Context) {
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
