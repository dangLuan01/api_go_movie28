package models

import (
	//"encoding/json"
	"fmt"
	"math"
	"github.com/dangLuan01/restapi_go/config"
	"github.com/dangLuan01/restapi_go/entities"
	"github.com/doug-martin/goqu/v9"
)

func buildMovieQueryFromTheme(theme entities.ThemeInfo) *goqu.SelectDataset {
    query := config.DB.From("movies").Select(
            goqu.I("movies.id"),
            goqu.I("movies.name"),
            goqu.I("movies.slug"),
            goqu.I("movies.type"),
            goqu.I("movies.release_date"),
            goqu.I("movies.rating"),
            goqu.Func("CONCAT", goqu.I("mi.path"), goqu.I("mi.image")).As("poster"),
            goqu.I("g.name").As("genre_name"),
        )
    // Thêm join với movie_genre nếu có genre_id
    if theme.Genre_id != nil {
        query = query.LeftJoin(
            goqu.T("movie_genres").As("mg"),
            goqu.On(goqu.I("mg.movie_id").Eq(goqu.I("movies.id"))),
        ).
		LeftJoin(
        	goqu.T("genres").As("g"),
        	goqu.On(goqu.I("mg.genre_id").Eq(goqu.I("g.id"))),
    	).
		 LeftJoin(
            goqu.T("movie_images").As("mi"),
            goqu.On(goqu.I("mi.movie_id").Eq(goqu.I("movies.id"))),
        ).
		Where(goqu.Ex{
			"mg.genre_id": *theme.Genre_id,
			"mi.is_thumbnail": 0,
		})
    }
    // Thêm điều kiện country nếu có
    if theme.Country_id != nil {
        //query = query.Where(goqu.Ex{"movies.country_id": *theme.Country_id})
        query = query.LeftJoin(
            goqu.T("movie_countries").As("mc"),
            goqu.On(goqu.I("mc.movie_id").Eq(goqu.I("movies.id"))),
        ).Where(goqu.Ex{
            "mc.country_id": *theme.Country_id,
        })
    }
    // Thêm điều kiện type nếu có
    if theme.Type != nil {
        query = query.Where(goqu.Ex{"movies.type": *theme.Type})
    }
    // Thêm điều kiện year nếu có
    if theme.Year != nil {
        query = query.Where(goqu.I("movies.release_date").Like(fmt.Sprintf("%d%%", *theme.Year)))
    }
    return query
}
func GetMoviesByTheme(theme entities.ThemeInfo, page, limit int) (entities.PaginatedMovies, error) {
    if page < 1 { page = 1 }
    if limit < 1 { limit = 10 }
    offset := (page - 1) * limit

    baseQuery := buildMovieQueryFromTheme(theme)
    
    totalCount, err := baseQuery.Count()
    if err != nil {
        return entities.PaginatedMovies{}, err
    }
    var movies []entities.MovieRaw
    err = baseQuery.
        Order(goqu.I("movies.updated_at").Desc()).
        Offset(uint(offset)).
        Limit(uint(limit)).
        ScanStructs(&movies)
    
    if err != nil {
        return entities.PaginatedMovies{}, err
    }
    
    resultMovies := make([]entities.Movie, 0, len(movies))
    for _, m := range movies {
        resultMovies = append(resultMovies, convertMovieRawToMovie(m))
    }
    return entities.PaginatedMovies{
        Movie:           resultMovies,
        Page:           page,
        PageSize:       limit,
        TotalPages:     int(math.Ceil(float64(totalCount)/float64(limit))),
    }, nil
}

func convertMovieRawToMovie(raw entities.MovieRaw) entities.Movie {
    return entities.Movie{
        Name:           raw.Name,
        Slug:           raw.Slug,
        Type:           raw.Type,
        Release_date:   raw.Release_date,
        Rating:         ConvertRating(float32(raw.Rating)),
        Image:          entities.Image{Poster: raw.Poster},
        Genres:         []entities.Genre{{Name: raw.Genre_name}},
    }
}
func GetAllThemesWithMovies(id, pageTheme, pageMovie, limit int) (entities.PagiateTheme, error) {
    if pageTheme < 1 { pageTheme = 1 }
    if limit < 1 { limit = 4 }
    offset := (pageTheme - 1) * limit
    var themes []entities.ThemeInfo
    
    ds := config.DB.From(goqu.T("themes").As("t")).Where(goqu.Ex{
        "t.status": 1,
    })
    if id != 0 {
        ds = ds.Where(goqu.Ex{
            "t.id": uint(id),
        })
    }
    if err := ds.Order(goqu.I("t.priority").Asc()).Offset(uint(offset)).Limit(uint(limit)).ScanStructs(&themes); err != nil {
        return entities.PagiateTheme{}, err
    }
    totalCount, _ := ds.Count()
    result := make([]entities.ThemeWithMovies, 0, len(themes))
    for _, theme := range themes {
        movies, err := GetMoviesByTheme(theme, pageMovie, theme.Limit)
        if err != nil {
            fmt.Printf("Error getting movies for theme %s: %v", theme.Name, err)
            continue
        }
        
        result = append(result, entities.ThemeWithMovies{
            ThemeInfo:          theme,
            PaginatedMovies:    movies,
        })
    }
    results := entities.PagiateTheme{
        ThemeWithMovies: result,
        Page: pageTheme,
        PageSize: limit,
        TotalPages: int(math.Ceil(float64(totalCount)/float64(limit))),
    }
    
    return results, nil
}