package movieapi

import (
	"encoding/json"
	"net/http"
	"strconv"
	//"github.com/dangLuan01/restapi_go/entities"
	"github.com/dangLuan01/restapi_go/models"
)
func responeWithJson(respone http.ResponseWriter, statusCode int, data interface{}) {
	result, err := json.Marshal(data)
	if err != nil {
		http.Error(respone, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	respone.Header().Set("Content-Type", "application/json")
	respone.WriteHeader(statusCode)
	respone.Write(result)
}
func GetMovieHot(respone http.ResponseWriter, request *http.Request) {
	movie := models.GetAllMovieHot()
	responeWithJson(respone, http.StatusOK, movie)
}

func GetCategory(respone http.ResponseWriter, request *http.Request)  {
	category := models.GetAllCategory()
	responeWithJson(respone, http.StatusOK, category)
}
func GetAllMovie(respone http.ResponseWriter, request *http.Request)  {
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
	movie := models.GetAllMovie(page, pageSize)
	responeWithJson(respone, http.StatusOK, movie)
}