data "cloudavenue_elb_service_engine_groups" "example" {
  edge_gateway_name = data.cloudavenue_edge_gateway.example.name
}

output "example" {
  value = data.cloudavenue_elb_service_engine_groups.example
}
