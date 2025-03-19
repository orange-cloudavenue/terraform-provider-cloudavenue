data "cloudavenue_edgegateway_network_routed" "example" {
  name              = "example"
  edge_gateway_name = cloudavenue_edgegateway.example.name
}
