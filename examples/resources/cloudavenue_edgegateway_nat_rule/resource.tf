resource "cloudavenue_edgegateway_nat_rule" "example" {
  edge_gateway_id = data.cloudavenue_edgegateways.example.edge_gateways[1].id

  name        = "example-snat"
  rule_type   = "SNAT"
  description = "description SNAT example"

  # Using primary_ip from edge gateway
  external_address         = data.cloudavenue_publicips.example.public_ips[2].public_ip
  internal_address         = "11.11.11.0/24"
  snat_destination_address = "8.8.8.8"

  priority = 10
}