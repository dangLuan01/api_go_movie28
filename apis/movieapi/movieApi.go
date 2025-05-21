package movieapi

import (
	"github.com/dangLuan01/restapi_go/models"
	"encoding/json"
	"net/http"
)
func GetMovie(respone http.ResponseWriter, request *http.Request) {
	movie := models.GetAllMovie()
	responeWithJson(respone, http.StatusOK, movie)
	
}
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
func GetCategory(respone http.ResponseWriter, request *http.Request)  {
	category := models.GetAllCategory()
	responeWithJson(respone, http.StatusOK, category)
}
