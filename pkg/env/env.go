package env

import (
	"os"
	"path/filepath"
)

func EnvBuild() {
	envDefault("CLOUDGOBRRR_ENV", "production")
	envDefault("HTTP_PORT", "8080")
	envDefault("TRUSTED_PROXIES", "nil")

	envDefault("DB_HOST", "localhost")
	envDefault("DB_PORT", "3306")
	envDefault("DB_USER", "cloudgobrrr")
	envDefault("DB_PASSWORD", "cloudgobrrr")
	envDefault("DB_NAME", "cloudgobrrr")

	// Temp and user data are not allowed to be on different drives / volumes
	envDefault("DATA_DIRECTORY", "./data")
	envDefault("USER_DIRECTORY", filepath.Join(os.Getenv("DATA_DIRECTORY"), "user"))
	envDefault("TEMP_DIRECTORY", filepath.Join(os.Getenv("DATA_DIRECTORY"), "tmp"))

	envDefault("ADMIN_USERNAME", "admin")
	envDefault("ADMIN_PASSWORD", "admin")
	envDefault("ADMIN_EMAIL", "admin@example.com")

	envDefault("SERVE_PUBLIC", "true")
	envDefault("PUBLIC_PATH", "./frontend/build")
	envDefault("PUBLIC_REGISTRATION", "false")
	envDefault("WEBDAV_ENABLED", "false")
}

func envDefault(EnvKey string, DefaultValue string) {
	if os.Getenv(EnvKey) == "" {
		os.Setenv(EnvKey, DefaultValue)
	}
}
