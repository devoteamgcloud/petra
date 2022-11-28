package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/devoteamgcloud/petra/routes"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var gcs_bucket string

func Init() error {
	gcs_bucket = os.Getenv("GCS_BUCKET")
	if gcs_bucket == "" {
		return fmt.Errorf("GCS_BUCKET env var must be set")
	}
	return nil
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
	r.Use(middleware.Heartbeat("/healthz"))

	r.Route("/v1/modules", routes.ModulesRouter)
	r.Route("/", routes.SDRouter)

	if err := http.ListenAndServe(":8080", r); err != nil {
		if err != http.ErrServerClosed {
			return err
		}
	}

	return nil
}
