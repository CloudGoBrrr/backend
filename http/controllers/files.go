package controllers

import (
	"cloudgobrrr/backend/pkg/helpers"
	"net/http"

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

	files, err := helpers.FilesList(path)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "error": "path not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok", "files": files})
}
