resource "cloudavenue_edgegateway" "example" {
  owner_name     = cloudavenue_vdcg.example.name
  tier0_vrf_name = data.cloudavenue_tier0_vrf.example.name
}
