package model

import (
	"gorm.io/gorm"

	"cloudgobrrr/backend/database"
	"cloudgobrrr/backend/pkg/helpers"
)

type DownloadSecret struct {
	gorm.Model
	Secret   string
	UserID   uint
	Path     string
	Filename string
}

// CREATE
func CreateDownloadSecret(userID uint, path string, filename string) (string, error) {
	db := database.GetDB()
	// create secret
	secret, err := helpers.GenerateRandomString(16)
	if err != nil {
		return "", err
	}
	// check if secret exists
	var count int64
	db.Model(&DownloadSecret{}).Where("secret = ?", secret).Count(&count)
	for count != 0 {
		secret, err = helpers.GenerateRandomString(16)
		if err != nil {
			return "", err
		}
		db.Model(&DownloadSecret{}).Where("secret = ?", secret).Count(&count)
	}
	db.Create(&DownloadSecret{Secret: secret, UserID: userID, Path: path, Filename: filename})
	return secret, nil
}

func GetDownloadSecretFromSecret(secret string) (*DownloadSecret, error) {
	db := database.GetDB()
	var downloadSecret DownloadSecret
	if err := db.Where("secret = ?", secret).First(&downloadSecret).Error; err != nil {
		return nil, err
	}
	return &downloadSecret, nil
}

// ToDo: delete secrets older than X days
