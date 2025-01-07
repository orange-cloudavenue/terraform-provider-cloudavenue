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
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type VMResourceModelDeployOS struct { //nolint:revive
	VappTemplateID   types.String `tfsdk:"vapp_template_id"`
	VMNameInTemplate types.String `tfsdk:"vm_name_in_template"`
	BootImageID      types.String `tfsdk:"boot_image_id"`
	AcceptAllEulas   types.Bool   `tfsdk:"accept_all_eulas"`
}

// attrTypes() returns the types of the attributes of the DeployOS attribute.
func (do *VMResourceModelDeployOS) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"vapp_template_id":    types.StringType,
		"vm_name_in_template": types.StringType,
		"boot_image_id":       types.StringType,
		"accept_all_eulas":    types.BoolType,
	}
}

// toAttrValues() returns the values of the attributes of the DeployOS attribute.
func (do *VMResourceModelDeployOS) toAttrValues() map[string]attr.Value {
	return map[string]attr.Value{
		"vapp_template_id":    do.VappTemplateID,
		"vm_name_in_template": do.VMNameInTemplate,
		"boot_image_id":       do.BootImageID,
		"accept_all_eulas":    do.AcceptAllEulas,
	}
}

// ToPlan returns the value of the DeployOS attribute, if set, as a types.Object.
func (do *VMResourceModelDeployOS) ToPlan(ctx context.Context) types.Object {
	if do == nil {
		return types.Object{}
	}

	return types.ObjectValueMust(do.AttrTypes(), do.toAttrValues())
}
