data "cloudavenue_s3_bucket_website_configuration" "example" {
  bucket = cloudavenue_s3_bucket.example.name
}
