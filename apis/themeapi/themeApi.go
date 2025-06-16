package themeapi

import (
	"net/http"
	"strconv"

	"github.com/dangLuan01/api_go_movie28/apis/utilapi"
	"github.com/dangLuan01/api_go_movie28/models"
)
func GetThemes(respone http.ResponseWriter, request *http.Request)  {
	query := request.URL.Query()
	idGet := query.Get("id")
	pageThemeGet := query.Get("page_theme")
	pageMovieGet := query.Get("page_movie")
	pageSizeGet := query.Get("page_size")

	id, _ := strconv.Atoi(idGet)

	pageTheme, err := strconv.Atoi(pageThemeGet)
	if err != nil || pageTheme < 1 {
		pageTheme = 1
	}
	pageMovie, err := strconv.Atoi(pageMovieGet)
	if err != nil || pageMovie < 1 {
		pageMovie = 1
	}
	pageSize, err := strconv.Atoi(pageSizeGet)
	if err != nil || pageSize < 1 {
		pageSize = 4
	}
	theme, err := models.GetAllThemesWithMovies(id, pageTheme, pageMovie, pageSize)
	if err != nil {
		utilapi.ResponseWithJson(respone, http.StatusOK, err)	
	}
	utilapi.ResponseWithJson(respone, http.StatusOK, theme)
}