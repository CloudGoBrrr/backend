package webdav

import (
	"cloudgobrrr/backend/http/middleware"
	"cloudgobrrr/backend/pkg/helpers"
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
		FileSystem: dav.Dir(os.Getenv(("USER_DIRECTORY"))),
		LockSystem: dav.NewMemLS(),
		Logger:     logger,
	}

	return func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, web.Prefix) {
			parsed := helpers.HttpSplitPath(c.Request.URL.Path)
			len := len(parsed)

			if len < 4 && c.Request.Method == "GET" {
				c.Data(http.StatusBadRequest, "text/plain; charset=utf-8", []byte("This is a webdav interface. Please use a webdav client to access the files."))
				c.Abort()
				return
			}

			if len < 3 {
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}

			if !middleware.ValidateBasic(c) {
				return
			}

			expectedPrefix := filepath.Join(web.Prefix, c.GetString("userName"))
			if !strings.HasPrefix(c.Request.URL.Path, expectedPrefix) {
				c.AbortWithStatus(http.StatusForbidden)
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
