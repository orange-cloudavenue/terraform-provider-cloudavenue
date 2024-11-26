resource "cloudavenue_vdcg_firewall" "example_with_app_port_profile" {
  vdc_group_name = cloudavenue_vdcg.example.name
  rules = [
    {
      action    = "ALLOW"
      name      = "allow all IPv4 traffic"
      direction = "IN_OUT"
      app_port_profile_ids = [
        data.cloudavenue_edgegateway_app_port_profile.example.id,
      ]
    }
  ]
}
