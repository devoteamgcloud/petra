package internal

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
)

func updateObjectMetadata(w io.Writer, bucket string, object string, newMetadata Metadata) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	o := client.Bucket(bucket).Object(object)

	// Update the object to set the metadata.
	objectAttrsToUpdate := storage.ObjectAttrsToUpdate{
		Metadata: map[string]string{
			"owner": newMetadata.Owner,
			"team":  newMetadata.Team,
		},
	}
	if _, err := o.Update(ctx, objectAttrsToUpdate); err != nil {
		return fmt.Errorf("ObjectHandle(%q).Update: %v", object, err)
	}
	fmt.Fprintf(w, "Updated custom metadata for object %v in bucket %v.\n", object, bucket)
	return nil
}

// moveFile moves an object into another location.
func moveFile(w io.Writer, bucket, object string, destinationObject string) error {
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

func updateMetadata(currentConfig *PetraConfig, flagConfig *PetraConfig) {
	if flagConfig.Metadata.Owner == "" {
		flagConfig.Metadata.Owner = currentConfig.Metadata.Owner
	}
	if flagConfig.Metadata.Team == "" {
		flagConfig.Metadata.Team = currentConfig.Metadata.Team
	}
}

func updateRequiredFields(currentConfig *PetraConfig, flagConfig *PetraConfig) {
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

func UpdateModule(bucket string, moduleDirectory string, flagConfig *PetraConfig) error {
	// 1. Get info from petra config file
	currentConf, err := getPetraConfig(moduleDirectory)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}

	// e.g.: main/rabbitmq/helm/0.0.1/main-rabbitmq-helm-0.0.1.tar.gz
	currentObject := getObjectPathFromConfig(currentConf)

	// 2.1 Change object's metadata:
	// - owner | team
	if flagConfig.Metadata.Owner != "" || flagConfig.Metadata.Team != "" {
		var buffer bytes.Buffer

		updateMetadata(currentConf, flagConfig)
		err = updateObjectMetadata(&buffer, bucket, currentObject, flagConfig.Metadata)
		if err != nil {
			return fmt.Errorf("error: %v", err)
		}
	}
	// 2.2 Move object if we want to update one of the following field:
	// - namespace | name | provider | version
	// because they're part of the object's path:
	// {namespace}/{module}/{provider}/{namespace}-{module}-{provider}-{version}.tar.gz
	if flagConfig.Namespace != "" || flagConfig.Name != "" || flagConfig.Provider != "" || flagConfig.Version != "" {
		var buffer bytes.Buffer

		updateRequiredFields(currentConf, flagConfig)

		destinationObject := getObjectPathFromConfig(flagConfig)
		fmt.Printf("destination object: %s\n", destinationObject)

		err = moveFile(&buffer, bucket, currentObject, destinationObject)
		if err != nil {
			return fmt.Errorf("error: %v", err)
		}
	}

	// 3. change petra config file
	err = editConfigFile(flagConfig, moduleDirectory)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}

	return nil
}
