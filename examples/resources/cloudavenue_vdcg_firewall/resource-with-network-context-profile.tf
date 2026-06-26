# Use a built-in SYSTEM profile (e.g. SSL) referenced by name via data source
data "cloudavenue_vdcg_network_context_profile" "ssl" {
  vdc_group_name = cloudavenue_vdcg.example.name
  name           = "SSL"
}

resource "cloudavenue_vdcg_firewall" "example_with_system_profile" {
  vdc_group_name = cloudavenue_vdcg.example.name
  rules = [
    {
      action      = "ALLOW"
      name        = "allow outbound SSL"
      direction   = "OUT"
      ip_protocol = "IPV4"

      network_context_profile_ids = [data.cloudavenue_vdcg_network_context_profile.ssl.id]
    }
  ]
}

