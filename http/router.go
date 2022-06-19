package http

import (
	"cloudgobrrr/backend/http/controllers"
	"cloudgobrrr/backend/http/middleware"
	"cloudgobrrr/backend/http/webdav"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

func newRouter(router *gin.Engine) {
	router.MaxMultipartMemory = 8 << 20 // 8 MiB
	router.Use(middleware.DefaultHeader)
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Content-Range", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Setup routes
	api := router.Group("api", middleware.ApiHeader)

	if os.Getenv("WEBDAV_ENABLED") == "true" {
		// WebDAV Support in /api/webdav
		router.Use(webdav.Handle())
	}

	v1 := api.Group("v1")
	{
		v1.GET("healthcheck", controllers.HttpHealthcheck)
		v1.GET("featureFlags", controllers.HttpFeatureFlag)

		v1Auth := v1.Group("auth")
		{
			v1Auth.GET("check", middleware.AuthenticateToken, controllers.HttpAuthCheck)
			v1Auth.POST("signin", controllers.HttpAuthSignin)
			v1Auth.POST("signup", controllers.HttpAuthSignup)
			v1Auth.POST("changepassword", middleware.AuthenticateToken, controllers.HttpAuthChangePassword)
			v1AuthToken := v1Auth.Group("/token")
			{
				v1AuthToken.POST("basic", middleware.AuthenticateToken, controllers.HttpAuthCreateBasicAuth)
				v1AuthToken.DELETE("", middleware.AuthenticateToken, controllers.HttpAuthDeleteAuthTokenWithID)
				v1AuthToken.GET("list", middleware.AuthenticateToken, controllers.HttpAuthListAuthTokens)
			}

		}

		v1Files := v1.Group("files")
		{
			v1Files.DELETE("", middleware.AuthenticateToken, controllers.HttpFileDelete)
			v1Files.PUT("upload", middleware.AuthenticateToken, controllers.HttpFileUpload)
			v1Files.POST("upload", middleware.AuthenticateToken, controllers.HttpFileUploadFinish)
			v1Files.GET("download", controllers.HttpFileDownloadWithSecret)
			v1Files.POST("download", middleware.AuthenticateToken, controllers.HttpFileDownloadCreateSecret)

			v1Files.GET("list", middleware.AuthenticateToken, controllers.HttpFilesList)
			v1Files.POST("folder", middleware.AuthenticateToken, controllers.HttpFolderCreate)
		}
	}

	if os.Getenv("SERVE_PUBLIC") == "true" {
		router.Use(static.Serve("/", static.LocalFile(os.Getenv("PUBLIC_PATH"), true)), middleware.PublicHeader)
		router.NoRoute(func(c *gin.Context) {
			c.File(os.Getenv("PUBLIC_PATH") + "/index.html")
		})
	}
}
