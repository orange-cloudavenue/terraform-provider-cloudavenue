data "cloudavenue_edgegateway_vpn_ipsec" "example" {
  edge_gateway_name = cloudavenue_edgegateway.example.name
  name              = "example"
}