# Petra


Putting it simply, Petra is a private registry for Terraform modules for GCP. 

## Why we made it ?

Petra was created to solve multiple problems encountered by most people that are using Terraform :

- **Internal Module Sharing** : Organizations often have their own set of standardized infrastructure modules that are used across multiple projects. With a private registry, these modules can be centrally stored and easily shared among different teams or projects within the organization. It ensures consistency, promotes reuse, and helps maintain best practices across infrastructure deployments.

- **Customization and Version Control** : A private registry allows organizations to customize and modify existing Terraform modules to fit their specific needs without impacting the original module source. This is particularly useful when making project-specific modifications or incorporating organization-specific configurations. Having version control in the private registry enables tracking and management of these customizations.

- **Offline Environments** : In certain scenarios, infrastructure deployments may be required in offline or air-gapped environments where internet access is limited or restricted. With a private registry, organizations can create a local copy of the registry and make it available within the restricted environment. This allows for self-contained infrastructure provisioning without relying on external networks.

- **Decorrelate Git repositories from Terraform modules**: People in need of a private terraform registry often try to use **Terraform Cloud** for their private terraform registry needs but find the limitation of having to create a git repository per module to be too cumbersome. 

## Usage

The server can be deployed using it's docker image on Kubernetes or even Cloud Run. All the server needs is rights to a Google Cloud Storage Bucket created to store the modules.

You can use your terraform modules like so :

![terraform-init](assets/terraform-init.gif)

!!! warning "The petra server endpoint has to support HTTPs or terraform won't allow you to init the module."

`petractl` needs a petra-config.yaml file to be present in the directories of the related terraform modules. You are then able to push modules directly to the bucket using the command `petractl push -b $MY_BUCKET .` 

![petractl](assets/petractl.gif)