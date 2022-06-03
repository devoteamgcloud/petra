package internal

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

const petraConfigFileName string = ".petra-config.yaml"

func editConfigFile(config *PetraConfig, modulePath string) error {
	data, err := yaml.Marshal(&config)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}

	fmt.Println(string(data))

	configFile := modulePath + petraConfigFileName

	err = ioutil.WriteFile(configFile, data, 0644)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}

	// f.Close()
	return nil
}

func GetPetraConfig(modulePath string) (*PetraConfig, error) {
	config := PetraConfig{}
	configPath := modulePath + petraConfigFileName

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

func GetObjectPathFromConfig(petraConf *PetraConfig) string {
	// {namespace}/{module}/{provider}/
	objectDirectory := petraConf.Namespace + "/" + petraConf.Name + "/" + petraConf.Provider + "/" + petraConf.Version + "/"
	// {namespace}-{module}-{provider}-{version}.zip
	object := petraConf.Namespace + "-" + petraConf.Name + "-" + petraConf.Provider + "-" + petraConf.Version + ".zip"

	// {namespace}/{module}/{provider}/{namespace}-{module}-{provider}-{version}.zip
	// e.g.: main/rabbitmq/helm/0.0.1/main-rabbitmq-helm-0.0.1.zip
	return objectDirectory + object
}
