package models

import (
	"github.com/dangLuan01/restapi_go/config"
	"github.com/dangLuan01/restapi_go/entities"
)

func GetAllMovie() []entities.Movie {
	var listHotMovie []entities.Movie
	err := config.DB.From("movies").Limit(10).ScanStructs(&listHotMovie)
	if err != nil {
		return nil
	}
	return listHotMovie
}