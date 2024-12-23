resource "cloudavenue_vdcg_firewall" "example_full" {
  vdc_group_name = cloudavenue_vdcg.example.name
  enabled        = true
  rules = [
    {
      action    = "ALLOW"
      name      = "allow all IPv4 traffic"
      direction = "IN"
      source_ids = [
        cloudavenue_vdcg_ip_set.example.id,
      ],
      destination_ids = [
        cloudavenue_vdcg_security_group.example.id,
      ],
      app_port_profile_ids = [
        data.cloudavenue_edgegateway_app_port_profile.example_provider_scope.id,
      ]
      source_groups_excluded      = true
      destination_groups_excluded = true
    },
    {
      action    = "DROP"
      name      = "drop IPv4 traffic"
      direction = "IN"
      source_ids = [
        cloudavenue_vdcg_ip_set.example.id,
      ],
      destination_ids = [
        cloudavenue_vdcg_security_group.example.id,
      ],
      app_port_profile_ids = [
        data.cloudavenue_edgegateway_app_port_profile.example_provider_scope.id,
        data.cloudavenue_edgegateway_app_port_profile.example_system_scope.id
      ]
      source_groups_excluded      = false
      destination_groups_excluded = false
    }
  ]
}
