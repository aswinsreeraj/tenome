package storage

import (
	"context"
	"database/sql"
	"strings"
	"tenome/internal/model"
)

type SQLiteStorage struct {
	db *sql.DB
}

func New(db *sql.DB) *SQLiteStorage {
	return &SQLiteStorage{db: db}
}

func (s *SQLiteStorage) Migrate(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, `
					CREATE TABLE IF NOT EXISTS pages (
						id INTEGER PRIMARY KEY AUTOINCREMENT,
						url TEXT UNIQUE,
						title TEXT,
						content TEXT
						)
				`)
	return err
}

func (s *SQLiteStorage) SavePage(ctx context.Context, page model.Page) (int64, error) {
	result, err := s.db.ExecContext(ctx, `
		INSERT INTO pages (
					url,
					title,
					content
				) VALUES (?, ?, ?)
	`, page.URL, page.Title, page.Content)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *SQLiteStorage) GetPagesByIDs(ctx context.Context, ids []int64) ([]model.Page, error) {
	if len(ids) == 0 {
		return []model.Page{}, nil
	}
	query := "SELECT id, url, title, content FROM pages WHERE id IN (" + strings.Repeat("?, ", len(ids)-1) + "?)"
	args := make([]any, len(ids))
	for i := range ids {
		args[i] = ids[i]
	}
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pages []model.Page
	for rows.Next() {
		var page model.Page

		if err := rows.Scan(
			&page.ID,
			&page.URL,
			&page.Title,
			&page.Content,
		); err != nil {
			return nil, err
		}
		pages = append(pages, page)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return pages, nil
}
