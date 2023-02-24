data "cloudavenue_edgegateways" "example" {}

output "gateways" {
  value = data.cloudavenue_edgegateways.example
}
