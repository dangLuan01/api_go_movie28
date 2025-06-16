package main

import (
	"net/http"
	"github.com/dangLuan01/api_go_movie28/config"
	"github.com/dangLuan01/api_go_movie28/router"
)
func main()  {
	config.InitDB()
	router := router.SetupRouter()

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		panic(err)
	}
}