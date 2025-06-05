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
func GetItemGenre(slug string, page, pageSize int) (entities.PaginatedMovies, error ) {
	var listMovie []entities.Movie
	var genreInfo entities.Genre
	var movie []entities.MovieRaw
	queryGenre := config.DB.From("genres").Where(
		goqu.Ex{"slug": slug},
	).Select(goqu.I("id"), goqu.I("name"))
	
	found, err := queryGenre.ScanStruct(&genreInfo)

	if err != nil {
		return entities.PaginatedMovies{}, err
	}
	if !found {
		return entities.PaginatedMovies{}, nil
	}

	queryMovie := config.DB.
	From("movies").
	LeftJoin(
		goqu.T("movie_images").As("mi"),
		goqu.On(
			goqu.I("movies.id").Eq(goqu.I("mi.movie_id")),
		),
	).
	LeftJoin(
		goqu.T("movie_genres").As("mg"),
		goqu.On(goqu.I("movies.id").Eq(goqu.I("mg.movie_id"))),
	).
	Where(
		goqu.Ex{"mg.genre_id":genreInfo.Id},
		goqu.Ex{"mi.is_thumbnail": 0},
	).
	Select(
		goqu.I("movies.name"),
		goqu.I("movies.slug"),
		goqu.I("movies.type"),
		goqu.I("movies.release_date"),
		goqu.I("movies.rating"),
		goqu.Func("CONCAT", goqu.I("mi.path"), goqu.I("mi.image")).As("poster"),
	).
	Order(goqu.I("movies.updated_at").Desc()).
	Limit(uint(pageSize)).Offset(uint((page - 1) * pageSize))

	if err := queryMovie.ScanStructs(&movie); err != nil {
		return entities.PaginatedMovies{}, err
	}
	for _, item := range movie {
		listMovie = append(listMovie, entities.Movie{
			Name: item.Name,
			Slug: item.Slug,
			Type: item.Type,
			Release_date: item.Release_date,
			Rating: ConvertRating(float32(item.Rating)),
			Image: entities.Image{
				Poster: item.Poster,
			},
			Genres: []entities.Genre{
				{
					Name: genreInfo.Name,
				},
			},
		})
	}
	
	return entities.PaginatedMovies{
		Data: listMovie,
		Page: page,
		PageSize: pageSize,
	}, nil
}