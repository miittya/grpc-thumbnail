package tests

import (
	thumbnailv1 "github.com/miittya/grpc-thumbnail/proto/gen/go/proto/thumbnail"
	"github.com/miittya/grpc-thumbnail/server/tests/suite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestThumbnail_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	videoURL := "https://www.youtube.com/watch?v=dQw4w9WgXcQ"
	resp, err := st.ThumbnailClient.Thumbnail(ctx, &thumbnailv1.ThumbnailRequest{VideoUrl: videoURL})
	require.NoError(t, err)
	assert.NotNil(t, resp.ThumbnailData)
	assert.True(t, len(resp.ThumbnailData) > 0)
}

func TestThumbnail_FailCases(t *testing.T) {
	ctx, st := suite.New(t)

	tests := []struct {
		name        string
		videoURL    string
		expectedErr string
	}{
		{
			name:        "empty url",
			videoURL:    "",
			expectedErr: "invalid video url",
		},
		{
			name:        "invalid video url",
			videoURL:    "https://www.youtube.com/watch?v=-1",
			expectedErr: "invalid video url",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := st.ThumbnailClient.Thumbnail(ctx, &thumbnailv1.ThumbnailRequest{VideoUrl: tt.videoURL})
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.expectedErr)
		})
	}
}
