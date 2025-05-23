---
page_title: "cloudavenue_vdc_acl Resource - cloudavenue"
subcategory: "vDC (Virtual Datacenter)"
description: |-
  Provides a Cloud Avenue vDC access control resource. This can be used to share vDC across users and/or groups.
---

# cloudavenue_vdc_acl (Resource)

Provides a Cloud Avenue vDC access control resource. This can be used to share vDC across users and/or groups.

## Example Usage

```terraform
resource "cloudavenue_vdc_acl" "example" {
  vdc                   = cloudavenue_vdc.example.name
  everyone_access_level = "ReadOnly"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `everyone_access_level` (String) Access level when the vApp is shared with everyone. Ensure that one and only one attribute from this collection is set : `shared_with`, `everyone_access_level`.
- `shared_with` (Attributes Set) One or more blocks defining the subjects with whom we are sharing. Ensure that one and only one attribute from this collection is set : `everyone_access_level`, `shared_with`. (see [below for nested schema](#nestedatt--shared_with))
- `vdc` (String) <i style="color:red;font-weight: bold">(ForceNew)</i> The name of vDC to use, optional if defined at provider level.

### Read-Only

- `id` (String) The ID of the acl rule.

<a id="nestedatt--shared_with"></a>
### Nested Schema for `shared_with`

Required:

- `access_level` (String) Access level for the user or group with whom we are sharing. Value must be one of : `ReadOnly`.

Optional:

- `group_id` (String) ID of the group with whom we are sharing. Ensure that one and only one attribute from this collection is set : `user_id`.
- `user_id` (String) ID of the user with whom we are sharing. Ensure that one and only one attribute from this collection is set : `group_id`.

Read-Only:

- `subject_name` (String) Name of the subject (group or user) with whom we are sharing.

## Import

Import is supported using the following syntax:
```shell
# use the vdc to import the resource
terraform import cloudavenue_vdc_acl.example vdc
```
