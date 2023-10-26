package models

import (
	"cloudgobrrr/utils"
	"time"

	"github.com/oklog/ulid/v2"
)

// Sessions are jwt refresh tokens for the web
type Session struct {
	ID          ulid.ULID `gorm:"primarykey"`
	UserID      uint
	Token       string `gorm:"primarykey"`
	Description string
	Remember    bool
	ExpiresAt   int64
	CreatedAt   int64
	UpdateAt    int64
}

// SessionCreate creates a session in the database
//
// Returns the identifier for the session
func SessionCreate(userId uint, description string, remember bool) (string, error) {
	var token string

	token, err := utils.GenerateRandomString(conf.GetInt("jwt.session.length"))
	if err != nil {
		return "", err
	}

	var expiresAt int64
	if remember {
		expiresAt = time.Now().Add(time.Duration(conf.GetUint("jwt.session.long")) * time.Minute).Unix()
	} else {
		expiresAt = time.Now().Add(time.Duration(conf.GetUint("jwt.session.default")) * time.Minute).Unix()
	}

	tx := db.Create(&Session{
		ID:          ulid.Make(),
		UserID:      userId,
		Token:       token,
		Description: description,
		Remember:    remember,
		ExpiresAt:   expiresAt,
		CreatedAt:   time.Now().Unix(),
		UpdateAt:    time.Now().Unix(),
	})

	if tx.Error != nil {
		return "", tx.Error
	}

	return token, nil
}

func SessionGetAllByUserID(userId uint) ([]Session, error) {
	var sessions []Session
	tx := db.Where("user_id = ?", userId).Find(&sessions)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return sessions, nil
}

func SessionGetByToken(token string) (*Session, error) {
	var session Session
	tx := db.Where("token = ?", token).First(&session)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &session, nil
}

func SessionGetByID(id ulid.ULID) (*Session, error) {
	var session Session
	tx := db.Where("id = ?", id).First(&session)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &session, nil
}

func SessionDeleteByToken(token string) error {
	tx := db.Where("token = ?", token).Delete(&Session{})
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func SessionDeleteByID(id ulid.ULID) error {
	tx := db.Where("id = ?", id).Delete(&Session{})
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
