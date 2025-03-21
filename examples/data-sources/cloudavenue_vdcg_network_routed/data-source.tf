data "cloudavenue_vdcg_network_routed" "example" {
  name              = "example"
  vdc_group_id      = cloudavenue_vdcg.example.id
  edge_gateway_name = cloudavenue_edgegateway.example.name
}
