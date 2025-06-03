package models

import (
	"log"
	"github.com/dangLuan01/restapi_go/config"
	"github.com/dangLuan01/restapi_go/entities"
	"github.com/doug-martin/goqu/v9"
)
func GetAllGenreHome() ([]entities.Genre, error)  {
	var listGenre []entities.Genre

	query := config.DB.From("genres").
	Select(
		goqu.I("name"), 
		goqu.I("slug")).
	Order(
		goqu.I("position").Asc())

	if err := query.ScanStructs(&listGenre); err != nil {
		return []entities.Genre{}, err
	}

	return listGenre, nil
}
func GetAllGenre() []entities.GenreWithMovies {
	var listGenre []entities.GenreWithMovies
	query := config.DB.From("genres").
		Select(
			goqu.I("name"),
			goqu.I("slug"),
			goqu.I("image"),
			goqu.COUNT(goqu.I("mg.movie_id")).As("total_movies"),
		).
		LeftJoin(goqu.T("movie_genres").As("mg"),goqu.On(goqu.Ex{
			"genres.id": goqu.I("mg.genre_id"),
		})).
		Where(
			goqu.Ex{"status": 1,
		}).
		GroupBy(
			"genres.name",
			"genres.slug",
			"genres.image",
		).
		Order(goqu.I("position").Asc())
	err := query.ScanStructs(&listGenre)
	if err != nil {
		log.Println("Error fetching genres:", err)
		return nil
	}
	return listGenre
}