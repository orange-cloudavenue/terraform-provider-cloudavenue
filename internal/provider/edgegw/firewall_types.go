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

package edgegw

import (
	"context"

	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

type firewallModel struct {
	ID              supertypes.StringValue                                `tfsdk:"id"`
	EdgeGatewayID   supertypes.StringValue                                `tfsdk:"edge_gateway_id"`
	EdgeGatewayName supertypes.StringValue                                `tfsdk:"edge_gateway_name"`
	Rules           supertypes.ListNestedObjectValueOf[firewallModelRule] `tfsdk:"rules"`
}

type firewallModelRule struct {
	ID                supertypes.StringValue        `tfsdk:"id"`
	Name              supertypes.StringValue        `tfsdk:"name"`
	Enabled           supertypes.BoolValue          `tfsdk:"enabled"`
	Direction         supertypes.StringValue        `tfsdk:"direction"`
	IPProtocol        supertypes.StringValue        `tfsdk:"ip_protocol"`
	Action            supertypes.StringValue        `tfsdk:"action"`
	Logging           supertypes.BoolValue          `tfsdk:"logging"`
	SourceIDs         supertypes.SetValueOf[string] `tfsdk:"source_ids"`
	DestinationIDs    supertypes.SetValueOf[string] `tfsdk:"destination_ids"`
	AppPortProfileIDs supertypes.SetValueOf[string] `tfsdk:"app_port_profile_ids"`
}

func (rm *firewallModel) Copy() *firewallModel {
	x := &firewallModel{}
	utils.ModelCopy(rm, x)
	return x
}

func (rm *firewallModel) rulesToNsxtFirewallRule(ctx context.Context) (nsxtFirewallRules []*govcdtypes.NsxtFirewallRule, diags diag.Diagnostics) {
	rules, d := rm.Rules.Get(ctx)
	if d.HasError() {
		diags.Append(d...)
		return nsxtFirewallRules, diags
	}

	nsxtFirewallRules = make([]*govcdtypes.NsxtFirewallRule, len(rules))
	for i, rule := range rules {
		nsxtFirewallRules[i] = &govcdtypes.NsxtFirewallRule{
			Name:                      rule.Name.Get(),
			ActionValue:               rule.Action.Get(),
			Enabled:                   rule.Enabled.Get(),
			IpProtocol:                rule.IPProtocol.Get(),
			Logging:                   rule.Logging.Get(),
			Direction:                 rule.Direction.Get(),
			Version:                   nil,
			SourceFirewallGroups:      nil,
			DestinationFirewallGroups: nil,
			ApplicationPortProfiles:   nil,
		}

		// ! If sourceIDs/destinationIDs is Null, it's an equivalent of any (source/destination)

		nsxtFirewallRules[i].SourceFirewallGroups, d = common.ToOpenAPIReferenceID(ctx, rule.SourceIDs)
		if d.HasError() {
			diags.Append(d...)
			return nsxtFirewallRules, diags
		}

		nsxtFirewallRules[i].DestinationFirewallGroups, d = common.ToOpenAPIReferenceID(ctx, rule.DestinationIDs)
		if d.HasError() {
			diags.Append(d...)
			return nsxtFirewallRules, diags
		}

		nsxtFirewallRules[i].ApplicationPortProfiles, d = common.ToOpenAPIReferenceID(ctx, rule.AppPortProfileIDs)
		if d.HasError() {
			diags.Append(d...)
			return nsxtFirewallRules, diags
		}
	}

	return nsxtFirewallRules, diags
}
