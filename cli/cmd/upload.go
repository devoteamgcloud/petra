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
	Short: "Upload a terraform module to a private registry",
	Long: `Compress a local Terraform module and upload to a private registry.
			1. Each module must have a .petra-config.yaml
			2. Read the values of the config file
			3. Compress all module's files, generate a {namespace}-{module}-{version}.zip file and upload it to the private registry.
			4. Path of the object in the Google Cloud Storage bucket: {namespace}-{module}-{version}/{namespace}-{module}-{version}.zip
		`,
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
