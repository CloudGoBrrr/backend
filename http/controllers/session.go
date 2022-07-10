package controllers

import (
	"cloudgobrrr/backend/database/model"
	"cloudgobrrr/backend/http/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HttpSessionChangeDescription(c *gin.Context) {
	var req json.ReqSessionChangeDescription
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, json.ResErrorInvalidRequest)
		return
	}

	oldDescription, err := model.SessionChangeDescription(req.Id, req.NewDescription)
	if err != nil {
		c.JSON(http.StatusInternalServerError, json.ResErrorInternalServerError)
		return
	}

	c.JSON(http.StatusOK, json.ResSessionChangeDescription{OldDescription: oldDescription})
}

func HttpSessionCreateBasicAuth(c *gin.Context) {
	var req json.ReqSessionCreateBasicAuth
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, json.ResErrorInvalidRequest)
		return
	}

	user, err := model.UserGetByID(c.MustGet("userID").(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, json.ResErrorInternalServerError)
		return
	}

	basicPassword, err := model.SessionCreateBasic(user, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, json.ResErrorInternalServerError)
		return
	}

	c.JSON(http.StatusOK, json.ResSessionCreateBasicAuth{Username: user.Username, Password: basicPassword})
}

func HttpSessionList(c *gin.Context) {
	sessions, err := model.SessionGetAll(c.MustGet("userID").(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "sessions": sessions})
}

type bindingSessionDelete struct {
	ID uint `form:"id" binding:"required"`
}

func HttpSessionDeleteWithID(c *gin.Context) {
	var query bindingSessionDelete
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
