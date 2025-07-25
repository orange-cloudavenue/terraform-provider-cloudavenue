---
page_title: "cloudavenue_vdc Resource - cloudavenue"
subcategory: "vDC (Virtual Datacenter)"
description: |-
  Provides a Cloud Avenue vDC (Virtual Data Center) resource. This can be used to create, update and delete vDC.
---

# cloudavenue_vdc (Resource)

Provides a Cloud Avenue vDC (Virtual Data Center) resource. This can be used to create, update and delete vDC.
 
 -> Note: For more information about Cloud Avenue vDC, please refer to the [Cloud Avenue documentation](https://cloud.orange-business.com/en/offres/infrastructure-iaas/cloud-avenue/wiki-cloud-avenue/services/virtual-datacenter/virtual-datacenter/).

 ~> **Warning**
 The VDC resource uses a complex validation system that is **incompatible** with the **Terraform module**. (See [Disable validation](#disable-validation))

## Example Usage

```terraform
resource "cloudavenue_vdc" "example" {
  name                  = "MyVDC"
  description           = "Example VDC created by Terraform"
  cpu_allocated         = 22000
  memory_allocated      = 30
  cpu_speed_in_mhz      = 2200
  billing_model         = "PAYG"
  disponibility_class   = "ONE-ROOM"
  service_class         = "STD"
  storage_billing_model = "PAYG"

  storage_profiles = [
    {
      class   = "gold"
      default = true
      limit   = 500
    },
    {
      class   = "silver"
      default = false
      limit   = 500
    },
  ]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `billing_model` (String) <i style="color:red;font-weight: bold">(ForceNew)</i> Choose Billing model of compute resources. The billing model available are different depending on the service class. See [Rules](https://registry.terraform.io/providers/orange-cloudavenue/cloudavenue/latest/docs/resources/vdc#rules) for more information. Value must be one of : `RESERVED`, `PAYG`, `DRAAS`.
- `cpu_allocated` (Number) CPU capacity in *MHz* that is committed to be available or used as a limit in PAYG mode. See [Rules](https://registry.terraform.io/providers/orange-cloudavenue/cloudavenue/latest/docs/resources/vdc#rules) for more information.
- `cpu_speed_in_mhz` (Number) Specifies the clock frequency, in Mhz, for any virtual CPU that is allocated to a VM. See [Rules](https://registry.terraform.io/providers/orange-cloudavenue/cloudavenue/latest/docs/resources/vdc#rules) for more information.
- `disponibility_class` (String) <i style="color:red;font-weight: bold">(ForceNew)</i> The disponibility class of the vDC. The disponibility class available are different depending on the service class. See [Rules](https://registry.terraform.io/providers/orange-cloudavenue/cloudavenue/latest/docs/resources/vdc#rules) for more information. Value must be one of : `ONE-ROOM`, `HA-DUAL-ROOM`, `DUAL-ROOM`.
- `memory_allocated` (Number) Memory capacity in Gb that is committed to be available or used as a limit in PAYG mode.
- `name` (String) <i style="color:red;font-weight: bold">(ForceNew)</i> The name of the vDC. String length must be between 2 and 27.
- `service_class` (String) <i style="color:red;font-weight: bold">(ForceNew)</i> The service class of the vDC. Value must be one of : `ECO`, `STD`, `HP`, `VOIP`.
- `storage_billing_model` (String) <i style="color:red;font-weight: bold">(ForceNew)</i> Choose Billing model of storage resources. The billing model available are different depending on the service class. See [Rules](https://registry.terraform.io/providers/orange-cloudavenue/cloudavenue/latest/docs/resources/vdc#rules) for more information. Value must be one of : `PAYG`, `RESERVED`.
- `storage_profiles` (Attributes Set) List of storage profiles for this vDC. Set must contain at least 1 elements. (see [below for nested schema](#nestedatt--storage_profiles))

### Optional

- `description` (String) A description of the vDC.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `id` (String) The ID of the vDC.

<a id="nestedatt--storage_profiles"></a>
### Nested Schema for `storage_profiles`

Required:

- `class` (String) The storage class of the storage profile. The storage class available are different depending on the service class. See [Rules](https://registry.terraform.io/providers/orange-cloudavenue/cloudavenue/latest/docs/resources/vdc#rules) for more information. Value must be one of `silver`, `silver_r1`, `silver_r2`, `gold`, `gold_r1`, `gold_r2`, `gold_hm`, `platinum3k`, `platinum3k_r1`, `platinum3k_r2`, `platinum3k_hm`, `platinum7k`, `platinum7k_r1`, `platinum7k_r2`, `platinum7k_hm` or custom storage profile class delivered by Cloud Avenue.
- `default` (Boolean) Set this storage profile as default for this vDC. Only one storage profile can be default per vDC.
- `limit` (Number) Max number in *Go* of units allocated for this storage profile.


<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

<!-- TABLE VDC ATTRIBUTES PARAMETERS -->
## Rules
More information about rules can be found [here](https://cloud.orange-business.com/en/offres/infrastructure-iaas/cloud-avenue/wiki-cloud-avenue/services/virtual-datacenter/virtual-datacenter/).All fields with a ** are editable.

### ServiceClass ECO
| BillingModels | StorageBillingModels | DisponibilityClasses | CPUInMhz (Mhz)          | CPUAllocated (Mhz)         | MemoryAllocated (Gb) |
| ------------- | -------------------- | -------------------- | ----------------------- | -------------------------- | -------------------- |
| RESERVED      | PAYG, RESERVED       | ONE-ROOM, DUAL-ROOM  | ** min: 1200, max: 2200 | ** min: 3000, max: 2500000 | ** min: 1, max: 5120 |
| PAYG          | PAYG, RESERVED       | ONE-ROOM, DUAL-ROOM  | equal: 2200             | ** min: 11000, max: 440000 | ** min: 1, max: 5120 |
| DRAAS         | PAYG, RESERVED       | ONE-ROOM, DUAL-ROOM  | equal: 2200             | ** min: 11000, max: 440000 | ** min: 1, max: 5120 |

### ServiceClass STD
| BillingModels | StorageBillingModels | DisponibilityClasses              | CPUInMhz (Mhz)          | CPUAllocated (Mhz)         | MemoryAllocated (Gb) |
| ------------- | -------------------- | --------------------------------- | ----------------------- | -------------------------- | -------------------- |
| RESERVED      | PAYG, RESERVED       | ONE-ROOM, HA-DUAL-ROOM, DUAL-ROOM | ** min: 1200, max: 2200 | ** min: 3000, max: 2500000 | ** min: 1, max: 5120 |
| PAYG          | PAYG, RESERVED       | ONE-ROOM, HA-DUAL-ROOM, DUAL-ROOM | equal: 2200             | ** min: 11000, max: 440000 | ** min: 1, max: 5120 |
| DRAAS         | PAYG, RESERVED       | ONE-ROOM, HA-DUAL-ROOM, DUAL-ROOM | equal: 2200             | ** min: 11000, max: 440000 | ** min: 1, max: 5120 |

### ServiceClass HP
| BillingModels | StorageBillingModels | DisponibilityClasses              | CPUInMhz (Mhz) | CPUAllocated (Mhz)         | MemoryAllocated (Gb) |
| ------------- | -------------------- | --------------------------------- | -------------- | -------------------------- | -------------------- |
| RESERVED      | PAYG, RESERVED       | ONE-ROOM, HA-DUAL-ROOM, DUAL-ROOM | equal: 2200    | ** min: 3000, max: 2500000 | ** min: 1, max: 5120 |
| PAYG          | PAYG, RESERVED       | ONE-ROOM, HA-DUAL-ROOM, DUAL-ROOM | equal: 2200    | ** min: 11000, max: 440000 | ** min: 1, max: 5120 |

### ServiceClass VOIP
| BillingModels | StorageBillingModels | DisponibilityClasses              | CPUInMhz (Mhz) | CPUAllocated (Mhz)         | MemoryAllocated (Gb) |
| ------------- | -------------------- | --------------------------------- | -------------- | -------------------------- | -------------------- |
| RESERVED      | PAYG, RESERVED       | ONE-ROOM, HA-DUAL-ROOM, DUAL-ROOM | equal: 3000    | ** min: 3000, max: 2500000 | ** min: 1, max: 5120 |


<!-- TABLE STORAGE PROFILES ATTRIBUTES PARAMETERS -->
## Storage Profiles
More information about storage profiles can be found [here](https://cloud.orange-business.com/en/offres/infrastructure-iaas/cloud-avenue/wiki-cloud-avenue/services/block-storage/).All fields with a ** are editable.

### ServiceClass ECO
| StorageProfileClass | SizeLimit (Go)          | IOPSLimit   | BillingModels  | DisponibilityClasses |
| ------------------- | ----------------------- | ----------- | -------------- | -------------------- |
| silver              | ** min: 100, max: 81920 | equal: 600  | PAYG, RESERVED | ONE-ROOM             |
| silver_r1           | ** min: 100, max: 81920 | equal: 600  | PAYG, RESERVED | DUAL-ROOM            |
| silver_r2           | ** min: 100, max: 81920 | equal: 600  | PAYG, RESERVED | DUAL-ROOM            |
| gold                | ** min: 100, max: 81920 | equal: 1000 | PAYG, RESERVED | ONE-ROOM             |
| gold_r1             | ** min: 100, max: 81920 | equal: 1000 | PAYG, RESERVED | DUAL-ROOM            |
| gold_r2             | ** min: 100, max: 81920 | equal: 1000 | PAYG, RESERVED | DUAL-ROOM            |

### ServiceClass STD
| StorageProfileClass | SizeLimit (Go)          | IOPSLimit   | BillingModels  | DisponibilityClasses |
| ------------------- | ----------------------- | ----------- | -------------- | -------------------- |
| silver              | ** min: 100, max: 81920 | equal: 600  | PAYG, RESERVED | ONE-ROOM             |
| silver_r1           | ** min: 100, max: 81920 | equal: 600  | PAYG, RESERVED | DUAL-ROOM            |
| silver_r2           | ** min: 100, max: 81920 | equal: 600  | PAYG, RESERVED | DUAL-ROOM            |
| gold                | ** min: 100, max: 81920 | equal: 1000 | PAYG, RESERVED | ONE-ROOM             |
| gold_r1             | ** min: 100, max: 81920 | equal: 1000 | PAYG, RESERVED | DUAL-ROOM            |
| gold_r2             | ** min: 100, max: 81920 | equal: 1000 | PAYG, RESERVED | DUAL-ROOM            |
| gold_hm             | ** min: 100, max: 81920 | equal: 1000 | PAYG, RESERVED | HA-DUAL-ROOM         |
| platinum3k          | ** min: 100, max: 81920 | equal: 3000 | PAYG, RESERVED | ONE-ROOM             |
| platinum3k_r1       | ** min: 100, max: 81920 | equal: 3000 | PAYG, RESERVED | DUAL-ROOM            |
| platinum3k_r2       | ** min: 100, max: 81920 | equal: 3000 | PAYG, RESERVED | DUAL-ROOM            |
| platinum3k_hm       | ** min: 100, max: 81920 | equal: 3000 | PAYG, RESERVED | HA-DUAL-ROOM         |
| platinum7k          | ** min: 100, max: 81920 | equal: 7000 | PAYG, RESERVED | ONE-ROOM             |
| platinum7k_r1       | ** min: 100, max: 81920 | equal: 7000 | PAYG, RESERVED | DUAL-ROOM            |
| platinum7k_r2       | ** min: 100, max: 81920 | equal: 7000 | PAYG, RESERVED | DUAL-ROOM            |
| platinum7k_hm       | ** min: 100, max: 81920 | equal: 7000 | PAYG, RESERVED | HA-DUAL-ROOM         |

### ServiceClass HP
| StorageProfileClass | SizeLimit (Go)          | IOPSLimit   | BillingModels  | DisponibilityClasses |
| ------------------- | ----------------------- | ----------- | -------------- | -------------------- |
| silver              | ** min: 100, max: 81920 | equal: 600  | PAYG, RESERVED | ONE-ROOM             |
| silver_r1           | ** min: 100, max: 81920 | equal: 600  | PAYG, RESERVED | DUAL-ROOM            |
| silver_r2           | ** min: 100, max: 81920 | equal: 600  | PAYG, RESERVED | DUAL-ROOM            |
| gold                | ** min: 100, max: 81920 | equal: 1000 | PAYG, RESERVED | ONE-ROOM             |
| gold_r1             | ** min: 100, max: 81920 | equal: 1000 | PAYG, RESERVED | DUAL-ROOM            |
| gold_r2             | ** min: 100, max: 81920 | equal: 1000 | PAYG, RESERVED | DUAL-ROOM            |
| gold_hm             | ** min: 100, max: 81920 | equal: 1000 | PAYG, RESERVED | HA-DUAL-ROOM         |
| platinum3k          | ** min: 100, max: 81920 | equal: 3000 | PAYG, RESERVED | ONE-ROOM             |
| platinum3k_r1       | ** min: 100, max: 81920 | equal: 3000 | PAYG, RESERVED | DUAL-ROOM            |
| platinum3k_r2       | ** min: 100, max: 81920 | equal: 3000 | PAYG, RESERVED | DUAL-ROOM            |
| platinum3k_hm       | ** min: 100, max: 81920 | equal: 3000 | PAYG, RESERVED | HA-DUAL-ROOM         |
| platinum7k          | ** min: 100, max: 81920 | equal: 7000 | PAYG, RESERVED | ONE-ROOM             |
| platinum7k_r1       | ** min: 100, max: 81920 | equal: 7000 | PAYG, RESERVED | DUAL-ROOM            |
| platinum7k_r2       | ** min: 100, max: 81920 | equal: 7000 | PAYG, RESERVED | DUAL-ROOM            |
| platinum7k_hm       | ** min: 100, max: 81920 | equal: 7000 | PAYG, RESERVED | HA-DUAL-ROOM         |

### ServiceClass VOIP
| StorageProfileClass | SizeLimit (Go)          | IOPSLimit   | BillingModels  | DisponibilityClasses |
| ------------------- | ----------------------- | ----------- | -------------- | -------------------- |
| silver              | ** min: 100, max: 81920 | equal: 600  | PAYG, RESERVED | ONE-ROOM             |
| silver_r1           | ** min: 100, max: 81920 | equal: 600  | PAYG, RESERVED | DUAL-ROOM            |
| silver_r2           | ** min: 100, max: 81920 | equal: 600  | PAYG, RESERVED | DUAL-ROOM            |
| gold                | ** min: 100, max: 81920 | equal: 1000 | PAYG, RESERVED | ONE-ROOM             |
| gold_r1             | ** min: 100, max: 81920 | equal: 1000 | PAYG, RESERVED | DUAL-ROOM            |
| gold_r2             | ** min: 100, max: 81920 | equal: 1000 | PAYG, RESERVED | DUAL-ROOM            |
| gold_hm             | ** min: 100, max: 81920 | equal: 1000 | PAYG, RESERVED | HA-DUAL-ROOM         |
| platinum3k          | ** min: 100, max: 81920 | equal: 3000 | PAYG, RESERVED | ONE-ROOM             |
| platinum3k_r1       | ** min: 100, max: 81920 | equal: 3000 | PAYG, RESERVED | DUAL-ROOM            |
| platinum3k_r2       | ** min: 100, max: 81920 | equal: 3000 | PAYG, RESERVED | DUAL-ROOM            |
| platinum3k_hm       | ** min: 100, max: 81920 | equal: 3000 | PAYG, RESERVED | HA-DUAL-ROOM         |
| platinum7k          | ** min: 100, max: 81920 | equal: 7000 | PAYG, RESERVED | ONE-ROOM             |
| platinum7k_r1       | ** min: 100, max: 81920 | equal: 7000 | PAYG, RESERVED | DUAL-ROOM            |
| platinum7k_r2       | ** min: 100, max: 81920 | equal: 7000 | PAYG, RESERVED | DUAL-ROOM            |
| platinum7k_hm       | ** min: 100, max: 81920 | equal: 7000 | PAYG, RESERVED | HA-DUAL-ROOM         |




## Disable validation

To disable the validation system, you can use the following environment variable:
```shell
export CLOUDAVENUE_VDC_VALIDATION=false
```

All checks will be skipped in the `terraform validate` sequence but will be running during the creation or an update of the resource. This is useful for terraform modules that are not compatible with the validation process.
The validation system is designed to ensure that the VDC resource is created with the correct parameters and configurations. However, in some cases, such as when using Terraform modules, the validation process may not be compatible.
The errors and warnings are returned during the creation of the resource, which can be confusing and time-consuming to troubleshoot.

The default value is `true`.

## Timeouts

The timeouts configuration allows you to specify the maximum amount of time that the provider will wait for a certain operation to complete. The following timeouts can be configured:

* `create` - 8 minutes.
* `update` - 8 minutes.
* `delete` - 5 minutes.

To configure the timeouts, use the following syntax:

```hcl
resource "cloudavenue_vdc" "example" {
  # ...
  timeouts {
    create = "10m"
    update = "10m"
    delete = "6m"
  }
}
```

## Import

Import is supported using the following syntax:
```shell
# VDC can be imported using the name.

terraform import cloudavenue_vdc.example name
```