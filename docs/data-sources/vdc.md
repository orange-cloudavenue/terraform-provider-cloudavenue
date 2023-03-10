---
page_title: "cloudavenue_vdc Data Source - cloudavenue"
subcategory: "vDC (Virtual Datacenter)"
description: |-
  Provides a Cloud Avenue Organization vDC data source. An Organization VDC can be used to reference a vDC and use its data within other resources or data sources.
---

# cloudavenue_vdc (Data Source)

Provides a Cloud Avenue Organization vDC data source. An Organization VDC can be used to reference a vDC and use its data within other resources or data sources.

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

- `name` (String) The name of the org vDC. It must be unique in the organization.
The length must be between 2 and 27 characters.

### Read-Only

- `billing_model` (String) Choose Billing model of compute resources. It can be `PAYG`, `DRAAS` or `RESERVED`.
- `cpu_allocated` (Number) CPU capacity in *MHz* that is committed to be available or used as a limit in PAYG mode.
It must be at least 5 * `cpu_speed_in_mhz` and at most 200 * `cpu_speed_in_mhz`.

 -> Note: Reserved capacity is automatically set according to the service class.
- `cpu_speed_in_mhz` (Number) Specifies the clock frequency, in Mhz, for any virtual CPU that is allocated to a VM.
It must be at least 1200.
- `description` (String) The description of the org vDC.
- `disponibility_class` (String) The disponibility class of the org vDC. It can be `ONE-ROOM`, `DUAL-ROOM` or `HA-DUAL-ROOM`.
- `id` (String) The ID of the vDC.
- `memory_allocated` (Number) Memory capacity in Gb that is committed to be available or used as a limit in PAYG mode.
It must be between 1 and 5000.
- `service_class` (String) The service class of the org vDC. It can be `ECO`, `STD`, `HP` or `VOIP`.
- `storage_billing_model` (String) Choose Billing model of storage resources. It can be `PAYG` or `RESERVED`.
- `storage_profiles` (Attributes List) List of storage profiles for this vDC. (see [below for nested schema](#nestedatt--storage_profiles))
- `vdc_group` (String) Name of an existing vDC group or a new one. This allows you to isolate your VDC.
VMs of vDCs which belong to the same vDC group can communicate together.

<a id="nestedatt--storage_profiles"></a>
### Nested Schema for `storage_profiles`

Read-Only:

- `class` (String) The storage class of the storage profile.
It can be `silver`, `silver_r1`, `silver_r2`, `gold`, `gold_r1`, `gold_r2`, `gold_hm`, `platinum3k`, `platinum3k_r1`, `platinum3k_r2`, `platinum3k_hm`, `platinum7k`, `platinum7k_r1`, `platinum7k_r2`, `platinum7k_hm`.
- `default` (Boolean) Set this storage profile as default for this vDC. Only one storage profile can be default per vDC.
- `limit` (Number) Max number of units allocated for this storage profile. In Gb. It must be between 500 and 10000.

