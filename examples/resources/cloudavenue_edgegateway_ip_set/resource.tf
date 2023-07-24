resource "cloudavenue_edgegateway_ip_set" "example" {
  name        = "example"
  description = "example of ip set"
  ip_addresses = [
    "12.12.12.1",            # IP Address
    "10.10.10.0/24",         # IP Address With CIDR
    "11.11.11.1-11.11.11.2", # IP Address Range
  ]
  edge_gateway_name = data.cloudavenue_edgegateway.example.name
}

data "cloudavenue_edgegateway" "example" {}
