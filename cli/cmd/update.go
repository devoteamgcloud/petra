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
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// updateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// updateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func update() error {
	err := internal.InitGCSBackend(flagGCSBucket)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}

	err = internal.UpdateModule(flagModuleDirectory, flagGCSBucket, &flagConfig)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}
	return nil
}
