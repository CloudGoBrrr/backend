package controllers

import (
	"cloudgobrrr/backend/http/binding"
	"cloudgobrrr/backend/pkg/env"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func HttpFeatureFlagGet(c *gin.Context) {
	c.JSON(http.StatusOK, binding.ResFeatureFlags{Version: env.VersionGet(), FeatureFlags: binding.FeatureFlags{
		PublicRegistration: os.Getenv("PUBLIC_REGISTRATION"),
		WebDav:             os.Getenv("WEBDAV_ENABLED"),
	}})
}
