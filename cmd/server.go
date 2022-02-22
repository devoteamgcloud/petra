package cmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/arthur-laurentdka/petra/module"
	"github.com/arthur-laurentdka/petra/provider"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "petra",
	Short: "private terraform registry",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := server(); err != nil {
			return err
		}
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// Declare Flags.
var (
	flagGCSBucket  string
	flagListenAddr string
)

func init() {
	rootCmd.PersistentFlags().StringVar(&flagGCSBucket, "gcs-bucket", "", "Name of the Google Cloud Storage bucket you want to use for storage")
	rootCmd.PersistentFlags().StringVar(&flagListenAddr, "listen-address", "3000", "Address to listen on")
}

const (
	prefixModules   = "/v1/modules"
	prefixProviders = "/v1/providers"
)

func server() error {
	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(middleware.GetHead)

	err := module.InitGCSBackend(flagGCSBucket)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}

	r.Use(middleware.Heartbeat("/is_alive"))

	r.Route(prefixModules, module.Routing)
	r.Route(prefixProviders, provider.Routing)

	if err := http.ListenAndServe(":"+flagListenAddr, r); err != nil {
		if err != http.ErrServerClosed {
			return err
		}
	}

	return nil
}
