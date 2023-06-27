data "cloudavenue_vdc" "example" {
  name = "VDC_Test"
}

resource "cloudavenue_edgegateway_app_port_profile" "example" {
  name        = "example-rule"
  description = "Application port profile for example"
  vdc         = data.cloudavenue_vdc.example.id

  app_ports = [
    {
      protocol = "ICMPv4"
    },
    {
      protocol = "TCP"
      ports = [
        "80",
        "443",
      ]
    },
  ]
}