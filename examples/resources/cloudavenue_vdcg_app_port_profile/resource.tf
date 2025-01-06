resource "cloudavenue_vdcg_app_port_profile" "example" {
  name         = "example-rule"
  description  = "Application port profile for example"
  vdc_group_id = cloudavenue_vdcg.example.id
  app_ports = [
    {
      protocol = "ICMPv4"
    },
    {
      protocol = "TCP"
      ports = [
        "80",
        "443",
        "8080-8090"
      ]
    },
  ]
}
