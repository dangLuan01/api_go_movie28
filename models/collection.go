package models

import (
	//"encoding/json"
	"github.com/dangLuan01/api_go_movie28/config"
	"github.com/dangLuan01/api_go_movie28/entities"
	"github.com/doug-martin/goqu/v9"
)

func GetAllCollection() ([]entities.Collection, error) {
	var listCollection []entities.Collection
	ds := config.DB.
	From(goqu.T("collections")).
	Where(
		goqu.Ex{
			"status": 1,
		}).
	Select(
		goqu.I("collections.name"), 
		goqu.I("collections.slug"),
		goqu.I("collections.image"),
	)
	if err := ds.ScanStructs(&listCollection); err != nil {
		return listCollection, nil
	}
	return listCollection, nil
}
func GetMovieCollection(slug string, page, pageSize int) (entities.PaginatedMovies, error) {
	var listMovie []entities.Movie
	var movie []entities.MovieRaw
	var totalPages int64
	ds := config.DB.From(goqu.T("collections")).
	LeftJoin(
		goqu.T("movie_collections").As("mc"),
		goqu.On(goqu.I("collections.id").Eq(goqu.I("mc.collection_id"))),
	).
	LeftJoin(
		goqu.T("movies").As("m"),
		goqu.On(goqu.I("mc.movie_id").Eq(goqu.I("m.id"))),
	).
	LeftJoin(
		goqu.T("movie_images").As("mi"),
		goqu.On(
			goqu.I("m.id").Eq(goqu.I("mi.movie_id")),
		),
	).
	LeftJoin(
        goqu.T("movie_genres").As("mg"),
        goqu.On(goqu.I("m.id").Eq(goqu.I("mg.movie_id"))),
    ).
	LeftJoin(
        goqu.T("genres").As("g"),
        goqu.On(goqu.I("mg.genre_id").Eq(goqu.I("g.id"))),
    ).
	Where(goqu.Ex{
		"collections.slug": slug,
		"mi.is_thumbnail": 0,
	}).
	Select(
		goqu.I("m.name"),
		goqu.I("m.slug"),
		goqu.I("m.type"),
		goqu.I("m.release_date"),
		goqu.I("m.rating"),
		goqu.Func("CONCAT", goqu.I("mi.path"), goqu.I("mi.image")).As("poster"),
		goqu.I("g.name").As("genre_name"),
	).Order(goqu.I("m.updated_at").Desc())
	
	count, err := ds.Count()
	if err != nil {
		return entities.PaginatedMovies{}, err
	}
	totalPages = count/int64(pageSize)
	if totalPages == 0 {
		totalPages = 1
	}
	
	if err := ds.Limit(uint(pageSize)).Offset(uint((page - 1) * pageSize)).ScanStructs(&movie); err != nil {
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
					Name: item.Genre_name,
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