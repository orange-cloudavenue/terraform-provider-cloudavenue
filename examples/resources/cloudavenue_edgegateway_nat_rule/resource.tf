resource "cloudavenue_edgegateway_nat_rule" "example" {
  edge_gateway_name = "myEdgeGateway"

  name        = "example-snat"
  rule_type   = "SNAT"
  description = "description SNAT example"

  # Using primary_ip from edge gateway
  external_address         = "89.32.25.10"
  internal_address         = "11.11.11.0/24"
  snat_destination_address = "8.8.8.8"

  priority = 10
}