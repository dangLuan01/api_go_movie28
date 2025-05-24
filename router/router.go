package router

import (
	"github.com/dangLuan01/restapi_go/apis/movieapi"
	"github.com/gorilla/mux"
)
func SetupRouter() *mux.Router{
	router := mux.NewRouter()
	
	router.HandleFunc("/api/v1/movie-hot", movieapi.GetMovieHot).Methods("GET")
	router.HandleFunc("/api/v1/category", movieapi.GetCategory).Methods("GET")
	router.HandleFunc("/api/v1/movies", movieapi.GetAllMovie).Methods("GET")
	return router
}