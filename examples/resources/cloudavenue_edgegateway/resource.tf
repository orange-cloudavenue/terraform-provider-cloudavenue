data "cloudavenue_tier0_vrfs" "example" {}

resource "cloudavenue_edgegateway" "example-with-vdc" {
  owner_name     = "MyVdc"
  tier0_vrf_name = data.cloudavenue_tier0_vrfs.example.names.0
  owner_type     = "vdc"
}

resource "cloudavenue_edgegateway" "example-with-group" {
  owner_name     = "MyVdcGroup"
  tier0_vrf_name = data.cloudavenue_tier0_vrfs.example.names.0
  owner_type     = "vdc-group"
}
