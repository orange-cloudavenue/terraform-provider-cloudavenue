resource "cloudavenue_network_isolated" "example" {
  name        = "rsx-example-isolated-network"
  description = "My isolated Org VDC network"

  gateway       = "1.1.1.1"
  prefix_length = 24

  dns1       = "8.8.8.8"
  dns2       = "8.8.4.4"
  dns_suffix = "example.com"

  static_ip_pool = [
    {
      start_address = "1.1.1.10"
      end_address   = "1.1.1.20"
    },
    {
      start_address = "1.1.1.100"
      end_address   = "1.1.1.103"
    }
  ]
}