package modules

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"github.com/go-chi/chi/v5"
)

func GetDownloadURL(w http.ResponseWriter, r *http.Request) {
	mod := module{
		Namespace: chi.URLParam(r, "namespace"),
		Name:      chi.URLParam(r, "name"),
		Provider:  chi.URLParam(r, "provider"),
		Version:   chi.URLParam(r, "version"),
	}
	ctx := context.Background()
	downloadURL, err := getModule(mod, ctx)
	if err != nil {
		// TODO add error handler
		fmt.Fprintln(os.Stderr, err)
	}

	w.Header().Set("X-Terraform-Get", downloadURL)
	fmt.Println("DownloadURL", downloadURL)
	w.WriteHeader(http.StatusNoContent)
}

func getModule(mod module, ctx context.Context) (string, error) {
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

		signedUrl, err := backend.Client.Bucket(backend.Bucket).SignedURL(modPath(mod), options)
		if err != nil {
			return "", fmt.Errorf("Bucket(%q).SignedURL: %v", backend.Bucket, err)
		}

		fmt.Println("Generated GET signed URL:")
		fmt.Printf("%q\n", signedUrl)
		return fmt.Sprintln(signedUrl), nil
	}

	object := backend.Client.Bucket(backend.Bucket).Object(modPath(mod))
	attrs, err := object.Attrs(ctx)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("gcs::https://www.googleapis.com/storage/v1/%s/%s", backend.Bucket, attrs.Name), nil
}
