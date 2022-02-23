package module

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func getDownloadURL(w http.ResponseWriter, r *http.Request) {
	mod := Module{
		Namespace: chi.URLParam(r, "namespace"),
		Name:      chi.URLParam(r, "name"),
		Provider:  chi.URLParam(r, "provider"),
		Version:   chi.URLParam(r, "version"),
	}
	ctx := context.Background()
	downloadURL, err := gcsBucket.getModule(mod, ctx)
	if err != nil {
		// TODO add error handler
		fmt.Fprintln(os.Stderr, err)
	}

	w.Header().Set("X-Terraform-Get", downloadURL)
	w.WriteHeader(http.StatusNoContent)
}

func (b *GCSBackend) getModule(mod Module, ctx context.Context) (string, error) {
	object := b.client.Bucket(b.bucket).Object(modPath(mod))
	attrs, err := object.Attrs(ctx)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("gcs::https://www.googleapis.com/storage/v1/%s/%s", b.bucket, attrs.Name), nil
}
