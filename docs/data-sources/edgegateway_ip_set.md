---
page_title: "cloudavenue_edgegateway_ip_set Data Source - cloudavenue"
subcategory: "Edge Gateway (Tier-1)"
description: |-
  The cloudavenue_edgegateway_ip_set data source allows you to retrieve information about an IP Set rule on an Edge Gateway.
---

# cloudavenue_edgegateway_ip_set (Data Source)

The `cloudavenue_edgegateway_ip_set` data source allows you to retrieve information about an IP Set rule on an Edge Gateway.

## Example Usage

```terraform
data "cloudavenue_edgegateway_ip_set" "example" {
  name            = "example"
  edge_gateway_id = cloudavenue_edgegateway.example.id
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The name of the IP Set.

### Optional

- `edge_gateway_id` (String) The ID of the Edge Gateway. Ensure that one and only one attribute from this collection is set : `edge_gateway_name`, `edge_gateway_id`.
- `edge_gateway_name` (String) The name of the Edge Gateway. Ensure that one and only one attribute from this collection is set : `edge_gateway_name`, `edge_gateway_id`.

### Read-Only

- `description` (String) The description of the IP Set.
- `id` (String) The ID of the IP Set.
- `ip_addresses` (Set of String) A set of IP address, CIDR or IP range.

