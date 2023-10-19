---
page_title: "cloudavenue_backup Resource - cloudavenue"
subcategory: "Backup"
description: |-
  The cloudavenue_backup resource allows you to manage backup strategy for vdc, vapp and vm from NetBackup solution. Please refer to the documentation for more information. https://wiki.cloudavenue.orange-business.com/wiki/Backup
---

# cloudavenue_backup (Resource)

The `cloudavenue_backup` resource allows you to manage backup strategy for `vdc`, `vapp` and `vm` from NetBackup solution. [Please refer to the documentation for more information.](https://wiki.cloudavenue.orange-business.com/wiki/Backup)

 ~> The credentials NetBackup are Requires to use this feature. [Please refer to the documentation for more information.](https://registry.terraform.io/providers/orange-cloudavenue/cloudavenue/latest/docs#netbackup-configuration)

## Examples
### Example Usage of a VDC Backup with 2 policy sets
```hcl
resource "cloudavenue_backup" "example-vdc" {
  type = "vdc"
  target_name = cloudavenue_vdc.example.name
  policies = [{
      policy_name = "D6"
    },
    {
      policy_name = "M3"
    }
  ]
}
```

### Example Usage of a VAPP Backup with a policy set
```hcl
resource "cloudavenue_backup" "example-vapp" {
  type = "vapp"
  target_name = cloudavenue_vapp.example.name
  policies = [{
      policy_name = "D6"
    }]
}
```

### Example Usage of a VM Backup with a policy set
```hcl
resource "cloudavenue_backup" "example-vm" {
  type = "vm"
  target_name = cloudavenue_vm.example.name
  policies = [{
      policy_name = "D6"
    }]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `policies` (Attributes Set) The backup policies of the target. Set must contain at least 1 elements. (see [below for nested schema](#nestedatt--policies))
- `type` (String) (ForceNew) Scope of the backup. Value must be one of : `vdc`, `vapp`, `vm`.

### Optional

- `target_id` (String) (ForceNew) The URN of the target. A target can be a VDC, a VApp or a VM. Ensure that one and only one attribute from this collection is set : `target_id`, `target_name`. Must be a valid URN.
- `target_name` (String) (ForceNew) The name of the target. A target can be a VDC, a VApp or a VM. Ensure that one and only one attribute from this collection is set : `target_id`, `target_name`.

### Read-Only

- `id` (Number) The ID of the backup.

<a id="nestedatt--policies"></a>
### Nested Schema for `policies`

Required:

- `policy_name` (String) The name of the backup policy. Each letter represent a strategy predefined: D = Daily, W = Weekly, M = Monthly, X = Replication, The number is the retention period. [Please refer to the documentation for more information.](https://wiki.cloudavenue.orange-business.com/wiki/Backup). Value must be one of : `D6`, `D30`, `D30NQ`, `D60`, `W4`, `M3`, `M12`, `XD6`, `XD30`, `XD60`, `XW4`, `XM3`, `XM12`.

Read-Only:

- `policy_id` (Number) The ID of the backup policy.

## Import

Import is supported using the following syntax:
```shell
terraform import cloudavenue_backup.example type.targetName
```