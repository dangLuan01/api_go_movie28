package movieapi

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/dangLuan01/api_go_movie28/apis/utilapi"
	"github.com/dangLuan01/api_go_movie28/entities"
	"github.com/dangLuan01/api_go_movie28/internal/cacheloader"
	"github.com/dangLuan01/api_go_movie28/models"
	"github.com/gorilla/mux"
)

func GetCategory(respone http.ResponseWriter, request *http.Request)  {
	category := models.GetAllCategory()
	utilapi.ResponseWithJson(respone, http.StatusOK, category)
}

func GetMovieHot(respone http.ResponseWriter, request *http.Request) {
	key 		:= "movie-hot"
	found 		:= false
	time 		:=  time.Duration(rand.Intn(350-250) + 250)
	movieCache 	:= cacheloader.GetCache(0, time)
	var data []entities.Movie
	if movieCache != nil && movieCache.Get(key, &data) {
		log.Println("Read from cache")
		utilapi.ResponseWithJson(respone, http.StatusOK, data)
		found = true
		return
	}
	if !found {
		movie := models.GetAllMovieHot()	
		data = movie
		if movieCache != nil {
			movieCache.Set(key, data)
		}
		log.Println("Read from db")
		utilapi.ResponseWithJson(respone, http.StatusOK, data)
	}
}

func GetAllMovie(respone http.ResponseWriter, request *http.Request)  {
	query 		:= request.URL.Query()
	pageGet 	:= query.Get("page")
	pageSizeGet := query.Get("page_size")
	page, err 	:= strconv.Atoi(pageGet)
	found 		:= false
	if err != nil || page < 1 {
		page = 1
	}
	pageSize, err := strconv.Atoi(pageSizeGet)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}
	var data entities.PaginatedMovies
	key := fmt.Sprintf("movies:page=%s:size=%s", pageGet, pageSizeGet)
	time 		:=  time.Duration(rand.Intn(300-250) + 250)
	movieCache := cacheloader.GetCache(0, time)
	if movieCache != nil && movieCache.Get(key, &data) {
		log.Println("Read from redis")
		utilapi.ResponseWithJson(respone, http.StatusOK, data)
		return
	}
	if !found {
		movie, err := models.GetAllMovie(page, pageSize)	
		if err != nil {
			utilapi.ResponseWithJson(respone, http.StatusOK, map[string]string{
				"error": "Không lấy được dữ liệu",
			})
		}
		data = movie
		if movieCache != nil {
			movieCache.Set(key, data)
		}
	}
	utilapi.ResponseWithJson(respone, http.StatusOK, data)
}
func GetMovieBySlug(response http.ResponseWriter, request *http.Request) {
	slug 		:= mux.Vars(request)["slug"]
	time 		:=  time.Duration(rand.Intn(900-700) + 700)
	movieCache 	:= cacheloader.GetCache(0, time)
	found 		:= false
	var movie entities.MovieDetail
	if movieCache != nil {
		log.Println("Read from redis")
		found = movieCache.Get(slug, &movie)
	}
	if !found {
		dbMovie, err := models.GetDetailMovie(slug)
		if err != nil {
			log.Printf("❌ DB error: %v", err)
			utilapi.ResponseWithJson(response, http.StatusInternalServerError, map[string]string{
				"error": "Không tìm thấy phim",
			})
			return
		}
		movie = dbMovie
		if movieCache != nil {
			movieCache.Set(slug, movie)
		}
	}

	utilapi.ResponseWithJson(response, http.StatusOK, movie)
}

