data "cloudavenue_vdcg_security_group" "example" {
  name           = "my-security-group"
  vdc_group_name = cloudavenue_vdcg.example.name
}
