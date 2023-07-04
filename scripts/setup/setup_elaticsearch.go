package scripts

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/shunsukenagashima/chat-api/pkg/domain/model"
)

func SetUpElasticsearch(user *model.User) error {
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

	checkExistsReq := esapi.IndicesExistsRequest{
		Index: []string{"users"},
	}

	res, err := checkExistsReq.Do(context.Background(), es)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode == 404 {
		mapping := `
		{
			"mappings": {
				"properties": {
					"userId": { "type": "keyword" },
					"userName": { "type": "text" },
					"email": { "type": "keyword" }
				}
			}
		}
		`
		res, err = es.Indices.Create(
			"users",
			es.Indices.Create.WithContext(context.Background()),
			es.Indices.Create.WithBody(strings.NewReader(mapping)),
		)

		defer res.Body.Close()

		if err != nil {
			return err
		}
		if res.IsError() {
			return fmt.Errorf("cannot create index: %s", res)
		}
	}

	userJson, err := json.Marshal(user)
	if err != nil {
		return err
	}

	createIndexReq := esapi.IndexRequest{
		Index:      "users",
		DocumentID: user.UserID,
		Body:       bytes.NewReader(userJson),
	}

	res, err = createIndexReq.Do(context.Background(), es)

	if err != nil {
		return err
	}
	if res.IsError() {
		return fmt.Errorf("cannot create index: %s", res.String())
	}

	defer res.Body.Close()

	return nil
}
