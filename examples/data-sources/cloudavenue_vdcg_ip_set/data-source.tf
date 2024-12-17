data "cloudavenue_vdcg_ip_set" "example" {
  name           = "example"
  vdc_group_name = cloudavenue_vdcg.name
}
