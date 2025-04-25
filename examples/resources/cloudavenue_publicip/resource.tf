resource "cloudavenue_publicip" "example" {
  edge_gateway_id = cloudavenue_edgegateway.example.id
}
