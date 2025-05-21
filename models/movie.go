package models

import (
	"RestAPI_Go_Basic/entities"
)
var listMovie = make([]entities.Movie, 0)
func GetAllMovie() []*entities.Movie {
	return listMovie
}