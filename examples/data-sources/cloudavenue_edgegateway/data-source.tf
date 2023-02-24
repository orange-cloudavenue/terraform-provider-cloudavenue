data "cloudavenue_edgegateway" "example" {
  edge_id = "cc1f35c2-90a2-48d1-9359-62794faf44ad"
}

output "gateway" {
  value = data.cloudavenue_edgegateway.example
}
