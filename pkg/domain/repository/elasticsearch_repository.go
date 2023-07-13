package repository

import (
	"context"
	"io"
)

//go:generate mockery --name=ElasticsearchRepository --output=mocks
type ElasticsearchRepository interface {
	Create(ctx context.Context, index, id string, body io.Reader) error
	Update(ctx context.Context, index, id string, body io.Reader) error
	Delete(ctx context.Context, index, id string) error
}
