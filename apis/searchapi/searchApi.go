package searchapi

import (
	"log"
	"net/http"

	"github.com/dangLuan01/api_go_movie28/apis/utilapi"
	"github.com/dangLuan01/api_go_movie28/models"
)
func SearchES(response http.ResponseWriter, request *http.Request) {
	query 			:= request.URL.Query()
	search 			:= query.Get("p")
	es, err 		:= models.Search(search)
	if err != nil {
		log.Printf("SearchES error: %v", err)
		utilapi.ResponseWithJson(response, http.StatusInternalServerError, "Failed to connect to Elasticsearch")
		return
	}
	utilapi.ResponseWithJson(response, http.StatusOK, es)
}
func Syn(response http.ResponseWriter, request *http.Request) {
	syn := models.SynES()
	utilapi.ResponseWithJson(response, http.StatusOK, syn)
}