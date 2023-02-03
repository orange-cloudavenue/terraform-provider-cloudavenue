data "cloudavenue_public_ip" "example" {}

output "public_ip" {
  value = data.cloudavenue_public_ip.example
}
