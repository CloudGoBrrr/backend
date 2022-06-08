package controllers

import (
	"cloudgobrrr/backend/pkg/helpers"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type bindingFilesList struct {
	Path string `form:"path" binding:"required"`
}

func HttpFilesList(c *gin.Context) {
	var query bindingFilesList
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "invalid request"})
		return
	}

	path, err := helpers.GetAndCheckPath(c.MustGet("userName").(string), query.Path)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err.Error()})
		return
	}

	files, err := helpers.ListFiles(path)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "error": "path not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok", "files": files})
}

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

type bindingFileDelete struct {
	Path string `form:"path" binding:"required"`
	Name string `form:"name" binding:"required"`
}

func HttpFileDelete(c *gin.Context) {
	var query bindingFileDelete
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "invalid request"})
		return
	}

	path, err := helpers.GetAndCheckPath(c.MustGet("userName").(string), query.Path)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err.Error()})
		return
	}

	removePath := filepath.Join(path, query.Name)

	if !helpers.FileExists(removePath) {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "error": "file not found"})
		return
	}

	if err := os.RemoveAll(removePath); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "an error occured while deleting file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
