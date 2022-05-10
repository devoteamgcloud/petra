## petracli update

Update one or multiple settings of a module and make changes in the .petra-config.yaml

```
petracli update [flags]
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

* [petracli](petracli.md)	 - private terraform registry cli

