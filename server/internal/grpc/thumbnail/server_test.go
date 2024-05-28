package thumbnailgrpc

import (
	"context"
	"errors"
	thumbnailv1 "github.com/miittya/grpc-thumbnail/proto/gen/go/proto/thumbnail"
	"github.com/miittya/grpc-thumbnail/server/internal/grpc/thumbnail/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestServerAPI_Thumbnail(t *testing.T) {
	type mockBehavior func(thumbnail *mocks.ThumbnailService, ctx context.Context, videoURL string)

	tests := []struct {
		name              string
		videoURL          string
		mockBehavior      mockBehavior
		expectedThumbnail []byte
		expectedError     error
	}{
		{
			name:     "success",
			videoURL: "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
			mockBehavior: func(thumbnail *mocks.ThumbnailService, ctx context.Context, videoURL string) {
				thumbnail.On("Thumbnail", ctx, videoURL).Return([]byte{0x1, 0x2, 0x3, 0x4, 0x5, 0x6}, nil)
			},
			expectedThumbnail: []byte{0x1, 0x2, 0x3, 0x4, 0x5, 0x6},
			expectedError:     nil,
		},
		{
			name:              "empty url",
			videoURL:          "",
			mockBehavior:      func(thumbnail *mocks.ThumbnailService, ctx context.Context, videoURL string) {},
			expectedThumbnail: nil,
			expectedError:     errors.New("invalid video url"),
		},
		{
			name:              "invalid video id",
			videoURL:          "https://www.youtube.com/watch?v=-1",
			mockBehavior:      func(thumbnail *mocks.ThumbnailService, ctx context.Context, videoURL string) {},
			expectedThumbnail: nil,
			expectedError:     errors.New("invalid video id"),
		},
		{
			name:              "invalid host",
			videoURL:          "www.google.com",
			mockBehavior:      func(thumbnail *mocks.ThumbnailService, ctx context.Context, videoURL string) {},
			expectedThumbnail: nil,
			expectedError:     errors.New("invalid video url"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			thumbnailService := mocks.NewThumbnailService(t)
			tt.mockBehavior(thumbnailService, context.Background(), tt.videoURL)
			srv := serverAPI{thumbnailService: thumbnailService}

			req := thumbnailv1.ThumbnailRequest{VideoUrl: tt.videoURL}
			resp, err := srv.Thumbnail(context.Background(), &req)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Nil(t, resp)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.expectedThumbnail, resp.ThumbnailData)
			}

			thumbnailService.AssertExpectations(t)
		})
	}
}
