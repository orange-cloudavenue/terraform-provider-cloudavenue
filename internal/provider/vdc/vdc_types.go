package vdc

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"

	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"
)

type vdcResourceModel struct {
	Timeouts               timeouts.Value            `tfsdk:"timeouts"`
	ID                     supertypes.StringValue    `tfsdk:"id"`
	Name                   supertypes.StringValue    `tfsdk:"name"`
	Description            supertypes.StringValue    `tfsdk:"description"`
	VDCServiceClass        supertypes.StringValue    `tfsdk:"service_class"`
	VDCDisponibilityClass  supertypes.StringValue    `tfsdk:"disponibility_class"`
	VDCBillingModel        supertypes.StringValue    `tfsdk:"billing_model"`
	VcpuInMhz2             supertypes.Int64Value     `tfsdk:"cpu_speed_in_mhz"`
	CPUAllocated           supertypes.Int64Value     `tfsdk:"cpu_allocated"`
	MemoryAllocated        supertypes.Int64Value     `tfsdk:"memory_allocated"`
	VDCStorageBillingModel supertypes.StringValue    `tfsdk:"storage_billing_model"`
	VDCStorageProfiles     supertypes.SetNestedValue `tfsdk:"storage_profiles"`
}

// * VDCStorageProfiles.
type vdcResourceModelVDCStorageProfiles []vdcResourceModelVDCStorageProfile

// * VDCStorageProfiles.
type vdcResourceModelVDCStorageProfile struct {
	Class   supertypes.StringValue `tfsdk:"class"`
	Limit   supertypes.Int64Value  `tfsdk:"limit"`
	Default supertypes.BoolValue   `tfsdk:"default"`
}

type vdcDataSourceModel struct {
	ID                     supertypes.StringValue    `tfsdk:"id"`
	Name                   supertypes.StringValue    `tfsdk:"name"`
	Description            supertypes.StringValue    `tfsdk:"description"`
	VDCServiceClass        supertypes.StringValue    `tfsdk:"service_class"`
	VDCDisponibilityClass  supertypes.StringValue    `tfsdk:"disponibility_class"`
	VDCBillingModel        supertypes.StringValue    `tfsdk:"billing_model"`
	VcpuInMhz2             supertypes.Int64Value     `tfsdk:"cpu_speed_in_mhz"`
	CPUAllocated           supertypes.Int64Value     `tfsdk:"cpu_allocated"`
	MemoryAllocated        supertypes.Int64Value     `tfsdk:"memory_allocated"`
	VDCStorageBillingModel supertypes.StringValue    `tfsdk:"storage_billing_model"`
	VDCStorageProfiles     supertypes.SetNestedValue `tfsdk:"storage_profiles"`
}

// * VDCStorageProfiles.
type vdcDataSourceModelVDCStorageProfiles []vdcDataSourceModelVDCStorageProfile

// * VDCStorageProfiles.
type vdcDataSourceModelVDCStorageProfile struct {
	Class   supertypes.StringValue `tfsdk:"class"`
	Limit   supertypes.Int64Value  `tfsdk:"limit"`
	Default supertypes.BoolValue   `tfsdk:"default"`
}

// GetVDCStorageProfiles returns a slice of vdcModelVDCStorageProfile from a vdcResourceModel.
func (rm *vdcResourceModel) GetVDCStorageProfiles(ctx context.Context) (values vdcResourceModelVDCStorageProfiles, d diag.Diagnostics) {
	values = make(vdcResourceModelVDCStorageProfiles, 0)
	d = rm.VDCStorageProfiles.Get(ctx, &values, false)
	return values, d
}

// GetVDCStorageProfiles returns a slice of vdcModelVDCStorageProfile from a vdcResourceModel.
func (dm *vdcDataSourceModel) GetVDCStorageProfiles(ctx context.Context) (values vdcDataSourceModelVDCStorageProfiles, d diag.Diagnostics) {
	values = make(vdcDataSourceModelVDCStorageProfiles, 0)
	d = dm.VDCStorageProfiles.Get(ctx, &values, false)
	return values, d
}
