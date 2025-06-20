resource "cloudavenue_s3_bucket_versioning_configuration" "example" {
  bucket = cloudavenue_s3_bucket.example.name
  status = "Enabled"
}