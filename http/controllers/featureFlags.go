package controllers

import (
	"cloudgobrrr/backend/pkg/env"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func HttpFeatureFlagGet(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok", "version": env.VersionGet(), "featureFlags": gin.H{"PUBLIC_REGISTRATION": os.Getenv("PUBLIC_REGISTRATION"), "WEBDAV_ENABLED": os.Getenv("WEBDAV_ENABLED")}})
}
