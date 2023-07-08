resource "cloudavenue_network_dhcp_binding" "example" {
  name           = "example"
  org_network_id = cloudavenue_network_dhcp.example.id
  mac_address    = "00:50:56:01:01:01"
  ip_address     = "192.168.1.231"
}

resource "cloudavenue_network_dhcp" "example" {
  org_network_id = cloudavenue_network_routed.example.id
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

