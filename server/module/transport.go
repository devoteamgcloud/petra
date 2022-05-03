package module

import "github.com/go-chi/chi/v5"

type Module struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Provider  string `json:"provider"`
	Version   string `json:"version"`
}

func Routing(r chi.Router) {
	r.Get("/{namespace}/{name}/{provider}/versions", getVersions)
	r.Get("/{namespace}/{name}/{provider}/{version}/download", getDownloadURL)
}
