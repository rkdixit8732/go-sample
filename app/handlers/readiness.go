package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	utils "ravi-test-example.com/app/utilities"
)

func ReadinessCheck(c *fiber.Ctx) error {
	// Check MongoDB connection status
	mongoDBStatus := "OK"
	if err := utils.CheckMongoDBConnection(); err != nil {
		mongoDBStatus = "Error"
	}
	// Check Redis connection status
	redisStatus := "OK"
	if err := utils.CheckRedisConnection(); err != nil {
		redisStatus = "Error"
	}
	// Determine overall app status based on connections
	appStatus := "OK"
	if mongoDBStatus == "Error" || redisStatus == "Error" {
		appStatus = "Error"
	}
	// Prepare response
	response := utils.ReadinessCheckResponse{
		MongoDBStatus: mongoDBStatus,
		RedisStatus:   redisStatus,
		AppStatus:     appStatus,
	}
	// Return readiness check response
	return c.Status(http.StatusOK).JSON(response)
}
