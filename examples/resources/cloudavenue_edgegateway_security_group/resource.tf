resource "cloudavenue_edgegateway_security_group" "example" {
  name        = "example"
  description = "This is an example security group"

  edge_gateway_id = cloudavenue_edgegateway.example.id

  member_org_network_ids = [
    cloudavenue_edgegateway_network_routed.example.id
  ]
}
