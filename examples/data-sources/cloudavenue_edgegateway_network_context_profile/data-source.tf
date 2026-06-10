data "cloudavenue_edgegateway_network_context_profile" "example" {
  edge_gateway_name = cloudavenue_edgegateway.example.name
  name              = "SSL"
}
