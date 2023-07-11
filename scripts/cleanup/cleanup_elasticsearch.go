package scripts

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

func CleanUpElasticsearch() error {
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{
			"https://localhost:9200",
		},
		Username: "elastic",
		Password: "password",
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	})
	if err != nil {
		return err
	}

	indexNames := []string{"users"}
	for _, indexName := range indexNames {
		deleteIndexReq := esapi.IndicesDeleteRequest{
			Index: []string{indexName},
		}

		res, err := deleteIndexReq.Do(context.Background(), es)
		if err != nil {
			return err
		}

		defer res.Body.Close()

		if res.StatusCode == 404 {
			continue
		}

		if res.IsError() {
			return fmt.Errorf("cannot delete index: %s", res)
		}
	}

	log.Print("Elasticsearch cleanup completed")

	return nil
}
