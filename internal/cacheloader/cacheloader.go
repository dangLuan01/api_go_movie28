package cacheloader

import (
	"os"
	"sync"
	"time"
	"github.com/dangLuan01/api_go_movie28/cache"
)

var (
	once        sync.Once
	dataCache  cache.DataCache
)
func GetCache(db int, exp time.Duration) cache.DataCache {
	host 	:= os.Getenv("REDIS_HOST")
	once.Do(func() {
		dataCache = cache.NewRedisCache(host, db, exp*time.Second)
	})
	return dataCache
}