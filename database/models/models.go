package models

import (
	"cloudgobrrr/config"
	"cloudgobrrr/filesystem"

	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var db *gorm.DB
var fs afero.Fs
var conf *viper.Viper

func Prepare(database *gorm.DB) {
	// load database connection
	db = database
	fs = filesystem.Get()
	conf = config.Get()
}
