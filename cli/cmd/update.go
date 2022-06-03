package cmd

import (
	"fmt"
	"os"

	"github.com/arthur-laurentdka/petra/cli/internal"
	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update one or multiple config values of a module.",
	Long: `Update one or multiple config values of a module.
			1. Get config values passed as arguments to petra.
			2. Make changes to {namespace}-{module}-{version}/{namespace}-{module}-{version}.zip in the Google Cloud Storage bucket.
			3. Make changes to the local .petra-config.yaml of the local module
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := update(); err != nil {
			return err
		}
		return nil
	},
}

var (
	flagConfig internal.PetraConfig
)

func init() {
	rootCmd.AddCommand(updateCmd)

	updateCmd.Flags().StringVar(&flagConfig.Namespace, "namespace", "", "Update module's namespace")
	updateCmd.Flags().StringVar(&flagConfig.Name, "name", "", "Update module's name")
	updateCmd.Flags().StringVar(&flagConfig.Provider, "provider", "", "Update module's provider")
	updateCmd.Flags().StringVar(&flagConfig.Version, "version", "", "Update module's version")
	updateCmd.Flags().StringVar(&flagConfig.Metadata.Owner, "owner", "", "Update module's owner")
	updateCmd.Flags().StringVar(&flagConfig.Metadata.Team, "team", "", "Update module's team")
}

func update() error {
	err := internal.UpdateModule(flagGCSBucket, flagModuleDirectory, flagConfig)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}
	return nil
}
