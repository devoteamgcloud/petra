package storage

import (
	"context"
	"fmt"

	"cloud.google.com/go/storage"
)

type GCSBackend struct {
	Client *storage.Client
	Bucket string
}

func InitGCSBackend(bucket string) (*GCSBackend, error) {
	ctx := context.Background()
	fmt.Println("\n" + bucket)

	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	if err != nil {
		return nil, err
	}

	backend := &GCSBackend{
		Client: client,
		Bucket: bucket,
	}

	return backend, nil
}
