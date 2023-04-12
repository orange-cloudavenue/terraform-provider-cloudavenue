data "cloudavenue_edgegateway" "example" {
  name = "tn01e02ocb0006205spt101"
}

resource "cloudavenue_network_routed" "example" {
  name            = "ExampleNetworkRouted"
  gateway         = "192.168.10.254"
  prefix_length   = 24
  edge_gateway_id = data.cloudavenue_edgegateway.example.id
  dns1            = "1.1.1.1"
  dns2            = "8.8.8.8"

  dns_suffix = "example"

  static_ip_pool = [
    {
      start_address = "192.168.10.10"
      end_address   = "192.168.10.20"
    }
  ]
}

data "cloudavenue_network_routed" "example" {
  name            = cloudavenue_network_routed.example.name
  edge_gateway_id = cloudavenue_network_routed.example.edge_gateway_id
}