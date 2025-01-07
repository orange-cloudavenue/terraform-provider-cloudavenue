/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

// Package edgegw provides a Terraform resource to manage edge gateways.
package edgegw

import (
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

type (
	edgeGatewaysDataSourceModel struct {
		ID           supertypes.StringValue                                                    `tfsdk:"id"`
		EdgeGateways supertypes.ListNestedObjectValueOf[edgeGatewayDataSourceModelEdgeGateway] `tfsdk:"edge_gateways"`
	}
	edgeGatewayDataSourceModelEdgeGateway struct {
		Tier0VrfName supertypes.StringValue `tfsdk:"tier0_vrf_name"`
		Name         supertypes.StringValue `tfsdk:"name"`
		ID           supertypes.StringValue `tfsdk:"id"`
		OwnerType    supertypes.StringValue `tfsdk:"owner_type"`
		OwnerName    supertypes.StringValue `tfsdk:"owner_name"`
		Description  supertypes.StringValue `tfsdk:"description"`
		LbEnabled    supertypes.BoolValue   `tfsdk:"lb_enabled"`
	}
)
