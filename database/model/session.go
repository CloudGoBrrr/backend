package model

import (
	"cloudgobrrr/backend/database"
	"cloudgobrrr/backend/pkg/helpers"
	"fmt"

	"gorm.io/gorm"
)

type Session struct {
	gorm.Model
	UserID      uint
	Identifier  string
	Description string
	Type        string `gorm:"default:token"` // migration
}

var whereToken = "identifier = ? AND type = \"token\""
var whereBasic = "identifier = ? AND type = \"basic\""

// CREATE
func SessionCreateToken(user User, description string) (string, error) {
	db := database.GetDB()
	token, err := helpers.GenerateRandomString(128)
	if err != nil {
		return "", err
	}
	var count int64
	db.Model(&Session{}).Where(whereToken, token).Count(&count)
	for count != 0 {
		token, err = helpers.GenerateRandomString(128)
		if err != nil {
			return "", err
		}
		db.Model(&Session{}).Where(whereToken, token).Count(&count)
	}
	db.Create(&Session{UserID: user.ID, Identifier: token, Description: description, Type: "token"})
	return token, nil
}

func SessionCreateBasic(user User, description string) (string, error) {
	db := database.GetDB()
	passwd, err := helpers.GenerateRandomString(32)
	if err != nil {
		return "", err
	}
	var basicToken string
	basicToken = user.Username + ":" + passwd
	var count int64
	db.Model(&Session{}).Where(whereBasic, basicToken).Count(&count)
	for count != 0 {
		passwd, err = helpers.GenerateRandomString(32)
		basicToken = user.Username + ":" + passwd
		if err != nil {
			return "", err
		}
		db.Model(&Session{}).Where(whereBasic, basicToken).Count(&count)
	}
	db.Create(&Session{UserID: user.ID, Identifier: basicToken, Description: description, Type: "basic"})
	return passwd, nil
}

// DELETE
func SessionDeleteAll(userID uint) error {
	db := database.GetDB()
	var sessions []Session
	db.Where("user_id = ?", userID).Find(&sessions)
	for _, session := range sessions {
		db.Delete(&session)
	}
	return nil
}

func SessionDeleteWithID(sessionID uint, userID uint) error {
	db := database.GetDB()
	var session Session
	db.Where("id = ?", sessionID).First(&session)
	if session.ID == 0 {
		return fmt.Errorf("session does not exist")
	}
	if session.UserID == userID {
		db.Delete(&session)
		return nil
	}
	return fmt.Errorf("token does not belong to user")
}

// GET
func SessionGetAll(userID uint) ([]map[string]interface{}, error) {
	db := database.GetDB()
	var results []map[string]interface{}
	db.Model(&Session{}).Select("id", "description", "created_at").Where("user_id = ?", userID).Find(&results)
	return results, nil
}

func SessionGetUserToken(token string) (User, Session, error) {
	db := database.GetDB()
	// ToDo: SQL join for better performance?
	var session Session
	db.Where(whereToken, token).First(&session)
	if session.ID == 0 {
		return User{}, Session{}, fmt.Errorf("token does not exist")
	}
	var user User
	db.Where("id = ?", session.UserID).First(&user)
	return user, session, nil
}

func SessionGetUserBasic(username string, password string) (User, Session, error) {
	db := database.GetDB()
	// ToDo: SQL join for better performance?
	token := username + ":" + password
	var session Session
	db.Where(whereBasic, token).First(&session)
	if session.ID == 0 {
		return User{}, Session{}, fmt.Errorf("session does not exist")
	}
	var user User
	db.Where("id = ?", session.UserID).First(&user)
	return user, session, nil
}

// CHANGE
func SessionChangeDescription(sessionID uint, newDescription string) (string, error) {
	db := database.GetDB()
	var session Session
	db.Where("id = ?", sessionID).First(&session)
	if session.ID == 0 {
		return "", fmt.Errorf("session does not exists")
	}
	oldDescription := session.Description
	session.Description = newDescription
	db.Save(&session)
	return oldDescription, nil
}
