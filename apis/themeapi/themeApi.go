package themeapi

import (
	"net/http"
	"github.com/dangLuan01/restapi_go/models"
	"github.com/dangLuan01/restapi_go/apis/utilapi"
)
func GetThemes(respone http.ResponseWriter, request *http.Request)  {
	theme, err := models.GetAllThemesWithMovies(1,24)
	if err != nil {
		utilapi.ResponseWithJson(respone, http.StatusOK, err)	
	}
	utilapi.ResponseWithJson(respone, http.StatusOK, theme)
}