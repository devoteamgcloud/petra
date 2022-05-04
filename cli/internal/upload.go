package internal

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
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

func Tar(moduleDirectory string) error {
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

func UploadModule(w io.Writer, zipFilePath string, petraConf *PetraConfig) error {
	ctx := context.Background()
	// Open local file.
	f, err := os.Open(zipFilePath)
	if err != nil {
		return fmt.Errorf("os.Open: %v", err)
	}
	defer f.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	object := getObjectPathFromConfig(petraConf)

	o := gcsBucket.client.Bucket(gcsBucket.bucket).Object(object)

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
