package models

import (
	"github.com/dangLuan01/restapi_go/entities"
)
//var listMovie = make([]entities.Movie, 0)

func GetAllMovie() []entities.Movie {
	listMovie := []entities.Movie{
		{
			Id:   1,
			Name: "The Shawshank Redemption",
			Slug: "the-shawshank-redemption",
			Year: 1994,
			Image: "https://example.com/image1.jpg",
		},
	} 
	return listMovie
}
func GetAllCategory() []entities.Category  {
	
	listCategory := []entities.Category{
		{
			Id:   1,
			Name: "Movie",
			Slug: "movie",
		},
		{
			Id:   2,
			Name: "TV Series",
			Slug: "tv-series",
		},
	}
	return listCategory
}