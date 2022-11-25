package storage

import (
	"context"
	"fmt"
	"io"
	"path"

	"cloud.google.com/go/storage"
)

type GCSBackend struct {
	client *storage.Client
	bucket string
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
		client: client,
		bucket: bucket,
	}

	return backend, nil
}

func (b *GCSBackend) UploadModule(namespace string, name string, provider string, version string, buff io.Reader) (string, error) {
	ctx := context.Background()
	var objPath = modPath(namespace, name, provider, version)
	wr := b.client.Bucket(b.bucket).Object(objPath).NewWriter(ctx)

	if _, err := io.Copy(wr, buff); err != nil {
		return "", fmt.Errorf("upload failed - %v", err.Error())
	}
	if err := wr.Close(); err != nil {
		return "", fmt.Errorf("upload failed - %v", err.Error())
	}

	return fmt.Sprintf("gcs::https://www.googleapis.com/storage/v1/%s/%s", b.bucket, objPath), nil
}

/* Utils */

func modPath(namespace string, name string, provider string, version string) string {
	return path.Join(
		namespace,
		name,
		provider,
		version,
		fmt.Sprintf("%s-%s-%s-%s.tar.gz", namespace, name, provider, version),
	)
}
