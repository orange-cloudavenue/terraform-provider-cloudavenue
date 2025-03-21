---
page_title: "cloudavenue_vapp_org_network Resource - cloudavenue"
subcategory: "vApp (Virtual Appliance)"
description: |-
  Provides a Cloud Avenue routed vAPP Org Network resource. This can be used to create, modify, and delete routed vAPP Network.
---

# cloudavenue_vapp_org_network (Resource)

Provides a Cloud Avenue routed vAPP Org Network resource. This can be used to create, modify, and delete routed vAPP Network.

!> **Warning on deleting resource:** Deleting a resource require **vApp to be in a powered OFF** state. 
If the vApp is in a powered on state, the resource will power OFF the vApp before deleting the resource and then power it back on.
On power **ALL** VMs in the vApp will be powered ON, regardless of their previous state.

## Example Usage

```terraform
data "cloudavenue_tier0_vrfs" "example" {}

resource "cloudavenue_edgegateway" "example" {
  owner_name     = "MyVDC"
  tier0_vrf_name = data.cloudavenue_tier0_vrfs.example.names.0
  owner_type     = "vdc"
}

resource "cloudavenue_network_routed" "example" {
  name        = "MyOrgNet"
  description = "This is an example Net"

  edge_gateway_id = cloudavenue_edgegateway.example.id

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

resource "cloudavenue_vapp" "example" {
  name        = "MyVapp"
  description = "This is an example vApp"
  vdc         = "MyVDC"
}

resource "cloudavenue_vapp_org_network" "example" {
  vapp_name    = cloudavenue_vapp.example.name
  network_name = cloudavenue_network_routed.example.name
  vdc          = "MyVDC"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `network_name` (String) <i style="color:red;font-weight: bold">(ForceNew)</i> Organization network name to which vApp network is connected to.

### Optional

- `vapp_id` (String) <i style="color:red;font-weight: bold">(ForceNew)</i> ID of the vApp. Ensure that one and only one attribute from this collection is set : `vapp_name`, `vapp_id`.
- `vapp_name` (String) <i style="color:red;font-weight: bold">(ForceNew)</i> Name of the vApp. Ensure that one and only one attribute from this collection is set : `vapp_id`, `vapp_name`.
- `vdc` (String) <i style="color:red;font-weight: bold">(ForceNew)</i> The name of vDC to use, optional if defined at provider level.

### Read-Only

- `id` (String) The ID of the network.

## Import

Import is supported using the following syntax:
```shell
# if vdc is not specified, the default vdc will be used
terraform import cloudavenue_vapp_org_network.example vapp_name.network_name

# if vdc is specified, the vdc will be used
terraform import cloudavenue_vapp_org_network.example vdc.vapp_name.network_name
```