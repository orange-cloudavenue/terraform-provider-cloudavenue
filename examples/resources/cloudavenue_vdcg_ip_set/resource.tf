resource "cloudavenue_vdcg_ip_set" "example" {
  name        = "example"
  description = "example of ip set"
  ip_addresses = [
    "12.12.12.1",            # IP Address
    "10.10.10.0/24",         # IP Address With CIDR
    "11.11.11.1-11.11.11.2", # IP Address Range
  ]
  vdc_group_name = cloudavenue_vdcg.name
}
