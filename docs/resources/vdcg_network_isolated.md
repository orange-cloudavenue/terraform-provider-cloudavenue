---
page_title: "cloudavenue_vdcg_network_isolated Resource - cloudavenue"
subcategory: "vDC Group (Virtual Datacenter Group)"
description: |-
  The cloudavenue_vdcg_network_isolated resource allows you to manage an isolated network in a VDC Group.
---

# cloudavenue_vdcg_network_isolated (Resource)

The `cloudavenue_vdcg_network_isolated` resource allows you to manage an isolated network in a `VDC Group`.
 
## Example Usage

```terraform
resource "cloudavenue_vdcg_network_isolated" "example" {
  name           = "my-isolated-network"
  vdc_group_name = cloudavenue_vdcg.example.name

  gateway       = "192.168.0.1"
  prefix_length = 24

  dns1       = "192.168.0.2"
  dns2       = "192.168.0.3"
  dns_suffix = "example.local"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `gateway` (String) (ForceNew) The gateway IP address for the network. This value define also the network IP range with the prefix length. Must be a valid IP with net.ParseIP.
- `name` (String) The name of the network. This value must be unique within the `VDC` that owns the network.
- `prefix_length` (Number) (ForceNew) The prefix length for the network. This value must be a valid prefix length for the network IP range. (e.g. /24 for netmask 255.255.255.0). For more information, see [CIDR notation](https://en.wikipedia.org/wiki/Classless_Inter-Domain_Routing). Value must be between 1 and 32.

### Optional

- `description` (String) A description of the network.
- `dns1` (String) The primary DNS server IP address for the network. Must be a valid IP with net.ParseIP.
- `dns2` (String) The secondary DNS server IP address for the network. Must be a valid IP with net.ParseIP.
- `dns_suffix` (String) The DNS suffix for the network.
- `guest_vlan_allowed` (Boolean) Indicates if the network allows guest VLANs. Value defaults to `false`.
- `static_ip_pool` (Attributes Set) A set of static IP pools to be used for this network. (see [below for nested schema](#nestedatt--static_ip_pool))
- `vdc_group_id` (String) (ForceNew) The ID of vDC group that owns the network. Ensure that at least one attribute from this collection is set: [vdc_group_name,vdc_group_id].
- `vdc_group_name` (String) (ForceNew) The name of vDC group that owns the network. Ensure that at least one attribute from this collection is set: [vdc_group_name,vdc_group_id].

### Read-Only

- `id` (String) The ID of the isolated network.

<a id="nestedatt--static_ip_pool"></a>
### Nested Schema for `static_ip_pool`

Required:

- `end_address` (String) The end address of the IP pool. This value must be a valid IP address in the network IP range. Must be a valid IP with net.ParseIP.
- `start_address` (String) The start address of the IP pool. This value must be a valid IP address in the network IP range. Must be a valid IP with net.ParseIP.

## Import

Import is supported using the following syntax:
```shell
# VDC Network isolated can be imported using the VDC Groupe name or ID and the network name or ID.
terraform import cloudavenue_vdcg_network_isolated.example vdcGroupNameOrId.networkNameOrId
```