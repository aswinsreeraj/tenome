package cache

import (
	"context"
	"tenome/internal/model"
)

type Cache interface {
	Get(ctx context.Context, key string) ([]model.Page, bool, error)
	Set(ctx context.Context, key string, pages []model.Page) error
}
