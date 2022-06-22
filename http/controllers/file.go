package controllers

import (
	"cloudgobrrr/backend/pkg/helpers"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

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
