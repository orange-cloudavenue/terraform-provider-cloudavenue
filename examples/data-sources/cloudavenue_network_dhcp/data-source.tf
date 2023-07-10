data "cloudavenue_network_dhcp" "example" {
  org_network_id = cloudavenue_network_routed.example.id
}


data "cloudavenue_edgegateways" "example" {}

resource "cloudavenue_network_routed" "example" {
  name        = "MyOrgNet"
  description = "This is an example Net"

  edge_gateway_id = data.cloudavenue_edgegateways.example.edge_gateways[0].id

  gateway       = "192.168.1.254"
  prefix_length = 24

  dns1 = "1.1.1.1"
  dns2 = "8.8.8.8"

  dns_suffix = "example"

  static_ip_pool = [
    {
      start_address = "192.168.1.10"
      end_address   = "192.168.1.20"
    }
  ]
}
