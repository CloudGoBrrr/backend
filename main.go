package main

import (
	"cloudgobrrr/config"
	"cloudgobrrr/http"
	"fmt"

	"github.com/rs/zerolog/log"
)

func main() {
	conf := config.Get()

	if conf.GetBool("development") {
		log.Warn().Msg("development mode is enabled - this is not recommended for production environments")
	}

	httpServer := http.Get()
	httpServer.Listen(fmt.Sprintf("%s:%s", conf.GetString("http.host"), conf.GetString("http.port")))
}
