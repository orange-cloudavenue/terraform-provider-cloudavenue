data "cloudavenue_iam_user" "example" {
  name = "example-user"
}

output "example_user_id" {
  value = data.cloudavenue_iam_user.example
}
