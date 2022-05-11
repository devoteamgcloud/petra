package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestExecuteCommand(t *testing.T) {
	rootCmd := &cobra.Command{
		Use:   "petra",
		Short: "Private terraform registry cli",
		Long:  "CLI to manage terraform modules in our private registry in a Google Cloud Storage bucket.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	rootCmd.Execute()
}