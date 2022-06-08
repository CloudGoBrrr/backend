package database

import (
	"database/sql"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

var SQLDB *sql.DB

func InitDB() (*gorm.DB, error) {
	var err error

	if SQLDB != nil {
		SQLDB.Close()
	}

	DB, err = NewDB()
	if err != nil {
		return nil, err
	}

	SQLDB, err = DB.DB()
	if err != nil {
		return nil, err
	}

	return DB, nil
}

func GetDB() *gorm.DB {
	return DB
}

func GetSQLDB() *sql.DB {
	return SQLDB
}

func NewDB() (*gorm.DB, error) {
	// open database connection
	dsn := os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") + "@tcp(" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") + ")/" + os.Getenv("DB_NAME") + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
