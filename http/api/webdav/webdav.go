package webdav

import (
	"cloudgobrrr/filesystem"
	"cloudgobrrr/http/middleware"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	dav "golang.org/x/net/webdav"
)

var app fiber.Router
var davHandler dav.Handler

func Setup(router fiber.Router) {
	app = router
	davHandler = dav.Handler{
		Prefix:     "/api/webdav",
		FileSystem: filesystem.WebdavFilesystem{Fs: filesystem.Get()},
		LockSystem: dav.NewMemLS(),
		Logger:     func(r *http.Request, err error) {},
	}

	app.All("/*",
		prePathCheck,
		middleware.Auth,
		postPathCheck,
		adaptor.HTTPHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			davHandler.ServeHTTP(w, r)
		}),
	)
}

func prePathCheck(c *fiber.Ctx) error {
	// check if path prefix is /api/webdav
	if c.Path() == "/api/webdav/" {
		return c.SendStatus(fiber.StatusNotFound)
	}

	return c.Next()
}

func postPathCheck(c *fiber.Ctx) error {
	// check if path prefix is /api/webdav/{username}
	if !strings.HasPrefix(c.Path(), "/api/webdav/"+c.Params("username")) {
		return c.SendStatus(fiber.StatusNotFound)
	}

	return c.Next()
}
