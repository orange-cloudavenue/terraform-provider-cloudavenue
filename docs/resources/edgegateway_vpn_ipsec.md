---
page_title: "cloudavenue_edgegateway_vpn_ipsec Resource - cloudavenue"
subcategory: "Edge Gateway (Tier-1)"
description: |-
  Provides a resource to manage an IPsec VPN Tunnel. You can configure a site-to-site connectivity between an Edge Gateway and remote site. The remote site must support IPSec protocol. The VPN is able to initiate and respond to incoming tunnel requests. The VPN tunnel is established only when both sides of the tunnel are configured. The VPN tunnel is terminated when one side of the tunnel is deleted or disabled. The VPN tunnel is re-established when the disabled side is enabled again.
---

# cloudavenue_edgegateway_vpn_ipsec (Resource)

Provides a resource to manage an IPsec VPN Tunnel. You can configure a site-to-site connectivity between an Edge Gateway and remote site. The remote site must support IPSec protocol. The VPN is able to initiate and respond to incoming tunnel requests. The VPN tunnel is established only when both sides of the tunnel are configured. The VPN tunnel is terminated when one side of the tunnel is deleted or disabled. The VPN tunnel is re-established when the disabled side is enabled again.

## Example Usage

Example Usage (IPsec VPN Tunnel with default Security Profile)
```terraform
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
```

Example Usage (IPsec VPN Tunnel with a custom Security Profile)
```terraform
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
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `local_ip_address` (String) An IPv4 Address for the local endpoint. This has to be a sub-allocated IP on the Edge Gateway. This endpoint must be reach by the remote endpoint. Must be a valid IP with net.ParseIP.
- `local_networks` (Set of String) Set of local networks in CIDR format. This local_networks will be exchanged between both sites in order to route ip traffic in VPN tunnel.
- `name` (String) The Name of the IPsec VPN Tunnel Configuration.
- `pre_shared_key` (String, Sensitive) The Pre-Shared Key (PSK) is an Authentication method. Is a complex password (ASCII) that will be exchanged between both sites in order to set up the IPsec tunnel.
- `remote_ip_address` (String) An IPv4 Address for the remote endpoint. This is your remote VPN endpoint you need to reach. Must be a valid IP with net.ParseIP.
- `remote_networks` (Set of String) Set of remote networks in CIDR format. This remote_networks will be exchanged between both sites in order to route ip traffic in VPN tunnel.

### Optional

- `description` (String) A description of the IPsec VPN Tunnel Configuration.
- `edge_gateway_id` (String) <i style="color:red;font-weight: bold">(ForceNew)</i> The ID of the Edge Gateway. Ensure that one and only one attribute from this collection is set : `edge_gateway_name`, `edge_gateway_id`.
- `edge_gateway_name` (String) <i style="color:red;font-weight: bold">(ForceNew)</i> The Name of the Edge Gateway. Ensure that one and only one attribute from this collection is set : `edge_gateway_name`, `edge_gateway_id`.
- `enabled` (Boolean) Enable or Disable the IPsec VPN Tunnel Configuration. Value defaults to `true`.
- `security_profile` (Attributes) Customization of your IPSec configuration. The configuration used must be symmetric for both endpoint VPN. (see [below for nested schema](#nestedatt--security_profile))

### Read-Only

- `id` (String) The ID of the IPsec VPN Tunnel Configuration.
- `security_type` (String) Type of Security Profile used for the IPsec VPN Tunnel.

<a id="nestedatt--security_profile"></a>
### Nested Schema for `security_profile`

Optional:

- `ike_dh_groups` (String) The Diffie-Hellman (DH) key exchange algorithm is a method used to make a shared encryption key available to two entities over an insecure communications channel. Value must be one of : `GROUP2`, `GROUP5`, `GROUP14`, `GROUP15`, `GROUP16`, `GROUP19`, `GROUP20`, `GROUP21`.
- `ike_digest_algorithm` (String) Secure hashing algorithms to use during the IKE negotiation. Value must be one of : `SHA1`, `SHA2_256`, `SHA2_384`, `SHA2_512`. If the value of [`<.ike_encryption_algorithm`](#<.ike_encryption_algorithm) attribute is one of `AES_GCM_128`, `AES_GCM_256` or `AES_GCM_512` this attribute is **NULL**.
- `ike_encryption_algorithm` (String) Encryption algorithms used by IKE. Value must be one of : `AES_128`, `AES_256`, `AES_GCM_128`, `AES_GCM_192`, `AES_GCM_256`.
- `ike_sa_lifetime` (Number) Security association lifetime in seconds. It is number of seconds before the IPsec tunnel ike part needs to reestablish. Value defaults to `86400`. Value must be between 21600 and 31536000.
- `ike_version` (String) IKE (Internet Key Exchange) is an encrypt protocol of your VPN data. Value must be one of: 
  - `IKE_V1` When you select this option, IPSec VPN initiates and responds to IKEv1 protocol only.
  - `IKE_V2` The default option. When you select this version, IPSec VPN initiates and responds to IKEv2 protocol only.
  - `IKE_FLEX` When you select this option, if the tunnel establishment fails with IKEv2 protocol, the source site does not fall back and initiate a connection with the IKEv1 protocol. Instead, if the remote site initiates a connection with the IKEv1 protocol, then the connection is accepted.
- `tunnel_df_policy` (String) Policy for handling defragmentation. Value defaults to `COPY`. Value must be one of: 
  - `COPY` Copies the defragmentation bit from the inner IP packet to the outer packet.
  - `CLEAR` Ignores the defragmentation bit present in the inner packet.
- `tunnel_dh_groups` (String) The Diffie-Hellman (DH) key exchange algorithm is a method used to make a shared encryption key available to two entities over an insecure communications channel. Value must be one of : `GROUP2`, `GROUP5`, `GROUP14`, `GROUP15`, `GROUP16`, `GROUP19`, `GROUP20`, `GROUP21`.
- `tunnel_digest_algorithms` (String) Digest algorithms to be used for message digest. Value must be one of : `SHA1`, `SHA2_256`, `SHA2_384`, `SHA2_512`. If the value of [`<.tunnel_encryption_algorithms`](#<.tunnel_encryption_algorithms) attribute is one of `AES_GCM_128`, `AES_GCM_256` or `AES_GCM_512` this attribute is **NULL**.
- `tunnel_dpd` (Number) Value in seconds of Dead Probe Detection interval. Value defaults to `60`. Value must be between 3 and 60.
- `tunnel_encryption_algorithms` (String) Encryption algorithms to use in IPSec tunnel establishment. Value must be one of : `AES_128`, `AES_256`, `AES_GCM_128`, `AES_GCM_192`, `AES_GCM_256`.
- `tunnel_pfs` (Boolean) PFS (Perfect Forward Secrecy) capacity enabled or disabled. It's generates unique private keys for each secure session. Value defaults to `true`.
- `tunnel_sa_lifetime` (Number) Security association lifetime in seconds. It is number of seconds before the IPsec tunnel needs to reestablish. Value defaults to `3600`. Value must be between 900 and 31536000.

## Import

Import is supported using the following syntax:
```shell
terraform import cloudavenue_edgegateway_vpn_ipsec.example edgeGatewayIDOrName.vpnIPSecNameOrID
```