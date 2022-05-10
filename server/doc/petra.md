## petra

Private terraform registry server

### Synopsis

Server to get versions of a terraform module in the private registry (Google Cloud Storage bucket) and to get a signed URL to download a module

```
petra [flags]
```

### Options

```
      --gcs-bucket string       Name of the Google Cloud Storage bucket you want to use for storage.
  -h, --help                    help for petra
      --listen-address string   Address to listen on (default "3000")
      --project-id string       Google Cloud project ID where the service account is stored in Secret Manager.
      --secret-id string        (Google Cloud Secret Manager) Secret ID of your service-account that allows you to generate signed URLs.
```

