resource "cloudavenue_iam_user" "example" {
  name      = "example"
  role_name = "Organization Administrator"
  password  = "Th!s1sSecur3P@ssword"
}

resource "cloudavenue_vapp" "example" {
  name        = "MyVapp"
  description = "This is an example vApp"
}

resource "cloudavenue_vapp_acl" "example" {
  vapp_name = cloudavenue_vapp.example.name
  shared_with = [{
    access_level = "ReadOnly"
    user_id      = cloudavenue_iam_user.example.id
  }]
}