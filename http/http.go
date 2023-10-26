package http

import (
	"cloudgobrrr/config"
	v1 "cloudgobrrr/http/api/v1"
	"cloudgobrrr/http/api/webdav"
	"cloudgobrrr/http/response"

	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var conf *viper.Viper
var app *fiber.App

var methods []string = append(fiber.DefaultMethods, "PROPFIND", "PROPPATCH", "MKCOL", "COPY", "MOVE", "LOCK", "UNLOCK")

func init() {
	conf = config.Get()
	app = fiber.New(fiber.Config{
		DisableStartupMessage:   true,         // disable startup message
		ErrorHandler:            errorHandler, // points to error handler
		Prefork:                 conf.GetBool("http.prefork"),
		ServerHeader:            conf.GetString("http.serverHeader"),
		AppName:                 "CloudGoBrrr",
		EnableTrustedProxyCheck: conf.GetBool("http.trustedProxy.enable"),
		TrustedProxies:          conf.GetStringSlice("http.trustedProxy.ips"),
		RequestMethods:          methods,
	})

	// run setup functions
	setupRecover()
	setupLogging()
	setupHookStartup()

	apiGroup := app.Group("/api")
	v1.Setup(apiGroup.Group("/v1"))
	webdav.Setup(apiGroup.Group("/webdav"))
}

// Error handler for failed requests
func errorHandler(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusBadRequest).JSON(response.Error{
		Success: false,
		Errors: []response.ErrorElement{
			{
				Message: err.Error(),
			},
		},
	})
}

// Setup recover middleware
//
// Recover middleware recovers from panics anywhere in the chain
func setupRecover() {
	if conf.GetBool("http.recover.enable") {
		app.Use(recover.New(recover.Config{
			EnableStackTrace: conf.GetBool("http.recover.stacktrace"),
		}))
	}
}

// Setup logging middleware
//
// Logging middleware logs the http requests
func setupLogging() {
	if conf.GetBool("logging.http") {
		app.Use(fiberzerolog.New(fiberzerolog.Config{
			Logger:   &log.Logger,
			Messages: []string{"server error", "client error", "success"},
		}))
	}
}

// Setup hook for startup
//
// Hook for startup logs the address the server is listening on
func setupHookStartup() {
	app.Hooks().OnListen(func(listenData fiber.ListenData) error {
		if fiber.IsChild() {
			return nil
		}

		// determine scheme
		scheme := "http"
		if listenData.TLS {
			scheme = "https"
		}

		// address to log
		addr := scheme + "://" + listenData.Host + ":" + listenData.Port
		log.Info().
			Str("addr", addr).
			Msg("server started")

		// localhost address to log
		addr = scheme + "://127.0.0.1:" + listenData.Port
		log.Debug().Str("addr", addr).Msg("localhost")

		return nil
	})
}

// Get returns the fiber app
func Get() *fiber.App {
	return app
}
