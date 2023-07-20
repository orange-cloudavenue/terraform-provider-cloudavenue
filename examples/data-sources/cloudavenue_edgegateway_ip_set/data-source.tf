data "cloudavenue_edgegateway_ip_set" "example" {
  name            = "example"
  edge_gateway_id = cloudavenue_edgegateway.example[0].id
}

data "cloudavenue_edgegateways" "example" {}
