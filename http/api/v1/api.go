package v1

import (
	"cloudgobrrr/config"
	"cloudgobrrr/http/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

var conf *viper.Viper
var val *validator.Validator
var app fiber.Router

func init() {
	conf = config.Get()
	val = validator.Get()
}

func Setup(router fiber.Router) {
	app = router

	setupDevelopment() // this is only run in development mode
	setupAuth()
	setupSecurity()
}
