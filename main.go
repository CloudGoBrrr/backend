package main

import (
	"cloudgobrrr/backend/database"
	"cloudgobrrr/backend/database/model"
	"cloudgobrrr/backend/http"
	"cloudgobrrr/backend/pkg/env"
	"cloudgobrrr/backend/pkg/migrator"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"gorm.io/gorm"
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
	db, err := database.InitDB()
	if err != nil {
		MainLogger.Fatal(err)
	}

	// Run Migrations
	migrator.RunMigrations()

	// Init DB
	initDb(db)

	// HTTP Boot
	http.NewHttpServer()

	// --
	// Shutdown
	// --
	database.GetSQLDB().Close()
	MainLogger.Println("Server stopped gracefully")
}

func initDb(db *gorm.DB) {
	// Run Migrations
	err := db.AutoMigrate(&model.User{}, &model.DownloadSecret{}, &model.Session{})
	if err != nil {
		MainLogger.Fatal(err)
	}

	// Create Admin user if none exists
	var count int64
	db.Model(&model.User{}).Count(&count)
	if count == 0 {
		fmt.Println("No users found, creating admin user")
		err := model.CreateUser("admin", "admin@example.com", "admin")
		if err != nil {
			panic(err)
		}
	}
}
