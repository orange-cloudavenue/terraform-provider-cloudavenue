resource "cloudavenue_network_dhcp" "example" {
  org_network_id = cloudavenue_edgegateway_network_routed.example.id
  mode           = "EDGE"
  pools = [
    {
      start_address = "192.168.1.30"
      end_address   = "192.168.1.100"
    }
  ]
  dns_servers = [
    "1.1.1.1",
    "1.0.0.1"
  ]
}
