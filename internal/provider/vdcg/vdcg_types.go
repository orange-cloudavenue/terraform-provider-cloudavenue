/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package vdcg

import (
	"context"

	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

type vdcgModel struct {
	ID          supertypes.StringValue        `tfsdk:"id"`
	Name        supertypes.StringValue        `tfsdk:"name"`
	Description supertypes.StringValue        `tfsdk:"description"`
	VDCIDs      supertypes.SetValueOf[string] `tfsdk:"vdc_ids"`
}

func (rm *vdcgModel) Copy() *vdcgModel {
	x := &vdcgModel{}
	utils.ModelCopy(rm, x)
	return x
}

func (rm *vdcgModel) fromSDK(ctx context.Context, data *types.ModelGetVdcGroup) (diags diag.Diagnostics) {
	if data == nil {
		return diags
	}

	rm.ID.Set(data.ID)
	rm.Name.Set(data.Name)
	rm.Description.Set(data.Description)
	if data.Vdcs != nil {
		vdcIDs := make([]string, len(data.Vdcs))
		for i, vdc := range data.Vdcs {
			vdcIDs[i] = vdc.ID
		}
		diags.Append(rm.VDCIDs.Set(ctx, vdcIDs)...)
	} else {
		rm.VDCIDs.SetNull(ctx)
	}

	return diags
}
