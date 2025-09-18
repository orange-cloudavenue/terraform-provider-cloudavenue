/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package edgegw

import (
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

type edgeGatewayResourceModel struct {
	ID          supertypes.StringValue `tfsdk:"id"`
	Name        supertypes.StringValue `tfsdk:"name"`
	Description supertypes.StringValue `tfsdk:"description"`
	T0Name      supertypes.StringValue `tfsdk:"t0_name"`
	OwnerName   supertypes.StringValue `tfsdk:"owner_name"`
	Bandwidth   supertypes.Int64Value  `tfsdk:"bandwidth"`

	// Deprecated
	Tier0VRFName supertypes.StringValue `tfsdk:"tier0_vrf_name"`

	// Read-Only
	OwnerID supertypes.StringValue `tfsdk:"owner_id"`
	T0ID    supertypes.StringValue `tfsdk:"t0_id"`
}

// Copy returns a copy of the edgeGatewayResourceModel.
func (rm *edgeGatewayResourceModel) Copy() *edgeGatewayResourceModel {
	x := &edgeGatewayResourceModel{}
	utils.ModelCopy(rm, x)
	return x
}

func (rm *edgeGatewayResourceModel) fromSDK(data *types.ModelEdgeGateway) {
	if data == nil {
		*rm = edgeGatewayResourceModel{}
		return
	}

	rm.ID.Set(data.ID)
	rm.Name.Set(data.Name)
	rm.Description.Set(data.Description)
	rm.OwnerName.Set(data.OwnerRef.Name)
	rm.OwnerID.Set(data.OwnerRef.ID)
	rm.T0ID.Set(data.UplinkT0.ID)
	rm.T0Name.Set(data.UplinkT0.Name)
	// Deprecated
	rm.Tier0VRFName.Set(data.UplinkT0.Name)
}
