---
page_title: "cloudavenue_vdc Data Source - cloudavenue"
subcategory: "vDC (Virtual Datacenter)"
description: |-
  Provides a Cloud Avenue vDC (Virtual Data Center) data source. This can be used to reference a vDC and use its data within other resources or data sources.
---

# cloudavenue_vdc (Data Source)

Provides a Cloud Avenue vDC (Virtual Data Center) data source. This can be used to reference a vDC and use its data within other resources or data sources.

 -> Note: For more information about Cloud Avenue vDC, please refer to the [Cloud Avenue documentation](https://cloud.orange-business.com/en/offres/infrastructure-iaas/cloud-avenue/wiki-cloud-avenue/services/virtual-datacenter/virtual-datacenter/).

## Example Usage

```terraform
data "cloudavenue_vdc" "example" {
  name = "VDC_Example"
}

output "example" {
  value = data.cloudavenue_vdc.example
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The name of the vDC. String length must be between 2 and 27.

### Read-Only

- `billing_model` (String) Choose Billing model of compute resources.
- `cpu_allocated` (Number) CPU capacity in *MHz* that is committed to be available or used as a limit in PAYG mode.
- `cpu_speed_in_mhz` (Number) Specifies the clock frequency, in Mhz, for any virtual CPU that is allocated to a VM.
- `description` (String) A description of the vDC.
- `disponibility_class` (String) The disponibility class of the vDC.
- `id` (String) The ID of the vDC.
- `memory_allocated` (Number) Memory capacity in Gb that is committed to be available or used as a limit in PAYG mode.
- `service_class` (String) The service class of the vDC.
- `storage_billing_model` (String) Choose Billing model of storage resources.
- `storage_profiles` (Attributes Set) List of storage profiles for this vDC. (see [below for nested schema](#nestedatt--storage_profiles))

<a id="nestedatt--storage_profiles"></a>
### Nested Schema for `storage_profiles`

Read-Only:

- `class` (String) The storage class of the storage profile.
- `default` (Boolean) Set this storage profile as default for this vDC. Only one storage profile can be default per vDC.
- `limit` (Number) Max number in *Go* of units allocated for this storage profile.

