resource "cloudavenue_edgegateway_app_port_profile" "example" {
  name            = "example-rule"
  description     = "Application port profile for example"
  edge_gateway_id = cloudavenue_edgegateway.example.id
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
