resource "cloudavenue_edgegateway_vpn_ipsec" "example" {
  edge_gateway_id = cloudavenue_edgegateway.example.id

  name        = "example"
  description = "example VPN IPSec"
  enabled     = false

  pre_shared_key = "my-preshared-key"

  local_ip_address = cloudavenue_publicip.example.public_ip
  local_networks   = ["10.10.10.0/24", "30.30.30.0/28"]

  remote_ip_address = "203.0.113.1"
  remote_networks   = ["192.168.1.0/24", "192.168.10.0/24"]

  security_profile = {
    ike_dh_groups                = "GROUP15"
    ike_digest_algorithm         = "SHA2_384"
    ike_encryption_algorithm     = "AES_128"
    ike_sa_lifetime              = 86400
    ike_version                  = "IKE_V2"
    tunnel_df_policy             = "COPY"
    tunnel_dh_groups             = "GROUP15"
    tunnel_digest_algorithms     = "SHA2_512"
    tunnel_dpd                   = 45
    tunnel_encryption_algorithms = "AES_128"
    tunnel_pfs                   = true
    tunnel_sa_lifetime           = 3600
  }
}