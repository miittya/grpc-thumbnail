package grpcclient

import (
	"context"
	"fmt"
	thumbnailv1 "github.com/miittya/grpc-thumbnail/proto/gen/go/proto/thumbnail"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type Client struct {
	api thumbnailv1.ThumbnailServiceClient
}

func New(
	addr string,
) (*Client, error) {
	op := "grpc.New"

	cc, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Client{
		api: thumbnailv1.NewThumbnailServiceClient(cc),
	}, nil
}

// DownloadThumbnail downloads thumbnail of YouTube video by its URL
func (c *Client) DownloadThumbnail(ctx context.Context, videoURL string) {
	resp, err := c.api.Thumbnail(ctx, &thumbnailv1.ThumbnailRequest{
		VideoUrl: videoURL,
	})
	if err != nil {
		log.Printf("Failed to get thumbnail for %s: %v", videoURL, err)
	}
	if err := saveThumbnail(videoURL, resp.ThumbnailData); err != nil {
		log.Printf("Failed to save thumbnail for %s: %v", videoURL, err)
	}
}

func (c *Client) DownloadThumbnails(ctx context.Context, videoURLs []string) {
	for _, videoURL := range videoURLs {
		c.DownloadThumbnail(ctx, videoURL)
	}
}

func (c *Client) DownloadThumbnailsAsync(ctx context.Context, videoURLs []string) {
	var wg sync.WaitGroup
	for _, videoURL := range videoURLs {
		wg.Add(1)
		go func(videoURL string) {
			defer wg.Done()
			c.DownloadThumbnail(ctx, videoURL)
		}(videoURL)
	}
	wg.Wait()
}

func saveThumbnail(videoURL string, thumbnail []byte) error {
	fileName, err := getFileNameFromURL(videoURL)
	if err != nil {
		return err
	}
	dir := "./thumbnails"
	if err := os.MkdirAll(dir, 0777); err != nil {
		return err
	}
	filePath := filepath.Join(dir, fileName)
	if err := os.WriteFile(filePath, thumbnail, 0644); err != nil {
		return err
	}

	return nil
}

func getFileNameFromURL(videoURL string) (string, error) {
	videoID, err := extractVideoID(videoURL)
	if err != nil {
		log.Printf("Failed to extract video ID for %s: %v", videoURL, err)
	}
	return fmt.Sprintf("%s_thumbnail.jpg", videoID), nil
}

func extractVideoID(videoURL string) (string, error) {
	parsedURL, err := url.Parse(videoURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse URL: %w", err)
	}

	queryParams := parsedURL.Query()
	videoID := queryParams.Get("v")
	if videoID == "" {
		if parsedURL.Host == "youtu.be" {
			videoID = strings.TrimPrefix(parsedURL.Path, "/")
		}
	}

	if videoID == "" {
		return "", fmt.Errorf("no video ID found in URL")
	}

	return videoID, nil
}
