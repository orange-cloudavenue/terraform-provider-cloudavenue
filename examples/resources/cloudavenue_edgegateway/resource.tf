data "cloudavenue_tier0_vrfs" "example_with_vdc" {}

resource "cloudavenue_edgegateway" "example_with_vdc" {
  owner_name     = "MyVDC"
  tier0_vrf_name = data.cloudavenue_tier0_vrfs.example_with_vdc.names.0
  owner_type     = "vdc"
  lb_enabled     = false
}