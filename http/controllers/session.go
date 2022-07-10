package controllers

import (
	"cloudgobrrr/backend/database/model"
	"cloudgobrrr/backend/http/binding"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HttpSessionChangeDescription(c *gin.Context) {
	var req binding.ReqSessionChangeDescription
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, binding.ResErrorInvalidRequest)
		return
	}

	oldDescription, err := model.SessionChangeDescription(req.Id, req.NewDescription)
	if err != nil {
		c.JSON(http.StatusInternalServerError, binding.ResErrorInternalServerError)
		return
	}

	c.JSON(http.StatusOK, binding.ResSessionChangeDescription{OldDescription: oldDescription})
}

func HttpSessionCreateBasicAuth(c *gin.Context) {
	var req binding.ReqSessionCreateBasicAuth
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, binding.ResErrorInvalidRequest)
		return
	}

	user, err := model.UserGetByID(c.MustGet("userID").(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, binding.ResErrorInternalServerError)
		return
	}

	basicPassword, err := model.SessionCreateBasic(user, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, binding.ResErrorInternalServerError)
		return
	}

	c.JSON(http.StatusOK, binding.ResSessionCreateBasicAuth{Username: user.Username, Password: basicPassword})
}

func HttpSessionList(c *gin.Context) {
	sessions, err := model.SessionGetAll(c.MustGet("userID").(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, binding.ResErrorInternalServerError)
		return
	}

	c.JSON(http.StatusOK, binding.ResSessionList{Sessions: sessions})
}

func HttpSessionDeleteWithID(c *gin.Context) {
	var query binding.ReqSessionDeleteWithID
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, binding.ResErrorInvalidRequest)
		return
	}

	if err := model.SessionDeleteWithID(query.ID, c.MustGet("userID").(uint)); err != nil {
		c.JSON(http.StatusInternalServerError, binding.ResErrorInternalServerError)
		return
	}

	c.JSON(http.StatusOK, binding.ResEmpty)
}
