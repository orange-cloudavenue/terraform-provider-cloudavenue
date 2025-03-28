---
page_title: "cloudavenue_edgegateway_ip_set Resource - cloudavenue"
subcategory: "Edge Gateway (Tier-1)"
description: |-
  The cloudavenue_edgegateway_ip_set resource allows you to manage an IP Set rule on an Edge Gateway. IP Sets are groups of objects to which the firewall rules apply. Combining multiple objects into IP Sets helps reduce the total number of firewall rules to be created.
---

# cloudavenue_edgegateway_ip_set (Resource)

The `cloudavenue_edgegateway_ip_set` resource allows you to manage an IP Set rule on an Edge Gateway. IP Sets are groups of objects to which the firewall rules apply. Combining multiple objects into IP Sets helps reduce the total number of firewall rules to be created.

## Example Usage

```terraform
resource "cloudavenue_edgegateway_ip_set" "example" {
  name        = "example"
  description = "example of ip set"
  ip_addresses = [
    "12.12.12.1",            # IP Address
    "10.10.10.0/24",         # IP Address With CIDR
    "11.11.11.1-11.11.11.2", # IP Address Range
  ]
  edge_gateway_name = "myEdgeName"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The name of the IP Set.

### Optional

- `description` (String) The description of the IP Set.
- `edge_gateway_id` (String) <i style="color:red;font-weight: bold">(ForceNew)</i> The ID of the Edge Gateway. Ensure that one and only one attribute from this collection is set : `edge_gateway_name`, `edge_gateway_id`.
- `edge_gateway_name` (String) <i style="color:red;font-weight: bold">(ForceNew)</i> The name of the Edge Gateway. Ensure that one and only one attribute from this collection is set : `edge_gateway_name`, `edge_gateway_id`.
- `ip_addresses` (Set of String) A set of IP address, CIDR or IP range.

### Read-Only

- `id` (String) The ID of the IP Set.

## Import

Import is supported using the following syntax:
```shell
terraform import cloudavenue_edgegateway_ip_set.example edgeGatewayIDOrName.ipSetName
```