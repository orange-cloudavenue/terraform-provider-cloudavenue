data "cloudavenue_elb_virtual_service" "example" {
  name            = "example"
  edge_gateway_id = data.cloudavenue_edgegateway.example.id
}
