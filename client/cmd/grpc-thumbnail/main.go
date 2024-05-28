package main

import (
	"context"
	"flag"
	grpcclient "github.com/miittya/grpc-thumbnail/client/internal/clients/thumbnail/grpc"
	"log"
)

func main() {
	async := flag.Bool("async", false, "download thumbnails asynchronously")
	flag.Parse()

	videoURLs := flag.Args()
	if len(videoURLs) == 0 {
		log.Fatal("no video url provided")
	}

	client, err := grpcclient.New("localhost:8080")
	if err != nil {
		log.Fatal(err)
	}

	if *async {
		client.DownloadThumbnailsAsync(context.Background(), videoURLs)
	} else {
		client.DownloadThumbnails(context.Background(), videoURLs)
	}
}
