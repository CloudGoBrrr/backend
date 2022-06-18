package model

import (
	"cloudgobrrr/backend/database"
	"cloudgobrrr/backend/pkg/helpers"
	"errors"

	"gorm.io/gorm"
)

type AuthToken struct {
	gorm.Model
	UserID      uint
	Token       string
	Description string
	Type        string `gorm:"default:token"` // migration
}

// CREATE
func CreateAuthToken(user User, description string) (string, error) {
	db := database.GetDB()
	token, err := helpers.GenerateRandomString(128)
	if err != nil {
		return "", err
	}
	var count int64
	db.Model(&AuthToken{}).Where("token = ?", token).Count(&count)
	for count != 0 {
		token, err = helpers.GenerateRandomString(128)
		if err != nil {
			return "", err
		}
		db.Model(&AuthToken{}).Where("token = ?", token).Count(&count)
	}
	db.Create(&AuthToken{UserID: user.ID, Token: token, Description: description, Type: "token"})
	return token, nil
}

func CreateBasicAuth(user User, description string) (string, error) {
	db := database.GetDB()
	passwd, err := helpers.GenerateRandomString(32)
	if err != nil {
		return "", err
	}
	var token string
	token = user.Username + ":" + passwd
	var count int64
	db.Model(&AuthToken{}).Where("token = ?", token).Count(&count)
	for count != 0 {
		passwd, err = helpers.GenerateRandomString(32)
		token = user.Username + ":" + passwd
		if err != nil {
			return "", err
		}
		db.Model(&AuthToken{}).Where("token = ?", token).Count(&count)
	}
	db.Create(&AuthToken{UserID: user.ID, Token: token, Description: description, Type: "basic"})
	return passwd, nil
}

// DELETE
func DeleteAuthToken(token string) error {
	db := database.GetDB()
	var count int64
	db.Model(&AuthToken{}).Where("token = ?", token).Count(&count)
	if count == 1 {
		db.Delete(&AuthToken{}, "token = ?", token)
		return nil
	}
	return errors.New("token does not exist")
}

func DeleteAllAuthTokensOfUser(userID uint) error {
	db := database.GetDB()
	var authTokens []AuthToken
	db.Where("user_id = ?", userID).Find(&authTokens)
	for _, authToken := range authTokens {
		db.Delete(&authToken)
	}
	return nil
}

func DeleteAuthTokenWithID(tokenID uint, userID uint) error {
	db := database.GetDB()
	var authToken AuthToken
	db.Where("id = ?", tokenID).First(&authToken)
	if authToken.ID == 0 {
		return errors.New("token does not exist")
	}
	if authToken.UserID == userID {
		db.Delete(&authToken)
		return nil
	}
	return errors.New("token does not belong to user")
}

// GET
func GetAllAuthTokensOfUserID(userID uint) ([]map[string]interface{}, error) {
	db := database.GetDB()
	var results []map[string]interface{}
	db.Model(&AuthToken{}).Select("id", "description", "created_at").Where("user_id = ?", userID).Find(&results)
	return results, nil
}

func GetUserFromAuthToken(token string) (User, AuthToken, error) {
	db := database.GetDB()
	// ToDo: SQL join for better performance?
	var authToken AuthToken
	db.Where("token = ? AND type = \"token\"", token).First(&authToken)
	if authToken.ID == 0 {
		return User{}, AuthToken{}, errors.New("token does not exist")
	}
	var user User
	db.Where("id = ?", authToken.UserID).First(&user)
	return user, authToken, nil
}

func GetUserFromBasicAuth(username string, password string) (User, AuthToken, error) {
	db := database.GetDB()
	// ToDo: SQL join for better performance?
	token := username + ":" + password
	var authToken AuthToken
	db.Where("token = ? AND type = \"basic\"", token).First(&authToken)
	if authToken.ID == 0 {
		return User{}, AuthToken{}, errors.New("token does not exist")
	}
	var user User
	db.Where("id = ?", authToken.UserID).First(&user)
	return user, authToken, nil
}
