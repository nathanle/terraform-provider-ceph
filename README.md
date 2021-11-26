# terraform-provider-ceph (S3)
A very simple Terraform provider to create/delete buckets via CEPH S3 API.

# Build and install
```
go build -o terraform-provider-ceph
mkdir -p ~/.terraform.d/plugins/github.com/modfin/ceph/0.1.0/linux_amd64
ln -s $(pwd)/terraform-provider-ceph ~/.terraform.d/plugins/github.com/modfin/ceph/0.1.0/linux_amd64/terraform-provider-ceph
```

If the provider have been rebuilt since last `terraform init` run, terraform
will bail on a checksum error. To fix that you can remove the `.terraform.lock.hcl`
from the terraform folder, and run `terraform init` to load the checksum for the
re-built binary.


# Terraform Example
A simple Terraform example is provided in the example folder, that first create a
new object storage instance, and then use this provider to create two buckets.