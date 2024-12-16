resource "cloudavenue_vdc_network_isolated" "example" {
  name = "my-isolated-network"
  vdc  = cloudavenue_vdc.example.name

  gateway       = "192.168.0.1"
  prefix_length = 24

  dns1       = "192.168.0.2"
  dns2       = "192.168.0.3"
  dns_suffix = "example.local"
}
