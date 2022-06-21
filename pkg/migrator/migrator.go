package migrator

import (
	"cloudgobrrr/backend/database"
	"cloudgobrrr/backend/database/model"
	"cloudgobrrr/backend/pkg/env"
	"fmt"
	"os"
	"path/filepath"
)

func RunMigrations() {
	versionFilePath := filepath.Join(os.Getenv("DATA_DIRECTORY"), "version.txt")

	// For future use maybe
	file, err := os.Create(versionFilePath)
	if err != nil {
		panic(err)
	}
	_, err = file.WriteString(env.GetVersion())
	if err != nil {
		panic(err)
	}
	file.Close()

	db := database.GetDB()

	// AutoMigrations
	err = db.AutoMigrate(&model.User{}, &model.DownloadSecret{}, &model.Session{})
	if err != nil {
		panic(err)
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

	// ToDo: write migration system
}
