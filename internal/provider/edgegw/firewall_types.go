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

package edgegw

import (
	"context"

	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

type firewallModel struct {
	ID              supertypes.StringValue                                `tfsdk:"id"`
	EdgeGatewayID   supertypes.StringValue                                `tfsdk:"edge_gateway_id"`
	EdgeGatewayName supertypes.StringValue                                `tfsdk:"edge_gateway_name"`
	Rules           supertypes.ListNestedObjectValueOf[firewallModelRule] `tfsdk:"rules"`
}

type firewallModelRule struct {
	ID                       supertypes.StringValue        `tfsdk:"id"`
	Name                     supertypes.StringValue        `tfsdk:"name"`
	Enabled                  supertypes.BoolValue          `tfsdk:"enabled"`
	Direction                supertypes.StringValue        `tfsdk:"direction"`
	IPProtocol               supertypes.StringValue        `tfsdk:"ip_protocol"`
	Action                   supertypes.StringValue        `tfsdk:"action"`
	Logging                  supertypes.BoolValue          `tfsdk:"logging"`
	SourceIDs                supertypes.SetValueOf[string] `tfsdk:"source_ids"`
	DestinationIDs           supertypes.SetValueOf[string] `tfsdk:"destination_ids"`
	AppPortProfileIDs        supertypes.SetValueOf[string] `tfsdk:"app_port_profile_ids"`
	NetworkContextProfileIDs supertypes.SetValueOf[string] `tfsdk:"network_context_profile_ids"`
}

func (rm *firewallModel) Copy() *firewallModel {
	x := &firewallModel{}
	utils.ModelCopy(rm, x)
	return x
}

func (rm *firewallModel) rulesToFirewallRule(ctx context.Context) (firewallRules []types.ParamsFirewallRule, diags diag.Diagnostics) {
	rules, d := rm.Rules.Get(ctx)
	if d.HasError() {
		diags.Append(d...)
		return firewallRules, diags
	}

	firewallRules = make([]types.ParamsFirewallRule, len(rules))
	for i, rule := range rules {
		firewallRules[i] = types.ParamsFirewallRule{
			Name:       rule.Name.Get(),
			Action:     rule.Action.Get(),
			Enabled:    rule.Enabled.Get(),
			IPProtocol: rule.IPProtocol.Get(),
			Logging:    rule.Logging.Get(),
			Direction:  rule.Direction.Get(),
		}

		// ! If sourceIDs/destinationIDs is Null, it's an equivalent of any (source/destination)

		sourceFirewallGroups := make([]types.ParamsFirewallGroupMember, 0)
		sourceIDs, d := rule.SourceIDs.Get(ctx)
		if d.HasError() {
			diags.Append(d...)
		}
		for _, id := range sourceIDs {
			sourceFirewallGroups = append(sourceFirewallGroups, types.ParamsFirewallGroupMember{
				ID: id,
			})
		}
		firewallRules[i].SourceFirewallGroups = sourceFirewallGroups

		destinationFirewallGroups := make([]types.ParamsFirewallGroupMember, 0)
		destinationIDs, d := rule.DestinationIDs.Get(ctx)
		if d.HasError() {
			diags.Append(d...)
		}
		for _, id := range destinationIDs {
			destinationFirewallGroups = append(destinationFirewallGroups, types.ParamsFirewallGroupMember{
				ID: id,
			})
		}
		firewallRules[i].DestinationFirewallGroups = destinationFirewallGroups

		applicationPortProfiles := make([]types.ParamsFirewallGroupMember, 0)
		applicationPortProfilesIDs, d := rule.AppPortProfileIDs.Get(ctx)
		if d.HasError() {
			diags.Append(d...)
		}
		for _, id := range applicationPortProfilesIDs {
			applicationPortProfiles = append(applicationPortProfiles, types.ParamsFirewallGroupMember{
				ID: id,
			})
		}
		firewallRules[i].ApplicationPortProfiles = applicationPortProfiles

		networkContextProfiles := make([]types.ParamsFirewallGroupMember, 0)
		networkContextProfilesIDs, d := rule.NetworkContextProfileIDs.Get(ctx)
		if d.HasError() {
			diags.Append(d...)
		}
		for _, id := range networkContextProfilesIDs {
			networkContextProfiles = append(networkContextProfiles, types.ParamsFirewallGroupMember{
				ID: id,
			})
		}
		firewallRules[i].NetworkContextProfiles = networkContextProfiles
		if d.HasError() {
			diags.Append(d...)
			return firewallRules, diags
		}
	}

	return firewallRules, diags
}
