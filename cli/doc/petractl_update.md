## petractl update

Update one or multiple config values of a module.

### Synopsis

Update one or multiple config values of a module.
			1. Get config values passed as arguments to petra.
			2. Make changes to {namespace}-{module}-{version}/{namespace}-{module}-{version}.zip in the Google Cloud Storage bucket.
			3. Make changes to the local .petra-config.yaml of the local module
	

```
petractl update [flags]
```

### Options

```
  -h, --help               help for update
      --name string        Update module's name
      --namespace string   Update module's namespace
      --owner string       Update module's owner
      --provider string    Update module's provider
      --team string        Update module's team
      --version string     Update module's version
```

### Options inherited from parent commands

```
      --gcs-bucket string         Name of the Google Cloud Storage bucket you want to use for storage (required) e.g.: my-bucket
      --module-directory string   Directory of the module you want to upload (required) e.g.: ./modules-example/rabbitmq/
```

### SEE ALSO

* [petractl](petractl.md)	 - Private terraform registry cli

