# Example NAT Rule SNAT (NAT out from network 11.11.11.0/24 to dest 8.8.8.8 translate in 89.32.25.10)
resource "cloudavenue_edgegateway_nat_rule" "example-snat" {
  edge_gateway_name = "myEdgeGateway"

  name        = "example-snat"
  rule_type   = "SNAT"
  description = "description SNAT example"

  external_address         = "89.32.25.10"
  internal_address         = "11.11.11.0/24"
  snat_destination_address = "8.8.8.8"

  priority = 10
}

# Example NAT Rule DNAT (Translate 89.32.25.10 to internal dest 4.11.11.11 on port 8080)
resource "cloudavenue_edgegateway_nat_rule" "example-dnat" {
  edge_gateway_name = "myEdgeGateway"

  name        = "example-dnat"
  rule_type   = "DNAT"
  description = "description DNAT example"

  external_address = "89.32.25.10"
  internal_address = "4.11.11.11"

  dnat_external_port = "8080"
}

# Example NAT Rule Reflexive (Nat in both way (in and out) external and internal on all port translated)
resource "cloudavenue_edgegateway_nat_rule" "example-reflexive" {
  edge_gateway_name = "myEdgeGateway"

  name        = "example-reflexive"
  rule_type   = "REFLEXIVE"
  description = "description REFLEXIVE example"

  external_address = "89.32.25.10"
  internal_address = "192.168.0.1"

  priority = 25
}