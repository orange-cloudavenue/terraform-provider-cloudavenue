data "cloudavenue_tier0_vrfs" "example" {}

data "cloudavenue_tier0_vrf" "example" {
  name = data.cloudavenue_tier0_vrfs.example.names.0
}

output "vrf" {
  value = data.cloudavenue_tier0_vrf.example
}
