package storage

import (
	"context"
	"fmt"
	"os"
	"time"

	"cloud.google.com/go/storage"
)

func (b *GCSBackend) GetModule(path string, ctx context.Context) (string, error) {

	enable_signed_url := os.Getenv("SIGNED_URL")
	if enable_signed_url == "true" {

		var options *storage.SignedURLOptions

		fmt.Println("Use existing credentials to create signed url")
		options = &storage.SignedURLOptions{
			Method:  "GET",
			Expires: time.Now().Add(2 * time.Minute),
		}

		signedUrl, err := b.Client.Bucket(b.Bucket).SignedURL(path, options)
		if err != nil {
			return "", fmt.Errorf("Bucket(%q).SignedURL: %v", b.Bucket, err)
		}

		fmt.Println("Generated GET signed URL:")
		fmt.Printf("%q\n", signedUrl)
		return fmt.Sprintln(signedUrl), nil
	}

	object := b.Client.Bucket(b.Bucket).Object(path)
	attrs, err := object.Attrs(ctx)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("gcs::https://www.googleapis.com/storage/v1/%s/%s", b.Bucket, attrs.Name), nil
}
