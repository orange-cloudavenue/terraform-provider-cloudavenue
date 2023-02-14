data "cloudavenue_public_ips" "example" {}

output "public_ips" {
  value = data.cloudavenue_public_ips.example
}
