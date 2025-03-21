package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	ErrCacheMiss = errors.New("cache miss: key not found")
	redisClient  *redis.Client
)

// SetupRedisClient initializes the Redis client and makes it available globally.
func InitCache(appRedisClient *redis.Client) {
	redisClient = appRedisClient
}

// GetRedisClient returns the Redis client instance.
func GetRedisClient() *redis.Client {
	if redisClient == nil {
		panic("Redis client not initialized. Call InitCache() first.")
	}
	return redisClient
}

// Set cache data with expiration
func Set(key string, value interface{}, expiration time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Convert any kind of data to JSON
	jsonData, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("error converting data: %v", err)
	}

	// Set data in Redis
	err = GetRedisClient().Set(ctx, key, jsonData, expiration).Err()
	if err != nil {
		return fmt.Errorf("error setting data in redis: %v", err)
	}

	return nil
}

// Get cache data and deserialize it into the target interface
func Get(key string, target interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	data, err := GetRedisClient().Get(ctx, key).Result()
	if err == redis.Nil {
		return ErrCacheMiss
	} else if err != nil {
		return fmt.Errorf("error getting data from redis: %v", err)
	}

	// Convert JSON string back to original data type
	err = json.Unmarshal([]byte(data), target)
	if err != nil {
		return fmt.Errorf("error converting data: %v", err)
	}
	return nil
}

// Delete key from cache
func Delete(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := GetRedisClient().Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("error deleting key from redis: %v", err)
	}

	return nil
}
