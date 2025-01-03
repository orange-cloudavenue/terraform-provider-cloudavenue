---
page_title: "cloudavenue_vdc_group Resource - cloudavenue"
subcategory: "vDC (Virtual Datacenter)"
description: |-
  The cloudavenue_vdc_group resource allows you to manage VDC Group.
  !> Resource deprecated The resource has renamed to cloudavenue_vdcg https://registry.terraform.io/providers/orange-cloudavenue/cloudavenue/latest/docs/resources/vdcg, it will be removed in the version v0.30.0 https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/milestone/18 of the provider. See the GitHub issue https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/869 for more information.
---

# cloudavenue_vdc_group (Resource)

The `cloudavenue_vdc_group` resource allows you to manage VDC Group. 

 !> **Resource deprecated** The resource has renamed to [`cloudavenue_vdcg`](https://registry.terraform.io/providers/orange-cloudavenue/cloudavenue/latest/docs/resources/vdcg), it will be removed in the version [`v0.30.0`](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/milestone/18) of the provider. See the [GitHub issue](https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/869) for more information.

## How to migrate existing resources

Original configuration:

```terraform
resource "cloudavenue_vdc_group" "example" {
  name = "example"
  vdc_ids = [
    cloudavenue_vdc.example.id,
  ]
}
```

Migrated configuration:

Rename the resource to `cloudavenue_vdcg` and add the `moved` block to the configuration:

```hcl
resource "cloudavenue_vdcg" "example" {
  name = "example"
  vdc_ids = [
    cloudavenue_vdc.example.id,
  ]
}

moved {
  from = cloudavenue_vdc_group.example
  to   = cloudavenue_vdcg.example
}
```

Run `terraform plan` and `terraform apply` to migrate the resource.

Example of terraform plan output:

```shell
Terraform will perform the following actions:

  # cloudavenue_vdc_group.example has moved to cloudavenue_vdcg.example
    resource "cloudavenue_vdcg" "example" {
        id      = "urn:vcloud:vdcGroup:xxxx-xxxx-xxxxx-xxxxxx"
        name    = "example"
        # (3 unchanged attributes hidden)
    }

Plan: 0 to add, 0 to change, 0 to destroy.
```

## Example Usage

```terraform
resource "cloudavenue_vdc_group" "example" {
  name = "example"
  vdc_ids = [
    cloudavenue_vdc.example.id,
  ]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The name of the VDC Group.
- `vdc_ids` (Set of String) The list of VDC IDs of the VDC Group. Set must contain at least 1 elements.

### Optional

- `description` (String) The description of the VDC Group.

### Read-Only

- `id` (String) The ID of the VDC Group.
- `status` (String) The status of the VDC Group. Value must be one of : `SAVING`, `SAVED`, `CONFIGURING`, `REALIZED`, `REALIZATION_FAILED`, `DELETING`, `DELETE_FAILED`, `OBJECT_NOT_FOUND`, `UNCONFIGURED`.
- `type` (String) The type of the VDC Group. Value must be one of : `LOCAL`, `UNIVERSAL`.

## Import

Import is supported using the following syntax:
```shell
terraform import cloudavenue_vdc_group.example vdcGroupNameOrID
```