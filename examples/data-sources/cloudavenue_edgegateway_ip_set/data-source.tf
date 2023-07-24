data "cloudavenue_edgegateway_ip_set" "example" {
  name            = "example"
  edge_gateway_id = data.cloudavenue_edgegateway.example.id
}

data "cloudavenue_edgegateway" "example" {}
