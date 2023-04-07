package routes

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

func SDRouter(r chi.Router) {
	r.Get("/.well-known/terraform.json", terraformServiceDiscovery)
}

func terraformServiceDiscovery(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write([]byte(fmt.Sprintf(`{"modules.v1": "%s/", "providers.v1": "%s/"}`, "/v1/modules", "/v1/providers")))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
