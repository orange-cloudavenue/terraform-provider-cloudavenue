resource "cloudavenue_vdcg_network_routed" "example" {
  name = "example"

  vdc_group_id    = cloudavenue_vdcg.example.id
  edge_gateway_id = cloudavenue_edgegateway.example.id

  gateway       = "192.168.1.254"
  prefix_length = 24

  dns1 = "1.1.1.1"
  dns2 = "8.8.8.8"

  dns_suffix = "example"

  static_ip_pool = [
    {
      start_address = "192.168.1.10"
      end_address   = "192.168.1.20"
    }
  ]
}
