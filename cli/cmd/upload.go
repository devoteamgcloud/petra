package cmd

import (
	"fmt"
	"os"

	"github.com/arthur-laurentdka/petra/cli/internal"
	"github.com/spf13/cobra"
)

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := upload(); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)
}

func upload() error {
	err := internal.UploadModule(flagGCSBucket, flagModuleDirectory)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}
	return nil
}
