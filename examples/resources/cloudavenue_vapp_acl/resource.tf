resource "cloudavenue_vapp_acl" "example" {
  vapp_name = cloudavenue_vapp.example.name
  shared_with = [{
    access_level = "ReadOnly"
    user_id      = cloudavenue_iam_user.example.id
  }]
}
