package route

import (
	"github.com/gofiber/fiber/v2"

	"UTS/helper"
)

type Dependencies struct {
	
}

func Register(app *fiber.App, deps Dependencies) {
	api := app.Group("/api/v1")

	api.Get("/health", func(c *fiber.Ctx) error {
		return helper.Success(c, fiber.StatusOK, "server berjalan", nil)
	})
}