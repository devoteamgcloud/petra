package module

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"github.com/go-chi/chi/v5"
	"golang.org/x/oauth2/google"
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
	fmt.Println("gcs bucket : ", gcsBucket)
	downloadURL, err := gcsBucket.getModule(mod, ctx)
	if err != nil {
		// TODO add error handler
		fmt.Fprintln(os.Stderr, err)
	}

	w.Header().Set("X-Terraform-Get", downloadURL)
	fmt.Println("DownloadURL", downloadURL)
	w.WriteHeader(http.StatusNoContent)
}

func (b *GCSBackend) getModule(mod Module, ctx context.Context) (string, error) {
	fmt.Println("mod :", modPath(mod))
	fmt.Println("context : ", ctx)

	sakeyFile := "./storage-sa.json"

	saKey, err := ioutil.ReadFile(sakeyFile)
	if err != nil {
		log.Fatalln(err)
	}

	cfg, err := google.JWTConfigFromJSON(saKey)
	if err != nil {
		log.Fatalln(err)
	}

	opts := &storage.SignedURLOptions{
		GoogleAccessID: cfg.Email,
		PrivateKey:     cfg.PrivateKey,
		Method:         "GET",
		Expires:        time.Now().Add(2 * time.Minute),
	}

	signedUrl, err := b.client.Bucket(b.bucket).SignedURL(modPath(mod), opts)
	if err != nil {
		return "", err
	}
	return fmt.Sprintln(signedUrl), nil
}
