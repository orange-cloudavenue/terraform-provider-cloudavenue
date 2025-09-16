/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package vrf

import (
	"context"

	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

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

func (m *tier0VrfDataSourceModel) fromAPI(ctx context.Context, api *types.ModelT0) (diags diag.Diagnostics) {
	m.ID.Set(utils.GenerateUUID(api.Name).String())
	m.Name.Set(api.Name)
	m.ClassService.Set(api.ClassOfService)

	bandwidth := &tier0VrfDataSourceModelBandwidth{
		Capacity:               supertypes.NewInt64Null(),
		Provisioned:            supertypes.NewInt64Null(),
		Remaining:              supertypes.NewInt64Null(),
		AllowedBandwidthValues: supertypes.NewListValueOfNull[int64](ctx),
		AllowUnlimited:         supertypes.NewBoolNull(),
	}
	bandwidth.Capacity.SetInt(api.Bandwidth.Capacity)
	bandwidth.Provisioned.SetInt(api.Bandwidth.Provisioned)
	bandwidth.Remaining.SetInt(api.Bandwidth.Remaining)
	allowedBandwidthValues := make([]int64, 0, len(api.Bandwidth.AllowedBandwidthValues))
	for _, bw := range api.Bandwidth.AllowedBandwidthValues {
		allowedBandwidthValues = append(allowedBandwidthValues, int64(bw))
	}
	diags.Append(bandwidth.AllowedBandwidthValues.Set(ctx, allowedBandwidthValues)...)
	bandwidth.AllowUnlimited.Set(api.Bandwidth.AllowUnlimited)

	diags.Append(m.Bandwidth.Set(ctx, bandwidth)...)

	edgegateways := make([]*tier0VrfDataSourceModelEdgeGateway, 0, len(api.EdgeGateways))
	for _, edgeGateway := range api.EdgeGateways {
		e := &tier0VrfDataSourceModelEdgeGateway{
			AllowedBandwidthValues: supertypes.NewListValueOfNull[int64](ctx),
		}
		e.ID.Set(edgeGateway.ID)
		e.Name.Set(edgeGateway.Name)
		e.Bandwidth.SetInt(edgeGateway.Bandwidth)

		allowedBandwidthValues := make([]int64, 0, len(edgeGateway.AllowedBandwidthValues))
		for _, bw := range edgeGateway.AllowedBandwidthValues {
			allowedBandwidthValues = append(allowedBandwidthValues, int64(bw))
		}
		diags.Append(e.AllowedBandwidthValues.Set(ctx, allowedBandwidthValues)...)
		edgegateways = append(edgegateways, e)
	}
	if diags.HasError() {
		return diags
	}

	diags.Append(m.EdgeGateways.Set(ctx, edgegateways)...)
	if diags.HasError() {
		return diags
	}

	return diags
}
