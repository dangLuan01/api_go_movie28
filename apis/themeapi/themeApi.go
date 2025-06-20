package themeapi

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/dangLuan01/api_go_movie28/apis/utilapi"
	"github.com/dangLuan01/api_go_movie28/entities"
	"github.com/dangLuan01/api_go_movie28/internal/cacheloader"
	"github.com/dangLuan01/api_go_movie28/models"
)
func GetThemes(respone http.ResponseWriter, request *http.Request) {
	query 			:= request.URL.Query()
	idGet 			:= query.Get("id")
	pageThemeGet 	:= query.Get("page_theme")
	pageMovieGet 	:= query.Get("page_movie")
	pageSizeGet 	:= query.Get("page_size")
	found 			:= false
	id, _ 			:= strconv.Atoi(idGet)
	var data entities.PagiateTheme
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
	key 		:= fmt.Sprintf("themes:page=%s:size=%s", pageThemeGet, pageSizeGet)
	time 		:=  time.Duration(rand.Intn(500-400) + 400)
	themeCache 	:= cacheloader.GetCache(0, time)
	if themeCache != nil && themeCache.Get(key, &data) {
		found = true
		log.Println("Read from redis")
		utilapi.ResponseWithJson(respone, http.StatusOK, data)
		return
	}
	if !found {
		themes, err := models.GetAllThemesWithMovies(id, pageTheme, pageMovie, pageSize)
		if err != nil {
			utilapi.ResponseWithJson(respone, http.StatusOK, err)	
		}
		data = themes
		if themeCache != nil {
			themeCache.Set(key, data)	
		}
	}
	utilapi.ResponseWithJson(respone, http.StatusOK, data)
}