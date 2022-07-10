package controllers

import (
	"cloudgobrrr/backend/http/binding"
	"cloudgobrrr/backend/pkg/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HttpFilesList(c *gin.Context) {
	var query binding.ReqFilesList
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, binding.ResErrorInvalidRequest)
		return
	}

	path, err := helpers.GetAndCheckPath(c.MustGet("userName").(string), query.Path)
	if err != nil {
		c.JSON(http.StatusBadRequest, binding.ResErrorInvalidPath)
		return
	}

	files, err := helpers.FilesList(path)
	if err != nil {
		c.JSON(http.StatusNotFound, binding.ResErrorNotFound)
		return
	}
	c.JSON(http.StatusOK, binding.ResFilesList{Files: files})
}
