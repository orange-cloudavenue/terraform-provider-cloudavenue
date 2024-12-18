data "cloudavenue_alb_service_engine_group" "example" {
  name              = "my-service-engine"
  edge_gateway_name = data.cloudavenue_edge_gateway.example.name
}

output "example" {
  value = data.cloudavenue_alb_service_engine_group.example
}