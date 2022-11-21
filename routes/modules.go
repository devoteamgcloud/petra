package routes

import (
	"fmt"
	"path"

	"github.com/go-chi/chi/v5"
)

/* Type Definitions */

type Module struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Provider  string `json:"provider"`
	Version   string `json:"version"`
}

/* Utils */

func modPath(mod Module) string {
	return path.Join(
		mod.Namespace,
		mod.Name,
		mod.Provider,
		mod.Version,
		fmt.Sprintf("%s-%s-%s-%s.zip", mod.Namespace, mod.Name, mod.Provider, mod.Version),
	)
}

func modPathPartial(mod Module) string {
	return path.Join(
		mod.Namespace,
		mod.Name,
		mod.Provider,
	)
}

/* <Router> */

func ModulesRouter(r chi.Router) {
	r.Get("/{namespace}/{name}/{provider}/versions", getVersions)
	r.Get("/{namespace}/{name}/{provider}/{version}/download", getDownloadURL)
}

/* </Router> */
