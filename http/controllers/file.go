package controllers

import (
	"cloudgobrrr/backend/http/binding"
	"cloudgobrrr/backend/pkg/helpers"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func HttpFileDelete(c *gin.Context) {
	var query binding.ReqFileDelete
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, binding.ResErrorInvalidRequest)
		return
	}

	path, err := helpers.GetAndCheckPath(c.MustGet("userName").(string), query.Path)
	if err != nil {
		c.JSON(http.StatusBadRequest, binding.ResErrorInvalidPath)
		return
	}

	removePath := filepath.Join(path, query.Name)

	if !helpers.FileExists(removePath) {
		c.JSON(http.StatusNotFound, binding.ResErrorNotFound)
		return
	}

	if err := os.RemoveAll(removePath); err != nil {
		c.JSON(http.StatusInternalServerError, binding.ResErrorInternalServerError)
		return
	}

	c.JSON(http.StatusOK, binding.ResEmpty)
}
