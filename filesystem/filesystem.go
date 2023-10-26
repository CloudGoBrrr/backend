package filesystem

import (
	"cloudgobrrr/config"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

var conf *viper.Viper
var fs afero.Fs

func init() {
	conf = config.Get()

	if conf.GetString("filesystem.backend") == "os" {
		log.Debug().
			Str("backend", "os").
			Msg("filesystem backend selected")
		path, err := filepath.Abs(conf.GetString("filesystem.os.path"))
		if err != nil {
			log.Fatal().Msg(err.Error())
		}
		fs = afero.NewBasePathFs(afero.NewOsFs(), path)
	} else if conf.GetString("filesystem.backend") == "memory" {
		log.Debug().
			Str("backend", "memory").
			Msg("filesystem backend selected")
		fs = afero.NewBasePathFs(afero.NewMemMapFs(), "/")
	} else {
		log.Fatal().Msg("invalid filesystem backend")
	}
}

func GetUserFs(username string, fs afero.Fs) afero.Fs {
	return afero.NewBasePathFs(fs, username)
}

// Get returns the filesystem
func Get() afero.Fs {
	return fs
}
