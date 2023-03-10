data "cloudavenue_tier0_vrfs" "example_with_vdc" {}

resource "cloudavenue_edgegateway" "example_with_vdc" {
  owner_name     = "MyVdc"
  tier0_vrf_name = data.cloudavenue_tier0_vrfs.example.names.0
  owner_type     = "vdc"
}

data "cloudavenue_tier0_vrfs" "example_with_group" {}

resource "cloudavenue_edgegateway" "example_with_group" {
  owner_name     = "MyVDCGroup"
  tier0_vrf_name = data.cloudavenue_tier0_vrfs.example.names.0
  owner_type     = "vdc-group"
}