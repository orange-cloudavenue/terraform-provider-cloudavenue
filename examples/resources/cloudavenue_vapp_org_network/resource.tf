resource "cloudavenue_vapp_org_network" "example" {
  vapp_name    = cloudavenue_vapp.example.name
  network_name = cloudavenue_edgegateway_network_routed.example.name
  vdc          = cloudavenue_vdc.example.name
}
