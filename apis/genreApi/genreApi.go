package genreapi

import (
	"net/http"
	"strconv"

	"github.com/dangLuan01/restapi_go/apis/utilapi"
	"github.com/dangLuan01/restapi_go/models"
	"github.com/gorilla/mux"
)
func GetAllGenre(respone http.ResponseWriter, request *http.Request) {
	genre := models.GetAllGenre()
	utilapi.ResponseWithJson(respone, http.StatusOK, genre)
}
func GetAllGenreHome(respone http.ResponseWriter, request *http.Request)  {
	genre, err := models.GetAllGenreHome()
	if err != nil {
		utilapi.ResponseWithJson(respone, http.StatusOK, err)	
	}
	utilapi.ResponseWithJson(respone, http.StatusOK, genre)
}
func GetGenreInfo(response http.ResponseWriter, request *http.Request)  {
	query := request.URL.Query()
	pageGet := query.Get("page")
	pageSizeGet := query.Get("page_size")
	page, err := strconv.Atoi(pageGet)
	if err != nil || page < 1 {
		page = 1
	}
	pageSize, err := strconv.Atoi(pageSizeGet)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}
	slug := mux.Vars(request)["slug"]
	genre, err := models.GetItemGenre(slug, page, pageSize)
	if err != nil {
		utilapi.ResponseWithJson(response, http.StatusOK, err)
	}
	utilapi.ResponseWithJson(response, http.StatusOK, genre)
}