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
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	_, err = client.Bucket(bucket).Attrs(ctx)
	if err != nil {
		return nil, fmt.Errorf("Bucket(%q): %v", bucket, err)
	}
	backend := &GCSBackend{
		Client: client,
		Bucket: bucket,
	}

	return backend, nil
}
