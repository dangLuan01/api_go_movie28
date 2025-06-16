package cache

import (
	"context"
	"encoding/json"
	"log"
	"time"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type redisCache struct {
	client    *redis.Client
	expires   time.Duration
	available bool
}

func NewRedisCache(host string, db int, exp time.Duration) DataCache {
	client := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: "",
		DB:       db,
	})
	available := true
	if _, err := client.Ping(ctx).Result(); err != nil {
		log.Printf("⚠️ Redis không khả dụng: %v", err)
		available = false
		return nil
	}

	return &redisCache{
		client:    client,
		expires:   exp,
		available: available,
	}
}

func (cache *redisCache) Set(key string, value any) {
	if !cache.available {
		return
	}
	data, err := json.Marshal(value)
	if err != nil {
		log.Printf("❌ Error marshaling JSON: %v", err)
		return
	}
	err = cache.client.Set(ctx, key, data, cache.expires*time.Second).Err()
	if err != nil {
		log.Printf("❌ Error setting cache: %v", err)
	}
}
func (cache *redisCache) Get(key string, dest any) bool {
	if !cache.available {
		return false
	}
	val, err := cache.client.Get(ctx, key).Result()
	if err != nil {
		return false
	}
	//var movie entities.Movie
	if err := json.Unmarshal([]byte(val), dest); err != nil {
		log.Printf("❌ Error unmarshaling JSON: %v", err)
		return false
	}
	return true
}
