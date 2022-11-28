package routes

import (
	"github.com/devoteamgcloud/petra/modules"
	"github.com/go-chi/chi/v5"
)

/* Type Definitions */

type Module struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Provider  string `json:"provider"`
	Version   string `json:"version"`
}

/* <Router> */

func ModulesRouter(r chi.Router) {
	r.Get("/{namespace}/{name}/{provider}/versions", modules.GetVersions)
	r.Get("/{namespace}/{name}/{provider}/{version}/download", modules.GetDownloadURL)
}

/* </Router> */
