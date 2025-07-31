/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package vrf

import supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

type tier0VrfDataSourceModel struct {
	ID supertypes.StringValue `tfsdk:"id"`

	Name            supertypes.StringValue `tfsdk:"name"`
	EdgeGatewayID   supertypes.StringValue `tfsdk:"edge_gateway_id"`
	EdgeGatewayName supertypes.StringValue `tfsdk:"edge_gateway_name"`

	ClassService supertypes.StringValue                                                 `tfsdk:"class_service"`
	Bandwidth    supertypes.SingleNestedObjectValueOf[tier0VrfDataSourceModelBandwidth] `tfsdk:"bandwidth"`
	EdgeGateways supertypes.ListNestedObjectValueOf[tier0VrfDataSourceModelEdgeGateway] `tfsdk:"edgegateways"`
}

type tier0VrfDataSourceModelBandwidth struct {
	Capacity               supertypes.Int64Value         `tfsdk:"capacity"`
	Provisioned            supertypes.Int64Value         `tfsdk:"provisioned"`
	Remaining              supertypes.Int64Value         `tfsdk:"remaining"`
	AllowedBandwidthValues supertypes.ListValueOf[int64] `tfsdk:"allowed_bandwidth_values"`
	AllowUnlimited         supertypes.BoolValue          `tfsdk:"allow_unlimited"`
}

type tier0VrfDataSourceModelEdgeGateway struct {
	ID                     supertypes.StringValue        `tfsdk:"id"`
	Name                   supertypes.StringValue        `tfsdk:"name"`
	Bandwidth              supertypes.Int64Value         `tfsdk:"bandwidth"`
	AllowedBandwidthValues supertypes.ListValueOf[int64] `tfsdk:"allowed_bandwidth_values"`
}
