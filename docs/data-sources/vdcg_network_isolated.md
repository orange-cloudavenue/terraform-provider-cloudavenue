---
page_title: "cloudavenue_vdcg_network_isolated Data Source - cloudavenue"
subcategory: "vDC Group (Virtual Datacenter Group)"
description: |-
  The cloudavenue_vdcg_network_isolated data source allows you to retrieve information about an isolated network in a VDC Group.
---

# cloudavenue_vdcg_network_isolated (Data Source)

The `cloudavenue_vdcg_network_isolated` data source allows you to retrieve information about an isolated network in a `VDC Group`.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The name of the network. This value must be unique within the `VDC` that owns the network.

### Optional

- `vdc_group_id` (String) The ID of vDC group that owns the network. Ensure that at least one attribute from this collection is set: [vdc_group_name,vdc_group_id]. This value must start with `urn:vcloud:vdcGroup:`.
- `vdc_group_name` (String) The name of vDC group that owns the network. Ensure that at least one attribute from this collection is set: [vdc_group_name,vdc_group_id].

### Read-Only

- `description` (String) A description of the network.
- `dns1` (String) The primary DNS server IP address for the network.
- `dns2` (String) The secondary DNS server IP address for the network.
- `dns_suffix` (String) The DNS suffix for the network.
- `gateway` (String) The gateway IP address for the network. This value define also the network IP range with the prefix length.
- `guest_vlan_allowed` (Boolean) Indicates if the network allows guest VLANs.
- `id` (String) The ID of the isolated network.
- `prefix_length` (Number) The prefix length for the network. This value must be a valid prefix length for the network IP range. (e.g. /24 for netmask 255.255.255.0). For more information, see [CIDR notation](https://en.wikipedia.org/wiki/Classless_Inter-Domain_Routing).
- `static_ip_pool` (Attributes Set) A set of static IP pools to be used for this network. (see [below for nested schema](#nestedatt--static_ip_pool))

<a id="nestedatt--static_ip_pool"></a>
### Nested Schema for `static_ip_pool`

Read-Only:

- `end_address` (String) The end address of the IP pool. This value must be a valid IP address in the network IP range.
- `start_address` (String) The start address of the IP pool. This value must be a valid IP address in the network IP range.
