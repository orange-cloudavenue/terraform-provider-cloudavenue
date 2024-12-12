data "cloudavenue_vdc_network_isolated" "example" {
  vdc  = cloudavenue_vdc.example.name
  name = "my-isolated-network"
}
