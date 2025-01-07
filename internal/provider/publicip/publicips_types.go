/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package publicip

import (
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

type publicIPDataSourceModel struct {
	ID        supertypes.StringValue                                         `tfsdk:"id"`
	PublicIPs supertypes.ListNestedObjectValueOf[publicIPNetworkConfigModel] `tfsdk:"public_ips"`
}

type publicIPNetworkConfigModel struct {
	ID              supertypes.StringValue `tfsdk:"id"`
	PublicIP        supertypes.StringValue `tfsdk:"public_ip"`
	EdgeGatewayName supertypes.StringValue `tfsdk:"edge_gateway_name"`
	EdgeGatewayID   supertypes.StringValue `tfsdk:"edge_gateway_id"`
}
