# Quickstart

## Setup the server

### 1. Create a Cloud Storage bucket and service account.   

```bash 
gcloud storage buckets create gs://$BUCKET_NAME --project=PROJECT_ID --location=$BUCKET_LOCATION --uniform-bucket-level-access
```
!!! node "`$BUCKET_NAME` needs to be a unique name"

Create a Service Account for the Petra server with the `Storage Object Admin` role on the `$BUCKET_NAME` bucket and `Service Account Token Creator` on the project if you want the server to run with the `SIGNED_URL` option.

```bash
gcloud iam service-accounts create $SA_NAME

# If you are not using the SIGNED_URL option, you will need to give each serviceAccount/user using the terraform modules read rights on the bucket. 
gcloud storage buckets add-iam-policy-binding gs://$BUCKET_NAME \
    --member=serviceAccount:$SA_NAME@$PROJECT_ID.iam.gserviceaccount.com \
    --role=roles/storage.objectAdmin

# (optional) If you want to run Petra server in the SIGNED_URL mode.
gcloud projects add-iam-policy-binding $PROJECT_ID \
    --member="serviceAccount:SA_NAME@PROJECT_ID.iam.gserviceaccount.com" \
    --role="roles/iam.serviceAccountTokenCreator"
```
!!! note "To push modules to the `$BUCKET_NAME` bucket, the serviceAccount/user running petractl will need write access to the bucket."

### 2. Build & Deploy to Cloud Run

Since `ghcr.io` images cannot be used in Cloud Run, you will have to host the docker image on either [Google's Artifact Registry](https://cloud.google.com/artifact-registry/docs/docker/store-docker-container-images) or {--Container Registry--} `(deprecated)`, or DockerHub. 
With Artifact Registry, you `$IMAGE_PATH` will be look like `[$REGION-]docker.pkg.dev/$PROJECT_ID/$REGISTRY_NAME/petra:$TAG`

=== "Rebuild and push with Cloud Build"

    ```yaml title="cloudbuild.yaml"
    steps:
    # Build the container image
    - name: 'gcr.io/cloud-builders/docker'
    args: ['build', '-t', '$IMAGE_PATH', '-f',  'Dockerfile.manual', '.']
    # Push the container image to Container Registry
    - name: 'gcr.io/cloud-builders/docker'
    args: ['push', '$IMAGE_PATH']
    # Deploy container image to Cloud Run
    - name: 'gcr.io/cloud-builders/gcloud'
    entrypoint: gcloud
    args: ['run', 'deploy', 'petra', '--image', '$IMAGE_PATH', '--platform', 'managed', '--region', '$REGION', '--allow-unauthenticated' '--set-env-vars', 'GCS_BUCKET=$BUCKET']
    images:
    -  $IMAGE_PATH
    ```

Here, the Cloud Run service created does not require authentication as it is not supported by terraform. The access is managed by managing the IAM access to the bucket itself. The Cloud Run is deployed with all traffic allowed, you can also deploy it and allow only internal traffic but you have to make sure that it exposes a HTTPs endpoint or terraform won't pull the modules.

## Pushing modules

### 1. Download Petractl

Go to the [latest release page](https://github.com/devoteamgcloud/petra/releases/latest) and one of the following files depending on you OS / CPU architecure :

=== "MacOS"

    - petra_x.x.x_darwin_amd64.tar.gz (Intel)
    - petra_x.x.x_darwin_arm64.tar.gz (Apple Silicon)

=== "Linux"

    - petra_0.4.1_linux_amd64.tar.gz (x86_64)
    - petra_0.4.1_linux_arm64.tar.gz (arm64)

### 2. Create a petra-config.yaml file for your module

```yaml title="petra-config.yaml"
namespace: production
name: my-module
provider: google
version: 1.0.3
```

!!! note "The file must be located at the root of your module"

### 3. Push the module

Make sure that you have the correct access rights on the `$BUCKET_NAME` bucket to be able to write new files. 

```bash
petractl push --bucket $BUCKET_NAME ./path/to/module
```

## Use your modules in Terraform

Make sure that the user or service account that will be running the terraform init command has access to the petra server and has read access rights to the `$BUCKET_NAME` bucket or that the `SIGNED_URL` option is enabled in order to use the modules.

Here's an example of Terraform code using a petra hosted module :

```hcl title="main.tf"
module "mod1" {
    source = "petra.example.com/production/my-module/google"
    version = "1.0.3"
}
```

!!! warning "`petra.example.com` has to support HTTPs for terraform to allow the init command"