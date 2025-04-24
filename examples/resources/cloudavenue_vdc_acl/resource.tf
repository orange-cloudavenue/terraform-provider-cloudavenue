resource "cloudavenue_vdc_acl" "example" {
  vdc                   = cloudavenue_vdc.example.name
  everyone_access_level = "ReadOnly"
}
