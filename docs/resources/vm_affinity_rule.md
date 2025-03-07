---
page_title: "cloudavenue_vm_affinity_rule Resource - cloudavenue"
subcategory: "VM (Virtual Machine)"
description: |-
  Provides a Cloud Avenue VM Affinity Rule. This can be used to create, modify and delete VM affinity and anti-affinity rules.
---

# cloudavenue_vm_affinity_rule (Resource)

Provides a Cloud Avenue VM Affinity Rule. This can be used to create, modify and delete VM affinity and anti-affinity rules.

~> **NOTE:** The CloudAvenue UI defines two different entities (`Affinity Rules` and `Anti-Affinity Rules`). This resource combines both entities: they are differentiated by the `polarity` property (see below).

## Example Usage

```terraform
resource "cloudavenue_vm_affinity_rule" "example" {
  name     = "example-affinity-rule"
  polarity = "Affinity"

  vm_ids = [
    cloudavenue_vm.example.id,
    cloudavenue_vm.example2.id,
  ]
}

resource "cloudavenue_vm" "example" {
  name      = "example-vm"
  vapp_name = cloudavenue_vapp.example.name
  deploy_os = {
    vapp_template_id = data.cloudavenue_catalog_vapp_template.example.id
  }
  settings = {
    customization = {
      auto_generate_password = true
    }
  }
  resource = {
  }

  state = {
  }
}

resource "cloudavenue_vm" "example2" {
  name      = "example-vm2"
  vapp_name = cloudavenue_vapp.example.name
  deploy_os = {
    vapp_template_id = data.cloudavenue_catalog_vapp_template.example.id
  }
  settings = {
    customization = {
      auto_generate_password = true
    }
  }
  resource = {
  }

  state = {
  }
}

data "cloudavenue_catalog_vapp_template" "example" {
  catalog_name  = "Orange-Linux"
  template_name = "debian_10_X64"
}

resource "cloudavenue_vapp" "example" {
  name        = "vapp_example"
  description = "This is a example vapp"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) VM affinity rule name.
- `polarity` (String) <i style="color:red;font-weight: bold">(ForceNew)</i> The polarity of the affinity rule. Value must be one of : `Affinity`, `Anti-Affinity`.
- `vm_ids` (Set of String) List of VM IDs to apply the affinity rule to. Set must contain at least 2 elements. Element value must satisfy all validations: must be a valid URN.

### Optional

- `enabled` (Boolean) `True` if this affinity rule is enabled. Value defaults to `true`.
- `required` (Boolean) `True` if this affinity rule is required. When a rule is mandatory, a host failover will not power on the VM if doing so would violate the rule. Value defaults to `true`.
- `vdc` (String) <i style="color:red;font-weight: bold">(ForceNew)</i> The name of vDC to use, optional if defined at provider level.

### Read-Only

- `id` (String) The ID of the affinity rule.

## Import

Import is supported using the following syntax:
```shell
# If `vDC` is not specified, the default `vDC` will be used
# The `affinityRuleIdentifier` can be either a name or an ID. If it is a name, it will succeed only if the name is unique.
terraform import cloudavenue_vm_affinity_rule.example affinityRuleIdentifier

# or you can specify the vDC
terraform import cloudavenue_vm_affinity_rule.example myVDC.affinityRuleIdentifier
```
