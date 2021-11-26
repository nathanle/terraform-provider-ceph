module "s3" {
  source = "./modules/objectstorage"
  description = "my test instance"
  buckets = ["first.test.bucket.modfin.se", "second.test.bucket.modfin.se"]
}

# Use `module.s3.access-key` and `module.s3.secret-key` as input to other
# modules/resources, that need those to access the bucket(s).