# terraform-provider-ceph example

To run this example you must have:
* Terraform, this have been tested with version 1.0.5
* Glesys credentials that can create an object storage instance
* terraform-provider-glesys, see [README.md](https://github.com/glesys/terraform-provider-glesys/blob/main/README.md) for installation notes
* terraform-provider-ceph, see [README.md](https://github.com/modfin/terraform-provider-ceph/blob/main/README.md) for installation notes

## Setup needed environment variables
```
export GLESYS_TOKEN="YOUR_TOKEN_WITH_OBJECT_STORAGE_PERMISSIONS"
export GLESYS_USERID="YOUR_USERID
```

## Initialize terraform project
```
terraform init
```

## Create S3 instance and buckets
```
terraform apply
```

## Destroy S3 instance and buckets
```
terraform destroy
```
