package cache

import (
	"context"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

type RedisCacheMemory struct {
	Cache *redis.Client
}

func ConnectCacheMemory() *RedisCacheMemory {
	log.Println("Starting cache connection...")
	client := redis.NewClient(&redis.Options{Addr: os.Getenv("CACHE_HOST") + ":" + os.Getenv("CACHE_PORT"), Password: os.Getenv("CACHE_PASSWORD"), DB: 0})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Failed to connect cache memory: %s", err.Error())
	}

	return &RedisCacheMemory{Cache: client}
}
