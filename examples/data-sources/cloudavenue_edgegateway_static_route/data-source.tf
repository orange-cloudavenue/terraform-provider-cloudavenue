data "cloudavenue_edgegateway_static_route" "example" {
  edge_gateway_name = cloudavenue_edgegateway.example.name
  name              = "example"
}
