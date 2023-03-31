data "cloudavenue_tier0_vrfs" "example" {}

resource "cloudavenue_edgegateway" "example" {
  owner_name     = "MyVDC"
  tier0_vrf_name = data.cloudavenue_tier0_vrfs.example.names.0
  owner_type     = "vdc"
  lb_enabled     = true
}

resource "cloudavenue_alb_pool" "example" {
  edge_gateway_name = cloudavenue_edgegateway.example.name
  name              = "Example"

  persistence_profile = {
    type = "CLIENT_IP"
  }

  members = [
    {
      ip_address = "192.168.1.1"
      port       = "80"
    },
    {
      ip_address = "192.168.1.2"
      port       = "80"
    },
    {
      ip_address = "192.168.1.3"
      port       = "80"
    }
  ]

  health_monitors = ["UDP", "TCP"]
}

data "cloudavenue_tier0_vrfs" "example" {}

resource "cloudavenue_edgegateway" "example" {
  owner_name     = "MyVDC"
  tier0_vrf_name = data.cloudavenue_tier0_vrfs.example.names.0
  owner_type     = "vdc"
  lb_enabled     = true
}

resource "cloudavenue_alb_pool" "example" {
  edge_gateway_name = cloudavenue_edgegateway.example.name
  name              = "Example"
}