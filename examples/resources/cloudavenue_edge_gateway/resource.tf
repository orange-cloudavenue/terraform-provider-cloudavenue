resource "cloudavenue_edge_gateway" "example" {
  vdc_name = "VDC_Frangipane"
  tier0_vrf_id = "vrf-1"
  owner_type = "vdc"
}
