resource "cloudavenue_vdcg_app_port_profile" "example" {
  name         = "MyApplication"
  description  = "Application port profile for my application"
  vdc_group_id = cloudavenue_vdcg.example.id
  app_ports = [
    {
      protocol = "TCP"
      ports = [
        "8080",
      ]
    },
  ]
}

resource "cloudavenue_vdcg_firewall" "example" {
  vdc_group_id = cloudavenue_vdcg.example.id
  rules = [{
    action               = "ALLOW"
    name                 = "From Internet to Application example"
    direction            = "IN"
    ip_protocol          = "IPV4"
    destination_ids      = [cloudavenue_vdcg_security_group.example.id]
    app_port_profile_ids = [cloudavenue_vdcg_app_port_profile.example.id]
  }]
}
