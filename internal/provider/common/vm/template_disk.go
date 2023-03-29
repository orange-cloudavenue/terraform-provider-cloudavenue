package vm

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vm/diskparams"
)

type TemplatesDiskModel []TemplateDiskModel

type TemplateDiskModel struct {
	BusType        types.String `tfsdk:"bus_type"`
	BusNumber      types.Int64  `tfsdk:"bus_number"`
	SizeInMb       types.Int64  `tfsdk:"size_in_mb"`
	StorageProfile types.String `tfsdk:"storage_profile"`
	UnitNumber     types.Int64  `tfsdk:"unit_number"`
}

// TemplateDiskSchema returns the schema for the template disk.
func TemplateDiskSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"bus_type":        diskparams.BusTypeAttributeRequired(),
		"size_in_mb":      diskparams.SizeInMBAttributeRequired(),
		"bus_number":      diskparams.BusNumberAttributeRequired(),
		"unit_number":     diskparams.UnitNumberAttributeRequired(),
		"storage_profile": diskparams.StorageProfileAttributeRequired(),
	}
}

// TemplateDiskAttrType returns the type map for the template disk.
func TemplateDiskAttrType() map[string]attr.Type {
	return map[string]attr.Type{
		"bus_type":        types.StringType,
		"size_in_mb":      types.Int64Type,
		"bus_number":      types.Int64Type,
		"unit_number":     types.Int64Type,
		"storage_profile": types.StringType,
	}
}

// ToAttrValue converts the model to an attr.Value.
func (m *TemplateDiskModel) ToAttrValue() map[string]attr.Value {
	return map[string]attr.Value{
		"bus_type":        m.BusType,
		"size_in_mb":      m.SizeInMb,
		"bus_number":      m.BusNumber,
		"unit_number":     m.UnitNumber,
		"storage_profile": m.StorageProfile,
	}
}

// ObjectType returns the type of the resource object.
func (m *TemplatesDiskModel) ObjectType() attr.Type {
	return types.ObjectType{
		AttrTypes: TemplateDiskAttrType(),
	}
}

// ResourceFromPlan converts a terraform plan to a TemplatesDiskModel struct.
func TemplatesDiskFromPlan(ctx context.Context, x types.Set) (*TemplatesDiskModel, diag.Diagnostics) {
	if x.IsNull() || x.IsUnknown() {
		return &TemplatesDiskModel{}, diag.Diagnostics{}
	}

	t := &TemplatesDiskModel{}
	return t, x.ElementsAs(ctx, t, false)
}

// ToPlan converts the resource struct to a plan.
func (m *TemplatesDiskModel) ToPlan(ctx context.Context) (basetypes.SetValue, diag.Diagnostics) {
	if m == nil {
		return types.SetNull(m.ObjectType()), diag.Diagnostics{}
	}

	return types.SetValueFrom(ctx, m.ObjectType(), m)
}
