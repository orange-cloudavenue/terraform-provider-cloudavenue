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

package vdcg

import (
	"context"

	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	v1 "github.com/orange-cloudavenue/cloudavenue-sdk-go/v1"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

type (
	NetworkRoutedModel struct {
		ID               supertypes.StringValue                                            `tfsdk:"id"`
		Name             supertypes.StringValue                                            `tfsdk:"name"`
		Description      supertypes.StringValue                                            `tfsdk:"description"`
		VDCGroupID       supertypes.StringValue                                            `tfsdk:"vdc_group_id"`
		VDCGroupName     supertypes.StringValue                                            `tfsdk:"vdc_group_name"`
		EdgeGatewayID    supertypes.StringValue                                            `tfsdk:"edge_gateway_id"`
		EdgeGatewayName  supertypes.StringValue                                            `tfsdk:"edge_gateway_name"`
		Gateway          supertypes.StringValue                                            `tfsdk:"gateway"`
		PrefixLength     supertypes.Int64Value                                             `tfsdk:"prefix_length"`
		DNS1             supertypes.StringValue                                            `tfsdk:"dns1"`
		DNS2             supertypes.StringValue                                            `tfsdk:"dns2"`
		DNSSuffix        supertypes.StringValue                                            `tfsdk:"dns_suffix"`
		StaticIPPool     supertypes.SetNestedObjectValueOf[NetworkRoutedModelStaticIPPool] `tfsdk:"static_ip_pool"`
		GuestVLANAllowed supertypes.BoolValue                                              `tfsdk:"guest_vlan_allowed"`
	}
	NetworkRoutedModelStaticIPPool struct {
		StartAddress supertypes.StringValue `tfsdk:"start_address"`
		EndAddress   supertypes.StringValue `tfsdk:"end_address"`
	}
)

func (rm *NetworkRoutedModel) Copy() *NetworkRoutedModel {
	x := &NetworkRoutedModel{}
	utils.ModelCopy(rm, x)
	return x
}

// ToSDKNetworkRoutedGroupModel converts the model to the SDK model.
func (rm *NetworkRoutedModel) ToSDKNetworkRoutedModel(ctx context.Context) (*v1.VDCNetworkRoutedModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	values := &v1.VDCNetworkRoutedModel{
		EdgeGatewayID:   rm.EdgeGatewayID.Get(),
		EdgeGatewayName: rm.EdgeGatewayName.Get(),
		VDCNetworkModel: v1.VDCNetworkModel{
			ID:                      rm.ID.Get(),
			Name:                    rm.Name.Get(),
			Description:             rm.Description.Get(),
			GuestVLANTaggingAllowed: rm.GuestVLANAllowed.GetPtr(),
			Subnet: func() v1.VDCNetworkModelSubnet {
				return v1.VDCNetworkModelSubnet{
					Gateway:      rm.Gateway.Get(),
					PrefixLength: rm.PrefixLength.GetInt(),
					DNSServer1:   rm.DNS1.Get(),
					DNSServer2:   rm.DNS2.Get(),
					DNSSuffix:    rm.DNSSuffix.Get(),
					IPRanges: func() v1.VDCNetworkModelSubnetIPRanges {
						var ipRanges v1.VDCNetworkModelSubnetIPRanges

						ipPools, d := rm.StaticIPPool.Get(ctx)
						if d.HasError() {
							diags.Append(d...)
							return ipRanges
						}

						for _, ipRange := range ipPools {
							ipRanges = append(ipRanges, v1.VDCNetworkModelSubnetIPRange{
								StartAddress: ipRange.StartAddress.Get(),
								EndAddress:   ipRange.EndAddress.Get(),
							})
						}
						return ipRanges
					}(),
				}
			}(),
		},
	}
	return values, diags
}
