package models

import (
	"cloudgobrrr/utils"
	"encoding/base64"
	"time"

	"github.com/oklog/ulid/v2"
)

type Token struct {
	ID          ulid.ULID `gorm:"primarykey"` // ToDo: replace with UUID
	UserID      uint
	Encoded     string `gorm:"primarykey"`
	Description string
	CreatedAt   int64
	UpdateAt    int64
}

// TokenCreate creates a token in the database
func TokenCreate(userId uint, description string) (ulid.ULID, string, error) {
	amount := conf.GetInt("auth.token.elementAmount")
	length := conf.GetInt("auth.token.elementLength")
	var generatedPassword string

	for i := 0; i < amount; i++ {
		if generatedPassword != "" {
			generatedPassword += "-"
		}
		element, err := utils.GenerateRandomStringNoDash(length)
		if err != nil {
			return utils.EmptyULID, "", err
		}
		generatedPassword += element
	}

	user, err := UserGetById(userId)
	if err != nil {
		return utils.EmptyULID, "", err
	}

	encoded := base64.StdEncoding.EncodeToString([]byte(user.Username + ":" + generatedPassword))

	token := Token{
		ID:          ulid.Make(),
		UserID:      userId,
		Encoded:     encoded,
		Description: description,
		CreatedAt:   time.Now().Unix(),
		UpdateAt:    time.Now().Unix(),
	}

	tx := db.Create(&token)
	if tx.Error != nil {
		return utils.EmptyULID, "", tx.Error
	}

	return token.ID, generatedPassword, nil
}

func TokenGetAllByUserID(userId uint) ([]Token, error) {
	var tokens []Token
	tx := db.Where("user_id = ?", userId).Find(&tokens)
	if tx.Error != nil {
		return []Token{}, tx.Error
	}

	return tokens, nil
}

// TokenGetByEncoded gets a token by its encoded value
func TokenGetByEncoded(encoded string) (*Token, error) {
	var token Token
	tx := db.Where("encoded = ?", encoded).Find(&token)
	if tx.Error != nil {
		return &Token{}, tx.Error
	}

	return &token, nil
}

// TokenGetByUsernameAndPassword gets a token by its username and password
func TokenGetByUsernameAndPassword(username string, password string) (*Token, error) {
	return TokenGetByEncoded(base64.StdEncoding.EncodeToString([]byte(username + ":" + password)))
}

// TokenGetByID gets a token by its id
func TokenGetByID(id ulid.ULID) (*Token, error) {
	var token Token
	tx := db.Where("id = ?", id).Find(&token)
	if tx.Error != nil {
		return &Token{}, tx.Error
	}

	return &token, nil
}

// TokenDeleteByID deletes a token by its id
func TokenDeleteByID(id ulid.ULID) error {
	tx := db.Where("id = ?", id).Delete(&Token{})
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}
