package v1

import (
	"cloudgobrrr/database/models"
	"cloudgobrrr/http/request"
	"cloudgobrrr/http/response"
	"cloudgobrrr/utils"
	"time"

	"github.com/gofiber/fiber/v2"
)

/*
 * Routes defined in this file:
 * - POST /auth/signin
 * - POST /auth/signup
 * - POST /auth/signout
 * - POST /auth/refresh
 */

var errInvalidCredentials response.Error = utils.ConvertStringsToErrorResponse("invalid credentials")
var errAlreadyTaken response.Error = utils.ConvertStringsToErrorResponse("user already taken")
var errSignupDisabled response.Error = utils.ConvertStringsToErrorResponse("signup disabled")
var errSessionNotFound response.Error = utils.ConvertStringsToErrorResponse("session not found")

// setupAuth sets up all routes for auth
func setupAuth() {
	auth := app.Group("/auth")

	auth.Post("/signin", authSignin)
	auth.Post("/signup", authSignup)
	auth.Post("/signout", authSignout)
	auth.Post("/refresh", authRefresh)
}

// authSignin handles the signin route
func authSignin(c *fiber.Ctx) error {
	req := new(request.AuthSignin)

	// parse request body
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ConvertErrorsToErrorResponse(err))
	}

	// validate request
	if errs := val.Validate(req); len(errs) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(val.ConvertToResponse(errs))
	}

	var user *models.User
	var token *models.Token
	var err error

	// check if req.User is email or username
	if utils.IsEmail(req.User) {
		// get user by email
		user, err = models.UserGetByEmail(req.User)
	} else {
		// get user by username
		user, err = models.UserGetByUsername(req.User)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(errInvalidCredentials)
		}
		token, err = models.TokenGetByUsernameAndPassword(req.User, req.Password)
	}

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(errInvalidCredentials)
	}

	// check if token exists
	if token != nil && token.ID != utils.EmptyULID {
		// get user by id
		user, err = models.UserGetById(token.UserID)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(errInvalidCredentials)
		}
	}

	// check if user exists
	if user.ID == utils.EmptyULID {
		return c.Status(fiber.StatusUnauthorized).JSON(errInvalidCredentials)
	}

	// check if password matches
	if match, err := utils.PasswordCheckHash(req.Password, user.Password); err != nil || !match {
		return c.Status(fiber.StatusUnauthorized).JSON(errInvalidCredentials)
	}

	// generate jwt token
	jwtToken, err := utils.GenerateJWT(user.Username, user.ID, user.IsAdmin)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ConvertErrorsToErrorResponse(err))
	}

	// create session
	httpAgent := c.Get("User-Agent")
	if httpAgent == "" {
		httpAgent = "Unknown"
	}
	session, err := models.SessionCreate(user.ID, httpAgent, req.Remember)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ConvertErrorsToErrorResponse(err))
	}

	// return response
	return c.JSON(response.Success{Success: true, Data: response.AuthToken{Token: jwtToken, Session: session}})
}

// authSignup handles the signup route
func authSignup(c *fiber.Ctx) error {
	if conf.GetBool("auth.signup") {
		return c.Status(fiber.StatusNotImplemented).JSON(errSignupDisabled)
	}

	req := new(request.AuthSignup)

	// parse request body
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ConvertErrorsToErrorResponse(err))
	}

	// validate request
	if errs := val.Validate(req); len(errs) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(val.ConvertToResponse(errs))
	}

	// check if username is taken
	if _, err := models.UserGetByUsername(req.Username); err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(errAlreadyTaken)
	}

	// check if email is taken
	if _, err := models.UserGetByEmail(req.Email); err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(errAlreadyTaken)
	}

	// create user
	if err := models.UserCreate(req.Username, req.Email, req.Password, false); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ConvertErrorsToErrorResponse(err))
	}

	// return response
	return c.JSON(response.Success{Success: true})
}

// authSignout handles the signout route
func authSignout(c *fiber.Ctx) error {
	req := new(request.AuthSignout)

	// parse request body
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ConvertErrorsToErrorResponse(err))
	}

	// validate request
	if errs := val.Validate(req); len(errs) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(val.ConvertToResponse(errs))
	}

	// delete session
	if err := models.SessionDeleteByToken(req.Session); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ConvertErrorsToErrorResponse(err))
	}

	// return response
	return c.JSON(response.Success{Success: true, Data: "signed out"})
}

// authRefresh handles the refresh route
func authRefresh(c *fiber.Ctx) error {
	req := new(request.AuthRefresh)

	// parse request body
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ConvertErrorsToErrorResponse(err))
	}

	// validate request
	if errs := val.Validate(req); len(errs) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(val.ConvertToResponse(errs))
	}

	// get session
	session, err := models.SessionGetByToken(req.Session)
	if err != nil {
		if err.Error() == "record not found" {
			return c.Status(fiber.StatusBadRequest).JSON(errSessionNotFound)
		}
		return c.Status(fiber.StatusBadRequest).JSON(utils.ConvertErrorsToErrorResponse(err))
	}

	// check if session is expired
	if session.ExpiresAt < time.Now().Unix() {
		// ToDo: delete session
		return c.Status(fiber.StatusBadRequest).JSON(errSessionNotFound)
	}

	// get user
	user, err := models.UserGetById(session.UserID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ConvertErrorsToErrorResponse(err))
	}

	// generate jwt token
	token, err := utils.GenerateJWT(user.Username, user.ID, user.IsAdmin)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ConvertErrorsToErrorResponse(err))
	}

	// update session
	newSessionToken, err := models.SessionUpdateToken(session)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ConvertErrorsToErrorResponse(err))
	}

	// return response
	return c.JSON(response.Success{Success: true, Data: response.AuthToken{Token: token, Session: newSessionToken}})
}
