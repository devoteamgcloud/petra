package cmd

import (
	"net/http"

	"github.com/arthur-laurentdka/petra/module"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "petra",
	Short: "private terraform registry",
	RunE: func(cmd *cobra.Command, args []string) error {
		server()
		return nil
	},
}

func Execute() error {
	return rootCmd.Execute()
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
	prefixProviders = "/v1/modules"
)

func server() {
	v := viper.New()
	v.SetEnvPrefix("PETRA")
	v.AutomaticEnv()

	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(middleware.GetHead)

	r.Use(middleware.Heartbeat("/is_alive"))

	r.Route(prefixModules, module.ModuleRouting)

	http.ListenAndServe(":"+flagListenAddr, r)
}
