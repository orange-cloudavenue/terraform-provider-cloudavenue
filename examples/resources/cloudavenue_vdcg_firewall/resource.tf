resource "cloudavenue_vdcg_firewall" "example" {
  vdc_group_name = cloudavenue_vdcg.example.name
  rules = [
    {
      action      = "ALLOW"
      name        = "allow all IPv4 traffic"
      direction   = "IN_OUT"
      ip_protocol = "IPV4"
    }
  ]
}
