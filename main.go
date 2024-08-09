package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	h "ravi-test-example.com/app/handlers"
	utils "ravi-test-example.com/app/utilities"
)

func main() {
	// Initialize MongoDB and Redis connections
	utils.InitMongoDB()
	utils.InitRedis()

	// Create Fiber app
	app := fiber.New()

	// Middleware
	app.Use(logger.New())
	app.Use(recover.New())

	// Routes
	app.Post("/players", h.CreatePlayer)
	app.Get("/players/:id", h.GetPlayer)
	app.Put("/players/:id/suspend", h.SuspendPlayer)
	app.Post("/players/:id/play", h.PlaySlotMachine)

	// Health Checks
	app.Get("/health", h.HealthCheck)
	app.Get("/liveness", h.LivenessCheck)
	app.Get("/readiness", h.ReadinessCheck)

	// Start server
	log.Fatal(app.Listen(":3000"))
}
