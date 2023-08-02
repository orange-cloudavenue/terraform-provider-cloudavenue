data "cloudavenue_edgegateway_nat_rule" "example" {
  edge_gateway_name = "myEdgeName"
  name              = "example-snat"
}