package genreapi
import (
	"net/http"
	"github.com/dangLuan01/restapi_go/models"
	"github.com/dangLuan01/restapi_go/apis/utilapi"
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