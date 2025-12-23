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
	"context"

	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

type VPNIPSecModel struct {
	Description     supertypes.StringValue       `tfsdk:"description"`
	EdgeGatewayID   supertypes.StringValue       `tfsdk:"edge_gateway_id"`
	EdgeGatewayName supertypes.StringValue       `tfsdk:"edge_gateway_name"`
	Enabled         supertypes.BoolValue         `tfsdk:"enabled"`
	ID              supertypes.StringValue       `tfsdk:"id"`
	LocalIPAddress  supertypes.StringValue       `tfsdk:"local_ip_address"`
	LocalNetworks   supertypes.SetValue          `tfsdk:"local_networks"`
	Name            supertypes.StringValue       `tfsdk:"name"`
	PreSharedKey    supertypes.StringValue       `tfsdk:"pre_shared_key"`
	RemoteID        supertypes.StringValue       `tfsdk:"remote_id"`
	RemoteIPAddress supertypes.StringValue       `tfsdk:"remote_ip_address"`
	RemoteNetworks  supertypes.SetValue          `tfsdk:"remote_networks"`
	SecurityProfile supertypes.SingleNestedValue `tfsdk:"security_profile"`
	SecurityType    supertypes.StringValue       `tfsdk:"security_type"`
}

type VPNIPSecModelLocalNetworks []supertypes.StringValue

type VPNIPSecModelRemoteNetworks []supertypes.StringValue

// * SecurityProfile.
type VPNIPSecModelSecurityProfile struct {
	IkeDhGroups                supertypes.StringValue `tfsdk:"ike_dh_groups"`
	IkeDigestAlgorithm         supertypes.StringValue `tfsdk:"ike_digest_algorithm"`
	IkeEncryptionAlgorithm     supertypes.StringValue `tfsdk:"ike_encryption_algorithm"`
	IkeSaLifetime              supertypes.Int64Value  `tfsdk:"ike_sa_lifetime"`
	IkeVersion                 supertypes.StringValue `tfsdk:"ike_version"`
	TunnelDfPolicy             supertypes.StringValue `tfsdk:"tunnel_df_policy"`
	TunnelDhGroups             supertypes.StringValue `tfsdk:"tunnel_dh_groups"`
	TunnelDigestAlgorithms     supertypes.StringValue `tfsdk:"tunnel_digest_algorithms"`
	TunnelDpd                  supertypes.Int64Value  `tfsdk:"tunnel_dpd"`
	TunnelEncryptionAlgorithms supertypes.StringValue `tfsdk:"tunnel_encryption_algorithms"`
	TunnelPfs                  supertypes.BoolValue   `tfsdk:"tunnel_pfs"`
	TunnelSaLifetime           supertypes.Int64Value  `tfsdk:"tunnel_sa_lifetime"`
}

const (
	vpnAuthentication string = "PSK"
	vpnModeInit       string = "INITIATOR"
	profileDefault    string = "DEFAULT"
	profileCustom     string = "CUSTOM"
)

func (rm *VPNIPSecModel) Copy() *VPNIPSecModel {
	x := &VPNIPSecModel{}
	utils.ModelCopy(rm, x)
	return x
}

// GetLocalNetworks returns the value of the LocalNetworks field.
func (rm *VPNIPSecModel) GetLocalNetworks(ctx context.Context) (values VPNIPSecModelLocalNetworks, diags diag.Diagnostics) {
	values = make(VPNIPSecModelLocalNetworks, 0)
	d := rm.LocalNetworks.Get(ctx, &values, false)
	return values, d
}

// GetRemoteNetworks returns the value of the RemoteNetworks field.
func (rm *VPNIPSecModel) GetRemoteNetworks(ctx context.Context) (values VPNIPSecModelRemoteNetworks, diags diag.Diagnostics) {
	values = make(VPNIPSecModelRemoteNetworks, 0)
	d := rm.RemoteNetworks.Get(ctx, &values, false)
	return values, d
}

// GetSecurityProfile returns the value of the SecurityProfile field.
func (rm *VPNIPSecModel) GetSecurityProfile(ctx context.Context) (values VPNIPSecModelSecurityProfile, diags diag.Diagnostics) {
	values = VPNIPSecModelSecurityProfile{}
	d := rm.SecurityProfile.Get(ctx, &values, basetypes.ObjectAsOptions{})
	return values, d
}

