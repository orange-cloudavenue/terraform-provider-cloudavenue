---
page_title: "cloudavenue_network_routed Data Source - cloudavenue"
subcategory: "Network"
description: |-
  Provides a Cloud Avenue vDC routed Network data source to read data or reference existing network
  !> Resource deprecated The resource has renamed to cloudavenue_edgegateway_network_routed https://registry.terraform.io/providers/orange-cloudavenue/cloudavenue/latest/docs/data-sources/edgegateway_network_routed, it will be removed in the version v0.38.0 https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/milestone/21 of the provider. See the GitHub issue https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/1020 for more information.
---

# cloudavenue_network_routed (Data Source)

Provides a Cloud Avenue vDC routed Network data source to read data or reference existing network 

 !> **Resource deprecated** The resource has renamed to [`cloudavenue_edgegateway_network_routed`](https://registry.terraform.io/providers/orange-cloudavenue/cloudavenue/latest/docs/data-sources/edgegateway_network_routed), it will be removed in the version [`v0.38.0`](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/milestone/21) of the provider. See the [GitHub issue](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/1020) for more information.

## Example Usage

```terraform
data "cloudavenue_edgegateway" "example" {
  name = "tn01e02ocb0006205spt101"
}

resource "cloudavenue_network_routed" "example" {
  name            = "ExampleNetworkRouted"
  gateway         = "192.168.10.254"
  prefix_length   = 24
  edge_gateway_id = data.cloudavenue_edgegateway.example.id
  dns1            = "1.1.1.1"
  dns2            = "8.8.8.8"

  dns_suffix = "example"

  static_ip_pool = [
    {
      start_address = "192.168.10.10"
      end_address   = "192.168.10.20"
    }
  ]
}

data "cloudavenue_network_routed" "example" {
  name            = cloudavenue_network_routed.example.name
  edge_gateway_id = cloudavenue_network_routed.example.edge_gateway_id
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The name of the network. This value must be unique within the `VDC` or `VDC Group` that owns the network.

### Optional

- `edge_gateway_id` (String) The ID of the edge gateway in which the routed network should be located.

### Read-Only

- `description` (String) A description of the network.
- `dns1` (String) The primary DNS server IP address for the network.
- `dns2` (String) The secondary DNS server IP address for the network.
- `dns_suffix` (String) The DNS suffix for the network.
- `edge_gateway_name` (String) The name of the edge gateway in which the routed network should be located.
- `gateway` (String) The gateway IP address for the network. This value define also the network IP range with the prefix length.
- `id` (String) The ID of the network.
- `interface_type` (String) An interface for the network.
- `prefix_length` (Number) The prefix length for the network. This value must be a valid prefix length for the network IP range. (e.g. /24 for netmask 255.255.255.0).
- `static_ip_pool` (Attributes Set) A set of static IP pools to be used for this network. (see [below for nested schema](#nestedatt--static_ip_pool))

<a id="nestedatt--static_ip_pool"></a>
### Nested Schema for `static_ip_pool`

Read-Only:

- `end_address` (String) The end address of the IP pool. This value must be a valid IP address in the network IP range.
- `start_address` (String) The start address of the IP pool. This value must be a valid IP address in the network IP range.

