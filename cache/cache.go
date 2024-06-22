package cache

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"os"
	"strconv"
	"time"
)

var ctx = context.Background()
var Cache *redis.Client

// InitCache initializes the Redis client using environment variables.
func InitCache() error {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		return fmt.Errorf("REDIS_ADDR environment variable is not set")
	}

	password := os.Getenv("REDIS_PASSWORD")
	dbStr := os.Getenv("REDIS_DB")
	db, err := strconv.Atoi(dbStr)
	if err != nil {
		return fmt.Errorf("invalid REDIS_DB value: %v", err)
	}

	Cache = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// Check if the connection is established
	_, err = Cache.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("could not connect to Redis: %v", err)
	}

	log.Println("Connected to Redis")
	return nil
}

// Set sets a key-value pair in Redis with an expiration time.
func Set(key string, value interface{}, expiration time.Duration) error {
	if Cache == nil {
		return fmt.Errorf("redis client is not initialized")
	}
	return Cache.Set(ctx, key, value, expiration).Err()
}

// Get retrieves the value for a given key from Redis.
func Get(key string) (string, error) {
	if Cache == nil {
		return "", fmt.Errorf("redis client is not initialized")
	}
	return Cache.Get(ctx, key).Result()
}

// Delete removes a key from Redis.
func Delete(key string) error {
	if Cache == nil {
		return fmt.Errorf("redis client is not initialized")
	}
	return Cache.Del(ctx, key).Err()
}
