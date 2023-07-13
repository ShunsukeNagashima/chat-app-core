package repository

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/shunsukenagashima/chat-api/pkg/domain/repository"
)

type ElasticsearchRepositoryImpl struct {
	es *elasticsearch.Client
}

func NewElasticsearchRepository(es *elasticsearch.Client) repository.ElasticsearchRepository {
	return &ElasticsearchRepositoryImpl{
		es,
	}
}

func (r *ElasticsearchRepositoryImpl) Create(ctx context.Context, index, documentId string, body io.Reader) error {
	req := esapi.IndexRequest{
		Index:      index,
		DocumentID: documentId,
		Body:       body,
		Refresh:    "true",
	}

	retryCount := 3
	var lastErr error
	for i := 0; i < retryCount; i++ {
		res, err := req.Do(ctx, r.es)
		if err != nil {
			lastErr = err
			log.Printf("Error when trying to create index for %s: %v", index, err)
			continue
		}
		defer res.Body.Close()

		if res.IsError() {
			lastErr = fmt.Errorf("unexpected status code: %d", res.StatusCode)
			log.Printf("%s", lastErr.Error())
			continue
		}

		lastErr = nil
		break
	}

	if lastErr != nil {
		return fmt.Errorf("failed to create index for %s: %s", index, lastErr)
	}

	return nil
}

func (r *ElasticsearchRepositoryImpl) Update(ctx context.Context, index, documentId string, body io.Reader) error {
	req := esapi.UpdateRequest{
		Index:      index,
		DocumentID: documentId,
		Body:       body,
		Refresh:    "true",
	}

	retryCount := 3
	var lastErr error
	for i := 0; i < retryCount; i++ {
		res, err := req.Do(ctx, r.es)
		if err != nil {
			lastErr = err
			log.Printf("Error when trying to update index for %s: %v", index, err)
			continue
		}
		defer res.Body.Close()

		if res.IsError() {
			lastErr = fmt.Errorf("unexpected status code: %d", res.StatusCode)
			log.Printf("%s", lastErr.Error())
			continue
		}

		lastErr = nil
		break
	}

	if lastErr != nil {
		return fmt.Errorf("failed to update index for %s: %s", index, lastErr)
	}

	return nil
}

func (r *ElasticsearchRepositoryImpl) Delete(ctx context.Context, index, documentId string) error {
	req := esapi.DeleteRequest{
		Index:      index,
		DocumentID: documentId,
		Refresh:    "true",
	}

	retryCount := 3
	var lastErr error
	for i := 0; i < retryCount; i++ {
		res, err := req.Do(ctx, r.es)
		if err != nil {
			lastErr = err
			log.Printf("Error when trying to delete index for %s: %v", index, err)
			continue
		}
		defer res.Body.Close()

		if res.IsError() {
			lastErr = fmt.Errorf("unexpected status code: %d", res.StatusCode)
			log.Printf("%s", lastErr.Error())
			continue
		}

		lastErr = nil
		break
	}

	if lastErr != nil {
		return fmt.Errorf("failed to delete index for %s: %s", index, lastErr)
	}

	return nil
}
