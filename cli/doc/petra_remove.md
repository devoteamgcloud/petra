## petra remove

Remove the module from a private registry

### Synopsis

Remove the module from a private registry.
			1. Read the value from .petra-config.yaml of the local module
			2. Remove the {namespace}-{module}-{version}/{namespace}-{module}-{version}-tar.gz from the Google Cloud Storage bucket.

```
petractl remove [flags]
```

### Options

```
  -h, --help   help for remove
```

### Options inherited from parent commands

```
      --gcs-bucket string         Name of the Google Cloud Storage bucket you want to use for storage (required) e.g.: my-bucket
      --module-directory string   Directory of the module you want to upload (required) e.g.: ./modules-example/rabbitmq/
```

### SEE ALSO

* [petractl](petra.md)	 - Private terraform registry cli

