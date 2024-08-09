package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func LivenessCheck(c *fiber.Ctx) error {
	// Check if application is live
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status": "ok",
	})
}
