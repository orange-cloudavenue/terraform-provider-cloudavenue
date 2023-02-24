data "cloudavenue_publicips" "example" {}

output "public_ips" {
  value = data.cloudavenue_publicips.example
}
