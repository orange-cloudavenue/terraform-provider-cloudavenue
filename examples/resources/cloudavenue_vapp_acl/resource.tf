# Shared wih everyone
resource "cloudavenue_vapp_acl" "example" {
  vdc                   = "MyVDC"
  vapp_name             = "MyVapp"
  everyone_access_level = "Change"
}

# Shared with user and/or group
data "cloudavenue_iam_user" "example_user" {
  user_name = "example-user"
}

data "cloudavenue_iam_group" "example_group" {
  group_name = "example-group"
}

resource "cloudavenue_vapp_acl" "example" {
  vdc       = "MyVDC"
  vapp_name = "MyVapp"
  shared_with = [
    {
      access_level = "FullControl"
      user_id      = data.cloudavenue_iam_user.example_user.id
    },
    {
      access_level = "ReadOnly"
      group_id     = data.cloudavenue_iam_group.example_group.id
  }]
}
