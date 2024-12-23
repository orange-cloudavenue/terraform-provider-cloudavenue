data "cloudavenue_vdcg_dynamic_security_group" "example" {
  name           = "my-dynamic-security-group"
  vdc_group_name = cloudavenue_vdcg.example.name
}
