package cmd

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"testing"
	"time"

	"cloud.google.com/go/storage"
	"github.com/arthur-laurentdka/petra/cli/internal"
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

func getMetadata(w io.Writer, bucket, object string) (*storage.ObjectAttrs, error) {
	// bucket := "bucket-name"
	// object := "object-name"
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	o := client.Bucket(bucket).Object(object)
	attrs, err := o.Attrs(ctx)
	if err != nil {
		return nil, fmt.Errorf("Object(%q).Attrs: %v", object, err)
	}
	fmt.Fprintf(w, "Bucket: %v\n", attrs.Bucket)
	fmt.Fprintf(w, "CacheControl: %v\n", attrs.CacheControl)
	fmt.Fprintf(w, "ContentDisposition: %v\n", attrs.ContentDisposition)
	fmt.Fprintf(w, "ContentEncoding: %v\n", attrs.ContentEncoding)
	fmt.Fprintf(w, "ContentLanguage: %v\n", attrs.ContentLanguage)
	fmt.Fprintf(w, "ContentType: %v\n", attrs.ContentType)
	fmt.Fprintf(w, "Crc32c: %v\n", attrs.CRC32C)
	fmt.Fprintf(w, "Generation: %v\n", attrs.Generation)
	fmt.Fprintf(w, "KmsKeyName: %v\n", attrs.KMSKeyName)
	fmt.Fprintf(w, "Md5Hash: %v\n", attrs.MD5)
	fmt.Fprintf(w, "MediaLink: %v\n", attrs.MediaLink)
	fmt.Fprintf(w, "Metageneration: %v\n", attrs.Metageneration)
	fmt.Fprintf(w, "Name: %v\n", attrs.Name)
	fmt.Fprintf(w, "Size: %v\n", attrs.Size)
	fmt.Fprintf(w, "StorageClass: %v\n", attrs.StorageClass)
	fmt.Fprintf(w, "TimeCreated: %v\n", attrs.Created)
	fmt.Fprintf(w, "Updated: %v\n", attrs.Updated)
	fmt.Fprintf(w, "Event-based hold enabled? %t\n", attrs.EventBasedHold)
	fmt.Fprintf(w, "Temporary hold enabled? %t\n", attrs.TemporaryHold)
	fmt.Fprintf(w, "Retention expiration time %v\n", attrs.RetentionExpirationTime)
	fmt.Fprintf(w, "Custom time %v\n", attrs.CustomTime)
	fmt.Fprintf(w, "\n\nMetadata\n")
	for key, value := range attrs.Metadata {
		fmt.Fprintf(w, "\t%v = %v\n", key, value)
	}
	return attrs, nil
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

// Upload in: staging/rabbitmq/helm/0.0.2/staging-rabbitmq-helm-0.0.2.tar.gz
func TestUploadModule(t *testing.T) {
	fmt.Println("=================================")
	fmt.Println("test: TestUploadModule")

	// Exec: petra upload --gcs-bucket=toltol-private-registry --module-directory=../modules-example/rabbitmq
	bucket := "toltol-private-registry"
	moduleDir := "../modules-example/rabbitmq/"
	_, err := executeCommand(rootCmd, "upload", "--gcs-bucket="+bucket, "--module-directory="+moduleDir)
	if err != nil {
		t.Errorf("error: %v", err)
	}

	// Get petra config in ../modules-example/rabbitmq/.petra-config.yaml
	conf, err := internal.GetPetraConfig(moduleDir)
	if err != nil {
		t.Errorf("error: %v", err)
	}
	fmt.Println(conf)

	// Get object path in Google Cloud Storage bucket from config
	// e.g.: main/rabbitmq/helm/0.0.1/main-rabbitmq-helm-0.0.1.tar.gz
	objectPath := internal.GetObjectPathFromConfig(conf)
	fmt.Println(objectPath)

	// Get object attributes from Google Cloud Storage bucket
	var buffer bytes.Buffer
	objectAttrs, err := getMetadata(&buffer, bucket, objectPath)
	if err != nil {
		t.Errorf("error: %v", err)
	}
	fmt.Println(objectAttrs)

	// Check attributes:
	// name, metadata.owner, metadata.team
	if objectPath != objectAttrs.Name || conf.Metadata.Owner != objectAttrs.Metadata["owner"] || conf.Metadata.Team != objectAttrs.Metadata["team"] {
		t.Errorf("error: attributes are not the same: %v", err)
	}

	fmt.Println("=================================")
	fmt.Printf("\n")
}

// Move
// from: staging/rabbitmq/helm/0.0.2/staging-rabbitmq-helm-0.0.2.tar.gz
// to: main/anotherModule/kubernetes/0.0.1/main-anotherModule-kubernetes-0.0.1.tar.gz
func TestUpdateModule(t *testing.T) {
	fmt.Println("=================================")
	fmt.Println("test: TestUpdateModule")

	bucket := "toltol-private-registry"
	moduleDir := "../modules-example/rabbitmq/"

	// Get previous petra config in ../modules-example/rabbitmq/.petra-config.yaml
	prevConf, err := internal.GetPetraConfig(moduleDir)
	if err != nil {
		t.Errorf("error: %v", err)
	}
	fmt.Printf("prevConf: %v\n", prevConf)

	name := "anotherModule"
	version := "0.0.1"
	provider := "kubernetes"
	namespace := "main"
	owner := "toto"
	team := "prod"
	// Exec: petra upload --gcs-bucket=toltol-private-registry --module-directory=../modules-example/rabbitmq
	_, err = executeCommand(rootCmd, "update", "--gcs-bucket="+bucket, "--module-directory="+moduleDir, "--name="+name, "--version="+version, "--provider="+provider, "--namespace="+namespace, "--owner="+owner, "--team="+team)
	if err != nil {
		t.Errorf("error: %v", err)
	}

	// Get new petra config in ../modules-example/rabbitmq/.petra-config.yaml
	newConf, err := internal.GetPetraConfig(moduleDir)
	if err != nil {
		t.Errorf("error: %v", err)
	}
	fmt.Printf("nextConf: %v\n", newConf)

	// 1. Check that object was moved to the new location
	// Get object path in Google Cloud Storage bucket from config
	// e.g.: main/rabbitmq/helm/0.0.1/main-rabbitmq-helm-0.0.1.tar.gz
	objectPath := internal.GetObjectPathFromConfig(newConf)
	fmt.Println(objectPath)

	// Get object attributes from Google Cloud Storage bucket
	var buffer bytes.Buffer
	objectAttrs, err := getMetadata(&buffer, bucket, objectPath)

	if err != nil {
		t.Errorf("error: %v", err)
	}
	fmt.Println(objectAttrs)

	// 2. Check changes in petra config file
	if newConf.Name != name || newConf.Namespace != namespace || newConf.Version != version || newConf.Provider != provider || newConf.Metadata.Owner != owner || newConf.Metadata.Team != team {
		t.Errorf("\nerror: the new config file had not the latest changes.\n")
	}

	fmt.Println("=================================")
	fmt.Printf("\n")
}

// Move
// from: main/anotherModule/kubernetes/0.0.1/main-anotherModule-kubernetes-0.0.1.tar.gz
// to: staging/rabbitmq/helm/0.0.2/staging-rabbitmq-helm-0.0.2.tar.gz
func TestUpdateModule_2(t *testing.T) {
	fmt.Println("=================================")
	fmt.Println("test: TestUpdateModule")

	bucket := "toltol-private-registry"
	moduleDir := "../modules-example/rabbitmq/"

	// Get previous petra config in ../modules-example/rabbitmq/.petra-config.yaml
	prevConf, err := internal.GetPetraConfig(moduleDir)
	if err != nil {
		t.Errorf("error: %v", err)
	}
	fmt.Printf("prevConf: %v\n", prevConf)

	name := "rabbitmq"
	version := "0.0.2"
	provider := "helm"
	namespace := "staging"
	owner := "Hyokil"
	team := "GCP"
	// Exec: petra upload --gcs-bucket=toltol-private-registry --module-directory=../modules-example/rabbitmq
	_, err = executeCommand(rootCmd, "update", "--gcs-bucket="+bucket, "--module-directory="+moduleDir, "--name="+name, "--version="+version, "--provider="+provider, "--namespace="+namespace, "--owner="+owner, "--team="+team)
	if err != nil {
		t.Errorf("error: %v", err)
	}

	// Get new petra config in ../modules-example/rabbitmq/.petra-config.yaml
	newConf, err := internal.GetPetraConfig(moduleDir)
	if err != nil {
		t.Errorf("error: %v", err)
	}
	fmt.Printf("nextConf: %v\n", newConf)

	// 1. Check that object was moved to the new location
	// Get object path in Google Cloud Storage bucket from config
	// e.g.: main/rabbitmq/helm/0.0.1/main-rabbitmq-helm-0.0.1.tar.gz
	objectPath := internal.GetObjectPathFromConfig(newConf)
	fmt.Println(objectPath)

	// Get object attributes from Google Cloud Storage bucket
	var buffer bytes.Buffer
	objectAttrs, err := getMetadata(&buffer, bucket, objectPath)

	if err != nil {
		t.Errorf("error: %v", err)
	}
	fmt.Println(objectAttrs)

	// 2. Check changes in petra config file
	if newConf.Name != name || newConf.Namespace != namespace || newConf.Version != version || newConf.Provider != provider || newConf.Metadata.Owner != owner || newConf.Metadata.Team != team {
		t.Errorf("\nerror: the new config file had not the latest changes.\n")
	}

	fmt.Println("=================================")
	fmt.Printf("\n")
}

// Remove: staging/rabbitmq/helm/0.0.2/staging-rabbitmq-helm-0.0.2.tar.gz
func TestRemoveModule(t *testing.T) {
	fmt.Println("=================================")
	fmt.Println("test: TestRemoveModule")

	// Exec: petra upload --gcs-bucket=toltol-private-registry --module-directory=../modules-example/rabbitmq
	bucket := "toltol-private-registry"
	moduleDir := "../modules-example/rabbitmq/"
	_, err := executeCommand(rootCmd, "remove", "--gcs-bucket="+bucket, "--module-directory="+moduleDir)
	if err != nil {
		t.Errorf("error: %v", err)
	}

	// Get petra config in ../modules-example/rabbitmq/.petra-config.yaml
	conf, err := internal.GetPetraConfig(moduleDir)
	if err != nil {
		t.Errorf("error: %v", err)
	}
	fmt.Println(conf)

	// Get object path in Google Cloud Storage bucket from config
	// e.g.: main/rabbitmq/helm/0.0.1/main-rabbitmq-helm-0.0.1.tar.gz
	objectPath := internal.GetObjectPathFromConfig(conf)
	fmt.Println(objectPath)

	// Get object attributes from Google Cloud Storage bucket
	var buffer bytes.Buffer
	_, err = getMetadata(&buffer, bucket, objectPath)

	expected := `Object("` + objectPath + `").Attrs: storage: object doesn't exist`

	// Must return a error: the object doesn't exist because we removed it
	if err == nil || err.Error() != expected {
		t.Errorf("\nExpected:\n %q\nGot:\n %q\n", expected, err.Error())
	}

	fmt.Println("=================================")
	fmt.Printf("\n")
}
