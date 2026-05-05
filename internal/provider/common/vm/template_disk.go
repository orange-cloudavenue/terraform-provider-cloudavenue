/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

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
		attrBusType:        diskparams.BusTypeAttributeRequired(),
		attrSizeInMB:       diskparams.SizeInMBAttributeRequired(),
		attrBusNumber:      diskparams.BusNumberAttributeRequired(),
		attrUnitNumber:     diskparams.UnitNumberAttributeRequired(),
		attrStorageProfile: diskparams.StorageProfileAttributeRequired(),
	}
}

// TemplateDiskAttrType returns the type map for the template disk.
func TemplateDiskAttrType() map[string]attr.Type {
	return map[string]attr.Type{
		attrBusType:        types.StringType,
		attrSizeInMB:       types.Int64Type,
		attrBusNumber:      types.Int64Type,
		attrUnitNumber:     types.Int64Type,
		attrStorageProfile: types.StringType,
	}
}

// ToAttrValue converts the model to an attr.Value.
func (m *TemplateDiskModel) ToAttrValue() map[string]attr.Value {
	return map[string]attr.Value{
		attrBusType:        m.BusType,
		attrSizeInMB:       m.SizeInMb,
		attrBusNumber:      m.BusNumber,
		attrUnitNumber:     m.UnitNumber,
		attrStorageProfile: m.StorageProfile,
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
