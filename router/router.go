package router

import (
	"github.com/dangLuan01/restapi_go/apis/collectionapi"
	"github.com/dangLuan01/restapi_go/apis/genreapi"
	"github.com/dangLuan01/restapi_go/apis/movieapi"
	"github.com/dangLuan01/restapi_go/apis/themeapi"
	"github.com/gorilla/mux"
)
func SetupRouter() *mux.Router{
	router := mux.NewRouter()
	
	router.HandleFunc("/api/v1/movie-hot", movieapi.GetMovieHot).Methods("GET")
	router.HandleFunc("/api/v1/category", movieapi.GetCategory).Methods("GET")
	router.HandleFunc("/api/v1/movies", movieapi.GetAllMovie).Methods("GET")
	router.HandleFunc("/api/v1/genre", genreapi.GetAllGenre).Methods("GET")
	router.HandleFunc("/api/v1/movie/{slug}", movieapi.GetMovieBySlug).Methods("GET")
	router.HandleFunc("/api/v1/genre-home", genreapi.GetAllGenreHome).Methods("GET")
	router.HandleFunc("/api/v1/genre/{slug}", genreapi.GetGenreInfo).Methods("GET")
	router.HandleFunc("/api/v1/collection", collectionapi.GetColletion).Methods("GET")
	router.HandleFunc("/api/v1/collection/{slug}", collectionapi.GetColletionBySlug).Methods("GET")
	router.HandleFunc("/api/v1/themes", themeapi.GetThemes).Methods("GET")

	return router
}