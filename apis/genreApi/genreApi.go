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