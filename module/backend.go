package module

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/storage"
)

var gcsBucket Backend

type Backend interface {
	getModule(mod Module, ctx context.Context) (string, error)
	// listModuleVersions(mod Module)
}

type GCSBackend struct {
	client *storage.Client
	bucket string
}

func InitGCSBackend(bckt string) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}

	gcsBucket := &GCSBackend{
		client: client,
		bucket: bckt,
	}

	attrs, err := gcsBucket.client.Bucket(gcsBucket.bucket).Attrs(ctx)
	if err == storage.ErrBucketNotExist {
		fmt.Fprintln(os.Stderr, "The", gcsBucket.bucket, "bucket does not exist")
		return err
	}
	if err != nil {
		// Other error to handle
		fmt.Fprintln(os.Stderr, err)
	}
	fmt.Println("The", gcsBucket.bucket, "bucket exists and has attributes:", attrs)
	return err
}
