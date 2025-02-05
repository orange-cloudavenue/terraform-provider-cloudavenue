/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package alb

import (
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

type serviceEngineGroupModel struct {
	ID                      supertypes.StringValue `tfsdk:"id"`
	Name                    supertypes.StringValue `tfsdk:"name"`
	EdgeGatewayID           supertypes.StringValue `tfsdk:"edge_gateway_id"`
	EdgeGatewayName         supertypes.StringValue `tfsdk:"edge_gateway_name"`
	MaxVirtualServices      supertypes.Int64Value  `tfsdk:"max_virtual_services"`
	ReservedVirtualServices supertypes.Int64Value  `tfsdk:"reserved_virtual_services"`
	DeployedVirtualServices supertypes.Int64Value  `tfsdk:"deployed_virtual_services"`
}
