package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"os"
	"strconv"
	"time"
)

var ctx = context.Background()
var Cache *redis.Client

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

	_, err = Cache.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("could not connect to Redis: %v", err)
	}

	log.Println("Connected to Redis")
	return nil
}

func Set(key string, value interface{}, expiration time.Duration) error {
	if Cache == nil {
		return fmt.Errorf("redis client is not initialized")
	}
	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return Cache.Set(ctx, key, jsonData, expiration).Err()
}

func Get(key string, dest interface{}) error {
	if Cache == nil {
		return fmt.Errorf("redis client is not initialized")
	}
	jsonData, err := Cache.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(jsonData), dest)
}

func Delete(key string) error {
	if Cache == nil {
		return fmt.Errorf("redis client is not initialized")
	}
	return Cache.Del(ctx, key).Err()
}

func FlushDB() error {
	if Cache == nil {
		return fmt.Errorf("redis client is not initialized")
	}
	return Cache.FlushDB(ctx).Err()
}

func FlushAll() error {
	if Cache == nil {
		return fmt.Errorf("redis client is not initialized")
	}
	return Cache.FlushAll(ctx).Err()
}
