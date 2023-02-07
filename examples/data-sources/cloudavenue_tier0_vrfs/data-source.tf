data "cloudavenue_tier0_vrfs" "example" {}

output "vrfs" {
  value = data.cloudavenue_tier0_vrfs.example
}
