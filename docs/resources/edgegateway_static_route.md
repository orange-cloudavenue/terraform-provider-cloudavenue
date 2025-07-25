---
page_title: "cloudavenue_edgegateway_static_route Resource - cloudavenue"
subcategory: "Edge Gateway (Tier-1)"
description: |-
  The cloudavenue_edgegateway_static_route resource allows you to create and manage static routes on an Edge Gateway.
---

# cloudavenue_edgegateway_static_route (Resource)

The `cloudavenue_edgegateway_static_route` resource allows you to create and manage static routes on an Edge Gateway.

## Example Usage

Example Simple Usage (Static Route with default next hop)
```terraform
resource "cloudavenue_edgegateway_static_route" "example" {
  name        = "example"
  description = "example description"

  edge_gateway_id = cloudavenue_edgegateway.example.id

  network_cidr = "192.168.2.0/24"
  next_hops = [
    {
      ip_address = "192.168.2.254"
    }
  ]
}
```

Example Advanced Usage (Static Route with 2 next hops)
```terraform
resource "cloudavenue_edgegateway_static_route" "example" {
  name        = "example"
  description = "example description"

  edge_gateway_id = cloudavenue_edgegateway.example.id

  network_cidr = "192.168.2.0/24"
  next_hops = [
    {
      ip_address = "192.168.2.254"
    },
    {
      ip_address     = "192.168.2.253"
      admin_distance = 2
    }
  ]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The name of the Static Route.
- `network_cidr` (String) The network CIDR of the Static Route. (e.g. 192.168.1.0/24).
- `next_hops` (Attributes Set) A set of next hops to use within the static route. Set must contain at least 1 elements. (see [below for nested schema](#nestedatt--next_hops))

### Optional

- `description` (String) The description of the Static Route.
- `edge_gateway_id` (String) <i style="color:red;font-weight: bold">(ForceNew)</i> The ID of the Edge Gateway. Ensure that one and only one attribute from this collection is set : `edge_gateway_name`, `edge_gateway_id`.
- `edge_gateway_name` (String) <i style="color:red;font-weight: bold">(ForceNew)</i> The name of the Edge Gateway. Ensure that one and only one attribute from this collection is set : `edge_gateway_name`, `edge_gateway_id`.

### Read-Only

- `id` (String) The ID of the Static Route.

<a id="nestedatt--next_hops"></a>
### Nested Schema for `next_hops`

Required:

- `ip_address` (String) IP address for next hop gateway IP Address for the Static Route. Must be a valid IP with net.ParseIP.

Optional:

- `admin_distance` (Number) Admin distance is used to choose which route to use when there are multiple routes for a specific network. The lower the admin distance, the higher the preference for the route. Value defaults to `1`. Value must be at least 1.

## Import

Import is supported using the following syntax:
```shell
terraform import cloudavenue_edgegateway_static_route.example edgeGatewayIDOrName.staticRouteIDOrName
```