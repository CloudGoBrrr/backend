package migrator

import (
	"cloudgobrrr/backend/database"
	"cloudgobrrr/backend/database/model"
	"cloudgobrrr/backend/pkg/env"
	"log"
	"os"
	"path/filepath"
)

func RunMigrations() {
	versionFilePath := filepath.Join(os.Getenv("DATA_DIRECTORY"), "version.txt")

	// For future use maybe
	file, err := os.Create(versionFilePath)
	if err != nil {
		log.Fatalf("Failed to create version file: %s", err)
	}
	_, err = file.WriteString(env.VersionGet())
	if err != nil {
		log.Fatalf("Failed to write version file: %s", err)
	}
	file.Close()

	db := database.GetDB()

	// AutoMigrations
	err = db.AutoMigrate(&model.User{}, &model.DownloadSecret{}, &model.Session{})
	if err != nil {
		log.Fatalf("Failed to run migrations: %s", err)
	}

	// Create Admin user if none exists
	var count int64
	db.Model(&model.User{}).Count(&count)
	if count == 0 {
		log.Println("No users found, creating admin user")
		err := model.UserCreate("admin", "admin@example.com", "admin")
		if err != nil {
			log.Fatalf("Failed to create admin user: %s", err)
		}
	}

	// ToDo: write migration system
}
