// SPDX-FileCopyrightText: Copyright (c) 2026 Orange
// SPDX-License-Identifier: MPL-2.0
//
// This software is distributed under the MPL-2.0 license.
// the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
// or see the "LICENSE" file for more details.

// SPDX-FileCopyrightText: Copyright (c) 2023-2026 Orange
// SPDX-License-Identifier: MPL-2.0
//
// This software is distributed under the MPL-2.0 license.
// the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
// or see the "LICENSE" file for more details.

package vapp

import (
	"cmp"
	"context"
	"fmt"

	"github.com/orange-cloudavenue/common-go/urn"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	sdktypes "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

// isolatedFenceMode is the fence mode of a vApp isolated network.
const isolatedFenceMode = "isolated"

type (
	isolatedNetworkModel struct {
		ID                 supertypes.StringValue                                              `tfsdk:"id"`
		VDC                supertypes.StringValue                                              `tfsdk:"vdc"`
		Name               supertypes.StringValue                                              `tfsdk:"name"`
		Description        supertypes.StringValue                                              `tfsdk:"description"`
		VAppName           supertypes.StringValue                                              `tfsdk:"vapp_name"`
		VAppID             supertypes.StringValue                                              `tfsdk:"vapp_id"`
		Netmask            supertypes.StringValue                                              `tfsdk:"netmask"`
		Gateway            supertypes.StringValue                                              `tfsdk:"gateway"`
		DNS1               supertypes.StringValue                                              `tfsdk:"dns1"`
		DNS2               supertypes.StringValue                                              `tfsdk:"dns2"`
		DNSSuffix          supertypes.StringValue                                              `tfsdk:"dns_suffix"`
		GuestVLANAllowed   supertypes.BoolValue                                                `tfsdk:"guest_vlan_allowed"`
		RetainIPMacEnabled supertypes.BoolValue                                                `tfsdk:"retain_ip_mac_enabled"`
		StaticIPPool       supertypes.SetNestedObjectValueOf[isolatedNetworkModelStaticIPPool] `tfsdk:"static_ip_pool"`
	}

	isolatedNetworkModelStaticIPPool struct {
		StartAddress supertypes.StringValue `tfsdk:"start_address"`
		EndAddress   supertypes.StringValue `tfsdk:"end_address"`
	}
)

// Copy returns a copy of the backupModel.
func (rm *isolatedNetworkModel) Copy() *isolatedNetworkModel {
	x := &isolatedNetworkModel{}
	utils.ModelCopy(rm, x)
	return x
}

// toSDK builds the SDK isolated network from the model.
// When current is not nil it is used as base, so the runtime fields not managed
// by the provider (IsConnected, IsDeployed, IsShared, PrefixLength...) are preserved on update.
func (rm *isolatedNetworkModel) toSDK(ctx context.Context, current *sdktypes.ModelVAppRuntimeNetworkConfig) (network sdktypes.ModelVAppRuntimeNetworkConfig, diags diag.Diagnostics) {
	staticIPPools, d := rm.StaticIPPool.Get(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return network, diags
	}

	ipRanges := make([]sdktypes.ModelVAppRuntimeIPRange, 0, len(staticIPPools))
	for _, staticIPPool := range staticIPPools {
		ipRanges = append(ipRanges, sdktypes.ModelVAppRuntimeIPRange{
			StartAddress: staticIPPool.StartAddress.Get(),
			EndAddress:   staticIPPool.EndAddress.Get(),
		})
	}

	scope := sdktypes.ModelVAppRuntimeIPScope{}
	if current != nil {
		network = *current
		if len(current.Configuration.IPScopes) > 0 {
			scope = current.Configuration.IPScopes[0]
		}
	}

	scope.Gateway = rm.Gateway.Get()
	scope.Netmask = rm.Netmask.Get()
	scope.DNS1 = rm.DNS1.Get()
	scope.DNS2 = rm.DNS2.Get()
	scope.DNSSuffix = rm.DNSSuffix.Get()
	scope.IPRanges = ipRanges

	network.NetworkName = rm.Name.Get()
	network.Description = rm.Description.Get()
	network.Configuration.FenceMode = isolatedFenceMode
	network.Configuration.ParentNetwork = nil
	network.Configuration.RetainNetInfoAcrossDeployments = rm.RetainIPMacEnabled.GetPtr()
	network.Configuration.IPScopes = []sdktypes.ModelVAppRuntimeIPScope{scope}
	network.Configuration.GuestVLANAllowed = rm.GuestVLANAllowed.GetPtr()

	return network, diags
}

// fromSDK fills the model from an SDK isolated network.
// VDC, vApp name and vApp ID are not part of the network and are set by the caller.
func (rm *isolatedNetworkModel) fromSDK(ctx context.Context, network *sdktypes.ModelVAppRuntimeNetworkConfig) (diags diag.Diagnostics) {
	networkID := urn.ExtractUUID(cmp.Or(network.ID, network.HREF))
	if networkID == "" {
		diags.AddError("Error creating vApp network ID", fmt.Sprintf("no UUID found in ID %q or HREF %q of network %q", network.ID, network.HREF, network.NetworkName))
		return diags
	}

	rm.ID.Set(urn.Normalize(urn.Network, networkID).String())
	rm.Name.Set(network.NetworkName)
	rm.Description.Set(network.Description)
	rm.RetainIPMacEnabled.SetPtr(network.Configuration.RetainNetInfoAcrossDeployments)
	rm.GuestVLANAllowed.SetPtr(network.Configuration.GuestVLANAllowed)

	if len(network.Configuration.IPScopes) == 0 {
		rm.StaticIPPool.SetNull(ctx)
		return diags
	}

	scope := network.Configuration.IPScopes[0]
	rm.Netmask.Set(scope.Netmask)
	rm.Gateway.Set(scope.Gateway)
	rm.DNS1.Set(scope.DNS1)
	rm.DNS2.Set(scope.DNS2)
	rm.DNSSuffix.Set(scope.DNSSuffix)

	if len(scope.IPRanges) == 0 {
		rm.StaticIPPool.SetNull(ctx)
		return diags
	}

	ipPool := make([]*isolatedNetworkModelStaticIPPool, 0, len(scope.IPRanges))
	for _, ipRange := range scope.IPRanges {
		ipPool = append(ipPool, &isolatedNetworkModelStaticIPPool{
			StartAddress: supertypes.NewStringValue(ipRange.StartAddress),
			EndAddress:   supertypes.NewStringValue(ipRange.EndAddress),
		})
	}
	diags.Append(rm.StaticIPPool.Set(ctx, ipPool)...)

	return diags
}
