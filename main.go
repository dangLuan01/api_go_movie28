package main

import (
	"net/http"
	"github.com/dangLuan01/restapi_go/apis/movieapi"
	"github.com/dangLuan01/restapi_go/config"
	"github.com/gorilla/mux"
)
func main()  {
	config.initDB()
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/movies", movieapi.GetMovie).Methods("GET")
	router.HandleFunc("/api/v1/categories", movieapi.GetCategory).Methods("GET")

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		panic(err)
	}
}