package module

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"cloud.google.com/go/storage"
	"github.com/go-chi/chi/v5"
	"google.golang.org/api/iterator"
)

type Version struct {
	Version string `json:"version"`
}

type listModuleVersions []struct {
	Source   string    `json:"source"`
	Versions []Version `json:"versions"`
}

func getVersions(w http.ResponseWriter, r *http.Request) {
	mod := Module{
		Namespace: chi.URLParam(r, "namespace"),
		Name:      chi.URLParam(r, "name"),
		Provider:  chi.URLParam(r, "provider"),
	}
	ctx := context.Background()
	fmt.Println("Enters the getVersions function")
	versions, err := gcsBucket.listModuleVersions(mod, ctx)
	if err != nil {
		// TODO add error handler
		fmt.Fprintln(os.Stderr, err)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(struct {
		Modules listModuleVersions `json:"modules"`
	}{
		listModuleVersions{
			{
				Source:   modPathPartial(mod),
				Versions: versions,
			},
		},
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func (b *GCSBackend) listModuleVersions(mod Module, ctx context.Context) ([]Version, error) {
	var versions []Version
	path := modPathPartial(mod)
	query := &storage.Query{
		Prefix: path,
	}
	it := b.client.Bucket(b.bucket).Objects(ctx, query)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return versions, err
		}
		if attrs.ContentType != "text/plain" {
			continue
		}
		version := strings.Split(attrs.Name, "/")[3]
		if version == "" {
			continue
		}
		vs := Version{
			Version: version,
		}
		versions = append(versions, vs)
	}
	return versions, nil
}
