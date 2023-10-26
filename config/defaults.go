package config

// defaults sets the default values for the config
func defaults() {
	// Development
	conf.SetDefault("development", false)

	// HTTP
	conf.SetDefault("http.host", "")
	conf.SetDefault("http.port", "8080")
	conf.SetDefault("http.prefork", false)
	conf.SetDefault("http.trustedProxy.enable", false)
	conf.SetDefault("http.trustedProxy.ips", []string{})
	conf.SetDefault("http.serverHeader", "CloudGoBrrr")
	conf.SetDefault("http.recover.enable", true)
	conf.SetDefault("http.recover.stacktrace", false)
	// Database
	conf.SetDefault("database.backend", "sqlite") // sqlite, mysql, memory or sqlite-memory
	conf.SetDefault("database.migrator.enable", true)
	conf.SetDefault("database.sqlite.path", "./cloudgobrrr.db")
	conf.SetDefault("database.mysql.host", "localhost")
	conf.SetDefault("database.mysql.port", "3306")
	conf.SetDefault("database.mysql.username", "root")
	conf.SetDefault("database.mysql.password", "")
	conf.SetDefault("database.mysql.name", "cloudgobrrr")
	// Frontend
	conf.SetDefault("frontend.enable", true)
	conf.SetDefault("frontend.mode", "static") // static, proxy
	conf.SetDefault("frontend.static.path", "./dist")
	conf.SetDefault("frontend.proxy.url", "http://localhost:3000")
	// Filesystem
	conf.SetDefault("filesystem.backend", "os") // os, memory
	conf.SetDefault("filesystem.os.path", "./storage")
	// Password
	conf.SetDefault("password.memory", 64) // in MB
	conf.SetDefault("password.iterations", 3)
	conf.SetDefault("password.parallelism", 2)
	conf.SetDefault("password.saltLength", 16)
	conf.SetDefault("password.keyLength", 32)
	// Logging
	conf.SetDefault("logging.format", "console") // json, console
	conf.SetDefault("logging.level", "info")     // trace, debug, info, warn, error, fatal, panic
	conf.SetDefault("logging.http", false)
	// Auth
	conf.SetDefault("auth.signup", false)
	// Security
	conf.SetDefault("security.token.groups", 4)
	conf.SetDefault("security.token.length", 8)
	// JWT
	conf.SetDefault("jwt.secret", "")
	conf.SetDefault("jwt.signingMethod", "HS256") // HS256
	conf.SetDefault("jwt.expiration", 10)         // in minutes
	conf.SetDefault("jwt.session.default", 60)    // in minutes
	conf.SetDefault("jwt.session.long", 43200)    // in minutes - 30 days
	conf.SetDefault("jwt.session.length", 64)     // in characters
}
