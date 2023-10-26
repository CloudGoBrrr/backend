package migrations

import (
	"cloudgobrrr/database/models"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

var db *gorm.DB

func Prepare(database *gorm.DB) {
	// load database connection
	db = database

	// check if migration table exists
	if !db.Migrator().HasTable("migrations") {
		log.Debug().Msg("creating migrations table")
		db.AutoMigrate(&models.Migration{})
	}
}
