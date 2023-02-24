data "cloudavenue_edgegateways" "example" {}

output "list_of_gateways" {
  value = data.cloudavenue_edgegateways.example
}
