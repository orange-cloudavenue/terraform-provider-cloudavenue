resource "cloudavenue_vdcg_security_group" "example" {
  vdc_group_id = cloudavenue_vdcg_network_isolated.example.vdc_group_id
  name         = "example"
  description  = "Example security group"
  member_org_network_ids = [
    cloudavenue_vdcg_network_isolated.example.id
  ]
}
