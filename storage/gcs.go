package storage

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/storage"
)

var GCSBucket *GCSBackend

type GCSBackend struct {
	Client *storage.Client
	Bucket string
}

func InitGCSBackend(bckt string) error {
	ctx := context.Background()
	fmt.Println("bucket name :", bckt)
	client, err := storage.NewClient(ctx)
	fmt.Println("Client : ", client)
	if err != nil {
		return err
	}

	GCSBucket = &GCSBackend{
		Client: client,
		Bucket: bckt,
	}

	attrs, err := GCSBucket.Client.Bucket(GCSBucket.Bucket).Attrs(ctx)
	if err == storage.ErrBucketNotExist {
		fmt.Fprintln(os.Stderr, "The", GCSBucket.Bucket, "bucket does not exist")
		return err
	}
	if err != nil {
		// Other error to handle
		fmt.Fprintln(os.Stderr, err)
	}
	fmt.Println("The", GCSBucket.Bucket, "bucket exists and has attributes:", attrs)
	return err
}
