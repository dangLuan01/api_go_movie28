package collectionapi

import (
	"net/http"
	"strconv"
	"github.com/dangLuan01/restapi_go/models"
	"github.com/dangLuan01/restapi_go/apis/utilapi"
	"github.com/gorilla/mux"
)
func GetColletion(respone http.ResponseWriter, request *http.Request)  {
	collection, err := models.GetAllCollection()
	if err != nil {
		utilapi.ResponseWithJson(respone, http.StatusOK, err)	
	}
	utilapi.ResponseWithJson(respone, http.StatusOK, collection)
}
func GetColletionBySlug(respone http.ResponseWriter, request *http.Request) {
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
	movie, err := models.GetMovieCollection(slug, page, pageSize)
	if err != nil {
		utilapi.ResponseWithJson(respone, http.StatusOK, err)
	}
	
	utilapi.ResponseWithJson(respone, http.StatusOK, movie)
}