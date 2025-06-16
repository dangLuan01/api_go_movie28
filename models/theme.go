package models

import (
	//"encoding/json"
	"fmt"	
	"math"
    "log"
	"sync"
	//"time"
	"github.com/dangLuan01/api_go_movie28/config"
	"github.com/dangLuan01/api_go_movie28/entities"
	"github.com/doug-martin/goqu/v9"
)

func buildMovieQueryFromTheme(theme entities.ThemeInfo) *goqu.SelectDataset {
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

    // Truy vấn chính
    query := config.DB.From(goqu.T("movies").As("m")).
        Select(
            goqu.I("m.id"),
            goqu.I("m.name"),
            goqu.I("m.origin_name"),
            goqu.I("m.slug"),
            goqu.I("m.type"),
            goqu.I("m.release_date"),
            goqu.I("m.rating"),
            posterSubquery.As("poster"),
            genreSubquery.As("genre_name"),
        ).
        Where(goqu.I("m.hot").Eq(0))
    // Điều kiện genre_id
    if theme.Genre_id != nil {
        query = query.Join(
            goqu.T("movie_genres").As("mg"),
            goqu.On(goqu.I("mg.movie_id").Eq(goqu.I("m.id"))),
        ).Where(goqu.I("mg.genre_id").Eq(*theme.Genre_id))
    }

    // Điều kiện country_id
    if theme.Country_id != nil {
        query = query.Join(
            goqu.T("movie_countries").As("mc"),
            goqu.On(goqu.I("mc.movie_id").Eq(goqu.I("m.id"))),
        ).Where(goqu.I("mc.country_id").Eq(*theme.Country_id))
    }

    // Điều kiện type
    if theme.Type != nil {
        query = query.Where(goqu.I("m.type").Eq(*theme.Type))
    }

    // Điều kiện year
    if theme.Year != nil {
       query = query.Where(goqu.I("m.release_year").Eq(*theme.Year))
    }
    query = query.Order(goqu.I("m.updated_at").Desc())
    return query
}

func GetMoviesByTheme(theme entities.ThemeInfo, page, limit int) (entities.PaginatedMovies, error) {
    // if page < 1 { page = 1 }
    // if limit < 1 { limit = 10 }
    offset      := (page - 1) * limit
    baseQuery   := buildMovieQueryFromTheme(theme)
    totalCount, err := baseQuery.Count()
    if err != nil {
        return entities.PaginatedMovies{}, fmt.Errorf("failed to count movies: %v", err)
    }
    var movies []entities.MovieRaw
    err = baseQuery.Offset(uint(offset)).
        Limit(uint(limit)).
        ScanStructs(&movies)
    
    if err != nil {
        return entities.PaginatedMovies{}, fmt.Errorf("failed to scan movies: %v", err)
    }
    
    resultMovies := make([]entities.Movie, 0, len(movies))
    for _, m := range movies {
        resultMovies = append(resultMovies, convertMovieRawToMovie(m))
    }
    return entities.PaginatedMovies{
        Movie:          resultMovies,
        Page:           page,
        PageSize:       limit,
        TotalPages:     int(math.Ceil(float64(totalCount)/float64(limit))),
    }, nil
}

func convertMovieRawToMovie(raw entities.MovieRaw) entities.Movie {
    return entities.Movie{
        Name:           raw.Name,
        Origin_name:    raw.Origin_name,
        Slug:           raw.Slug,
        Type:           raw.Type,
        Release_date:   raw.Release_date,
        Rating:         ConvertRating(float32(raw.Rating)),
        Image:          entities.Image{Poster: raw.Poster},
        Genres:         []entities.Genre{{Name: raw.Genre_name}},
    }
}
// func GetAllThemesWithMovies(id, pageTheme, pageMovie, limit int) (entities.PagiateTheme, error) {
//     if pageTheme < 1 { pageTheme = 1 }
//     if limit < 1 { limit = 4 }
//     offset := (pageTheme - 1) * limit
//     var themes []entities.ThemeInfo
    
//     ds := config.DB.From(goqu.T("themes").As("t")).Where(goqu.Ex{
//         "t.status": 1,
//     })
//     if id != 0 {
//         ds = ds.Where(goqu.Ex{
//             "t.id": uint(id),
//         })
//     }
//     totalCount, _ := ds.Count()
//     if err := ds.Order(goqu.I("t.priority").Asc()).Offset(uint(offset)).Limit(uint(limit)).ScanStructs(&themes); err != nil {
//         return entities.PagiateTheme{}, fmt.Errorf("failed: %v", err)
//     }
    
