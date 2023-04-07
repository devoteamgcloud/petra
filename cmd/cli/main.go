package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/devoteamgcloud/petra/internal/modules"
	"github.com/devoteamgcloud/petra/internal/storage"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "petractl",
	Short: "package a terraform module and upload it to a backend bucket",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

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

	pushCmd.PersistentFlags().StringVarP(&flagBucket, "bucket", "b", "", "Name of the destination GCS bucket.")
	pushCmd.PersistentFlags().BoolVarP(&flagRecursive, "recursive", "r", false, "Also process files in subdirectories. By default, only the given directory (or current directory) is processed.")
}
