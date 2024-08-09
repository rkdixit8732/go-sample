package utilities

import (
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	redisDB *redis.Client
)

// InitRedis initializes the Redis client.
func InitRedis() {
	redisDB = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Replace with your Redis server address
		Password: "",               // No password set
		DB:       0,                // Use default DB
	})

	// Example: Ping the Redis server to test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := redisDB.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
}

// CheckRedisConnection checks the connection status of Redis.
func CheckRedisConnection() error {
	if redisDB == nil {
		return errors.New("Redis client is not initialized")
	}
	// Ping Redis server
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if _, err := redisDB.Ping(ctx).Result(); err != nil {
		return err
	}

	return nil
}

func GetRadisClient() *redis.Client {
	return redisDB
}