//     result := make([]entities.ThemeWithMovies, 0, len(themes))
//     for _, theme := range themes {
//         movies, err := GetMoviesByTheme(theme, pageMovie, theme.Limit)
//         if err != nil {
//             fmt.Printf("Error getting movies for theme %s: %v", theme.Name, err)
//             continue
//         }
        
//         result = append(result, entities.ThemeWithMovies{
//             ThemeInfo:          theme,
//             PaginatedMovies:    movies,
//         })
//     }
//     results := entities.PagiateTheme{
//         ThemeWithMovies: result,
//         Page: pageTheme,
//         PageSize: limit,
//         TotalPages: int(math.Ceil(float64(totalCount)/float64(limit))),
//     }
    
//     return results, nil
// }

// ## Run concurrency ##
func GetAllThemesWithMovies(id, pageTheme, pageMovie, limit int) (entities.PagiateTheme, error) {
    if pageTheme < 1 {
        pageTheme = 1
    }
    if limit < 1 {
        limit = 2
    }
    offset := (pageTheme - 1) * limit

    // Lấy danh sách theme
    var themes []entities.ThemeInfo
    ds := config.DB.From(goqu.T("themes").As("t")).Where(goqu.Ex{
        "t.status": 1,
    })
    if id != 0 {
        ds = ds.Where(goqu.Ex{
            "t.id": uint(id),
        })
    }
    totalCount, err := ds.Count()
    if err != nil {
        return entities.PagiateTheme{}, fmt.Errorf("failed to count themes: %v", err)
    }
    if err := ds.Order(goqu.I("t.priority").Asc()).Offset(uint(offset)).Limit(uint(limit)).ScanStructs(&themes); err != nil {
        return entities.PagiateTheme{}, fmt.Errorf("failed to scan themes: %v", err)
    }

    // Kênh để thu thập kết quả và lỗi
    type resultStruct struct {
        index           int
        themeWithMovies entities.ThemeWithMovies
        err             error
    }
    resultsChan := make(chan resultStruct, len(themes))
    var wg sync.WaitGroup

    // Chạy goroutines cho mỗi theme
    // start := time.Now()
    for i, theme := range themes {
        wg.Add(1)
        go func(idx int, t entities.ThemeInfo) {
            defer wg.Done()
            movies, err := GetMoviesByTheme(t, pageMovie, t.Limit)
            if err != nil {
                resultsChan <- resultStruct{
                    index: idx, 
                    err: fmt.Errorf("failed to get movies for theme %s: %v", t.Name, err),
                }
                return
            }
            resultsChan <- resultStruct{
                index: idx,
                themeWithMovies: entities.ThemeWithMovies{
                    ThemeInfo:       t,
                    PaginatedMovies: movies,
                },
            }
        }(i, theme)
    }

    // Đóng kênh sau khi tất cả goroutines hoàn thành
    go func() {
        wg.Wait()
        close(resultsChan)
    }()

    // Thu thập kết quả theo thứ tự
    resultSlice := make([]entities.ThemeWithMovies, len(themes))
    var errors []error
    for res := range resultsChan {
        if res.err != nil {
            errors = append(errors, res.err)
            continue
        }
        resultSlice[res.index] = res.themeWithMovies
    }

    // // Ghi log thời gian
    // log.Printf("GetAllThemesWithMovies took %v for %d themes", time.Since(start), len(themes))

    // Kiểm tra lỗi
    if len(errors) > 0 {
        log.Printf("Encountered %d errors while fetching movies: %v", len(errors), errors)
        // Nếu tất cả theme đều lỗi, trả về lỗi
        if len(errors) == len(themes) {
            return entities.PagiateTheme{}, fmt.Errorf("failed to fetch movies for all themes: %v", errors)
        }
    }

    // Loại bỏ các theme rỗng (nếu có)
    result := make([]entities.ThemeWithMovies, 0, len(themes))
    for _, r := range resultSlice {
        if r.ThemeInfo.Id != 0 { // Chỉ thêm các theme có dữ liệu hợp lệ
            result = append(result, r)
        }
    }

    return entities.PagiateTheme{
        ThemeWithMovies: result,
        Page:            pageTheme,
        PageSize:        limit,
        TotalPages:      int(math.Ceil(float64(totalCount)/float64(limit))),
    }, nil
}
