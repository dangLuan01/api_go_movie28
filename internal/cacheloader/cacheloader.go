package cacheloader

import (
	"sync"
	"github.com/dangLuan01/api_go_movie28/cache"
)

var (
	once        sync.Once
	movieCache  cache.MovieCache
)
func GetMovieCache() cache.MovieCache {
	once.Do(func() {
		movieCache = cache.NewRedisCache("localhost:6379", 0, 300) // TTL = 300 giây
	})
	return movieCache
}