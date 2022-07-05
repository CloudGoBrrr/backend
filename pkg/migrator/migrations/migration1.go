package migrations

import (
	"cloudgobrrr/backend/database"
	"cloudgobrrr/backend/database/model"
	"log"
	"os"
)

func Migration1() {
	db := database.GetDB()
	user, err := model.UserGetByUsername(os.Getenv("ADMIN_USERNAME"))
	if err != nil {
		log.Fatalf("Failed to get admin user: %s", err)
	}
	if !user.IsAdmin {
		user.IsAdmin = true
		db.Save(&user)
		log.Println("Migration 1 complete")
	}
}
