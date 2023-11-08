data "cloudavenue_s3_bucket_acl" "example" {
  bucket = cloudavenue_s3_bucket.example.name
}
