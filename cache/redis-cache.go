package cache

import (
	"context"
	"encoding/json"
	"log"
	"time"
	"github.com/dangLuan01/api_go_movie28/entities"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type redisCache struct {
	host    string
	db      int
	expires time.Duration
}
func NewRedisCache(host string, db int, exp time.Duration) MovieCache {
	return &redisCache{
		host: 		host,
		db: 		db,
		expires: 	exp,
	}
}

func (cache *redisCache) getClient() *redis.Client  {
	return redis.NewClient(&redis.Options{
		Addr: 		cache.host,
		Password: 	"",
		DB: 		cache.db,
	})
}

func (cache *redisCache) Set(key string, value *entities.Movie){
	
	client 		:= cache.getClient()
	json, err 	:= json.Marshal(value)
	if err != nil {
		log.Printf("Error json:%v", err)
		return
	}
	error := client.Set(ctx, key, json, cache.expires*time.Second).Err()
	if error != nil {
		log.Printf("Error Set cache:%v", error)
		return
	}
}

func (cache *redisCache) Get(key string) *entities.Movie{
	client 		:= cache.getClient()
	val, err 	:= client.Get(ctx, key).Result()
	if err != nil {
		return nil
	}
	movie := entities.Movie{}
	err = json.Unmarshal([]byte(val), &movie)
	if err != nil {
		log.Printf("Error get json cache:%v", err)
		return nil
	}
	return &movie
}
