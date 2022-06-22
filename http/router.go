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
		v1.GET("featureFlags", controllers.HttpFeatureFlagGet)

		v1Auth := v1.Group("auth")
		{
			v1Auth.GET("details", middleware.AuthenticateToken, controllers.HttpAuthDetails)
			v1Auth.POST("signin", controllers.HttpAuthSignin)
			v1Auth.POST("signup", controllers.HttpAuthSignup)
			v1Auth.POST("changepassword", middleware.AuthenticateToken, controllers.HttpAuthChangePassword)
			v1AuthSession := v1Auth.Group("session")
			{
				v1AuthSession.POST("basic", middleware.AuthenticateToken, controllers.HttpSessionCreateBasicAuth)
				v1AuthSession.GET("list", middleware.AuthenticateToken, controllers.HttpSessionList)
				v1AuthSession.PUT("description", middleware.AuthenticateToken, controllers.HttpSessionChangeDescription)
				v1AuthSession.DELETE("", middleware.AuthenticateToken, controllers.HttpSessionDeleteWithID)
			}

		}

		v1File := v1.Group("file")
		{
			v1File.PUT("upload", middleware.AuthenticateToken, controllers.HttpFileChunkedUpload)
			v1File.POST("upload", middleware.AuthenticateToken, controllers.HttpFileChunkedUploadFinish)
			v1File.GET("download", controllers.HttpFileDownloadWithSecret)
			v1File.POST("download", middleware.AuthenticateToken, controllers.HttpFileDownloadCreateSecret)
			v1File.DELETE("", middleware.AuthenticateToken, controllers.HttpFileDelete)
		}

		v1Folder := v1.Group("folder")
		{
			v1Folder.POST("", middleware.AuthenticateToken, controllers.HttpFolderCreate)
		}

		v1Files := v1.Group("files")
		{
			v1Files.GET("", middleware.AuthenticateToken, controllers.HttpFilesList)
		}
	}

	if os.Getenv("SERVE_PUBLIC") == "true" {
		router.Use(static.Serve("/", static.LocalFile(os.Getenv("PUBLIC_PATH"), true)), middleware.PublicHeader)
		router.NoRoute(func(c *gin.Context) {
			c.File(os.Getenv("PUBLIC_PATH") + "/index.html")
		})
	}
}
