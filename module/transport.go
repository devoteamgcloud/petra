package module

import "github.com/go-chi/chi/v5"

func ModuleRouting(r chi.Router) {
	r.Get("/{namespace}/{name}/{provider}/versions", moduleGetVersions)
	r.Get("/{namespace}/{name}/{provider}/{version}/download", moduleGetDownloadURL)
}
