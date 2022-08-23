package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/arthur-laurentdka/petra/server/module"
	"github.com/arthur-laurentdka/petra/server/provider"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Config struct {
	bucket  string
	address string
}

var config Config

func Init() error {
	bucket := os.Getenv("GCS_BUCKET")
	if bucket == "" {
		return fmt.Errorf("GCS_BUCKET env var must be set")
	}
	address := os.Getenv("LISTEN_ADDRESS")
	if address == "" {
		return fmt.Errorf("LISTEN_ADDRESS env var must be set")
	}

	config = Config{
		bucket:  bucket,
		address: address,
	}
	return nil
}

const (
	prefixModules   = "/v1/modules"
	prefixProviders = "/v1/providers"
)

func getSD(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write([]byte(fmt.Sprintf(`{"modules.v1": "%s/", "providers.v1": "%s/"}`, prefixModules, prefixProviders)))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func Run() error {
	err := Init()
	if err != nil {
		return err
	}
	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(middleware.GetHead)

	err = module.InitGCSBackend(config.bucket)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}

	r.Use(middleware.Heartbeat("/is_alive"))

	// Basic Service Discovery Handler
	r.Get("/.well-known/terraform.json", getSD)

	// Modules Handler
	r.Route(prefixModules, module.Routing)

	// Providers Handler
	r.Route(prefixProviders, provider.Routing)

	if err := http.ListenAndServe(":"+config.address, r); err != nil {
		if err != http.ErrServerClosed {
			return err
		}
	}

	return nil
}
