package models

import (
	"github.com/dangLuan01/restapi_go/config"
	"github.com/dangLuan01/restapi_go/entities"
	"github.com/doug-martin/goqu/v9"
)
type MovieRaw struct {
	Id    			int    `json:"id"`
	Name  			string `json:"name"`
	Slug  			string `json:"slug"`
	Type  			string `json:"type"`
	Release_date 	int    `json:"release_date"`
	Rating			float64 `json:"rating"`
	Thumb 			string `json:"thumb"`
	Poster			string `json:"poster"`
}

func GetAllMovieHot() []entities.Movie {
	var listHotMovie []entities.Movie
	err := config.DB.From("movies").
	LeftJoin(
		goqu.T("movie_images").As("mi"),
		goqu.On(
			goqu.I("movies.id").Eq(goqu.I("mi.movie_id")),
		),
	).
	Where(
		goqu.Ex{
			"movies.hot": 1,
			"mi.is_thumbnail": 1,
		},
	).
	Select(
		goqu.I("movies.id"),
		goqu.I("movies.name"),
		goqu.I("movies.slug"),
		goqu.I("movies.type"),
		goqu.I("movies.release_date"),
		goqu.I("movies.rating"),
		goqu.Func("CONCAT", goqu.I("mi.path"), goqu.I("mi.image")).As("thumb"),
	).
	Limit(10)

	var hot []MovieRaw
	if err := err.ScanStructs(&hot); err != nil {
		return nil
	}

	for _, item := range hot {
		listHotMovie = append(listHotMovie, entities.Movie{
			Id: item.Id,
			Name: item.Name,
			Slug: item.Slug,
			Type: item.Type,
			Release_date: item.Release_date,
			Rating: item.Rating,
			Image: entities.Image{
				Thumb: item.Thumb,
			},
		})
	}
	return listHotMovie
}
func GetAllMovie(page, pageSize int) (entities.PaginatedMovies) {
	var listMovie []entities.Movie
	err := config.DB.From("movies").
	LeftJoin(
		goqu.T("movie_images").As("mi"),
		goqu.On(
			goqu.I("movies.id").Eq(goqu.I("mi.movie_id")),
		),
	).
	Where(
		goqu.Ex{
			"movies.hot": nil,
			"mi.is_thumbnail": 0,
		},
	).
	Select(
		goqu.I("movies.id"),
		goqu.I("movies.name"),
		goqu.I("movies.slug"),
		goqu.I("movies.type"),
		goqu.I("movies.release_date"),
		goqu.I("movies.rating"),
		goqu.Func("CONCAT", goqu.I("mi.path"), goqu.I("mi.image")).As("poster"),
	).
	Limit(uint(pageSize)).Offset(uint((page - 1) * pageSize))

	var movie []MovieRaw
	if err := err.ScanStructs(&movie); err != nil {
		return entities.PaginatedMovies{}
	}
	for _, item := range movie {
		listMovie = append(listMovie, entities.Movie{
			Id: item.Id,
			Name: item.Name,
			Slug: item.Slug,
			Type: item.Type,
			Release_date: item.Release_date,
			Rating: item.Rating,
			Image: entities.Image{
				Poster: item.Poster,
			},
		})
	}
	return entities.PaginatedMovies{
		Data: listMovie,
		Page: page,
		PageSize: pageSize,
	}	
}