package brokers

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
)

// Global variable for RedisPubSubService instance
var redisPubSubInstance *RedisPubSubService

// RedisPubSubService manages Redis Pub/Sub.
type RedisPubSubService struct {
	client *redis.Client
	ctx    context.Context
}

// InitRedisPubSub initializes the Redis Pub/Sub service (called from main.go)
func InitRedisPubSub(redisClient *redis.Client) {
	if redisPubSubInstance == nil {
		redisPubSubInstance = &RedisPubSubService{
			client: redisClient,
			ctx:    context.Background(),
		}
		log.Println("Redis Pub/Sub Service initialized")
	}
}

// GetRedisPubSubService provides the singleton instance.
func GetRedisPubSubService() *RedisPubSubService {
	if redisPubSubInstance == nil {
		log.Fatal("Redis Pub/Sub Service is not initialized. Call InitRedisPubSub() in main.go first.")
	}
	return redisPubSubInstance
}

// Publish sends a message to a Redis channel.
func (r *RedisPubSubService) Publish(channel, message string) error {
	err := r.client.Publish(r.ctx, channel, message).Err()
	if err != nil {
		return err
	}
	return nil
}

// Subscribe listens for messages on a Redis channel.
func (r *RedisPubSubService) Subscribe(channel string, handler func(msg string)) {
	pubsub := r.client.Subscribe(r.ctx, channel)
	ch := pubsub.Channel()

	log.Printf("Subscribed to channel: %s", channel)
	for msg := range ch {
		handler(msg.Payload)
	}
}

// Close closes the Redis Pub/Sub service.
func (r *RedisPubSubService) Close() {
	r.client.Close()
	log.Println("Redis Pub/Sub Service closed")
}
