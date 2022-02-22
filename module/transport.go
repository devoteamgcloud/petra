package module

import "github.com/go-chi/chi/v5"

type Module struct {
	Namespace string
	Name      string
	Provider  string
	Version   string
}

func Routing(r chi.Router) {
	// r.Get("/{namespace}/{name}/{provider}/versions", getVersions)
	r.Get("/{namespace}/{name}/{provider}/{version}/download", getDownloadURL)
}
