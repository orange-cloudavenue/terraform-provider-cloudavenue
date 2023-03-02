data "cloudavenue_edgegateway" "example" {
  name = "myEdgeName"
}

output "gateway" {
  value = data.cloudavenue_edgegateway.example
}
