package models

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"os"
	"github.com/dangLuan01/api_go_movie28/entities"
	"github.com/dangLuan01/api_go_movie28/internal/elasticsearch"
)
var (
	r map[string]interface{}
	index = os.Getenv("ES_INDEX")
	
)
func Search(search string) (entities.SearchResult, error) {
	movie := []entities.Movie{}
	es := elasticsearch.ConnectES()
	if es == nil {
		log.Printf("Model Elasticsearch: %v", es)
		return entities.SearchResult{}, nil
	}
	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  search,
				"fields": []string{"name^2", "origin_name"},
				"fuzziness": "AUTO",
			},
		},
		"size": 17,
		"from": 0,
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}
	
	res, err := es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex(index),
		es.Search.WithBody(&buf),
		es.Search.WithTrackTotalHits(true),
		es.Search.WithPretty(),
	)
	if err != nil {
		log.Printf("Elasticsearch search error: %v", err)
		return entities.SearchResult{}, err 
	}
	defer res.Body.Close()
	if res.IsError() {
		log.Printf("Elasticsearch response error: %s", res.String())
		
	}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Printf("Error parsing the response body: %s", err)
	}
	
	hits := r["hits"].(map[string]interface{})["hits"].([]interface{})

	for _, hit := range hits {
		hitMap := hit.(map[string]interface{})
		if source, ok := hitMap["_source"].(map[string]interface{}); ok {
			movie = append(movie, entities.Movie{
				Name: source["name"].(string),
				Origin_name: source["origin_name"].(string),
				Slug: source["slug"].(string),
				Image: entities.Image{
					Poster: source["poster"].(string),
				},	
				Type: source["type"].(string),
				Age: source["age"].(string),
				Release_date: int(source["release_date"].(float64)),
				Runtime: source["runtime"].(string),
			})
		}
	}
	return entities.SearchResult{
		Movies: movie,
	}, nil
}