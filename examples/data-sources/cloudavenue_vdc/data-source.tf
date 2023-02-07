data "cloudavenue_vdc" "example" {
  name = "VDC_Example"
}

output "example" {
  value = data.cloudavenue_vdc.example
}
