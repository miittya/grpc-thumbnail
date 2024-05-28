package suite

import (
	"context"
	thumbnailv1 "github.com/miittya/grpc-thumbnail/proto/gen/go/proto/thumbnail"
	"github.com/miittya/grpc-thumbnail/server/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"strconv"
	"testing"
	"time"
)

type Suite struct {
	*testing.T
	Cfg             *config.Config
	ThumbnailClient thumbnailv1.ThumbnailServiceClient
}

const (
	grpcHost         = "localhost"
	maxRetryAttempts = 10
	retryDelay       = 2 * time.Second
)

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	cfg := config.MustLoadPath("../config/local.yaml")

	ctx, cancelCtx := context.WithTimeout(context.Background(), cfg.GRPC.Timeout)

	t.Cleanup(func() {
		t.Helper()
		cancelCtx()
	})

	var cc *grpc.ClientConn
	var err error
	for att := 0; att < maxRetryAttempts; att++ {
		cc, err = grpc.NewClient(
			net.JoinHostPort(grpcHost, strconv.Itoa(cfg.GRPC.Port)),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err == nil {
			break
		}
		t.Logf("grpc server connection attempt %d failed: %v", att+1, err)
		time.Sleep(retryDelay)
	}

	if err != nil {
		t.Fatalf("grpc server connection failed after %d attempts: %v", maxRetryAttempts, err)
	}

	return ctx, &Suite{
		T:               t,
		Cfg:             cfg,
		ThumbnailClient: thumbnailv1.NewThumbnailServiceClient(cc),
	}
}
