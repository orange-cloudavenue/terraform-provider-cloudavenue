package vm

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

/*
Disks is a list of Disk.
*/
type Disks []Disk

/*
Disk is independent disk attached to a VM.
This disk is always detached disk type.
*/
type Disk struct {
	ID             types.String `tfsdk:"id"`
	VAppName       types.String `tfsdk:"vapp_name"`
	VAppID         types.String `tfsdk:"vapp_id"`
	VMName         types.String `tfsdk:"vm_name"`
	VMID           types.String `tfsdk:"vm_id"`
	VDC            types.String `tfsdk:"vdc"`
	Name           types.String `tfsdk:"name"`
	SizeInMb       types.Int64  `tfsdk:"size_in_mb"`
	StorageProfile types.String `tfsdk:"storage_profile"`
	IsDetachable   types.Bool   `tfsdk:"is_detachable"`

	BusType    types.String `tfsdk:"bus_type"`
	BusNumber  types.Int64  `tfsdk:"bus_number"`
	UnitNumber types.Int64  `tfsdk:"unit_number"`
}

/*
ToAttrValue

converts the Disk struct to a terraform plan

attr.Value is a map[string]attr.Value
  - "id"
  - "vdc"
  - "vapp_name"
  - "vapp_id"
  - "vm_name"
  - "vm_id"
  - "name"
  - "size_in_mb"
  - "storage_profile"
  - "is_detachable"
  - "bus_type"
  - "bus_number"
  - "unit_number"
*/
func (d *Disk) ToAttrValue() map[string]attr.Value {
	return map[string]attr.Value{
		"id":              d.ID,
		"vapp_name":       d.VAppName,
		"vapp_id":         d.VAppID,
		"vm_name":         d.VMName,
		"vm_id":           d.VMID,
		"vdc":             d.VDC,
		"name":            d.Name,
		"size_in_mb":      d.SizeInMb,
		"storage_profile": d.StorageProfile,
		"is_detachable":   d.IsDetachable,
		"bus_type":        d.BusType,
		"bus_number":      d.BusNumber,
		"unit_number":     d.UnitNumber,
	}
}

/*
DiskAttrType

returns the attr.Type map for the disk

attr.Type is a map[string]attr.Type
  - "id" 				(types.StringType)
  - "vapp_name" 		(types.StringType)
  - "vapp_id" 			(types.StringType)
  - "vm_name" 			(types.StringType)
  - "vm_id" 			(types.StringType)
  - "vdc" 				(types.StringType)
  - "name" 				(types.StringType)
  - "size_in_mb" 		(types.Int64Type)
  - "storage_profile" 	(types.StringType)
  - "bus_type" 			(types.StringType)
  - "bus_number" 		(types.Int64Type)
  - "unit_number" 		(types.Int64Type)
*/
func DiskAttrType() map[string]attr.Type {
	return map[string]attr.Type{
		"id":              types.StringType,
		"vapp_name":       types.StringType,
		"vapp_id":         types.StringType,
		"vm_name":         types.StringType,
		"vm_id":           types.StringType,
		"vdc":             types.StringType,
		"name":            types.StringType,
		"bus_type":        types.StringType,
		"size_in_mb":      types.Int64Type,
		"storage_profile": types.StringType,
		"is_detachable":   types.BoolType,
		"bus_number":      types.Int64Type,
		"unit_number":     types.Int64Type,
	}
}

// DisksFromPlan converts the terraform plan to a OLDDisk struct.
func DisksFromPlan(ctx context.Context, x types.Set) (*Disk, diag.Diagnostics) {
	if x.IsNull() || x.IsUnknown() {
		return &Disk{}, diag.Diagnostics{}
	}

	c := &Disk{}

	d := x.ElementsAs(ctx, c, false)

	return c, d
}

/*
ElementType

return the attr.Type for the disk.
*/
func (d *Disk) ElementType() attr.Type {
	return types.ObjectType{AttrTypes: DiskAttrType()}
}

func (d *Disks) ElementType() attr.Type {
	return types.ObjectType{AttrTypes: DiskAttrType()}
}

// ToPlan converts the disk struct to a terraform plan.
func (d *Disks) ToPlan(ctx context.Context) (basetypes.SetValue, diag.Diagnostics) {
	if d == nil {
		return types.SetNull(d.ElementType()), diag.Diagnostics{}
	}

	return types.SetValueFrom(ctx, d.ElementType(), d)
}
