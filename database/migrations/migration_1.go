package migrations

import (
	"cloudgobrrr/database/models"
)

func Migration_1() {
	db.AutoMigrate(&models.User{}, &models.Session{}, &models.Token{})

	// Create admin user
	models.UserCreate("admin", "admin@example.com", "changeme", true)
}
