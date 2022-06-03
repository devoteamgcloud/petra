package internal

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
)

func removeFile(w io.Writer, bucket, modulePath string) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	// Get info from petra config file
	petraConf, err := GetPetraConfig(modulePath)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}

	// Get full path of the .zip object in the bucket
	object := GetObjectPathFromConfig(petraConf)
	fmt.Printf("Object to delete: %v\n", object)

	o := client.Bucket(bucket).Object(object)

	// Optional: set a generation-match precondition to avoid potential race
	// conditions and data corruptions. The request to upload is aborted if the
	// object's generation number does not match your precondition.
	attrs, err := o.Attrs(ctx)
	if err != nil {
		return fmt.Errorf("object.Attrs: %v", err)
	}
	o = o.If(storage.Conditions{GenerationMatch: attrs.Generation})

	if err := o.Delete(ctx); err != nil {
		return fmt.Errorf("Object(%q).Delete: %v", object, err)
	}
	fmt.Fprintf(w, "Blob %v deleted.\n", object)
	return nil
}

func RemoveModule(bucket string, modulePath string) error {
	var buffer bytes.Buffer

	err := removeFile(&buffer, bucket, modulePath)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	return nil
}
