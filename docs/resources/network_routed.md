---
page_title: "cloudavenue_network_routed Resource - cloudavenue"
subcategory: "Network"
description: |-
  Provides a CloudAvenue Org VDC routed Network. This can be used to create, modify, and delete routed VDC networks.
---

# cloudavenue_network_routed (Resource)

Provides a CloudAvenue Org VDC routed Network. This can be used to create, modify, and delete routed VDC networks.

## Example Usage

```terraform
data "cloudavenue_edgegateway" "example" {
  name = "tn01e02ocb0006205spt101"
}

resource "cloudavenue_network_routed" "example" {
  name        = "OrgNetExample"
  description = "Org Net Example"

  edge_gateway_id = data.cloudavenue_edgegateway.example.id

  gateway       = "192.168.1.254"
  prefix_length = 24

  dns1 = "1.1.1.1"
  dns2 = "8.8.8.8"

  dns_suffix = "example"

  static_ip_pool = [
    {
      start_address = "192.168.1.10"
      end_address   = "192.168.1.20"
    }
  ]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `edge_gateway_id` (String) Edge gateway ID in which Routed network should be located.
- `gateway` (String) Gateway IP address.
- `name` (String) Network name.
- `prefix_length` (Number) Network prefix length.

### Optional

- `description` (String) Network description.
- `dns1` (String) DNS server 1.
- `dns2` (String) DNS server 2.
- `dns_suffix` (String) DNS suffix.
- `interface_type` (String) Optional interface type (only for NSX-V networks). One of `INTERNAL` (default), `DISTRIBUTED`, `SUBINTERFACE`
- `static_ip_pool` (Attributes Set) IP ranges used for static pool allocation in the network. (see [below for nested schema](#nestedatt--static_ip_pool))

### Read-Only

- `id` (String) The ID of the routed network.

<a id="nestedatt--static_ip_pool"></a>
### Nested Schema for `static_ip_pool`

Required:

- `end_address` (String) End address of the IP range.
- `start_address` (String) Start address of the IP range.

## Import

Import is supported using the following syntax:
```shell
terraform import vdc-or-vdc-group-name.NetworkName
```
