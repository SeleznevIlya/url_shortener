package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"url_shortener/internal/storage"

	"github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

// New creates a new instance of the PostgreSQL storage.
func New(storagePath string) (*Storage, error) {
	const op = "storage.postgres.New"

	db, err := sql.Open("postgres", storagePath)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(urlToSave string, alias string) (int64, error) {
	const op = "storage.postgres.SaveURL"

	stmt, err := s.db.Prepare("INSERT INTO url(url_name, alias) VALUES ($1, $2) RETURNING id")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	var id int64

	err = stmt.QueryRow(urlToSave, alias).Scan(&id)
	if err != nil {
		var postgresErr *pq.Error
		if errors.As(err, &postgresErr) && postgresErr.Code == pq.ErrorCode("23505") {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrUrlExists)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.postgres.GetURL"

	stmt, err := s.db.Prepare("SELECT url_name FROM url WHERE alias=$1")
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	var url string

	err = stmt.QueryRow(alias).Scan(&url)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("%s: %w", op, storage.ErrUrlNotFound)
		}
		return "", fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return url, nil
}

func (s *Storage) DeleteURL(alias string) error {
	const op = "storage.postgres.DeleteURL"

	stmt, err := s.db.Prepare("DELETE FROM url WHERE alias=$1")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(alias)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%s: %w", op, storage.ErrUrlNotFound)
		}
		return fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return nil
}
