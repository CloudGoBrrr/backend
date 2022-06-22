package main

import (
	"cloudgobrrr/backend/database"
	"cloudgobrrr/backend/http"
	"cloudgobrrr/backend/pkg/env"
	"cloudgobrrr/backend/pkg/migrator"
	"log"
	"math/rand"
	"os"
	"time"
)

func main() {
	var err error

	// --
	// Pre boot
	// --
	rand.Seed(time.Now().Unix())
	env.EnvBuild()
	// Creating necessary folders
	err = os.MkdirAll(os.Getenv("DATA_DIRECTORY"), 0755)
	if err != nil {
		log.Fatalf("Failed to create data directory: %s", err)
	}
	err = os.MkdirAll(os.Getenv("USER_DIRECTORY"), 0755)
	if err != nil {
		log.Fatalf("Failed to create user directory: %s", err)
	}
	err = os.MkdirAll(os.Getenv("TEMP_DIRECTORY"), 0755)
	if err != nil {
		log.Fatalf("Failed to create temp directory: %s", err)
	}

	// --
	// Boot
	// --
	// DB Connect
	err = database.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %s", err)
	}

	// Run Migrations
	migrator.RunMigrations()

	// HTTP Boot
	http.NewHttpServer()

	// --
	// Shutdown
	// --
	database.GetSQLDB().Close()
	log.Println("Server stopped gracefully")
}
