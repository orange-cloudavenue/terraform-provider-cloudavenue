resource "cloudavenue_edgegateway_vpn_ipsec" "example" {
  edge_gateway_id = cloudavenue_edgegateway.example.id

  name        = "example"
  description = "example VPN IPSec"
  enabled     = true

  pre_shared_key = "my-preshared-key"

  local_ip_address = cloudavenue_publicip.example.public_ip
  local_networks   = ["10.10.10.0/24", "30.30.30.0/28"]

  remote_ip_address = "203.0.113.1"
  remote_networks   = ["192.168.1.0/24", "192.168.10.0/24"]
}
