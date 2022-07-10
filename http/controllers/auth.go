package controllers

import (
	"cloudgobrrr/backend/database"
	"cloudgobrrr/backend/database/model"
	"cloudgobrrr/backend/http/binding"
	"cloudgobrrr/backend/pkg/helpers"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func HttpAuthSignin(c *gin.Context) {
	var json binding.ReqAuthSignin
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, binding.ResErrorInvalidRequest)
		return
	}

	user, err := model.UserGetByUsername(json.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, binding.ResErrorUnauthorized)
		c.Abort()
		return
	}

	if user.Username == "" {
		c.JSON(http.StatusUnauthorized, binding.ResErrorInvalidLogin)
		return
	}

	if !helpers.PasswordCheckHash(json.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, binding.ResErrorInvalidLogin)
		return
	}

	token, err := model.SessionCreateToken(user, json.Description)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, binding.ResErrorInternalServerError)
		return
	}

	c.JSON(http.StatusOK, binding.ResAuthSignin{Token: token})
}

func HttpAuthSignup(c *gin.Context) {
	if os.Getenv("PUBLIC_REGISTRATION") != "true" {
		c.JSON(http.StatusForbidden, binding.ResErrorForbidden)
		return
	}

	var json binding.ReqAuthSignup
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, binding.ResErrorInvalidRequest)
		return
	}

	var user model.User
	database.GetDB().First(&user, "username = ?", json.Username)

	if user.Username != "" {
		c.JSON(http.StatusConflict, binding.ResError{Error: "username already exists"})
		return
	}

	err := model.UserCreate(json.Username, json.Email, json.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, binding.ResErrorInternalServerError)
		return
	}

	c.JSON(http.StatusOK, binding.ResEmpty)

}

func HttpAuthChangePassword(c *gin.Context) {
	var json binding.ReqAuthChangePassword
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, binding.ResErrorInvalidRequest)
		return
	}

	user, err := model.UserGetByID(c.MustGet("userID").(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, binding.ResErrorInternalServerError)
		return
	}

	if !helpers.PasswordCheckHash(json.OldPassword, user.Password) {
		c.JSON(http.StatusUnauthorized, binding.ResErrorUnauthorized)
		return
	}

	err = model.UserChangePassword(user.ID, json.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, binding.ResErrorInternalServerError)
		return
	}

	err = model.SessionDeleteAll(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, binding.ResErrorInternalServerError)
		return
	}

	c.JSON(http.StatusOK, binding.ResEmpty)
}

func HttpAuthDetails(c *gin.Context) {
	c.JSON(http.StatusOK, binding.ResAuthDetails{
		UserDetails: binding.UserDetails{
			ID:        c.MustGet("userID").(uint),
			Username:  c.MustGet("userName").(string),
			Email:     c.MustGet("userEmail").(string),
			SessionID: c.MustGet("sessionId").(uint),
			IsAdmin:   c.MustGet("isAdmin").(bool),
		},
	})
}
