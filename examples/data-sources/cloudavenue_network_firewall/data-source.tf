data "cloudavenue_edgegateways" "example" {}

data "cloudavenue_network_firewall" "example" {
  edge_gateway_id = data.cloudavenue_edgegateways.example.edge_gateways[0].id
}
