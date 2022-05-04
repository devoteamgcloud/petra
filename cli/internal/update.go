package internal

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
)

// func UpdateObject() error {

// }

// moveFile moves an object into another location.
func moveFile(w io.Writer, bucket, object string, destinationObject string) error {
	// bucket := "bucket-name"
	// object := "object-name"
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	src := client.Bucket(bucket).Object(object)
	dst := client.Bucket(bucket).Object(destinationObject)

	// Optional: set a generation-match precondition to avoid potential race
	// conditions and data corruptions. The request to upload is aborted if the
	// object's generation number does not match your precondition.
	// For a dst object that does not yet exist, set the DoesNotExist precondition.
	dst = dst.If(storage.Conditions{DoesNotExist: true})
	// If the destination object already exists in your bucket, set instead a
	// generation-match precondition using its generation number.
	// attrs, err := dst.Attrs(ctx)
	// if err != nil {
	//      return fmt.Errorf("object.Attrs: %v", err)
	// }
	// dst = dst.If(storage.Conditions{GenerationMatch: attrs.Generation})

	if _, err := dst.CopierFrom(src).Run(ctx); err != nil {
		return fmt.Errorf("Object(%q).CopierFrom(%q).Run: %v", destinationObject, object, err)
	}
	if err := src.Delete(ctx); err != nil {
		return fmt.Errorf("Object(%q).Delete: %v", object, err)
	}
	fmt.Fprintf(w, "Blob %v moved to %v.\n", object, destinationObject)
	return nil
}

func update(currentConfig *PetraConfig, flagConfig *PetraConfig) {
	if flagConfig.Namespace == "" {
		flagConfig.Namespace = currentConfig.Namespace
	}
	if flagConfig.Name == "" {
		flagConfig.Name = currentConfig.Name
	}
	if flagConfig.Provider == "" {
		flagConfig.Provider = currentConfig.Provider
	}
	if flagConfig.Version == "" {
		flagConfig.Version = currentConfig.Version
	}
}

// 1. Read petra config file
// 2. Read object
// 3. Update object
// 4. Update field in petra config file
func UpdateModule(modulePath string, bucket string, flagConfig *PetraConfig) error {
	// 1. Get info from petra config file
	currentConf, err := GetPetraConfig(modulePath)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}

	// 2.1 Move object if we want to update one of the following field:
	// - namespace | name | provider | version
	// because they're part of the object's path:
	// {namespace}/{module}/{provider}/{namespace}-{module}-{provider}-{version}.tar.gz
	if flagConfig.Namespace != "" || flagConfig.Name != "" || flagConfig.Provider != "" || flagConfig.Version != "" {
		var buffer bytes.Buffer

		currentObject := GetObjectPathFromConfig(currentConf)
		update(currentConf, flagConfig)
		destinationObject := GetObjectPathFromConfig(flagConfig)

		fmt.Println(currentObject)
		fmt.Println(destinationObject)
		err = moveFile(&buffer, bucket, currentObject, destinationObject)
		if err != nil {
			return fmt.Errorf("error: %v", err)
		}
	}

	// 2.2 Otherwise we change the metadata of the object for:
	// - owner | team
	return nil
}
