resource "cloudavenue_network_dhcp_binding" "example" {
  name = "example"

  org_network_id = cloudavenue_network_dhcp.example.id

  mac_address = "00:50:56:01:01:01"
  ip_address  = "192.168.1.231"
}
