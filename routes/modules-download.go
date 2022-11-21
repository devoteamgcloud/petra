package routes

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/storage"
	petra_storage "github.com/devoteamgcloud/petra/storage"
	"github.com/go-chi/chi/v5"
)

func getDownloadURL(w http.ResponseWriter, r *http.Request) {
	mod := Module{
		Namespace: chi.URLParam(r, "namespace"),
		Name:      chi.URLParam(r, "name"),
		Provider:  chi.URLParam(r, "provider"),
		Version:   chi.URLParam(r, "version"),
	}
	fmt.Println("Enters the getDownloadURL function")
	ctx := context.Background()
	fmt.Println("gcs bucket : ", petra_storage.GCSBucket)
	downloadURL, err := getModule(mod, ctx)
	if err != nil {
		// TODO add error handler
		fmt.Fprintln(os.Stderr, err)
	}

	w.Header().Set("X-Terraform-Get", downloadURL)
	fmt.Println("DownloadURL", downloadURL)
	w.WriteHeader(http.StatusNoContent)
}

func getModule(mod Module, ctx context.Context) (string, error) {
	fmt.Println("mod :", modPath(mod))
	fmt.Println("context : ", ctx)

	enable_signed_url := os.Getenv("SIGNED_URL")
	if enable_signed_url == "true" {

		var options *storage.SignedURLOptions

		fmt.Println("Use existing credentials to create signed url")
		options = &storage.SignedURLOptions{
			Method:  "GET",
			Expires: time.Now().Add(2 * time.Minute),
		}

		signedUrl, err := petra_storage.GCSBucket.Client.Bucket(petra_storage.GCSBucket.Bucket).SignedURL(modPath(mod), options)
		if err != nil {
			return "", fmt.Errorf("Bucket(%q).SignedURL: %v", petra_storage.GCSBucket.Bucket, err)
		}

		fmt.Println("Generated GET signed URL:")
		fmt.Printf("%q\n", signedUrl)
		return fmt.Sprintln(signedUrl), nil
	}

	object := petra_storage.GCSBucket.Client.Bucket(petra_storage.GCSBucket.Bucket).Object(modPath(mod))
	attrs, err := object.Attrs(ctx)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("gcs::https://www.googleapis.com/storage/v1/%s/%s", petra_storage.GCSBucket.Bucket, attrs.Name), nil
}
