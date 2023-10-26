package utils

import (
	"cloudgobrrr/structs"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secret []byte

var ErrInvalidSigningMethod = errors.New("invalid signing method")

// GenerateJWT generates a JWT token and returns it
func GenerateJWT(username string, userId uint, isAdmin bool) (string, error) {
	var token *jwt.Token
	if conf.GetString("jwt.signingMethod") == "HS256" {
		token = jwt.New(jwt.SigningMethodHS256)
	} else {
		return "", ErrInvalidSigningMethod
	}

	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = username
	claims["userId"] = userId
	claims["exp"] = time.Now().Add(time.Duration(conf.GetUint("jwt.expiration")) * time.Minute).Unix()
	claims["isAdmin"] = isAdmin

	return token.SignedString(secret)
}

// DecodeJWT decodes a JWT token and returns the userId and isAdmin
func DecodeJWT(token string) (*structs.User, error) {
	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return nil, err
	}

	claims := t.Claims.(jwt.MapClaims)
	username := claims["username"].(string)
	userId := uint(claims["userId"].(float64))
	isAdmin := claims["isAdmin"].(bool)
	exp := int64(claims["exp"].(float64))

	if exp < time.Now().Unix() {
		err = errors.New("token expired")
		return nil, err
	}

	return &structs.User{ID: userId, Username: username, IsAdmin: isAdmin}, nil
}
