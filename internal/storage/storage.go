package storage

import (
	"context"
	"tenome/internal/model"
)

type Storage interface {
	Migrate(ctx context.Context) error
	SavePage(ctx context.Context, page model.Page) (int64, error)
}
