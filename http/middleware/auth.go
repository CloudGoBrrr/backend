package middleware

import (
	"cloudgobrrr/database/models"
	"cloudgobrrr/structs"
	"cloudgobrrr/utils"
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
)

var errUnauthorized = utils.ConvertStringsToErrorResponse("unauthorized")

func Auth(c *fiber.Ctx) error {
	// get auth header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		c.Set("WWW-Authenticate", "Basic realm=\"User Visible Realm\", charset=\"UTF-8\"")
		return c.Status(fiber.StatusUnauthorized).JSON(errUnauthorized)
	}

	// split auth header
	authHeaderSplit := strings.Split(authHeader, " ")
	if len(authHeaderSplit) != 2 {
		return c.Status(fiber.StatusUnauthorized).JSON(errUnauthorized)
	}

	// get token type
	tokenType := strings.ToLower(authHeaderSplit[0])

	// var declarations
	var user *structs.User
	var err error

	if tokenType == "bearer" {
		// bearer token
		user, err = bearerAuth(authHeaderSplit[1])
	} else if tokenType == "basic" {
		// basic token
		user, err = basicAuth(authHeaderSplit[1])
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON(errUnauthorized)
	}

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(errUnauthorized)
	}

	// set locals
	c.Locals("user", user)

	return c.Next()
}

func bearerAuth(token string) (*structs.User, error) {
	return utils.DecodeJWT(token)
}

func basicAuth(token string) (*structs.User, error) {
	// get token from database
	tokenModel, err := models.TokenGetByEncoded(token)
	if err != nil {
		return nil, err
	}
	if tokenModel.ID == utils.EmptyULID {
		return nil, errors.New("token not found")
	}

	// get user by id
	user, err := models.UserGetById(tokenModel.UserID)
	if err != nil {
		return nil, err
	}
	if user.ID == utils.EmptyULID {
		return nil, errors.New("user not found")
	}

	return &structs.User{ID: user.ID, Username: user.Username, IsAdmin: user.IsAdmin}, nil
}
