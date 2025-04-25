data "cloudavenue_edgegateway_security_group" "example" {
  name            = "example"
  edge_gateway_id = cloudavenue_edgegateway.example.id
}
