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
	Short: "Remove the .tar.gz file of a Terraform module in the bucket",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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
