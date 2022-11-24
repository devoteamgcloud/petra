package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/devoteamgcloud/petra/petractl/modules"
	"github.com/devoteamgcloud/petra/petractl/storage"
	"github.com/spf13/cobra"
)

// Declare Flags
var (
	flagBucket    string
	flagRecursive bool
)

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "package a terraform module and upload it to a backend bucket",
	Long: `CLI tool to package and push terraform modules to a Petra terraform module registry backend bucket
	`,
	RunE: func(cmd *cobra.Command, args []string) error {

		backend, err := storage.InitGCSBackend(flagBucket)
		if err != nil {
			return fmt.Errorf("failed to setup storage: %v", err)
		}

		if len(args) == 0 {
			return fmt.Errorf("missing argument")
		}

		if _, err := os.Stat(args[0]); errors.Is(err, os.ErrNotExist) {
			return err
		}

		err = modules.PackageModules(args[0], flagRecursive, backend)
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)

	pushCmd.PersistentFlags().StringVarP(&flagBucket, "bucket", "b", "", "")
	pushCmd.PersistentFlags().BoolVarP(&flagRecursive, "recursive", "r", false, "Also process files in subdirectories. By default, only the given directory (or current directory) is processed.")
}
