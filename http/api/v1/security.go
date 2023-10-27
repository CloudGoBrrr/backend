package v1

import (
	"cloudgobrrr/database/models"
	"cloudgobrrr/http/middleware"
	"cloudgobrrr/http/request"
	"cloudgobrrr/http/response"
	"cloudgobrrr/structs"
	"cloudgobrrr/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/oklog/ulid/v2"
)

func setupSecurity() {
	security := app.Group("/security")

	security.Get("/token", middleware.Auth, securityTokenGetAll)
	security.Post("/token", middleware.Auth, securityTokenCreate)
	security.Delete("/token/:identifier", middleware.Auth, securityTokenDelete)
	security.Get("/session", middleware.Auth, securitySessionGetAll)
	security.Delete("/session/:identifier", middleware.Auth, securitySessionDelete)
}

func securityTokenGetAll(c *fiber.Ctx) error {
	var user *structs.User = c.Locals("user").(*structs.User)

	getTokens, err := models.TokenGetAllByUserID(user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ConvertErrorsToErrorResponse(err))
	}

	// convert to response
	var tokens []response.SecurityTokenElement = make([]response.SecurityTokenElement, 0)
	for _, token := range getTokens {
		tokens = append(tokens, response.SecurityTokenElement{
			ID:          token.ID,
			Description: token.Description,
			CreatedAt:   token.CreatedAt,
			UpdatedAt:   token.UpdateAt,
		})
	}

	return c.JSON(response.Success{
		Success: true,
		Data:    response.SecurityTokenGetAll{Tokens: tokens},
	})
}

// securityTokenCreate handles the token creation route
func securityTokenCreate(c *fiber.Ctx) error {
	req := new(request.SecurityTokenCreate)

	// parse request body
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ConvertErrorsToErrorResponse(err))
	}

	// validate request
	if errs := val.Validate(req); len(errs) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(val.ConvertToResponse(errs))
	}

	var user *structs.User = c.Locals("user").(*structs.User)

	// create token
	id, password, err := models.TokenCreate(user.ID, req.Description)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ConvertErrorsToErrorResponse(err))
	}

	return c.JSON(response.Success{
		Success: true,
		Data:    response.SecurityTokenCreate{ID: id, Password: password},
	})
}

// securityTokenDelete handles the token deletion route
func securityTokenDelete(c *fiber.Ctx) error {
	// get id from params and convert to ulid.ULID
	id, err := ulid.Parse(c.Params("identifier"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ConvertStringsToErrorResponse("invalid identifier"))
	}

	// get token by id
	token, err := models.TokenGetByID(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ConvertStringsToErrorResponse("invalid identifier"))
	}

	// compare token user id with user id from locals
	user := c.Locals("user").(*structs.User)
	if token.UserID != user.ID {
		return c.Status(fiber.StatusUnauthorized).JSON(utils.ConvertStringsToErrorResponse("unauthorized"))
	}

	// delete token
	err = models.TokenDeleteByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ConvertErrorsToErrorResponse(err))
	}

	return c.JSON(response.Success{Success: true, Data: "token deleted"})
}

// securitySessionGetAll handles the session get all route
func securitySessionGetAll(c *fiber.Ctx) error {
	var user *structs.User = c.Locals("user").(*structs.User)

	getSessions, err := models.SessionGetAllByUserID(user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ConvertErrorsToErrorResponse(err))
	}

	// convert to response
	var sessions []response.SecuritySessionElement = make([]response.SecuritySessionElement, 0)
	for _, session := range getSessions {
		sessions = append(sessions, response.SecuritySessionElement{
			ID:          session.ID,
			Description: session.Description,
			CreatedAt:   session.CreatedAt,
			UpdatedAt:   session.UpdateAt,
		})
	}

	return c.JSON(response.Success{
		Success: true,
		Data:    response.SecuritySessionGetAll{Sessions: sessions},
	})
}

// securitySessionDelete handles the session deletion route
func securitySessionDelete(c *fiber.Ctx) error {
	// get id from params and convert to ulid.ULID
	id, err := ulid.Parse(c.Params("identifier"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ConvertStringsToErrorResponse("invalid identifier"))
	}

	// note: we are always returning success here, even if the token does not exist

	// get token by id
	token, err := models.SessionGetByID(id)
	if err != nil {
		return c.JSON(response.Success{Success: true, Data: "token deleted"})
	}

	// compare token user id with user id from locals
	user := c.Locals("user").(*structs.User)
	if token.UserID != user.ID {
		return c.JSON(response.Success{Success: true, Data: "token deleted"})
	}

	// delete token
	err = models.SessionDeleteByID(id)
	if err != nil {
		return c.JSON(response.Success{Success: true, Data: "token deleted"})
	}

	return c.JSON(response.Success{Success: true, Data: "token deleted"})
}
