package controllers

import (
	"cloudgobrrr/backend/database/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type bindingSessionChangeDescription struct {
	Id             uint   `json:"sessionId" binding:"required"`
	NewDescription string `json:"newDescription" binding:"required"`
}

func HttpSessionChangeDescription(c *gin.Context) {
	var json bindingSessionChangeDescription
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "invalid request"})
		return
	}

	oldDescription, err := model.SessionChangeDescription(json.Id, json.NewDescription)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "oldDescription": oldDescription})
}

type bindingSessionCreateBasicAuth struct {
	Description string `json:"description" binding:"required"`
}

func HttpSessionCreateBasicAuth(c *gin.Context) {
	var json bindingSessionCreateBasicAuth
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "invalid request"})
		return
	}

	user, err := model.UserGetByID(c.MustGet("userID").(uint))
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
