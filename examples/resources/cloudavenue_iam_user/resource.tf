resource "cloudavenue_iam_user" "example" {
  name              = "example"
  role_name         = "Organization Administrator"
  password          = "Th!s1sSecur3P@ssword"
  enabled           = true # Default true
  email             = "foo@bar.com"
  telephone         = "1234567890"
  full_name         = "Test User"
  take_ownership    = true # Default true
  deployed_vm_quota = 10   # Default 0
  stored_vm_quota   = 5    # Default 0
}
