package yt

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"
)

type Client struct {
	httpClient *http.Client
}

func New(httpClient *http.Client) *Client {
	return &Client{
		httpClient: httpClient,
	}
}

var videoIDRegex = regexp.MustCompile(`(?:v=|\/)([0-9A-Za-z_-]{10}[048AEIMQUYcgkosw])`)

// extractVideoID extracts video ID from URL
func extractVideoID(videoURL string) (string, error) {
	matches := videoIDRegex.FindStringSubmatch(videoURL)
	if len(matches) < 2 {
		return "", fmt.Errorf("invalid video URL: %s", videoURL)
	}
	return matches[1], nil
}

// Thumbnail gets thumbnail from YouTube video
func (c *Client) Thumbnail(ctx context.Context, videoURL string) ([]byte, error) {
	op := "yt.Thumbnail"
	videoID, err := extractVideoID(videoURL)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	thumbnailURL := fmt.Sprintf("https://img.youtube.com/vi/%s/default.jpg", videoID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, thumbnailURL, nil)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", op, fmt.Errorf(http.StatusText(resp.StatusCode)))
	}

	thumbnail, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return thumbnail, nil
}
