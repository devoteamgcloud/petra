package cmd

import (
	"bytes"
	"fmt"
	"os"

	"github.com/arthur-laurentdka/petra/cli/internal"
	"github.com/spf13/cobra"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "A brief description of your command",
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

	// removeCmd.Flags().StringVar(&flagConfigFile, "config-file", "", "Path to your petra config file (required)")
	// removeCmd.MarkFlagRequired("config-file")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// removeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// removeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func remove() error {
	var buffer bytes.Buffer

	err := internal.DeleteFile(&buffer, flagGCSBucket, flagModuleDirectory)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}
	return nil
}
