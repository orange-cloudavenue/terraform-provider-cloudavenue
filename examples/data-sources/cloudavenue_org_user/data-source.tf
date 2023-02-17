data "cloudavenue_org_user" "example" {
	user_name = "example-user"
}

output "example_user_id" {
  value = data.cloudavenue_org_user.example.id
}
