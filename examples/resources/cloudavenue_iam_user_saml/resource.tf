resource "cloudavenue_iam_user_saml" "example" {
  user_name         = "example"
  role_name         = "Organization Administrator"
  enabled           = true # Default true
  take_ownership    = true # Default true
  deployed_vm_quota = 10   # Default 0
  stored_vm_quota   = 5    # Default 0
}
