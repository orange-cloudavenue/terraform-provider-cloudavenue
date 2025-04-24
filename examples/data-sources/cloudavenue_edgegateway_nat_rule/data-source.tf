data "cloudavenue_edgegateway_nat_rule" "example" {
  name            = "example-snat"
  edge_gateway_id = cloudavenue_edgegateway.example.id
}
