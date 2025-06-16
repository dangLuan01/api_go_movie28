package models

import (
	"encoding/json"
	"fmt"
	"github.com/dangLuan01/api_go_movie28/config"
	"github.com/dangLuan01/api_go_movie28/entities"
	"github.com/doug-martin/goqu/v9"
)
type MovieRaw struct {
	Id   			int    	`json:"id"`
	Name  			string 	`json:"name"`
	Origin_name		string	`json:"origin_name"`
	Slug  			string 	`json:"slug"`
	Type  			string 	`json:"type"`
	Release_date 	int    	`json:"release_date"`
	Rating			float64 `json:"rating"`
	Content 		string 	`json:"content,omitempty"`
	Runtime 		string 	`json:"runtime,omitempty"`
	Age 			string 	`json:"age,omitempty"`
	Trailer 		string 	`json:"trailer,omitempty"`
	Thumb 			string 	`json:"thumb"`
	Poster			string 	`json:"poster"`
	Genre_name 		string
}

func ConvertRating(rating float32) float32 {
	return float32(int(rating * 10)) / 10
}

func GetAllMovieHot() []entities.Movie {
	var listHotMovie []entities.Movie
	// Subquery cho genre_name
	genreSubquery := config.DB.From(goqu.T("genres").As("g")).
		Join(goqu.T("movie_genres").As("mg"), goqu.On(goqu.I("g.id").Eq(goqu.I("mg.genre_id")))).
		Where(goqu.I("mg.movie_id").Eq(goqu.I("m.id"))).
		Select(goqu.I("g.name")).
		Limit(1)
		// Subquery cho poster
	thumbSubquery := config.DB.From(goqu.T("movie_images").As("mi")).
		Where(
			goqu.I("mi.movie_id").Eq(goqu.I("m.id")),
			goqu.I("mi.is_thumbnail").Eq(1),
		).
		Select(goqu.Func("CONCAT", goqu.I("mi.path"), goqu.I("mi.image"))).
		Limit(1)
	err := config.DB.From(goqu.T("movies").As("m")).
	Where(goqu.I("m.hot").Eq(1)).
	Select(
		goqu.I("m.id"),
		goqu.I("m.name"),
		goqu.I("m.origin_name"),
		goqu.I("m.slug"),
		goqu.I("m.type"),
		goqu.I("m.release_date"),
		goqu.I("m.rating"),
		genreSubquery.As("genre_name"),
		thumbSubquery.As("thumb"),
	).
	Limit(10)
	var hot []MovieRaw
	if err := err.ScanStructs(&hot); err != nil {
		fmt.Println(err)
		return nil
	}

	for _, item := range hot {
		listHotMovie = append(listHotMovie, entities.Movie{
			Name: item.Name,
			Origin_name: item.Origin_name,
			Slug: item.Slug,
			Type: item.Type,
			Release_date: item.Release_date,
			Rating: ConvertRating(float32(item.Rating)),
			Image: entities.Image{
				Thumb: item.Thumb,
			},
			Genres: []entities.Genre{
				{
					Name: item.Genre_name,
				},
			},
		})
	}
	return listHotMovie
}
func GetAllMovie(page, pageSize int) (entities.PaginatedMovies, error) {
	var listMovie []entities.Movie
	var movie []MovieRaw
	var totalPages int64
	// Subquery cho genre_name
	genreSubquery := config.DB.From(goqu.T("genres").As("g")).
		Join(goqu.T("movie_genres").As("mg"), goqu.On(goqu.I("g.id").Eq(goqu.I("mg.genre_id")))).
		Where(goqu.I("mg.movie_id").Eq(goqu.I("m.id"))).
		Select(goqu.I("g.name")).
		Limit(1)

	// Subquery cho poster
	posterSubquery := config.DB.From(goqu.T("movie_images").As("mi")).
		Where(
			goqu.I("mi.movie_id").Eq(goqu.I("m.id")),
			goqu.I("mi.is_thumbnail").Eq(0),
		).
		Select(goqu.Func("CONCAT", goqu.I("mi.path"), goqu.I("mi.image"))).
		Limit(1)
	ds := config.DB.From(goqu.T("movies").As("m")).
		Where(goqu.I("m.hot").Eq(0)).
		Select(
			goqu.I("m.name"),
			goqu.I("m.origin_name"),
			goqu.I("m.slug"),
			goqu.I("m.type"),
			goqu.I("m.release_date"),
			goqu.I("m.rating"),
			genreSubquery.As("genre_name"),
			posterSubquery.As("poster"),
		).
		Order(goqu.I("m.updated_at").Desc()).Limit(uint(pageSize)).Offset(uint((page - 1) * pageSize))
	if err := ds.ScanStructs(&movie); err != nil {
		return entities.PaginatedMovies{}, nil
	}
	
	for _, item := range movie {
		listMovie = append(listMovie, entities.Movie{
			Name: item.Name,
			Origin_name: item.Origin_name,
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
// func GetAllMovie(page, pageSize int) (entities.PaginatedMovies, error) {
// 	var movies []entities.Movie

// 	// Subquery để lấy thumbnail
// 	thumbnailSubq := config.DB.
// 		From("movie_images").
// 		Select(
// 			goqu.I("movie_id"),
// 			goqu.Func("CONCAT", goqu.I("path"), goqu.I("image")).As("poster"),
// 		).
// 		Where(goqu.Ex{"is_thumbnail": 0}).
// 		GroupBy("movie_id")

// 	// Main query
// 	ds := config.DB.
// 		From("movies").
// 		LeftJoin(
// 			goqu.T("movie_genres").As("mg"),
// 			goqu.On(goqu.I("movies.id").Eq(goqu.I("mg.movie_id"))),
// 		).
// 		LeftJoin(
// 			goqu.T("genres").As("g"),
// 			goqu.On(goqu.I("mg.genre_id").Eq(goqu.I("g.id"))),
// 		).
// 		LeftJoin(
// 			goqu.L("(?)", thumbnailSubq).As("mi"),
// 			goqu.On(goqu.I("movies.id").Eq(goqu.I("mi.movie_id"))),
// 		).
// 		Where(goqu.Ex{"movies.hot": 0}).
// 		Select(
// 			goqu.I("movies.name"),
// 			goqu.I("movies.slug"),
// 			goqu.I("movies.type"),
// 			goqu.I("movies.release_date"),
// 			goqu.I("movies.rating"),
// 			goqu.I("mi.poster").As("image_raw"),
// 			// Gom genres thành JSON array
// 			goqu.L(`JSON_ARRAYAGG(
// 				CASE 
// 					WHEN g.id IS NULL THEN NULL
// 					ELSE JSON_OBJECT('name', g.name)
// 				END
// 			)`).As("genres_raw"),
// 		).
// 		GroupBy(goqu.I("movies.id")).
// 		Order(goqu.I("movies.updated_at").Desc()).
// 		Limit(uint(pageSize)).
// 		Offset(uint((page - 1) * pageSize))
	
// 	//Scan
// 	if err := ds.ScanStructs(&movies); err != nil {
// 		fmt.Println(err)
// 		return entities.PaginatedMovies{}, err
// 	}
	
// 	// Sau khi ScanStructs(&movies)
// 	for i := range movies {
// 		if movies[i].GenresRaw != "" {
// 			_ = json.Unmarshal([]byte(movies[i].GenresRaw), &movies[i].Genres)
// 		}
// 		if movies[i].ImageRaw != "" {
//         	movies[i].Image.Poster = movies[i].ImageRaw
//     	}
// 	}

// 	return entities.PaginatedMovies{
// 		Data:     movies,
// 		Page:     page,
// 		PageSize: pageSize,
// 	}, nil
// }

func GetDetailMovie(slug string) (entities.Movie, error) {
	var row MovieRaw
	var genres []entities.Genre
	ds := config.DB.Select(
		"movies.id",
		"movies.name",
		"movies.origin_name",
		"movies.slug",
		"movies.type",
		"movies.release_date",
		"movies.rating",
		goqu.Func("IFNULL", goqu.I("movies.content"), "").As("content"),
		goqu.Func("IFNULL", goqu.I("movies.runtime"), "").As("runtime"),
		goqu.Func("IFNULL", goqu.I("movies.age"), "").As("age"),
		goqu.Func("IFNULL", goqu.I("movies.trailer"), "").As("trailer"),
		//"movies.trailer",
		goqu.Func("CONCAT", goqu.I("mi.path"), goqu.I("mi.image")).As("thumb"),
	
	).From("movies").
	LeftJoin(
		goqu.T("movie_images").As("mi"),
		goqu.On(
			goqu.I("movies.id").Eq(goqu.I("mi.movie_id")),
		),
	).Where(goqu.Ex{"slug": slug,"mi.is_thumbnail": 1})
	_, err := ds.ScanStruct(&row)

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
		fmt.Println(err)
		return entities.Movie{}, err
	}

	movie := entities.Movie{
		Name: row.Name,
		Origin_name: row.Origin_name,
		Slug: row.Slug,
		Type: row.Type,
		Release_date: row.Release_date,
		Rating: ConvertRating(float32(row.Rating)),
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
	return movie, nil
}
// func GetServerWithEpisodes(movieId int) ([]entities.Server, error) {
	
// 	var servers []entities.Server
// 	var episodes []entities.Episode
	
// 	err := config.DB.From("movie_servers").
// 		Join(goqu.T("episodes"), goqu.On(goqu.I("movie_servers.id").Eq(goqu.I("episodes.server_id")))).
// 		Where(goqu.Ex{"episodes.movie_id": uint(movieId)}).
// 		Select("movie_servers.id", "movie_servers.name").
// 		GroupBy("movie_servers.id").ScanStructs(&servers)
	
// 	if err != nil {
// 		return nil, err
// 	}

// 	err = config.DB.From("episodes").
// 		Where(goqu.Ex{"movie_id": uint(movieId)}).
// 		Select("server_id", "episode", "hls").
// 		ScanStructs(&episodes)
	
// 	if err != nil {
// 		return nil, err
// 	}
// 	serverMap := make(map[int]*entities.Server)
// 	for i := range servers {
// 		serverMap[servers[i].Id] = &servers[i]
// 	}
	
// 	for _, ep := range episodes {
// 		if server, ok := serverMap[ep.Server_id]; ok {
// 			server.Episodes = append(server.Episodes, ep)
// 		}
// 	}
// 	return servers, nil
// }
func GetServerWithEpisodes(movieId int) ([]entities.Server, error) {
    var servers []struct {
        Id       int             `db:"id"`
        Name     string          `db:"name"`
        Episodes json.RawMessage `db:"episodes"`
    }
    err := config.DB.From("movie_servers").
        LeftJoin(goqu.T("episodes"), 
            goqu.On(goqu.I("movie_servers.id").Eq(goqu.I("episodes.server_id")))).
        Where(goqu.Ex{
            "episodes.movie_id": uint(movieId),
        }).
        Select(
            "movie_servers.id",
            "movie_servers.name",
            goqu.L(`JSON_ARRAYAGG(
                CASE 
                    WHEN episodes.server_id IS NULL THEN NULL
                    ELSE JSON_OBJECT(
                        'server_id', episodes.server_id,
                        'episode', episodes.episode,
                        'hls', episodes.hls
                    )
                END
            )`).As("episodes"),
        ).
        GroupBy("movie_servers.id").
        ScanStructs(&servers)

    if err != nil {
        return nil, err
    }
    result := make([]entities.Server, 0, len(servers))
    for _, s := range servers {
        server := entities.Server{
            Id:   s.Id,
            Name: s.Name,
        }

        // Parse JSON episodes nếu có
        if len(s.Episodes) > 0 && string(s.Episodes) != "null" {
            var episodes []entities.Episode
            if err := json.Unmarshal(s.Episodes, &episodes); err != nil {
                return nil, fmt.Errorf("failed to unmarshal episodes: %v", err)
            }
            server.Episodes = episodes
        }

        result = append(result, server)
    }

    return result, nil
}
