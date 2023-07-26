data "cloudavenue_edgegateway_nat_rule" "example" {
  depends_on      = [cloudavenue_edgegateway_nat_rule.example]
  edge_gateway_id = data.cloudavenue_edgegateway.main.id
  name            = "example-snat"
}