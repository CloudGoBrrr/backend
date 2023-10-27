package models

import (
	"cloudgobrrr/filesystem"
	"cloudgobrrr/utils"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

type User struct {
	ID        ulid.ULID `gorm:"primarykey"`
	Username  string    `gorm:"unique"`
	Email     string    `gorm:"unique"`
	Password  string
	IsAdmin   bool `gorm:"default:false"`
	CreatedAt int64
	UpdatedAt int64
}

// UserCreate creates a user in the database
func UserCreate(username, email, plainTextPassword string, isAdmin bool) error {
	log.Debug().Str("username", username).Str("email", email).Msg("creating new user")
	// hash password
	hashedPassword, err := utils.PasswordHash(plainTextPassword)
	if err != nil {
		return err
	}

	// create user in database
	tx := db.Create(&User{
		ID:        ulid.Make(),
		Username:  username,
		Email:     email,
		Password:  hashedPassword,
		IsAdmin:   isAdmin,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	})
	if tx.Error != nil {
		return tx.Error
	}

	userFs := filesystem.GetUserFs(username, fs)
	// create user directory
	err = userFs.MkdirAll("/", 0755)
	if err != nil {
		return err
	}

	// create user files
	readmeFile, err := userFs.Create("README.md")
	if err != nil {
		return err
	}
	_, err = readmeFile.WriteString("Welcome to CloudGoBrrr!\nThis is your user directory. You can upload files here.\n\nYour username is: " + username)
	if err != nil {
		return err
	}

	return nil
}

// UserGetByUsername gets a user by their username
func UserGetByUsername(username string) (*User, error) {
	var user User
	tx := db.Where("username = ?", username).Find(&user)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &user, nil
}

// UserGetByEmail gets a user by their email
func UserGetByEmail(email string) (*User, error) {
	var user User
	tx := db.Where("email = ?", email).Find(&user)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &user, nil
}

// UserGetById gets a user by their id
func UserGetById(id ulid.ULID) (*User, error) {
	var user User
	tx := db.Where("id = ?", id).Find(&user)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &user, nil
}
