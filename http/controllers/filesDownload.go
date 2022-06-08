package controllers

import (
	"cloudgobrrr/backend/database"
	"cloudgobrrr/backend/database/model"
	"cloudgobrrr/backend/pkg/helpers"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type bindingFileDownloadCreateSecret struct {
	Path string `json:"path" binding:"required"`
	Name string `json:"name" binding:"required"`
}

func HttpFileDownloadCreateSecret(c *gin.Context) {
	var json bindingFileDownloadCreateSecret
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "invalid request"})
		return
	}

	path, err := helpers.GetAndCheckPath(c.MustGet("userName").(string), json.Path)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err.Error()})
		return
	}

	if !helpers.FileExists(path + "/" + json.Name) {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "error": "file not found"})
		return
	}

	secret, err := model.CreateDownloadSecret(c.MustGet("userID").(uint), path, json.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "secret": secret})
}

type bindingFileDownloadWithSecret struct {
	Secret string `form:"secret" binding:"required"`
}

func HttpFileDownloadWithSecret(c *gin.Context) {
	var query bindingFileDownloadWithSecret
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "invalid request"})
		return
	}

	downloadSecret, err := model.GetDownloadSecretFromSecret(query.Secret)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err.Error()})
		return
	}

	// delete used secret
	if err := database.GetDB().Delete(&downloadSecret).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "an error occured while deleting secret"})
		return
	}

	filePath := filepath.Join(downloadSecret.Path, downloadSecret.Filename)
	c.Writer.Header().Set("Content-Type", "application/octet-stream")
	c.Writer.Header().Set("Content-Disposition", "attachment; filename=\""+downloadSecret.Filename+"\"")
	c.File(filePath)
}
