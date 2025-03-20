package redis

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

var rdb *redis.Client

func InitRedis(appConfig AppConfig) {
	rdb = redis.NewClient(&redis.Options{Addr: appConfig.RedisAddress, Password: appConfig.RedisPassword, DB: appConfig.RedisDb, PoolSize: 10, MinIdleConns: 2, IdleTimeout: 5 * time.Minute})

	ctx := context.Background()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("could not connect to redis: %v", err)
	}

	log.Println("Connected to redis")

	// Healthcheck monitor
	go monitorHealth()
}

// Get Client
func GetRedisClient() *redis.Client {
	return rdb
}

func monitorHealth() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		_, err := rdb.Ping(ctx).Result()
		cancel()
		if err != nil {
			log.Printf("Redis health check failed")
		}
	}
}
