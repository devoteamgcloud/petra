package modules

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

type Version struct {
	Version string `json:"version"`
}

type listModuleVersions []struct {
	Source   string    `json:"source"`
	Versions []Version `json:"versions"`
}

func GetVersions(w http.ResponseWriter, r *http.Request) {
	mod := module{
		Namespace: chi.URLParam(r, "namespace"),
		Name:      chi.URLParam(r, "name"),
		Provider:  chi.URLParam(r, "provider"),
	}
	ctx := context.Background()
	path := modPathPartial(mod)
	retrievedVersions, err := backend.ModuleVersions(path, ctx)
	if err != nil {
		// TODO add error handler
		fmt.Fprintln(os.Stderr, err)
	}
	versions := []Version{}
	for _, v := range retrievedVersions {
		versions = append(versions, Version{v})
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(struct {
		Modules listModuleVersions `json:"modules"`
	}{
		listModuleVersions{
			{
				Source:   path,
				Versions: versions,
			},
		},
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
