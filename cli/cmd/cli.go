package cmd

import (
	"fmt"
	"os"

	"github.com/arthur-laurentdka/petra/cli/internal"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "petra",
	Short: "private terraform registry",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := cli(); err != nil {
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
	flagGCSBucket       string
	flagModuleDirectory string
)

func init() {
	rootCmd.PersistentFlags().StringVar(&flagGCSBucket, "gcs-bucket", "", "Name of the Google Cloud Storage bucket you want to use for storage")
	rootCmd.PersistentFlags().StringVar(&flagModuleDirectory, "module-directory", "", "Directory of the module you want to upload")
}

func cli() error {
	err := module.InitGCSBackend(flagGCSBucket)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}

	err = module.Tar(flagModuleDirectory)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}

	petraConf, err := module.GetPetraConfig(flagModuleDirectory)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}

	err = module.UploadModule("./module.zip", petraConf)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}
	return nil
}
