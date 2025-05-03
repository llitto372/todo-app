package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"todo/internal/config/storage"
)

type Storage struct {
	db *sql.DB
}

func New(connStr string) (*Storage, error) {
	const op = "storage.postgres.New"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS url(
    	id SERIAL PRIMARY KEY,
   		alias TEXT UNIQUE NOT NULL,
    	url TEXT NOT NULL);
    CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(URLtoSave string, alias string) error {
	const op = "storage.postgres.SaveURL"

	stmt, err := s.db.Prepare("INSERT INTO url(alias, url) VALUES ($1, $2)")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(alias, URLtoSave)
	var pqErr *pq.Error
	if err != nil {
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return fmt.Errorf("%s: %w", op, storage.ErrURLExists)
		}
	}
	return nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.postgres.GetURL"
	stmt, err := s.db.Prepare("SELECT url FROM url WHERE alias = $1")
	if err != nil {
		return "", storage.ErrURLNotFound
	}
	defer stmt.Close()
	var resultUrl string
	err = stmt.QueryRow(alias).Scan(&resultUrl)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return resultUrl, nil
}

func (s *Storage) DeleteURL(alias string) error {
	const op = "storage.postgres."
	stmt, err := s.db.Prepare("DELETE FROM url WHERE alias = $1")
	if err != nil {
		fmt.Printf("%s: %s\n", op, err)
	}
	_, err = stmt.Exec(alias)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return err
}
