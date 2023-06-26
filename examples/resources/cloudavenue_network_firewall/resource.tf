data "cloudavenue_edgegateways" "example" {}

resource "cloudavenue_network_firewall" "example" {

  edge_gateway_id = data.cloudavenue_edgegateways.example.edge_gateways[0].id
  rules = [
    {
      action      = "ALLOW"
      name        = "allow all IPv4 traffic"
      direction   = "IN_OUT"
      ip_protocol = "IPV4"
    }
  ]
}