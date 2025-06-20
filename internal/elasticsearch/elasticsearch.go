package elasticsearch

import (
	"log"
	"os"

	"github.com/elastic/go-elasticsearch/v7"
)
var host = os.Getenv("ES_HOST")
func ConnectES() *elasticsearch.Client {
    es, err := elasticsearch.NewClient(elasticsearch.Config{
        Addresses: []string{
           host,
        },
    }) 
    if err != nil {
        log.Printf("Error creating Elasticsearch client: %s", err)
    }
    return es
}
