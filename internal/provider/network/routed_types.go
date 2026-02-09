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

package network

import (
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

type (
	RoutedModel struct {
		ID              supertypes.StringValue                                     `tfsdk:"id"`
		Name            supertypes.StringValue                                     `tfsdk:"name"`
		Description     supertypes.StringValue                                     `tfsdk:"description"`
		EdgeGatewayID   supertypes.StringValue                                     `tfsdk:"edge_gateway_id"`
		EdgeGatewayName supertypes.StringValue                                     `tfsdk:"edge_gateway_name"`
		InterfaceType   supertypes.StringValue                                     `tfsdk:"interface_type"`
		Gateway         supertypes.StringValue                                     `tfsdk:"gateway"`
		PrefixLength    supertypes.Int64Value                                      `tfsdk:"prefix_length"`
		DNS1            supertypes.StringValue                                     `tfsdk:"dns1"`
		DNS2            supertypes.StringValue                                     `tfsdk:"dns2"`
		DNSSuffix       supertypes.StringValue                                     `tfsdk:"dns_suffix"`
		StaticIPPool    supertypes.SetNestedObjectValueOf[RoutedModelStaticIPPool] `tfsdk:"static_ip_pool"`
	}
	RoutedModelStaticIPPool struct {
		StartAddress supertypes.StringValue `tfsdk:"start_address"`
		EndAddress   supertypes.StringValue `tfsdk:"end_address"`
	}
)

func (rm *RoutedModel) Copy() *RoutedModel {
	x := &RoutedModel{}
	utils.ModelCopy(rm, x)
	return x
}
