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

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/v1/edgegateway"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

type firewallModel struct {
	ID              supertypes.StringValue                               `tfsdk:"id"`
	EdgeGatewayID   supertypes.StringValue                               `tfsdk:"edge_gateway_id"`
	EdgeGatewayName supertypes.StringValue                               `tfsdk:"edge_gateway_name"`
	Rules           supertypes.SetNestedObjectValueOf[firewallModelRule] `tfsdk:"rules"`
}

type firewallModelRule struct {
	ID                     supertypes.StringValue        `tfsdk:"id"`
	Name                   supertypes.StringValue        `tfsdk:"name"`
	Enabled                supertypes.BoolValue          `tfsdk:"enabled"`
	Direction              supertypes.StringValue        `tfsdk:"direction"`
	IPProtocol             supertypes.StringValue        `tfsdk:"ip_protocol"`
	Priority               supertypes.Int64Value         `tfsdk:"priority"`
	Action                 supertypes.StringValue        `tfsdk:"action"`
	Logging                supertypes.BoolValue          `tfsdk:"logging"`
	SourceIDs              supertypes.SetValueOf[string] `tfsdk:"source_ids"`
	SourceIPAddresses      supertypes.SetValueOf[string] `tfsdk:"source_ip_addresses"`
	DestinationIDs         supertypes.SetValueOf[string] `tfsdk:"destination_ids"`
	DestinationIPAddresses supertypes.SetValueOf[string] `tfsdk:"destination_ip_addresses"`
	AppPortProfileIDs      supertypes.SetValueOf[string] `tfsdk:"app_port_profile_ids"`
}

func (rm *firewallModel) Copy() *firewallModel {
	x := &firewallModel{}
	utils.ModelCopy(rm, x)
	return x
}

func (rm *firewallModel) ToSDK(ctx context.Context) (rules edgegateway.FirewallModelRules, diags diag.Diagnostics) {
	diags = diag.Diagnostics{}

	if rm == nil {
		return edgegateway.FirewallModelRules{}, diags
	}

	rulesRM, d := rm.Rules.Get(ctx)
	if d.HasError() {
		diags.Append(d...)
		return edgegateway.FirewallModelRules{}, diags
	}

	sdkRules := make([]*edgegateway.FirewallModelRule, len(rulesRM))
	for i, rule := range rulesRM {
		sdkRules[i] = &edgegateway.FirewallModelRule{
			ID:         rule.ID.Get(),
			Name:       rule.Name.Get(),
			Enabled:    rule.Enabled.Get(),
			Direction:  rule.Direction.Get(),
			IPProtocol: rule.IPProtocol.Get(),
			Priority:   rule.Priority.GetIntPtr(),
			Action:     rule.Action.Get(),
			Logging:    rule.Logging.Get(),
			SourceIPAddresses: func() []string {
				ips, d := rule.SourceIPAddresses.Get(ctx)
				if d.HasError() {
					diags.Append(d...)
					return nil
				}
				return ips
			}(),
			SourceFirewallGroups: func() []govcdtypes.OpenApiReference {
				ids, d := common.ToOpenAPIReferenceID(ctx, rule.SourceIDs)
				if d.HasError() {
					diags.Append(d...)
					return nil
				}
				return ids
			}(),
			DestinationFirewallGroups: func() []govcdtypes.OpenApiReference {
				ids, d := common.ToOpenAPIReferenceID(ctx, rule.DestinationIDs)
				if d.HasError() {
					diags.Append(d...)
					return nil
				}
				return ids
			}(),
			DestinationIPAddresses: func() []string {
				ips, d := rule.DestinationIPAddresses.Get(ctx)
				if d.HasError() {
					diags.Append(d...)
					return nil
				}
				return ips
			}(),
			ApplicationPortProfiles: func() []govcdtypes.OpenApiReference {
				ids, d := common.ToOpenAPIReferenceID(ctx, rule.AppPortProfileIDs)
				if d.HasError() {
					diags.Append(d...)
					return nil
				}
				return ids
			}(),
		}
	}

	return edgegateway.FirewallModelRules{Rules: sdkRules}, diags
}
