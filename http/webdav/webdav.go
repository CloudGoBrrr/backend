package webdav

import (
	"cloudgobrrr/backend/http/middleware"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	dav "golang.org/x/net/webdav"
)

func Handle() gin.HandlerFunc {

	web := dav.Handler{
		Prefix:     "/api/webdav",
		FileSystem: dav.Dir(os.Getenv(("DATA_DIRECTORY"))),
		LockSystem: dav.NewMemLS(),
		Logger:     logger,
	}

	return func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, web.Prefix) {
			if !middleware.ValidateBasic(c) {
				return
			}

			expectedPrefix := filepath.Join(web.Prefix, c.GetString("userName"))
			if !strings.HasPrefix(c.Request.URL.Path, expectedPrefix) {
				c.AbortWithStatus(403)
				return
			}

			c.Status(200)
			web.ServeHTTP(c.Writer, c.Request)
			c.Abort()
		}
	}
}

func logger(req *http.Request, err error) {
	if err != nil {
		log.Println("webdav error: " + err.Error())
	}
}
