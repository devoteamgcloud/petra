package internal

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"cloud.google.com/go/storage"
	"gopkg.in/yaml.v2"
)

type GCSBackend struct {
	client *storage.Client
	bucket string
}

func initGCSBackend(bckt string) error {
	ctx := context.Background()
	fmt.Println("bucket name :", bckt)
	client, err := storage.NewClient(ctx)
	fmt.Println("Client : ", client)
	if err != nil {
		return err
	}

	gcsBucket = &GCSBackend{
		client: client,
		bucket: bckt,
	}

	attrs, err := gcsBucket.client.Bucket(gcsBucket.bucket).Attrs(ctx)
	if err == storage.ErrBucketNotExist {
		fmt.Fprintln(os.Stderr, "The", gcsBucket.bucket, "bucket does not exist")
		return err
	}
	if err != nil {
		// Other error to handle
		fmt.Fprintln(os.Stderr, err)
	}
	fmt.Println("The", gcsBucket.bucket, "bucket exists and has attributes:", attrs)
	return err
}

func editConfigFile(config *PetraConfig, modulePath string) error {
	data, err := yaml.Marshal(&config)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}

	fmt.Println(string(data))

	configFile := modulePath + ".petra-config.yaml"

	err = ioutil.WriteFile(configFile, data, 0644)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}

	// f.Close()
	return nil
}

func getPetraConfig(modulePath string) (*PetraConfig, error) {
	config := PetraConfig{}
	configPath := modulePath + ".petra-config.yaml"

	fmt.Println(configPath)

	yamlFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return nil, fmt.Errorf("error: %v", err)
	}
	fmt.Printf("%+v\n", config)

	// check required fields
	if config.Namespace == "" {
		return nil, fmt.Errorf("error: required field (namespace) is missing in the config file")
	}
	if config.Name == "" {
		return nil, fmt.Errorf("error: required field (name) is missing in the config file")
	}
	if config.Provider == "" {
		return nil, fmt.Errorf("error: required field (provider) is missing in the config file")
	}
	if config.Version == "" {
		return nil, fmt.Errorf("error: required field (version) is missing in the config file")
	}

	return &config, nil
}

func getObjectPathFromConfig(petraConf *PetraConfig) string {
	// {namespace}/{module}/{provider}/
	objectDirectory := petraConf.Namespace + "/" + petraConf.Name + "/" + petraConf.Provider + "/"
	// {namespace}-{module}-{provider}-{version}.tar.gz
	object := petraConf.Namespace + "-" + petraConf.Name + "-" + petraConf.Provider + "-" + petraConf.Version + ".tar.gz"

	// {namespace}/{module}/{provider}/{namespace}-{module}-{provider}-{version}.tar.gz
	// e.g.: main/rabbitmq/helm/0.0.1/main-rabbitmq-helm-0.0.1.tar.gz
	return objectDirectory + object
}
