package model

import (
	"cloudgobrrr/backend/database"
	"cloudgobrrr/backend/pkg/helpers"
	"fmt"
	"os"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string
	Email    string
	Password string
	IsAdmin  bool `gorm:"default:false"`
}

// CREATE
func UserCreate(username string, email string, plainTextPassword string) error {
	db := database.GetDB()
	// check if user exists
	var count int64
	db.Model(&User{}).Where("username = ?", username).Count(&count)
	if count == 0 {
		// create user in directory
		err := os.Mkdir(os.Getenv("USER_DIRECTORY")+"/"+username, 0755)
		if err != nil {
			return err
		}
		// ---

		// add sample files (needs a better solution)
		file, err := os.Create(os.Getenv("USER_DIRECTORY") + "/" + username + "/README.txt")
		if err != nil {
			return err
		}
		_, err = file.WriteString("Hello World")
		if err != nil {
			return err
		}
		file.Close()
		err = os.Mkdir(os.Getenv("USER_DIRECTORY")+"/"+username+"/test_folder", 0755)
		if err != nil {
			return err
		}
		file, err = os.Create(os.Getenv("USER_DIRECTORY") + "/" + username + "/test_folder" + "/README.txt")
		if err != nil {
			return err
		}
		_, err = file.WriteString("Hello World")
		if err != nil {
			return err
		}
		file.Close()
		err = os.Mkdir(os.Getenv("USER_DIRECTORY")+"/"+username+"/empty_folder", 0755)
		if err != nil {
			return err
		}

		// ---
		// create user in database
		hash, err := helpers.PasswordHash(plainTextPassword)
		if err != nil {
			return err
		}
		db.Create(&User{Username: username, Email: email, Password: hash})
		return nil
	}
	return fmt.Errorf("user already exists")
}

// ToDo: add quota
// ToDo: add template system for createuser

func UserChangePassword(userID uint, plainTextPassword string) error {
	db := database.GetDB()
	var user User
	db.Where("id = ?", userID).First(&user)
	if user.ID == 0 {
		return fmt.Errorf("user does not exist")
	}
	hash, err := helpers.PasswordHash(plainTextPassword)
	if err != nil {
		return err
	}
	user.Password = hash
	db.Save(&user)
	return nil
}

func UserGetByUsername(username string) (User, error) {
	db := database.GetDB()
	var user User
	db.Where("username = ?", username).First(&user)
	if user.ID == 0 {
		return User{}, fmt.Errorf("user does not exist")
	}
	return user, nil
}

func UserGetByID(userID uint) (User, error) {
	db := database.GetDB()
	var user User
	db.Where("id = ?", userID).First(&user)
	if user.ID == 0 {
		return User{}, fmt.Errorf("user does not exist")
	}
	return user, nil
}
