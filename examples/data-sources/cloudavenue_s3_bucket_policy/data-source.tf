data "cloudavenue_s3_bucket_policy" "example" {
  bucket = cloudavenue_s3_bucket.example.name
}
