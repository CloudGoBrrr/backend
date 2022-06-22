package controllers

import (
	"cloudgobrrr/backend/pkg/helpers"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type bindingFolderCreate struct {
	Path string `json:"path" binding:"required"`
	Name string `json:"name" binding:"required"`
}

func HttpFolderCreate(c *gin.Context) {
	var json bindingFolderCreate
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "invalid request"})
		return
	}

	path, err := helpers.GetAndCheckPath(c.MustGet("userName").(string), json.Path)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err.Error()})
		return
	}

	folderPath := filepath.Join(path, json.Name)

	if helpers.FileExists(folderPath) {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "folder already exists"})
		return
	}

	if err := os.MkdirAll(folderPath, os.ModePerm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "an error occured while creating folder"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "created"})
}
