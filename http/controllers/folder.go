package controllers

import (
	"cloudgobrrr/backend/http/binding"
	"cloudgobrrr/backend/pkg/helpers"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func HttpFolderCreate(c *gin.Context) {
	var json binding.ReqFolderCreate
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, binding.ResErrorInvalidRequest)
		return
	}

	path, err := helpers.GetAndCheckPath(c.MustGet("userName").(string), json.Path)
	if err != nil {
		c.JSON(http.StatusBadRequest, binding.ResErrorInvalidPath)
		return
	}

	folderPath := filepath.Join(path, json.Name)

	if helpers.FileExists(folderPath) {
		c.JSON(http.StatusBadRequest, binding.ResErrorFolderAlreadyExists)
		return
	}

	if err := os.MkdirAll(folderPath, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, binding.ResErrorInternalServerError)
		return
	}

	c.JSON(http.StatusOK, binding.ResEmpty)
}
