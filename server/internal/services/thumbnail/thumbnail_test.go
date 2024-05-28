package thumbnail

import (
	"context"
	"errors"
	"github.com/miittya/grpc-thumbnail/server/internal/lib/slogdiscard"
	"github.com/miittya/grpc-thumbnail/server/internal/services/thumbnail/mocks"
	storageErrors "github.com/miittya/grpc-thumbnail/server/internal/storage"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestThumbnailService_Thumbnail(t *testing.T) {
	type mockBehavior func(storage *mocks.Storage, client *mocks.Client, ctx context.Context, videoURL string)

	tests := []struct {
		name              string
		videoURL          string
		mockBehavior      mockBehavior
		expectedThumbnail []byte
		expectedErr       error
	}{
		{
			name:     "success from cache",
			videoURL: "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
			mockBehavior: func(storage *mocks.Storage, client *mocks.Client, ctx context.Context, videoURL string) {
				storage.On("Thumbnail", ctx, videoURL).Return([]byte{0x1, 0x2, 0x3, 0x4, 0x5, 0x6}, nil)
			},
			expectedThumbnail: []byte{0x1, 0x2, 0x3, 0x4, 0x5, 0x6},
			expectedErr:       nil,
		},
		{
			name:     "success from YouTube",
			videoURL: "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
			mockBehavior: func(storage *mocks.Storage, client *mocks.Client, ctx context.Context, videoURL string) {
				storage.On("Thumbnail", ctx, videoURL).Return(nil, storageErrors.ErrNotFound)
				client.On("Thumbnail", ctx, videoURL).Return([]byte{0x1, 0x2, 0x3, 0x4, 0x5, 0x6}, nil)
				storage.On("SaveThumbnail", ctx, videoURL, []byte{0x1, 0x2, 0x3, 0x4, 0x5, 0x6}).Return(nil)
			},
			expectedThumbnail: []byte{0x1, 0x2, 0x3, 0x4, 0x5, 0x6},
			expectedErr:       nil,
		},
		{
			name:     "YouTube fetch error",
			videoURL: "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
			mockBehavior: func(storage *mocks.Storage, client *mocks.Client, ctx context.Context, videoURL string) {
				storage.On("Thumbnail", ctx, videoURL).Return(nil, storageErrors.ErrNotFound)
				client.On("Thumbnail", ctx, videoURL).Return(nil, errors.New("failed to fetch thumbnail"))
			},
			expectedThumbnail: nil,
			expectedErr:       errors.New("failed to fetch thumbnail"),
		},
		{
			name:     "cache save error",
			videoURL: "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
			mockBehavior: func(storage *mocks.Storage, client *mocks.Client, ctx context.Context, videoURL string) {
				storage.On("Thumbnail", ctx, videoURL).Return(nil, storageErrors.ErrNotFound)
				client.On("Thumbnail", ctx, videoURL).Return([]byte{0x1, 0x2, 0x3, 0x4, 0x5, 0x6}, nil)
				storage.On("SaveThumbnail", ctx, videoURL, []byte{0x1, 0x2, 0x3, 0x4, 0x5, 0x6}).Return(errors.New("failed to save thumbnail in cache"))
			},
			expectedThumbnail: []byte{0x1, 0x2, 0x3, 0x4, 0x5, 0x6},
			expectedErr:       nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := &mocks.Storage{}
			client := &mocks.Client{}
			tt.mockBehavior(storage, client, context.Background(), tt.videoURL)

			service := New(slogdiscard.NewDiscardLogger(), storage, client)
			thumbnail, err := service.Thumbnail(context.Background(), tt.videoURL)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Nil(t, thumbnail)
				assert.Contains(t, err.Error(), tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, thumbnail)
				assert.Equal(t, tt.expectedThumbnail, thumbnail)
			}

			storage.AssertExpectations(t)
			client.AssertExpectations(t)
		})
	}
}
