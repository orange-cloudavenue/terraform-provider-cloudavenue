data "cloudavenue_s3_bucket_ownership_controls" "example" {
  bucket = cloudavenue_s3_bucket.example.name
}
