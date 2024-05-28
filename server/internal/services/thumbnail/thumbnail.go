package thumbnail

import (
	"context"
	"errors"
	"fmt"
	"github.com/miittya/grpc-thumbnail/server/internal/lib/sl"
	"github.com/miittya/grpc-thumbnail/server/internal/storage"
	"log/slog"
)

type ThumbnailService struct {
	log      *slog.Logger
	ytClient Client
	storage  Storage
}

//go:generate go run github.com/vektra/mockery/v2@v2.43.0 --name=Storage
type Storage interface {
	SaveThumbnail(ctx context.Context, videoURL string, thumbnail []byte) error
	Thumbnail(ctx context.Context, videoURL string) ([]byte, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.43.0 --name=Client
type Client interface {
	Thumbnail(ctx context.Context, videoURL string) ([]byte, error)
}

func New(
	log *slog.Logger,
	storage Storage,
	client Client,
) *ThumbnailService {
	return &ThumbnailService{
		log:      log,
		storage:  storage,
		ytClient: client,
	}
}

// Thumbnail gets thumbnail by URL
func (s *ThumbnailService) Thumbnail(ctx context.Context, videoURL string) ([]byte, error) {
	op := "thumbnail.Thumbnail"

	log := s.log.With(
		slog.String("op", op),
		slog.String("videoURL", videoURL),
	)

	log.Info("getting thumbnail")

	thumbnail, err := s.storage.Thumbnail(ctx, videoURL)
	if err != nil {
		if !errors.Is(err, storage.ErrNotFound) {
			log.Info(err.Error())
			log.Info(storage.ErrNotFound.Error())
			log.Info("hui", slog.Bool("hui", errors.Is(err, storage.ErrNotFound)))
			log.Error("failed to get thumbnail from db", sl.Err(err))
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		log.Warn("thumbnail not found in cache, fetching from youtube")

		thumbnail, err = s.ytClient.Thumbnail(ctx, videoURL)
		if err != nil {
			log.Error("failed to fetch thumbnail from youtube", sl.Err(err))
			return nil, fmt.Errorf("%s: %w", "failed to fetch thumbnail from youtube", err)
		}

		log.Info("successfully fetched thumbnail from youtube")
		log.Info("saving thumbnail in cache...")

		err = s.storage.SaveThumbnail(ctx, videoURL, thumbnail)
		if err != nil {
			log.Warn("failed to save thumbnail in cache", sl.Err(err))
		}
	}
	return thumbnail, nil
}
