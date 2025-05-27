package models

import (
	"github.com/dangLuan01/restapi_go/config"
	"github.com/dangLuan01/restapi_go/entities"
	"github.com/doug-martin/goqu/v9"
)
type MovieRaw struct {
	Id   			int    `json:"id"`
	Name  			string `json:"name"`
	Slug  			string `json:"slug"`
	Type  			string `json:"type"`
	Release_date 	int    `json:"release_date"`
	Rating			float64 `json:"rating"`
	Content 		string `json:"content,omitempty"`
	Runtime 		string `json:"runtime,omitempty"`
	Age 			string `json:"age,omitempty"`
	Trailer 		string `json:"trailer,omitempty"`
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
func GetDetailMovie(slug string)  entities.Movie {
	var row MovieRaw
	var genres []entities.Genre
	ds := config.DB.Select(
		"movies.id",
		"movies.name",
		"movies.slug",
		"movies.type",
		"movies.release_date",
		"movies.rating",
		goqu.Func("IFNULL", goqu.I("movies.content"), "").As("content"),
		"movies.runtime",
		goqu.Func("IFNULL", goqu.I("movies.age"), "").As("age"),
		"movies.trailer",
		goqu.Func("CONCAT", goqu.I("mi.path"), goqu.I("mi.image")).As("thumb"),
	
	).From("movies").
	LeftJoin(
		goqu.T("movie_images").As("mi"),
		goqu.On(
			goqu.I("movies.id").Eq(goqu.I("mi.movie_id")),
		),
	).Where(goqu.Ex{"slug": slug,"mi.is_thumbnail": 1})
	found, err := ds.ScanStruct(&row)

	err2 := config.DB.
	From("genres").
	Join(
		goqu.T("movie_genres").As("mg"),
		goqu.On(goqu.I("genres.id").Eq(goqu.I("mg.genre_id"))),
	).
	Where(goqu.Ex{"mg.movie_id": row.Id}).
	Select("genres.name", "genres.slug").
	ScanStructs(&genres)

	if err != nil {
		return entities.Movie{}
	}
	if !found {
		return entities.Movie{}
	}

	movie := entities.Movie{
		Name: row.Name,
		Slug: row.Slug,
		Type: row.Type,
		Release_date: row.Release_date,
		Rating: row.Rating,
		Content: row.Content,
		Runtime: row.Runtime,
		Age: row.Age,
		Trailer: row.Trailer,
		Image: entities.Image{
			Thumb: row.Thumb,
		},
	}
	if err2 == nil {
		movie.Genres = genres
	}
	server, err := GetServerWithEpisodes(row.Id)
	if err != nil {
		server = []entities.Server{}
	}
	movie.Servers = server
	return movie
}
func GetServerWithEpisodes(movieId int) ([]entities.Server, error) {
	
	var servers []entities.Server
	var episodes []entities.Episode
	
	err := config.DB.From("movie_servers").
		Join(goqu.T("episodes"), goqu.On(goqu.I("movie_servers.id").Eq(goqu.I("episodes.server_id")))).
		Where(goqu.Ex{"episodes.movie_id": uint(movieId)}).
		Select("movie_servers.id", "movie_servers.name").
		GroupBy("movie_servers.id").ScanStructs(&servers)
	
	if err != nil {
		return nil, err
	}

	err = config.DB.From("episodes").
		Where(goqu.Ex{"movie_id": uint(movieId)}).
		Select("server_id", "episode", "hls").
		ScanStructs(&episodes)
	
	if err != nil {
		return nil, err
	}
	serverMap := make(map[int]*entities.Server)
	for i := range servers {
		serverMap[servers[i].Id] = &servers[i]
	}
	
	for _, ep := range episodes {
		if server, ok := serverMap[ep.Server_id]; ok {
			server.Episodes = append(server.Episodes, ep)
		}
	}
	return servers, nil
}