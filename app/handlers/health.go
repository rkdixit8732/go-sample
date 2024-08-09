package handlers

import (
	"github.com/gofiber/fiber/v2"
)

// Health Checks of MongoDB, Redis, and overall system
func HealthCheck(c *fiber.Ctx) error {
	return ReadinessCheck(c)
}
