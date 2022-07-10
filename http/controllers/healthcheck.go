package controllers

import (
	"cloudgobrrr/backend/http/binding"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HttpHealthcheck(c *gin.Context) {
	c.JSON(http.StatusOK, binding.ResEmpty)
}
