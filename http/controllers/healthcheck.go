package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HttpHealthcheck(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{"status": "ok"})
}
