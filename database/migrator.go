package database

import (
	"cloudgobrrr/database/migrations"
	"cloudgobrrr/database/models"

	"github.com/rs/zerolog/log"
)

func runMigrator() {
	// check if migrator is enabled
	if !conf.GetBool("database.migrator.enable") {
		log.Info().Msg("migrator disabled")
		return
	}
	log.Info().Msg("running migrations")

	// prepare migrator
	migrations.Prepare(db)

	// run migrations
	runMigration("migration_1", migrations.Migration_1)
}

func runMigration(name string, fn func()) {
	// Check if migration has already been run
	var migration models.Migration
	db.Where("name = ?", name).Find(&migration)
	if migration.Name == name {
		log.Debug().Str("name", name).Msg("migration already run")
		return
	}

	// Run migration
	log.Debug().Str("name", name).Msg("running migration")
	fn()

	// Create migration
	db.Create(&models.Migration{Name: name})
	log.Debug().Str("name", name).Msg("migration created")
}
