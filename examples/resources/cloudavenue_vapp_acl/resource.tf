resource "cloudavenue_iam_user" "example" {
  user_name   = "exampleuser"
  description = "An example user"
  role        = "Organization Administrator"
  password    = "Th!s1sSecur3P@ssword"
}

resource "cloudavenue_iam_group" "example" {
  name        = "examplegroup"
  role        = "Organization Administrator"
  description = "An example group"
}

resource "cloudavenue_vapp" "example" {
  name        = "MyVapp"
  description = "This is an example vApp"
}

resource "cloudavenue_vapp_acl" "example" {
  vdc       = "MyVDC" # Optional
  vapp_name = cloudavenue_vapp.example.name
  shared_with = [{
    access_level = "ReadOnly"
    user_id      = cloudavenue_iam_user.example.id
    },
    {
      access_level = "FullControl"
      group_id     = cloudavenue_iam_group.example.id
  }]
}