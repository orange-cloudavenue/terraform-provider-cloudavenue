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

package catalog

import (
	"context"
	"fmt"

	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

type VAPPTemplateModel struct {
	CatalogID    supertypes.StringValue `tfsdk:"catalog_id"`
	CatalogName  supertypes.StringValue `tfsdk:"catalog_name"`
	CreatedAt    supertypes.StringValue `tfsdk:"created_at"`
	Description  supertypes.StringValue `tfsdk:"description"`
	ID           supertypes.StringValue `tfsdk:"id"`
	TemplateID   supertypes.StringValue `tfsdk:"template_id"`
	VMNames      supertypes.SetValue    `tfsdk:"vm_names"`
	TemplateName supertypes.StringValue `tfsdk:"template_name"`
}

type VAPPTemplateModelVMNames []supertypes.StringValue

func NewCatalogVappTemplate(t any) *VAPPTemplateModel {
	switch x := t.(type) {
	case tfsdk.State: //nolint:dupl
		return &VAPPTemplateModel{
			CatalogID:    supertypes.NewStringUnknown(),
			CatalogName:  supertypes.NewStringUnknown(),
			CreatedAt:    supertypes.NewStringUnknown(),
			Description:  supertypes.NewStringUnknown(),
			ID:           supertypes.NewStringUnknown(),
			TemplateID:   supertypes.NewStringNull(),
			TemplateName: supertypes.NewStringNull(),
			VMNames:      supertypes.NewSetUnknown(x.Schema.GetAttributes()["vm_names"].GetType().(supertypes.SetType).ElementType()),
		}

	case tfsdk.Plan: //nolint:dupl
		return &VAPPTemplateModel{
			CatalogID:    supertypes.NewStringUnknown(),
			CatalogName:  supertypes.NewStringUnknown(),
			CreatedAt:    supertypes.NewStringUnknown(),
			Description:  supertypes.NewStringUnknown(),
			ID:           supertypes.NewStringUnknown(),
			TemplateID:   supertypes.NewStringNull(),
			TemplateName: supertypes.NewStringNull(),
			VMNames:      supertypes.NewSetUnknown(x.Schema.GetAttributes()["vm_names"].GetType().(supertypes.SetType).ElementType()),
		}

	case tfsdk.Config: //nolint:dupl
		return &VAPPTemplateModel{
			CatalogID:    supertypes.NewStringUnknown(),
			CatalogName:  supertypes.NewStringUnknown(),
			CreatedAt:    supertypes.NewStringUnknown(),
			Description:  supertypes.NewStringUnknown(),
			ID:           supertypes.NewStringUnknown(),
			TemplateID:   supertypes.NewStringNull(),
			TemplateName: supertypes.NewStringNull(),
			VMNames:      supertypes.NewSetUnknown(x.Schema.GetAttributes()["vm_names"].GetType().(supertypes.SetType).ElementType()),
		}

	default:
		panic(fmt.Sprintf("unexpected type %T", t))
	}
}

func (rm *VAPPTemplateModel) Copy() *VAPPTemplateModel {
	x := &VAPPTemplateModel{}
	utils.ModelCopy(rm, x)
	return x
}

// GetVmNames returns the value of the VmNames field.
func (rm *VAPPTemplateModel) GetVMNames(ctx context.Context) (values VAPPTemplateModelVMNames, diags diag.Diagnostics) {
	values = make(VAPPTemplateModelVMNames, 0)
	d := rm.VMNames.Get(ctx, &values, false)
	return values, d
}
