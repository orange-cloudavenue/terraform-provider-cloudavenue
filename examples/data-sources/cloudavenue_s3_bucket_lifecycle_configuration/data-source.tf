data "cloudavenue_s3_bucket_lifecycle_configuration" "example" {
  bucket = cloudavenue_s3_bucket.example.name
}
