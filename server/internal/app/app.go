package app

import (
	grpcapp "github.com/miittya/grpc-thumbnail/server/internal/app/grpc"
	"github.com/miittya/grpc-thumbnail/server/internal/clients/yt"
	"github.com/miittya/grpc-thumbnail/server/internal/services/thumbnail"
	"github.com/miittya/grpc-thumbnail/server/internal/storage/sqlite"
	"log/slog"
	"net/http"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(
	log *slog.Logger,
	gRPCPort int,
	storagePath string,
) *App {
	storage, err := sqlite.New(storagePath)
	if err != nil {
		panic("failed to open storage")
	}

	client := yt.New(http.DefaultClient)

	thumbnailService := thumbnail.New(log, storage, client)
	grpcApp := grpcapp.New(log, thumbnailService, gRPCPort)

	return &App{
		GRPCSrv: grpcApp,
	}
}
