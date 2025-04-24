resource "cloudavenue_edgegateway_dhcp_forwarding" "example" {
  edge_gateway_id = cloudavenue_edgegateway.example.id

  dhcp_servers = [
    "192.168.10.10"
  ]
}
