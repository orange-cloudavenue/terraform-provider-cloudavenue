data "cloudavenue_edge_gateways" "example" {}

resource "cloudavenue_public_ip" "example" {
    edge_id = data.cloudavenue_edge_gateways.example.edge_gateways[0].edge_id
}
