data "cloudavenue_elb_pool" "example" {
  name            = "example"
  edge_gateway_id = cloudavenue_edge_gateway.example.id
}
