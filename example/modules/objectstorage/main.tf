# versions.tf
terraform {
  required_version = ">= 0.13"
  required_providers {
    glesys = {
      source = "github.com/glesys/glesys"
      version = "1.0.0"
    }
    ceph = {
      source = "github.com/modfin/ceph"
      version = "0.1.0"
    }
  }
}

# provider.tf
provider "ceph" {
  access_key = glesys_objectstorage_instance.s3.accesskey
  secret_key = glesys_objectstorage_instance.s3.secretkey
  region = var.datacenter
}

# input.tf
variable "datacenter" {
  type = string
  default = "dc-sto1"
}

variable "description" {
  type = string
}

variable "buckets" {
  type = set(string)
}

# main.tf
resource "glesys_objectstorage_instance" "s3" {
  datacenter = var.datacenter
  description = var.description
}

resource "ceph_s3_bucket" "bucket" {
  for_each = var.buckets
  name = each.value
}

#output.tf
output "access-key" {
  value = glesys_objectstorage_instance.s3.accesskey
}

output "secret-key" {
  value = glesys_objectstorage_instance.s3.secretkey
}
