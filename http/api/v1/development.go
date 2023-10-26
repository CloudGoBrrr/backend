package v1

import (
	"cloudgobrrr/database"
	"cloudgobrrr/http/middleware"
	"cloudgobrrr/http/request"

	"github.com/gofiber/fiber/v2"
)

/*
 * Routes defined in this file:
 * - GET /dev/ping
 * - GET /dev/ping_database
 * - POST /dev/test_validation
 */

// ToDo: remove this later

// setupDevelopment sets up all routes for development
func setupDevelopment() {
	if !conf.GetBool("development") {
		return
	}

	dev := app.Group("/dev")

	dev.Get("/ping", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "pong",
		})
	})

	dev.Get("/ping_database", func(c *fiber.Ctx) error {
		db, _ := database.Get().DB()
		err := db.Ping()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "database connection failed",
			})
		}
		return c.JSON(fiber.Map{
			"message": "database connection successful",
		})
	})

	dev.Post("/test_validation", func(c *fiber.Ctx) error {
		req := new(request.TestPost)

		if err := c.BodyParser(req); err != nil {
			return err
		}

		if errs := val.Validate(req); len(errs) > 0 {
			return c.Status(fiber.StatusBadRequest).JSON(val.ConvertToResponse(errs))
		}

		return c.JSON(fiber.Map{
			"message": "success",
		})
	})

	dev.Get("/test_auth", middleware.Auth, func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "success",
		})
	})
}
