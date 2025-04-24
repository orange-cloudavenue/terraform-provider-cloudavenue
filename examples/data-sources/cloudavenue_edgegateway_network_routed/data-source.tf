data "cloudavenue_edgegateway_network_routed" "example" {
  name            = "example"
  edge_gateway_id = cloudavenue_edgegateway.example.id
}