func (rm *VPNIPSecModel) GetNsxtIPSecVPNTunnelSecurityProfile(ctx context.Context) (values *govcdtypes.NsxtIpSecVpnTunnelSecurityProfile, diags diag.Diagnostics) {
	values = &govcdtypes.NsxtIpSecVpnTunnelSecurityProfile{}
	if !rm.SecurityProfile.IsKnown() {
		return values, diags
	}

	securityProfile, d := rm.GetSecurityProfile(ctx)
	diags.Append(d...)
	if d.HasError() {
		return values, diags
	}

	if securityProfile.IkeDhGroups.IsKnown() {
		values.IkeConfiguration.DhGroups = append(values.IkeConfiguration.DhGroups, securityProfile.IkeDhGroups.Get())
	}
	if securityProfile.IkeDigestAlgorithm.IsKnown() {
		values.IkeConfiguration.DigestAlgorithms = append(values.IkeConfiguration.DigestAlgorithms, securityProfile.IkeDigestAlgorithm.Get())
	}
	if securityProfile.IkeEncryptionAlgorithm.IsKnown() {
		values.IkeConfiguration.EncryptionAlgorithms = append(values.IkeConfiguration.EncryptionAlgorithms, securityProfile.IkeEncryptionAlgorithm.Get())
	}
	if securityProfile.IkeSaLifetime.IsKnown() {
		values.IkeConfiguration.SaLifeTime = securityProfile.IkeSaLifetime.GetIntPtr()
	}
	if securityProfile.IkeVersion.IsKnown() {
		values.IkeConfiguration.IkeVersion = securityProfile.IkeVersion.Get()
	}
	if securityProfile.TunnelDfPolicy.IsKnown() {
		values.TunnelConfiguration.DfPolicy = securityProfile.TunnelDfPolicy.Get()
	}
	if securityProfile.TunnelDhGroups.IsKnown() {
		values.TunnelConfiguration.DhGroups = append(values.TunnelConfiguration.DhGroups, securityProfile.TunnelDhGroups.Get())
	}
	if securityProfile.TunnelDigestAlgorithms.IsKnown() {
		values.TunnelConfiguration.DigestAlgorithms = append(values.TunnelConfiguration.DigestAlgorithms, securityProfile.TunnelDigestAlgorithms.Get())
	}
	if securityProfile.TunnelEncryptionAlgorithms.IsKnown() {
		values.TunnelConfiguration.EncryptionAlgorithms = append(values.TunnelConfiguration.EncryptionAlgorithms, securityProfile.TunnelEncryptionAlgorithms.Get())
	}
	if securityProfile.TunnelPfs.IsKnown() {
		values.TunnelConfiguration.PerfectForwardSecrecyEnabled = securityProfile.TunnelPfs.Get()
	}
	if securityProfile.TunnelSaLifetime.IsKnown() {
		values.TunnelConfiguration.SaLifeTime = securityProfile.TunnelSaLifetime.GetIntPtr()
	}
	if securityProfile.TunnelDpd.IsKnown() {
		values.DpdConfiguration.ProbeInterval = securityProfile.TunnelDpd.GetInt()
	}
	return values, diags
}

func (rm VPNIPSecModelLocalNetworks) Get() []string {
	return utils.SuperSliceTypesStringToSliceString(rm)
}

func (rm VPNIPSecModelRemoteNetworks) Get() []string {
	return utils.SuperSliceTypesStringToSliceString(rm)
}

func (rm *VPNIPSecModel) ToNsxtIPSecVPNTunnel(ctx context.Context) (values *govcdtypes.NsxtIpSecVpnTunnel, diags diag.Diagnostics) {
	values = &govcdtypes.NsxtIpSecVpnTunnel{
		Name:                    rm.Name.Get(),
		Description:             rm.Description.Get(),
		Enabled:                 rm.Enabled.Get(),
		PreSharedKey:            rm.PreSharedKey.Get(),
		Logging:                 false,             // not available on cloudavenue
		AuthenticationMode:      vpnAuthentication, // Only PSK is supported on cloudavenue
		ConnectorInitiationMode: vpnModeInit,
	}

	if rm.ID.IsKnown() {
		values.ID = rm.ID.Get()
	}

	// Get local networks
	localNet, diags := rm.GetLocalNetworks(ctx)
	diags.Append(diags...)
	if diags.HasError() {
		return values, diags
	}

	values.LocalEndpoint.LocalId = rm.LocalIPAddress.Get()
	values.LocalEndpoint.LocalAddress = rm.LocalIPAddress.Get()
	values.LocalEndpoint.LocalNetworks = localNet.Get()

	// Get remote networks
	remoteNet, diags := rm.GetRemoteNetworks(ctx)
	diags.Append(diags...)
	if diags.HasError() {
		return values, diags
	}

	values.RemoteEndpoint.RemoteId = rm.RemoteIPAddress.Get()

	if rm.RemoteID.IsKnown() {
		values.RemoteEndpoint.RemoteId = rm.RemoteID.Get()
	}

	values.RemoteEndpoint.RemoteAddress = rm.RemoteIPAddress.Get()
	values.RemoteEndpoint.RemoteNetworks = remoteNet.Get()

	return values, diags
}
