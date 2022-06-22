package controllers

import (
	"cloudgobrrr/backend/database"
	"cloudgobrrr/backend/database/model"
	"cloudgobrrr/backend/pkg/helpers"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type bindingAuthSignin struct {
	Username    string `json:"username"    binding:"required"`
	Password    string `json:"password"    binding:"required"`
	Description string `json:"description" binding:"required"`
}

func HttpAuthSignin(c *gin.Context) {
	var json bindingAuthSignin
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "invalid request"})
		return
	}

	user, err := model.UserGetByUsername(json.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		c.Abort()
		return
	}

	if user.Username == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "error": "invalid username or password"})
		return
	}

	if !helpers.PasswordCheckHash(json.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "error": "invalid username or password"})
		return
	}

	token, err := model.SessionCreateToken(user, json.Description)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "token": token})
}

type bindingAuthSignup struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email"    binding:"required"`
	Password string `json:"password" binding:"required"`
}

func HttpAuthSignup(c *gin.Context) {
	if os.Getenv("PUBLIC_REGISTRATION") != "true" {
		c.JSON(http.StatusForbidden, gin.H{"status": "forbidden"})
		return
	}

	var json bindingAuthSignup
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "invalid request"})
		return
	}

	var user model.User
	database.GetDB().First(&user, "username = ?", json.Username)

	if user.Username != "" {
		c.JSON(http.StatusConflict, gin.H{"status": "error", "error": "username already exists"})
		return
	}

	err := model.UserCreate(json.Username, json.Email, json.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})

}

type bindingAuthChangePassword struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required"`
}

func HttpAuthChangePassword(c *gin.Context) {
	var json bindingAuthChangePassword
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "invalid request"})
		return
	}

	user, err := model.UserGetByID(c.MustGet("userID").(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": "internal server error"})
		return
	}

	if !helpers.PasswordCheckHash(json.OldPassword, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "error": "invalid old password"})
		return
	}

	err = model.UserChangePassword(user.ID, json.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": "internal server error"})
		return
	}

	err = model.SessionDeleteAll(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func HttpAuthDetails(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"userDetails": gin.H{
			"id":        c.MustGet("userID").(uint),
			"username":  c.MustGet("userName").(string),
			"email":     c.MustGet("userEmail").(string),
			"sessionId": c.MustGet("sessionId").(uint),
		},
	})
}
