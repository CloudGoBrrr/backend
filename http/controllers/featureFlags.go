package controllers

import (
	"cloudgobrrr/backend/pkg/env"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func HttpFeatureFlag(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{"status": "ok", "version": env.GetVersion(), "featureFlags": gin.H{"PUBLIC_REGISTRATION": os.Getenv("PUBLIC_REGISTRATION")}})
}
