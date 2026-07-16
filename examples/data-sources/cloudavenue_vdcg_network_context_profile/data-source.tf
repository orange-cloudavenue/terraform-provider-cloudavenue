data "cloudavenue_vdcg_network_context_profile" "example" {
  vdc_group_name = cloudavenue_vdcg.example.name
  name           = "SSL"
}
