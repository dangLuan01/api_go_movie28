package movieapi

import (
	"log"
	"net/http"
	"strconv"

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
	movie := models.GetAllMovieHot()
	utilapi.ResponseWithJson(respone, http.StatusOK, movie)
}

func GetAllMovie(respone http.ResponseWriter, request *http.Request)  {
	query 		:= request.URL.Query()
	pageGet 	:= query.Get("page")
	pageSizeGet := query.Get("page_size")
	page, err 	:= strconv.Atoi(pageGet)
	if err != nil || page < 1 {
		page = 1
	}
	pageSize, err := strconv.Atoi(pageSizeGet)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}
	movie, err := models.GetAllMovie(page, pageSize)
	if err != nil {
		utilapi.ResponseWithJson(respone, http.StatusOK, err)
	}
	utilapi.ResponseWithJson(respone, http.StatusOK, movie)
}
func GetMovieBySlug(respone http.ResponseWriter, request *http.Request)  {
	movieCache 	:= cacheloader.GetMovieCache()
	slug 		:= mux.Vars(request)["slug"]
	var movie *entities.Movie = movieCache.Get(slug)
	if movie == nil {
		movie, err := models.GetDetailMovie(slug)
		if err != nil {
			log.Printf("Err:%v", err)
			utilapi.ResponseWithJson(respone, http.StatusOK, err)
		}
		movieCache.Set(slug, &movie)
		utilapi.ResponseWithJson(respone, http.StatusOK, movie)
	}else {
		utilapi.ResponseWithJson(respone, http.StatusOK, movie)
	}
	
}
