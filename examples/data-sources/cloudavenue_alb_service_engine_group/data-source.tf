data "cloudavenue_alb_service_engine_group" "example" {}

output "example" {
  value = data.cloudavenue_alb_service_engine_group.example
}