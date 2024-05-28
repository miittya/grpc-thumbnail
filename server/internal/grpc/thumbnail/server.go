package thumbnailgrpc

import (
	"context"
	"errors"
	thumbnailv1 "github.com/miittya/grpc-thumbnail/proto/gen/go/proto/thumbnail"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/url"
	"regexp"
)

// Regular expression for validating YouTube video ID
var (
	videoIDRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]{11}$`)
)

//go:generate go run github.com/vektra/mockery/v2@v2.43.0 --name=ThumbnailService
type ThumbnailService interface {
	Thumbnail(ctx context.Context, videoURL string) ([]byte, error)
}

type serverAPI struct {
	thumbnailv1.UnimplementedThumbnailServiceServer
	thumbnailService ThumbnailService
}

func Register(gRPCServer *grpc.Server, thumbnailService ThumbnailService) {
	thumbnailv1.RegisterThumbnailServiceServer(gRPCServer, &serverAPI{thumbnailService: thumbnailService})
}

func (s *serverAPI) Thumbnail(
	ctx context.Context,
	req *thumbnailv1.ThumbnailRequest,
) (*thumbnailv1.ThumbnailResponse, error) {
	if err := validateURL(req.GetVideoUrl()); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	thumbnail, err := s.thumbnailService.Thumbnail(ctx, req.GetVideoUrl())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &thumbnailv1.ThumbnailResponse{ThumbnailData: thumbnail}, nil
}

func validateURL(videoURL string) error {
	parsedURL, err := url.ParseRequestURI(videoURL)
	if err != nil {
		return errors.New("invalid video url")
	}

	if parsedURL.Host != "www.youtube.com" && parsedURL.Host != "youtube.com" {
		return errors.New("invalid video url")
	}

	queryParams := parsedURL.Query()
	videoID := queryParams.Get("v")

	if !videoIDRegex.MatchString(videoID) {
		return errors.New("invalid video id")
	}
	return nil
}
