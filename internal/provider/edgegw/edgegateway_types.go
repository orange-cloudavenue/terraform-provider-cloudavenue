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

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

type edgeGatewayResourceModel struct {
	Timeouts    timeouts.Value         `tfsdk:"timeouts"`
	ID          supertypes.StringValue `tfsdk:"id"`
	Tier0VrfID  supertypes.StringValue `tfsdk:"tier0_vrf_name"`
	Name        supertypes.StringValue `tfsdk:"name"`
	OwnerType   supertypes.StringValue `tfsdk:"owner_type"`
	OwnerName   supertypes.StringValue `tfsdk:"owner_name"`
	Description supertypes.StringValue `tfsdk:"description"`
	Bandwidth   supertypes.Int64Value  `tfsdk:"bandwidth"`
}

type edgeGatewayDatasourceModel struct {
	ID          supertypes.StringValue `tfsdk:"id"`
	Tier0VrfID  supertypes.StringValue `tfsdk:"tier0_vrf_name"`
	Name        supertypes.StringValue `tfsdk:"name"`
	OwnerType   supertypes.StringValue `tfsdk:"owner_type"`
	OwnerName   supertypes.StringValue `tfsdk:"owner_name"`
	Description supertypes.StringValue `tfsdk:"description"`
	Bandwidth   supertypes.Int64Value  `tfsdk:"bandwidth"`
}

// Copy returns a copy of the edgeGatewayResourceModel.
func (rm *edgeGatewayResourceModel) Copy() *edgeGatewayResourceModel {
	x := &edgeGatewayResourceModel{}
	utils.ModelCopy(rm, x)
	return x
}

// Copy returns a copy of the edgeGatewayDatasourceModel.
func (dm *edgeGatewayDatasourceModel) Copy() *edgeGatewayDatasourceModel {
	x := &edgeGatewayDatasourceModel{}
	utils.ModelCopy(dm, x)
	return x
}
