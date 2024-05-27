package thumbnailgrpc

import (
	"context"
	"errors"
	thumbnailv1 "github.com/miittya/grpc-thumbnail/proto/gen/go/proto/thumbnail"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/url"
)

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
	_, err := url.ParseRequestURI(videoURL)
	if err != nil {
		return errors.New("invalid video url")
	}
	return nil
}
