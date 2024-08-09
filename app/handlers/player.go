package handlers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/exp/rand"
	"ravi-test-example.com/app/models"
	utils "ravi-test-example.com/app/utilities"
)

const (
	tableName      = "players"
	rtpStatisticID = "rtp_statistic"

	databaseName        = "mongodb_data"
	maxPlaysBeforeCheck = 1000000000 // 1 billion plays
	maxRTP              = 0.975      // 97.5%
)

// Handlers for CreatePlayer
func CreatePlayer(c *fiber.Ctx) error {
	// Parse request body
	var req utils.CreatePlayerRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Bad Request",
		})
	}

	//set default value
	req.Status = "Active"

	// Insert player into MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	playerCollection := utils.GetMongoClient().Database(databaseName).Collection(tableName)
	result, err := playerCollection.InsertOne(ctx, req)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create player",
		})
	}

	// Return success response with created player ID
	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message":  "Player created successfully",
		"status":   "Success",
		"playerID": result.InsertedID,
	})
}

// Handlers for retrieve player details from MongoDB
func GetPlayer(c *fiber.Ctx) error {
	var player models.Player
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Convert playerID string to ObjectID
	objID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid player ID",
		})
	}

	playerCollection := utils.GetMongoClient().Database(databaseName).Collection(tableName)
	if err := playerCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&player); err != nil {
		if err == errors.New("player not found") {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{
				"error": "Player not found",
			})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get player",
		})
	}
	return c.Status(http.StatusOK).JSON(player)
}

// Handlers for SuspendPlayer
func SuspendPlayer(c *fiber.Ctx) error {
	objID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid player ID",
		})
	}

	// Update player in MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	playerCollection := utils.GetMongoClient().Database(databaseName).Collection(tableName)
	result, err := playerCollection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": bson.M{"status": "Suspended"}})
	if err != nil {
		log.Printf("Failed to suspend player: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to suspend player",
		})
	}

	// Check if player was found and updated
	if result.ModifiedCount == 0 {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Player not found",
		})
	}

	// Return success response
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Player suspended successfully",
	})
}

// Handlers for PlaySlotMachine
func PlaySlotMachine(c *fiber.Ctx) error {
	objID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid player ID",
		})
	}

	// Get player details from MongoDB
	var player models.Player
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	playerCollection := utils.GetMongoClient().Database(databaseName).Collection(tableName)
	err = playerCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&player)
	if err != nil {
		// Check if player not found
		if err == mongo.ErrNoDocuments {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{
				"error": "Player not found",
			})
		}
		// Handle other errors
		log.Printf("Failed to find player: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to find player",
		})
	}

	// Simulate slot machine gameplay
	betAmount := 10
	symbols := []string{"cherry", "bell", "lemon", "orange", "star", "red"}

	// Randomly select 3 symbols
	var resultSymbols []string
	for i := 0; i < 3; i++ {
		resultSymbols = append(resultSymbols, symbols[rand.Intn(len(symbols))])
	}

	// Calculate win amount based on result (simplified for example)
	winAmount := 0
	if resultSymbols[0] == resultSymbols[1] && resultSymbols[1] == resultSymbols[2] {
		winAmount = betAmount * 10 // Example: Win 10 times the bet amount for three matching symbols
	} else {
		winAmount = 0 // No win for this example
	}

	// Update player credits and log game outcome
	updatedCredits := player.Credits - betAmount + winAmount
	update := bson.M{
		"$set": bson.M{
			"credits": updatedCredits,
		},
	}

	// Update player credits in MongoDB
	_, err = playerCollection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		log.Printf("Failed to update player credits: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update player credits",
		})
	}

	// Update RTP statistics
	err = updateRTPStatistics(ctx, winAmount, betAmount)
	if err != nil {
		log.Printf("Failed to update RTP statistics: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update RTP statistics",
		})
	}

	// Return JSON response with game outcome
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message":         "Slot machine game played successfully",
		"player_id":       c.Params("id"),
		"result":          resultSymbols,
		"win_amount":      winAmount,
		"updated_credits": updatedCredits,
	})
}

func updateRTPStatistics(ctx context.Context, winAmount, betAmount int) error {
	// Get current RTP statistics from MongoDB
	var rtpStat models.RTPStatistic
	rtpCollection := utils.GetMongoClient().Database(databaseName).Collection(rtpStatisticID)
	err := rtpCollection.FindOne(ctx, bson.M{"_id": rtpStatisticID}).Decode(&rtpStat)

	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}

	// Initialize RTP statistics if it doesn't exist
	if err == mongo.ErrNoDocuments {
		rtpStat = models.RTPStatistic{
			ID:          rtpStatisticID,
			TotalPlays:  0,
			TotalWins:   0,
			TotalLosses: 0,
			TotalPayout: 0,
			LastUpdated: time.Now(),
		}
	}

	// Update RTP statistics based on game outcome
	rtpStat.TotalPlays++
	if winAmount > 0 {
		rtpStat.TotalWins++
		rtpStat.TotalPayout += winAmount
	} else {
		rtpStat.TotalLosses++
	}

	// Update last updated timestamp
	rtpStat.LastUpdated = time.Now()

	// Update RTP statistics in MongoDB
	_, err = rtpCollection.ReplaceOne(ctx, bson.M{"_id": rtpStatisticID}, rtpStat, options.Replace().SetUpsert(true))
	if err != nil {
		return err
	}

	// Check and adjust RTP if necessary
	if rtpStat.TotalPlays >= maxPlaysBeforeCheck {
		currentRTP := float64(rtpStat.TotalPayout) / float64(rtpStat.TotalPlays*betAmount)

		if currentRTP > maxRTP {
			// Implement logic to adjust RTP (e.g., reduce payouts)
			log.Printf("RTP exceeds %.2f%%. Adjusting payouts.", maxRTP*100)
			// Example: Reduce payout ratios or take corrective action
		}
	}
	return nil
}
