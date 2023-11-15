data "cloudavenue_s3_user" "example" {
  user_name = cloudavenue_iam_user.example.name
}
