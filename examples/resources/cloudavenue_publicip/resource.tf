data "cloudavenue_edgegateways" "example" {}

resource "cloudavenue_publicip" "example" {
  edge_id = data.cloudavenue_edgegateways.example.edge_gateways[0].id
}
