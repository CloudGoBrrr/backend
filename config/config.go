package config

import (
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var conf *viper.Viper

func init() {
	conf = viper.New()

	// For config file
	conf.SetConfigName("cloudgobrrr")
	conf.SetConfigType("toml")
	conf.AddConfigPath(".")
	conf.AddConfigPath("/config")

	// For environment variables
	conf.AutomaticEnv()
	conf.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Default values
	defaults()

	// Find and read the config file
	err := conf.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			setupLogging()
			log.Warn().Msg("config file not found, using defaults and environment")
		} else {
			panic(err)
		}
	} else {
		setupLogging()
		log.Debug().Msg("config file found and loaded")
	}
}

// setupLogging sets up the logging format and level
func setupLogging() {
	if conf.GetString("logging.format") == "console" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	switch conf.GetString("logging.level") {
	case "trace":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "panic":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	}

	log.Debug().
		Str("logging_level", conf.GetString("logging.level")).
		Str("logging_format", conf.GetString("logging.format")).
		Msg("logging setup finished")
}

// Get returns the config object
func Get() *viper.Viper {
	return conf
}
