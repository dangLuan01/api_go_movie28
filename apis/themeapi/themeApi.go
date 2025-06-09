package themeapi

import (
	"net/http"
	"strconv"

	"github.com/dangLuan01/restapi_go/apis/utilapi"
	"github.com/dangLuan01/restapi_go/models"
)
func GetThemes(respone http.ResponseWriter, request *http.Request)  {
	theme, err := models.GetAllThemesWithMovies(1,24)
	if err != nil {
		utilapi.ResponseWithJson(respone, http.StatusOK, err)	
	}
	utilapi.ResponseWithJson(respone, http.StatusOK, theme)
}
func GetThemeId(respone http.ResponseWriter, request *http.Request)  {
	query := request.URL.Query()
	pageGet := query.Get("page")
	pageSizeGet := query.Get("page_size")
	idGet := query.Get("id")
	page, err := strconv.Atoi(pageGet)
	if err != nil || page < 1 {
		page = 1
	}
	pageSize, err := strconv.Atoi(pageSizeGet)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}
	id, _ := strconv.Atoi(idGet)

	theme, err := models.GetThemeById(id, page, pageSize)
	if err != nil {
		utilapi.ResponseWithJson(respone, http.StatusOK, err)	
	}
	utilapi.ResponseWithJson(respone, http.StatusOK, theme)	
}