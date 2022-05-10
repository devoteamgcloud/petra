## petracli upload

Compress a Terraform module as a .tar.gz file and upload it to a bucket

```
petracli upload [flags]
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

* [petracli](petracli.md)	 - private terraform registry cli

