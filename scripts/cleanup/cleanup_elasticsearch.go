package scripts

import (
	"context"
	"crypto/tls"
	"encoding/json"
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

	catIndexReq := esapi.CatIndicesRequest{
		Format: "json",
	}
	res, err := catIndexReq.Do(context.Background(), es)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	var indices []map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&indices)
	if err != nil {
		return err
	}

	for _, index := range indices {
		indexName := index["index"].(string)
		fmt.Println("Deleting index:", indexName)

		deleteReq := esapi.IndicesDeleteRequest{
			Index: []string{indexName},
		}

		deleteRes, err := deleteReq.Do(context.Background(), es)
		if err != nil {
			continue
		}
		defer deleteRes.Body.Close()

		if deleteRes.IsError() {
			return fmt.Errorf("error deleting index %s: %s", indexName, deleteRes.String())
		}
	}

	log.Print("Elasticsearch cleanup completed")

	return nil
}
