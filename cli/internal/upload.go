package internal

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"cloud.google.com/go/storage"
)

var gcsBucket *GCSBackend

type Metadata struct {
	Owner string
	Team  string
}

type PetraConfig struct {
	Namespace string
	Name      string
	Provider  string
	Version   string
	Metadata  Metadata
}

func tar(moduleDirectory string) error {
	file, err := os.Create("module.zip")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	w := zip.NewWriter(file)
	defer w.Close()

	walker := func(path string, info os.FileInfo, err error) error {
		fmt.Printf("Crawling: %#v\n", path)
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		// Ensure that `path` is not absolute; it should not start with "/".
		// This snippet happens to work because I don't use
		// absolute paths, but ensure your real-world code
		// transforms path into a zip-root relative path.
		f, err := w.Create(path)
		if err != nil {
			return err
		}

		_, err = io.Copy(f, file)
		if err != nil {
			return err
		}

		return nil
	}
	err = filepath.Walk(moduleDirectory, walker)
	if err != nil {
		panic(err)
	}
	return nil
}

func uploadFile(w io.Writer, bucket string, zipFilePath string, petraConf *PetraConfig) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	f, err := os.Open(zipFilePath)
	if err != nil {
		return fmt.Errorf("os.Open: %v", err)
	}
	defer f.Close()

	object := getObjectPathFromConfig(petraConf)

	o := client.Bucket(bucket).Object(object)

	wc := o.NewWriter(ctx)

	// Update the object to set the metadata:
	// - owner
	// - team
	objectAttrs := map[string]string{
		"owner": petraConf.Metadata.Owner,
		"team":  petraConf.Metadata.Team,
	}
	wc.ObjectAttrs.Metadata = objectAttrs

	// Upload an object with storage.Writer.
	if _, err = io.Copy(wc, f); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}
	fmt.Fprintf(w, "Blob %v uploaded.\n", object)
	return nil
}

func UploadModule(bucket string, modulePath string) error {
	var buffer bytes.Buffer

	err := tar(modulePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}

	petraConf, err := getPetraConfig(modulePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}

	err = uploadFile(&buffer, bucket, "./module.zip", petraConf)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}
	return nil
}
