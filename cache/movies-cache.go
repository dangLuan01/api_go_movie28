package cache

import "github.com/dangLuan01/api_go_movie28/entities"

type MovieCache interface {
	Set(key string, value *entities.Movie)
	Get(key string) *entities.Movie
}