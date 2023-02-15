resource "cloudavenue_org_user" "example" {
  user_name   = "exampleuser"
  description = "A example user"
  role        = "Organization Administrator"
  password    = "Th!s1sSecur3P@ssword"
}