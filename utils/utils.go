package utils

import (
	"cloudgobrrr/config"

	"github.com/oklog/ulid/v2"
	"github.com/spf13/viper"
)

var conf *viper.Viper

var EmptyULID ulid.ULID = ulid.ULID{}

func init() {
	conf = config.Get()

	params = paramsStruct{
		memory:      conf.GetUint32("password.memory") * 1024,
		iterations:  conf.GetUint32("password.iterations"),
		parallelism: uint8(conf.GetUint("password.parallelism")),
		saltLength:  conf.GetUint32("password.saltLength"),
		keyLength:   conf.GetUint32("password.keyLength"),
	}

	secret = []byte(conf.GetString("jwt.secret"))
}
