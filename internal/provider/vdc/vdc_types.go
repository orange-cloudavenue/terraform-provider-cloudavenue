package vdc

import (
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
)

type vdcDataSourceModel struct {
	ID                     types.String             `tfsdk:"id"`
	Name                   types.String             `tfsdk:"name"`
	Description            types.String             `tfsdk:"description"`
	VDCServiceClass        types.String             `tfsdk:"service_class"`
	VDCDisponibilityClass  types.String             `tfsdk:"disponibility_class"`
	VDCBillingModel        types.String             `tfsdk:"billing_model"`
	VcpuInMhz2             types.Float64            `tfsdk:"cpu_speed_in_mhz"`
	CPUAllocated           types.Float64            `tfsdk:"cpu_allocated"`
	MemoryAllocated        types.Float64            `tfsdk:"memory_allocated"`
	VDCStorageBillingModel types.String             `tfsdk:"storage_billing_model"`
	VDCStorageProfiles     []vdcStorageProfileModel `tfsdk:"storage_profiles"`
	VDCGroup               types.String             `tfsdk:"vdc_group"`
}

type vdcResourceModel struct {
	Timeouts               timeouts.Value           `tfsdk:"timeouts"`
	ID                     types.String             `tfsdk:"id"`
	Name                   types.String             `tfsdk:"name"`
	Description            types.String             `tfsdk:"description"`
	VDCServiceClass        types.String             `tfsdk:"service_class"`
	VDCDisponibilityClass  types.String             `tfsdk:"disponibility_class"`
	VDCBillingModel        types.String             `tfsdk:"billing_model"`
	VcpuInMhz2             types.Float64            `tfsdk:"cpu_speed_in_mhz"`
	CPUAllocated           types.Float64            `tfsdk:"cpu_allocated"`
	MemoryAllocated        types.Float64            `tfsdk:"memory_allocated"`
	VDCStorageBillingModel types.String             `tfsdk:"storage_billing_model"`
	VDCStorageProfiles     []vdcStorageProfileModel `tfsdk:"storage_profiles"`
	VDCGroup               types.String             `tfsdk:"vdc_group"`
}
