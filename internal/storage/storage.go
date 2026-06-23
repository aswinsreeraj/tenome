package storage

type Storage interface {
	SavePage(ctx context.Context, page Page) error
}
