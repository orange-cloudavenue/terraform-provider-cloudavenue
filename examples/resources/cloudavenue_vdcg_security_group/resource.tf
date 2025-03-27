resource "cloudavenue_vdcg_security_group" "example" {
  name        = "example"
  description = "Example security group"

  vdc_group_id = cloudavenue_vdcg.example.id

  member_org_network_ids = [
    cloudavenue_vdcg_network_isolated.example.id,
    cloudavenue_vdcg_network_routed.example.id
  ]
}
