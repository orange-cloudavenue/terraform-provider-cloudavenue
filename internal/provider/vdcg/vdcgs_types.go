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

type (
	vdcgsModel struct {
		ID         supertypes.StringValue                            `tfsdk:"id"`
		FilterID   supertypes.StringValue                            `tfsdk:"filter_by_id"`
		FilterName supertypes.StringValue                            `tfsdk:"filter_by_name"`
		VDCGs      supertypes.SetNestedObjectValueOf[vdcgsModelVDCG] `tfsdk:"vdc_groups"`
	}

	vdcgsModelVDCG struct {
		ID           supertypes.StringValue                           `tfsdk:"id"`
		Name         supertypes.StringValue                           `tfsdk:"name"`
		Description  supertypes.StringValue                           `tfsdk:"description"`
		NumberOfVDCs supertypes.Int64Value                            `tfsdk:"number_of_vdcs"`
		VDCs         supertypes.SetNestedObjectValueOf[vdcgsModelVDC] `tfsdk:"vdcs"`
	}

	vdcgsModelVDC struct {
		ID   supertypes.StringValue `tfsdk:"id"`
		Name supertypes.StringValue `tfsdk:"name"`
	}
)

func (rm *vdcgsModel) fromSDK(ctx context.Context, data *types.ModelListVdcGroup) (diags diag.Diagnostics) {
	if data == nil {
		return
	}

	// listOfVDCGIDs is used to create an ID for the datasource
	listOfVDCGIDs := make([]string, 0, len(data.VdcGroups))

	listOfVDCgs := make([]*vdcgsModelVDCG, 0, len(data.VdcGroups))

	for _, vdcg := range data.VdcGroups {
		item := &vdcgsModelVDCG{}
		item.ID.Set(vdcg.ID)
		item.Name.Set(vdcg.Name)
		item.Description.Set(vdcg.Description)
		item.NumberOfVDCs.SetInt(vdcg.NumberOfVdcs)

		listOfVDCGIDs = append(listOfVDCGIDs, vdcg.ID)

		listOfVDCs := make([]*vdcgsModelVDC, 0, len(vdcg.Vdcs))
		for _, vdc := range vdcg.Vdcs {
			vdcItem := vdcgsModelVDC{}
			vdcItem.ID.Set(vdc.ID)
			vdcItem.Name.Set(vdc.Name)
			listOfVDCs = append(listOfVDCs, &vdcItem)
		}
		diags.Append(item.VDCs.Set(ctx, listOfVDCs)...)
		if diags.HasError() {
			return
		}

		listOfVDCgs = append(listOfVDCgs, item)
	}

	// Generate a stable ID based on the list of VDCG IDs
	// If the list is empty inject a constant to avoid generating a random UUID
	if len(listOfVDCGIDs) == 0 {
		listOfVDCGIDs = append(listOfVDCGIDs, "no-vdcg")
	}
	rm.ID.Set(utils.GenerateUUID(listOfVDCGIDs...).ValueString())

	diags.Append(rm.VDCGs.Set(ctx, listOfVDCgs)...)

	return
}
