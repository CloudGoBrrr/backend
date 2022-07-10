package controllers

import (
	"cloudgobrrr/backend/database"
	"cloudgobrrr/backend/database/model"
	"cloudgobrrr/backend/http/binding"
	"cloudgobrrr/backend/pkg/helpers"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func HttpFileDownloadCreateSecret(c *gin.Context) {
	var json binding.ReqFileDownloadCreateSecret
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, binding.ResErrorInvalidRequest)
		return
	}

	path, err := helpers.GetAndCheckPath(c.MustGet("userName").(string), json.Path)
	if err != nil {
		c.JSON(http.StatusBadRequest, binding.ResErrorInvalidPath)
		return
	}

	if !helpers.FileExists(path + "/" + json.Name) {
		c.JSON(http.StatusNotFound, binding.ResErrorNotFound)
		return
	}

	secret, err := model.DownloadSecretCreate(c.MustGet("userID").(uint), path, json.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, binding.ResErrorInternalServerError)
		return
	}

	c.JSON(http.StatusOK, binding.ResFileDownloadCreateSecret{Secret: secret})
}

func HttpFileDownloadWithSecret(c *gin.Context) {
	var query binding.ReqFileDownloadWithSecret
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, binding.ResErrorInvalidRequest)
		return
	}

	downloadSecret, err := model.DownloadSecretGetBySecret(query.Secret)
	if err != nil {
		c.JSON(http.StatusNotFound, binding.ResErrorNotFound)
		return
	}

	// delete used secret
	if err := database.GetDB().Delete(&downloadSecret).Error; err != nil {
		c.JSON(http.StatusInternalServerError, binding.ResErrorInternalServerError)
		return
	}

	filePath := filepath.Join(downloadSecret.Path, downloadSecret.Filename)
	c.Writer.Header().Set("Content-Type", "application/octet-stream")
	c.Writer.Header().Set("Content-Disposition", "attachment; filename=\""+downloadSecret.Filename+"\"")
	c.File(filePath)
}
