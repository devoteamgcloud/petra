package cmd

import (
	"fmt"
	"os"

	"github.com/arthur-laurentdka/petra/cli/internal"
	"github.com/spf13/cobra"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove the module from a private registry",
	Long: `Remove the module from a private registry.\n
			1. Read the value from .petra-config.yaml of the local module\n
			2. Remove the {namespace}-{module}-{version}/{namespace}-{module}-{version}-tar.gz from the Google Cloud Storage bucket`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := remove(); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}

func remove() error {
	err := internal.RemoveModule(flagGCSBucket, flagModuleDirectory)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}
	return nil
}
