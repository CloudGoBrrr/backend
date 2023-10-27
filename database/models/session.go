package models

import (
	"cloudgobrrr/utils"
	"time"

	"github.com/oklog/ulid/v2"
)

// Sessions are jwt refresh tokens for the web
type Session struct {
	ID          ulid.ULID `gorm:"primarykey"`
	UserID      ulid.ULID
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
func SessionCreate(userId ulid.ULID, description string, remember bool) (string, error) {
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

func SessionUpdateToken(session *Session) (string, error) {
	var expiresAt int64
	if session.Remember {
		expiresAt = time.Now().Add(time.Duration(conf.GetUint("jwt.session.long")) * time.Minute).Unix()
	} else {
		expiresAt = time.Now().Add(time.Duration(conf.GetUint("jwt.session.default")) * time.Minute).Unix()
	}
	newSessionToken, err := utils.GenerateRandomString(conf.GetInt("jwt.session.length"))
	if err != nil {
		return "", err
	}
	updatedFields := Session{
		ExpiresAt: expiresAt,
		UpdateAt:  time.Now().Unix(),
		Token:     newSessionToken,
	}
	if err := SessionUpdate(session, updatedFields); err != nil {
		return "", err
	}
	return newSessionToken, nil
}

func SessionGetAllByUserID(userId ulid.ULID) ([]Session, error) {
	var sessions []Session
	tx := db.Where("user_id = ?", userId).Find(&sessions)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return sessions, nil
}

func SessionGetByToken(token string) (*Session, error) {
	var session Session
	tx := db.Where("token = ?", token).Find(&session)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &session, nil
}

func SessionGetByID(id ulid.ULID) (*Session, error) {
	var session Session
	tx := db.Where("id = ?", id).Find(&session)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &session, nil
}

func SessionUpdate(session *Session, updatedFields Session) error {
	tx := db.Model(&session).Updates(updatedFields)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
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
