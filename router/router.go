package router

import (
	"github.com/dangLuan01/restapi_go/apis/movieapi"
	"github.com/gorilla/mux"
)
func SetupRouter() *mux.Router{
	router := mux.NewRouter()
	
	router.HandleFunc("/api/v1/movies", movieapi.GetMovie).Methods("GET")
	router.HandleFunc("/api/v1/categories", movieapi.GetCategory).Methods("GET")

	return router
}