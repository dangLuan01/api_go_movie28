package models

import (
	//"fmt"

	"bytes"
	"encoding/json"
	"log"
	"os"

	"github.com/dangLuan01/api_go_movie28/config"
	"github.com/doug-martin/goqu/v9"
)

func SynES() any {
	posterSubquery := config.DB.From(goqu.T("movie_images").As("mi")).
		Where(
			goqu.I("mi.movie_id").Eq(goqu.I("movies.id")),
			goqu.I("mi.is_thumbnail").Eq(0),
		).
		Select(goqu.Func("CONCAT", goqu.I("mi.path"), goqu.I("mi.image"))).
		Limit(1)
	found := config.DB.From("movies").Select(
		goqu.I("movies.id"),
		goqu.I("movies.name"),
		goqu.I("movies.origin_name"),
		goqu.I("movies.slug"),
		goqu.I("movies.type"),
		goqu.I("movies.age"),
		goqu.I("movies.release_date"),
		goqu.I("movies.runtime"),
		posterSubquery.As("poster"),
	)
	type Movie struct {
		Id          	int    `json:"id"`
		Name        	string `json:"name"`
		Origin_name  	string `json:"origin_name"`
		Slug        	string `json:"slug"`
		Type        	string `json:"type"`
		Age         	string `json:"age"`
		Runtime     	string `json:"runtime"`
		Release_date 	int    `json:"release_date"`
		Poster      	string `json:"poster"`
	}
	var movie []Movie
	if err := found.ScanStructs(&movie); err != nil {
		log.Printf("Error scanning structs: %v", err)
		panic(err)
	}
	file, err := os.Create("movies.ndjson")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	for _, m := range movie {
		var buf bytes.Buffer
		meta := map[string]map[string]interface{}{
			"index": {
				"_index": "movies",
				"_id":    m.Id,
			},
		}
		if err := json.NewEncoder(&buf).Encode(meta); err != nil {
			log.Printf("Error encoding meta: %v", err)
			continue
		}
		data := map[string]interface{}{
			"name":         m.Name,
			"origin_name":  m.Origin_name,
			"slug":         m.Slug,
			"type":         m.Type,
			"release_date": m.Release_date,
			"age":          m.Age,
			"runtime":      m.Runtime,
			"poster":       m.Poster,
		}
		if err := json.NewEncoder(&buf).Encode(data); err != nil {
			log.Printf("Error encoding data: %v", err)
			continue
		}
		file.Write(buf.Bytes())
	}
	return "Syn completed successfully"
}