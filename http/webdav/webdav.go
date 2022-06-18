package webdav

import (
	"cloudgobrrr/backend/database/model"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	dav "golang.org/x/net/webdav"
)

var prefix = "/remote/webdav"

func Handle() gin.HandlerFunc {

	web := dav.Handler{
		Prefix:     prefix,
		FileSystem: dav.Dir(os.Getenv(("DATA_DIRECTORY"))),
		LockSystem: dav.NewMemLS(),
		Logger:     logger,
	}

	return func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, prefix) {
			if !validator(c) {
				c.Writer.Header().Set("WWW-Authenticate", "Basic realm=\"User Visible Realm\", charset=\"UTF-8\"")
				c.AbortWithStatus(401)
				return
			}
			c.Status(200)
			web.ServeHTTP(c.Writer, c.Request)
			c.Abort()
		}
	}
}

func validator(c *gin.Context) bool {
	username, password, hasAuth := c.Request.BasicAuth()
	if !hasAuth {
		return false
	}

	user, _, err := model.GetUserFromBasicAuth(username, password)
	if err != nil {
		log.Println(err)
		return false
	}

	if user.Username == username {
		t := filepath.Join(prefix, username)
		return strings.HasPrefix(c.Request.URL.Path, t)
	}
	return false
}

func logger(req *http.Request, err error) {
	if err != nil {
		log.Println("webdav error: " + err.Error())
	}
}
