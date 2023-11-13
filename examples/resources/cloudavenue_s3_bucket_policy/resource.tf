data "cloudavenue_s3_bucket" "example" {
  bucket = "example-bucket"
}

resource "cloudavenue_s3_bucket_policy" "example" {
  bucket = data.cloudavenue_s3_bucket.example.name
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect    = "Allow"
        Principal = "*"
        Action = [
          "s3:DeleteObject",
          "s3:GetObject",
          "s3:ListBucketVersions",
        ]
        Resource = [
          data.cloudavenue_s3_bucket.example.arn,
          "${data.cloudavenue_s3_bucket.example.arn}/*",
        ]
      }
    ]
  })
}