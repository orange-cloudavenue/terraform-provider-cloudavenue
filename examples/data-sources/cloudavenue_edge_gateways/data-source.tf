data "cloudavenue_edge_gateways" "example" {}

output "gateways" {
  value = data.cloudavenue_edge_gateways.example
}
