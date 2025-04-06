package config

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

var Cache *redis.Client

func InitCache() {
	Cache = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // adjust if necessary
		Password: "",               // no password set
		DB:       0,                // use default DB
	})
	// Check Redis connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := Cache.Ping(ctx).Err(); err != nil {
		Logger.Fatalf("Could not connect to Redis: %v", err)
	}
}
