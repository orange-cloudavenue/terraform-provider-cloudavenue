---
page_title: "cloudavenue_edgegateway Data Source - cloudavenue"
subcategory: "Edge Gateway (Tier-1)"
description: |-
  The Edge Gateway data source allows you to show the details of an Edge Gateways in Cloud Avenue.
---

# cloudavenue_edgegateway (Data Source)

The Edge Gateway data source allows you to show the details of an Edge Gateways in Cloud Avenue.

## Example Usage

```terraform
data "cloudavenue_edgegateway" "example" {
  name = "myEdgeName"
}

output "gateway" {
  value = data.cloudavenue_edgegateway.example
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The name of the Edge Gateway.

### Read-Only

- `bandwidth` (Number) The bandwidth in `Mbps` of the Edge Gateway.
- `description` (String) The description of the Edge Gateway.
- `id` (String) The ID of the Edge Gateway.
- `owner_name` (String) The name of the Edge Gateway owner. It can be a VDC or a VDC Group name.
- `tier0_vrf_name` (String) The name of the Tier-0 VRF to which the Edge Gateway is attached.

