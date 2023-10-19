---
page_title: "cloudavenue_backup Data Source - cloudavenue"
subcategory: "Backup"
description: |-
  The cloudavenue_backup data source allows you to retrieve information about a backup of NetBackup solution.
---

# cloudavenue_backup (Data Source)

The `cloudavenue_backup` data source allows you to retrieve information about a backup of NetBackup solution.

 ~> The credentials NetBackup are Requires to use this feature. [Please refer to the documentation for more information.](https://registry.terraform.io/providers/orange-cloudavenue/cloudavenue/latest/docs#netbackup-configuration)

## Example Usage

```terraform
data "cloudavenue_backup" "example" {
  type        = "vdc"
  target_name = data.cloudavenue_vdc.example.name
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `type` (String) Scope of the backup. Value must be one of : `vdc`, `vapp`, `vm`.

### Optional

- `id` (Number) The ID of the backup.
- `target_id` (String) The URN of the target. A target can be a VDC, a VApp or a VM. Ensure that one and only one attribute from this collection is set : `target_id`, `target_name`. Must be a valid URN.
- `target_name` (String) The name of the target. A target can be a VDC, a VApp or a VM. Ensure that one and only one attribute from this collection is set : `target_id`, `target_name`.

### Read-Only

- `policies` (Attributes Set) The backup policies of the target. (see [below for nested schema](#nestedatt--policies))

<a id="nestedatt--policies"></a>
### Nested Schema for `policies`

Read-Only:

- `policy_id` (Number) The ID of the backup policy.
- `policy_name` (String) The name of the backup policy. Each letter represent a strategy predefined: D = Daily, W = Weekly, M = Monthly, X = Replication, The number is the retention period. [Please refer to the documentation for more information.](https://wiki.cloudavenue.orange-business.com/wiki/Backup).
