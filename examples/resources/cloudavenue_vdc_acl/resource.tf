resource "cloudavenue_vdc_acl" "example" {
  vdc                   = "VDC_Test" # Optional
  everyone_access_level = "ReadOnly"
}