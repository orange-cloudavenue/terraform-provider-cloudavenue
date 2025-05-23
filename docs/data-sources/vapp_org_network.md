---
page_title: "cloudavenue_vapp_org_network Data Source - cloudavenue"
subcategory: "vApp (Virtual Appliance)"
description: |-
  Provides a Cloud Avenue routed vAPP Org Network data source to read data or reference existing network.
---

# cloudavenue_vapp_org_network (Data Source)

Provides a Cloud Avenue routed vAPP Org Network data source to read data or reference existing network.

## Example Usage

```terraform
data "cloudavenue_vapp_org_network" "example" {
  vapp_name    = cloudavenue_vapp.example.name
  network_name = cloudavenue_edgegateway_network_routed.example.name
  vdc          = cloudavenue_vapp.example.vdc
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `network_name` (String) Organization network name to which vApp network is connected to.

### Optional

- `vapp_id` (String) ID of the vApp. Ensure that one and only one attribute from this collection is set : `vapp_name`, `vapp_id`.
- `vapp_name` (String) Name of the vApp. Ensure that one and only one attribute from this collection is set : `vapp_id`, `vapp_name`.
- `vdc` (String) The name of vDC to use, optional if defined at provider level.

### Read-Only

- `id` (String) The ID of the network.

