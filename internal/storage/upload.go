package storage

import (
	"context"
	"fmt"
	"io"
)

func (b *GCSBackend) UploadModule(path string, buff io.Reader) (string, error) {
	ctx := context.Background()
	wr := b.Client.Bucket(b.Bucket).Object(path).NewWriter(ctx)

	if _, err := io.Copy(wr, buff); err != nil {
		return "", fmt.Errorf("upload failed - %v", err.Error())
	}
	if err := wr.Close(); err != nil {
		return "", fmt.Errorf("upload failed - %v", err.Error())
	}

	return fmt.Sprintf("gcs::https://www.googleapis.com/storage/v1/%s/%s", b.Bucket, path), nil
}
