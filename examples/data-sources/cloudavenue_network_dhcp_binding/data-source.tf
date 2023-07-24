data "cloudavenue_network_dhcp_binding" "example" {
  name           = "example"
  org_network_id = cloudavenue_network_routed.example.id
}
