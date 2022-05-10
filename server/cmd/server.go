package cmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/arthur-laurentdka/petra/server/module"
	"github.com/arthur-laurentdka/petra/server/provider"
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
	flagProjetID   string
	flagSecretID   string
)

func init() {
	rootCmd.PersistentFlags().StringVar(&flagGCSBucket, "gcs-bucket", "", "Name of the Google Cloud Storage bucket you want to use for storage.")
	rootCmd.PersistentFlags().StringVar(&flagProjetID, "project-id", "", "Google Cloud project ID where the service account is stored in Secret Manager.")
	rootCmd.PersistentFlags().StringVar(&flagSecretID, "secret-id", "", "(Google Cloud Secret Manager) Secret ID of your service-account that allows you to generate signed URLs.")

	rootCmd.PersistentFlags().StringVar(&flagListenAddr, "listen-address", "3000", "Address to listen on")

	// Making Flags required
	err := rootCmd.MarkPersistentFlagRequired("project-id")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	err = rootCmd.MarkPersistentFlagRequired("secret-id")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
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

	module.InitSecretManagerInfo(flagProjetID, flagSecretID)

	r.Use(middleware.Heartbeat("/is_alive"))

	// Basic Service Discovery Handler
	r.Get("/.well-known/terraform.json", getSD)

	// Modules Handler
	r.Route(prefixModules, module.Routing)

	// Providers Handler
	r.Route(prefixProviders, provider.Routing)

	if err := http.ListenAndServe(":"+flagListenAddr, r); err != nil {
		if err != http.ErrServerClosed {
			return err
		}
	}

	return nil
}
