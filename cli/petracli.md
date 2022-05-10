## petracli

private terraform registry cli

```
petracli [flags]
```

### Options

```
      --gcs-bucket string         Name of the Google Cloud Storage bucket you want to use for storage (required) e.g.: my-bucket
  -h, --help                      help for petracli
      --module-directory string   Directory of the module you want to upload (required) e.g.: ./modules-example/rabbitmq/
```

### SEE ALSO

* [petracli remove](petracli_remove.md)	 - Remove the .tar.gz file of a Terraform module in the bucket
* [petracli update](petracli_update.md)	 - Update one or multiple settings of a module and make changes in the .petra-config.yaml
* [petracli upload](petracli_upload.md)	 - Compress a Terraform module as a .tar.gz file and upload it to a bucket

