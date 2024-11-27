package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
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
