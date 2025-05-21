package main

import (
	"net/http"
	"github.com/dangLuan01/restapi_go/config"
	"github.com/dangLuan01/restapi_go/router"
)
func main()  {
	config.InitDB()
	router := router.SetupRouter()

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		panic(err)
	}
}