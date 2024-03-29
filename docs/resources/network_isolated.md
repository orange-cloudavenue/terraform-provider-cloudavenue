---
page_title: "cloudavenue_network_isolated Resource - cloudavenue"
subcategory: "Network"
description: |-
  Provides a Cloud Avenue VDC isolated Network. This can be used to create, modify, and delete VDC isolated networks.
---

# cloudavenue_network_isolated (Resource)

Provides a Cloud Avenue VDC isolated Network. This can be used to create, modify, and delete VDC isolated networks.

## Example Usage

```terraform
resource "cloudavenue_network_isolated" "example" {
  vdc         = "VDC_Test"
  name        = "rsx-example-isolated-network"
  description = "My isolated Org VDC network"

  gateway       = "1.1.1.1"
  prefix_length = 24

  dns1       = "8.8.8.8"
  dns2       = "8.8.4.4"
  dns_suffix = "example.com"

  static_ip_pool = [
    {
      start_address = "1.1.1.10"
      end_address   = "1.1.1.20"
    },
    {
      start_address = "1.1.1.100"
      end_address   = "1.1.1.103"
    }
  ]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `gateway` (String) (ForceNew) The gateway IP address for the network. This value define also the network IP range with the prefix length. Must be a valid IP with net.ParseIP.
- `name` (String) The name of the network. This value must be unique within the `VDC` or `VDC Group` that owns the network.
- `prefix_length` (Number) (ForceNew) The prefix length for the network. This value must be a valid prefix length for the network IP range. (e.g. /24 for netmask 255.255.255.0). Value must be between 1 and 32.

### Optional

- `description` (String) A description of the network.
- `dns1` (String) The primary DNS server IP address for the network. Must be a valid IP with net.ParseIP.
- `dns2` (String) The secondary DNS server IP address for the network. Must be a valid IP with net.ParseIP.
- `dns_suffix` (String) The DNS suffix for the network.
- `static_ip_pool` (Attributes Set) A set of static IP pools to be used for this network. Set must contain at least 1 elements. (see [below for nested schema](#nestedatt--static_ip_pool))
- `vdc` (String) (ForceNew) The name of vDC to use, optional if defined at provider level.

### Read-Only

- `id` (String) The ID of the network.

<a id="nestedatt--static_ip_pool"></a>
### Nested Schema for `static_ip_pool`

Required:

- `end_address` (String) The end address of the IP pool. This value must be a valid IP address in the network IP range. Must be a valid IP with net.ParseIP.
- `start_address` (String) The start address of the IP pool. This value must be a valid IP address in the network IP range. Must be a valid IP with net.ParseIP.

## Import

Import is supported using the following syntax:
```shell
terraform import cloudavenue_network_isolated.example vdc-or-vdc-group-name.NetworkName
```
