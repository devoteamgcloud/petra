## petra upload

Upload a terraform module to a private registry

### Synopsis

Compress a local Terraform module and upload to a private registry.
			1. Each module must have a .petra-config.yaml
			2. Read the values of the config file
			3. Compress all module's files, generate a {namespace}-{module}-{version}-tar.gz file and upload it to the private registry.
			4. Path of the object in the Google Cloud Storage bucket: {namespace}-{module}-{version}/{namespace}-{module}-{version}-tar.gz
		

```
petra upload [flags]
```

### Options

```
  -h, --help   help for upload
```

### Options inherited from parent commands

```
      --gcs-bucket string         Name of the Google Cloud Storage bucket you want to use for storage (required) e.g.: my-bucket
      --module-directory string   Directory of the module you want to upload (required) e.g.: ./modules-example/rabbitmq/
```

### SEE ALSO

* [petra](petra.md)	 - Private terraform registry cli

