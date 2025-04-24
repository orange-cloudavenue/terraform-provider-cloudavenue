data "cloudavenue_network_dhcp" "example" {
  org_network_id = cloudavenue_edgegateway_network_routed.example.id
}
