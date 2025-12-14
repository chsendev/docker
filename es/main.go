package main

import (
	"context"
	"log"

	"github.com/elastic/go-elasticsearch/v9"
)

func main() {
	es, _ := elasticsearch.NewTypedClient(elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
	})
	defer es.Close(context.Background())
	log.Println(es.Info().Do(context.Background()))
}
