package database

import (
	"cloudgobrrr/config"
	"cloudgobrrr/database/models"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB
var conf *viper.Viper

func init() {
	var err error
	conf = config.Get()

	if conf.GetString("database.backend") == "mysql" {
		log.Debug().
			Str("backend", "mysql").
			Msg("database backend selected")
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			conf.GetString("database.mysql.username"),
			conf.GetString("database.mysql.password"),
			conf.GetString("database.mysql.host"),
			conf.GetString("database.mysql.port"),
			conf.GetString("database.mysql.name"))
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	} else if conf.GetString("database.backend") == "sqlite" {
		log.Debug().
			Str("backend", "sqlite").
			Msg("database backend selected")
		db, err = gorm.Open(sqlite.Open(conf.GetString("database.sqlite.path")), &gorm.Config{})
	} else if conf.GetString("database.backend") == "sqlite-memory" || conf.GetString("database.backend") == "memory" {
		log.Debug().
			Str("backend", "sqlite-memory").
			Msg("database backend selected")
		db, err = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	} else {
		log.Fatal().
			Str("backend", conf.GetString("database.backend")).
			Msg("invalid db backend")
	}

	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	models.Prepare(db)

	log.Debug().Msg("database connection established")
	runMigrator()
}

// Get returns the database connection
func Get() *gorm.DB {
	return db
}
