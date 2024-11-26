resource "cloudavenue_vdcg_firewall" "example_with_source_ids" {
  vdc_group_name = cloudavenue_vdcg.example.name
  enabled        = true
  rules = [
    {
      action    = "ALLOW"
      name      = "allow in IPv4 traffic"
      direction = "IN"
      source_ids = [
        cloudavenue_vdcg_ip_set.example.id,
        cloudavenue_vdcg_security_group.example.id
      ]
    }
  ]
}
