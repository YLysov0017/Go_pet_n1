package sqlite

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/YLysov0017/go_pet_n1/internal/config/storage"
	"github.com/mattn/go-sqlite3" // init driver
)

type Storage struct {
	db *sql.DB
}

const errorf string = "%s: %w"

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf(errorf, op, err) // ошибка случилась при создании соединения с sqlite
	}

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS url(
		ID INTEGER PRIMARY KEY,
		alias TEXT NOT NULL UNIQUE,
		url TEXT NOT NULL UNIQUE);
	CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
	`)

	if err != nil {
		return nil, fmt.Errorf(errorf, op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf(errorf, op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(urlToSave, alias string) (int64, error) {
	const op = "storage.sqlite.SaveURL"
	stmt, err := s.db.Prepare(`SELECT COUNT(*) FROM url WHERE alias = ?`)
	if err != nil {
		return 0, fmt.Errorf(errorf, op, err)
	}
	var resURL int

	_ = stmt.QueryRow(alias).Scan(&resURL)
	if resURL != 0 {
		return 0, fmt.Errorf(errorf, op, storage.ErrURLExists)
	}

	stmt, err = s.db.Prepare(`INSERT INTO url(url, alias) VALUES(?, ?)`)
	if err != nil {
		return 0, fmt.Errorf(errorf, op, err)
	}

	res, err := stmt.Exec(urlToSave, alias)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf(errorf, op, storage.ErrURLExists)
		}
		return 0, fmt.Errorf(errorf, op, err)
	}

	id, err := res.LastInsertId() // Not available in every SQL Database
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get last insert id: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.sqlite.GetURL"
	stmt, err := s.db.Prepare(`SELECT url FROM url WHERE alias = ?`)
	if err != nil {
		return "", fmt.Errorf(errorf, op, err)
	}
	var resURL string

	err = stmt.QueryRow(alias).Scan(&resURL)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return "", storage.ErrURLNotFound
	case err != nil:
		return "", fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return resURL, nil
}

func (s *Storage) DeleteURL(alias string) error {
	const op = "storage.sqlite.DeleteURL"
	stmt, err := s.db.Prepare(`DELETE FROM url WHERE alias = ?`)
	if err != nil {
		return fmt.Errorf(errorf, op, err)
	}
	_, err = stmt.Exec(alias)

	if err != nil {
		return fmt.Errorf(errorf, op, err)
	}

	return nil
}
