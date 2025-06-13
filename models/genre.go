package models

import (
	"fmt"
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
        goqu.I("genres.name"),
        goqu.I("genres.slug"),
        goqu.I("genres.image"),
        goqu.COUNT(goqu.I("mg.movie_id")).As("total_movies"),
    ).
    InnerJoin(goqu.T("movie_genres").As("mg"), goqu.On(goqu.I("genres.id").Eq(goqu.I("mg.genre_id")))).
    Where(goqu.I("genres.status").Eq(1)).
    GroupBy(goqu.I("genres.id")).
    Order(goqu.I("genres.position").Asc()).Limit(70)

	err := query.ScanStructs(&listGenre)
	if err != nil {
		fmt.Println("Error fetching genres:", err)
		return nil
	}
	return listGenre
}
func GetItemGenre(slug string, page, pageSize int) (entities.PaginatedMovies, error ) {
	var listMovie []entities.Movie
	var genreInfo entities.Genre
	var movie []entities.MovieRaw
	var totalPages int64
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
	posterSubquery := config.DB.From(goqu.T("movie_images").As("mi")).
		Where(
			goqu.I("mi.movie_id").Eq(goqu.I("m.id")),
			goqu.I("mi.is_thumbnail").Eq(0),
		).
		Select(goqu.Func("CONCAT", goqu.I("mi.path"), goqu.I("mi.image"))).
		Limit(1)
	queryMovie := config.DB.
	From(goqu.T("movies").As("m")).
	LeftJoin(
		goqu.T("movie_genres").As("mg"),
		goqu.On(goqu.I("m.id").Eq(goqu.I("mg.movie_id"))),
	).
	Where(
		goqu.Ex{"mg.genre_id":genreInfo.Id},
	).
	Select(
		goqu.I("m.name"),
		goqu.I("m.slug"),
		goqu.I("m.type"),
		goqu.I("m.release_date"),
		goqu.I("m.rating"),
		posterSubquery.As("poster"),
	).
	Order(goqu.I("m.updated_at").Desc())
	
	count, err := queryMovie.Count()
	if err != nil {
		return entities.PaginatedMovies{}, err
	}
	totalPages = count/int64(pageSize)
	if totalPages == 0 {
		totalPages = 1
	}
	if err := queryMovie.Limit(uint(pageSize)).Offset(uint((page - 1) * pageSize)).ScanStructs(&movie); err != nil {
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
		Movie: listMovie,
		Page: page,
		PageSize: pageSize,
		TotalPages: int(totalPages),
	}, nil
}