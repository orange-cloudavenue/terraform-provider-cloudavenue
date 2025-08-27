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
		ID           supertypes.StringValue                                                     `tfsdk:"id"`
		EdgeGateways supertypes.ListNestedObjectValueOf[edgeGatewaysDataSourceModelEdgeGateway] `tfsdk:"edge_gateways"`
	}
	edgeGatewaysDataSourceModelEdgeGateway struct {
		ID          supertypes.StringValue `tfsdk:"id"`
		Name        supertypes.StringValue `tfsdk:"name"`
		Description supertypes.StringValue `tfsdk:"description"`

		T0Name supertypes.StringValue `tfsdk:"t0_name"`
		T0ID   supertypes.StringValue `tfsdk:"t0_id"`

		OwnerName supertypes.StringValue `tfsdk:"owner_name"`
		OwnerID   supertypes.StringValue `tfsdk:"owner_id"`

		// Deprecated
		Tier0VRFName supertypes.StringValue `tfsdk:"tier0_vrf_name"`
	}
)
