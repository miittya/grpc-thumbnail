package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
	"github.com/miittya/grpc-thumbnail/server/internal/storage"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	op := "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &Storage{db: db}, nil
}

func (s *Storage) Stop() error {
	return s.db.Close()
}

// SaveThumbnail saves thumbnail in cache
func (s *Storage) SaveThumbnail(ctx context.Context, videoURL string, thumbnail []byte) error {
	op := "storage.sqlite.SaveThumbnail"

	_, err := s.db.ExecContext(
		ctx,
		"INSERT INTO thumbnails (video_url, thumbnail) VALUES (?, ?)",
		videoURL,
		string(thumbnail),
	)

	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return fmt.Errorf("%s: %w", op, storage.ErrUrlExists)
		}
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

// Thumbnail gets thumbnail from cache
func (s *Storage) Thumbnail(ctx context.Context, videoURL string) ([]byte, error) {
	op := "storage.sqlite.Thumbnail"

	stmt, err := s.db.Prepare("SELECT thumbnail FROM thumbnails WHERE video_url = ?")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	row := stmt.QueryRowContext(ctx, videoURL)

	var thumbnail []byte
	err = row.Scan(&thumbnail)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return thumbnail, nil
}
