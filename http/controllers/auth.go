package controllers

import (
	"cloudgobrrr/backend/database"
	"cloudgobrrr/backend/database/model"
	"cloudgobrrr/backend/pkg/helpers"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type bindingSignin struct {
	Username    string `json:"username"    binding:"required"`
	Password    string `json:"password"    binding:"required"`
	Description string `json:"description" binding:"required"`
}

func HttpAuthSignin(c *gin.Context) {
	var json bindingSignin
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "invalid request"})
		return
	}

	user, err := model.GetUserByUsername(json.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		c.Abort()
		return
	}

	if user.Username == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "error": "invalid username or password"})
		return
	}

	if !helpers.CheckPasswordHash(json.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "error": "invalid username or password"})
		return
	}

	token, err := model.SessionCreateToken(user, json.Description)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "token": token})
}

type bindingSignup struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email"    binding:"required"`
	Password string `json:"password" binding:"required"`
}

func HttpAuthSignup(c *gin.Context) {
	if os.Getenv("PUBLIC_REGISTRATION") != "true" {
		c.JSON(http.StatusForbidden, gin.H{"status": "forbidden"})
		return
	}

	var json bindingSignup
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

	err := model.CreateUser(json.Username, json.Email, json.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})

}

type bindingChangePassword struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required"`
}

func HttpAuthChangePassword(c *gin.Context) {
	var json bindingChangePassword
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "invalid request"})
		return
	}

	user, err := model.GetUserByID(c.MustGet("userID").(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": "internal server error"})
		return
	}

	if !helpers.CheckPasswordHash(json.OldPassword, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "error": "invalid old password"})
		return
	}

	err = model.ChangePassword(user.ID, json.NewPassword)
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

type bindingCreateBasicAuth struct {
	Description string `json:"description" binding:"required"`
}

func HttpAuthCreateBasicAuth(c *gin.Context) {
	var json bindingCreateBasicAuth
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "invalid request"})
		return
	}

	user, err := model.GetUserByID(c.MustGet("userID").(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": "internal server error"})
		return
	}

	basicPassword, err := model.SessionCreateBasic(user, json.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "username": user.Username, "password": basicPassword})
}

func HttpAuthListAuthTokens(c *gin.Context) {
	tokens, err := model.SessionGetAll(c.MustGet("userID").(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "tokens": tokens})
}

type bindingDeleteAuthToken struct {
	ID uint `form:"id" binding:"required"`
}

func HttpAuthDeleteAuthTokenWithID(c *gin.Context) {
	var query bindingDeleteAuthToken
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "invalid request"})
		return
	}

	if err := model.SessionDeleteWithID(query.ID, c.MustGet("userID").(uint)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func HttpAuthCheck(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{
		"status": "ok",
		"userDetails": gin.H{
			"id":       c.MustGet("userID").(uint),
			"username": c.MustGet("userName").(string),
			"email":    c.MustGet("userEmail").(string),
			"tokenID":  c.MustGet("tokenID").(uint),
		},
	})
}
