package model

import (
	"time"

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
func DownloadSecretCreate(userID uint, path string, filename string) (string, error) {
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

func DownloadSecretGetBySecret(secret string) (*DownloadSecret, error) {
	db := database.GetDB()
	var downloadSecret DownloadSecret
	if err := db.Where("secret = ?", secret).First(&downloadSecret).Error; err != nil {
		return nil, err
	}
	return &downloadSecret, nil
}

// CLEANUP
func DownloadSecretCleanup() error {
	//delete all secrets older than 2 hour
	db := database.GetDB()
	createdAt := time.Now().Add(-2 * time.Hour)
	db.Where("created_at < ?", createdAt).Unscoped().Delete(&DownloadSecret{})
	return nil
}
