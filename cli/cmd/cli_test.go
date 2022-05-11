package cmd

import (
	"bytes"
	"fmt"
	"testing"
	// "os"

	"github.com/spf13/cobra"
)

// Utils
func executeCommand(root *cobra.Command, args ...string) (output string, err error) {
	_, output, err = executeCommandC(root, args...)
	return output, err
}

func executeCommandC(root *cobra.Command, args ...string) (c *cobra.Command, output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)

	c, err = root.ExecuteC()

	return c, buf.String(), err
}

// Tests
func TestRootExecuteUnknownCommand(t *testing.T) {
	fmt.Println("=================================")
	fmt.Println("test: TestRootExecuteUnknownCommand")
	output, err := executeCommand(rootCmd, "unknown")

	fmt.Println(err)
	expected := "Error: unknown command \"unknown\" for \"petra\"\nRun 'petra --help' for usage.\n"

	if output != expected {
		t.Errorf("\nExpected:\n %q\nGot:\n %q\n", expected, output)
	}
	fmt.Println("=================================")
	fmt.Printf("\n")
}

func TestUploadSubCmdNoFlag(t *testing.T) {
	fmt.Println("=================================")
	fmt.Println("test: TestUploadSubCmdNoFlag")
	_, err := executeCommand(rootCmd, "upload")

	fmt.Println(err)
	expected := "required flag(s) \"gcs-bucket\", \"module-directory\" not set"

	if err.Error() != expected {
		t.Errorf("\nExpected:\n %q\nGot:\n %q\n", expected, err.Error())
	}
	fmt.Println("=================================")
	fmt.Printf("\n")
}

func TestRemoveSubCmdNoFlag(t *testing.T) {
	fmt.Println("=================================")
	fmt.Println("test: TestRemoveSubCmdNoFlag")
	_, err := executeCommand(rootCmd, "remove")

	fmt.Println(err)
	expected := "required flag(s) \"gcs-bucket\", \"module-directory\" not set"

	if err.Error() != expected {
		t.Errorf("\nExpected:\n %q\nGot:\n %q\n", expected, err.Error())
	}
	fmt.Println("=================================")
	fmt.Printf("\n")
}

func TestUpdateSubCmdNoFlag(t *testing.T) {
	fmt.Println("=================================")
	fmt.Println("test: TestUpdateSubCmdNoFlag")
	_, err := executeCommand(rootCmd, "update")

	fmt.Println(err)
	expected := "required flag(s) \"gcs-bucket\", \"module-directory\" not set"

	if err.Error() != expected {
		t.Errorf("\nExpected:\n %q\nGot:\n %q\n", expected, err.Error())
	}
	fmt.Println("=================================")
	fmt.Printf("\n")
}
