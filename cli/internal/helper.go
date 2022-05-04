package internal

func GetObjectPathFromConfig(petraConf *PetraConfig) string {
	// {namespace}/{module}/{provider}/
	objectDirectory := petraConf.Namespace + "/" + petraConf.Name + "/" + petraConf.Provider + "/"
	// {namespace}-{module}-{provider}-{version}.tar.gz
	object := petraConf.Namespace + "-" + petraConf.Name + "-" + petraConf.Provider + "-" + petraConf.Version + ".tar.gz"

	// {namespace}/{module}/{provider}/{namespace}-{module}-{provider}-{version}.tar.gz
	// e.g.: main/rabbitmq/helm/0.0.1/main-rabbitmq-helm-0.0.1.tar.gz
	return objectDirectory + object
}
