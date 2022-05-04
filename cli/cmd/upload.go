package cmd

import (
	"bytes"
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// uploadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// uploadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func upload() error {
	var buffer bytes.Buffer

	err := internal.InitGCSBackend(flagGCSBucket)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}

	err = internal.Tar(flagModuleDirectory)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}

	petraConf, err := internal.GetPetraConfig(flagModuleDirectory)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}

	err = internal.UploadModule(&buffer, "./module.zip", petraConf)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}
	return nil
}
