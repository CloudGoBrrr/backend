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

var MainLogger = log.New(os.Stdout, "[APP] ", log.Ldate|log.Ltime)

func main() {
	var err error

	// --
	// Pre boot
	// --
	rand.Seed(time.Now().Unix())
	env.BuildEnv()
	// Creating necessary folders
	err = os.MkdirAll(os.Getenv("DATA_DIRECTORY"), 0755)
	if err != nil {
		MainLogger.Fatal(err)
	}
	err = os.MkdirAll(os.Getenv("USER_DIRECTORY"), 0755)
	if err != nil {
		MainLogger.Fatal(err)
	}
	err = os.MkdirAll(os.Getenv("TEMP_DIRECTORY"), 0755)
	if err != nil {
		MainLogger.Fatal(err)
	}

	// --
	// Boot
	// --
	// DB Connect
	err = database.InitDB()
	if err != nil {
		MainLogger.Fatal(err)
	}

	// Run Migrations
	migrator.RunMigrations()

	// HTTP Boot
	http.NewHttpServer()

	// --
	// Shutdown
	// --
	database.GetSQLDB().Close()
	MainLogger.Println("Server stopped gracefully")
}
