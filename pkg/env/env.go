package env

import (
	"os"
)

func BuildEnv() {
	defaultEnv("CLOUDGOBRRR_ENV", "production")
	defaultEnv("HTTP_PORT", "8080")
	defaultEnv("TRUSTED_PROXIES", "nil")

	defaultEnv("DB_HOST", "localhost")
	defaultEnv("DB_PORT", "3306")
	defaultEnv("DB_USER", "cloudgobrrr")
	defaultEnv("DB_PASSWORD", "cloudgobrrr")
	defaultEnv("DB_NAME", "cloudgobrrr")

	// Temp and user data are not allowed to be on different drives / volumes
	defaultEnv("DATA_DIRECTORY", "./data/user")
	defaultEnv("TEMP_DIRECTORY", "./data/tmp")

	defaultEnv("SERVE_PUBLIC", "true")
	defaultEnv("PUBLIC_PATH", "./frontend/build")
	defaultEnv("PUBLIC_REGISTRATION", "false")
}

func defaultEnv(EnvKey string, DefaultValue string) {
	if os.Getenv(EnvKey) == "" {
		os.Setenv(EnvKey, DefaultValue)
	}
}
