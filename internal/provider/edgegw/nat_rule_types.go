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

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

type NATRuleModel struct {
	AppPortProfileID supertypes.StringValue `tfsdk:"app_port_profile_id"`
	Description      supertypes.StringValue `tfsdk:"description"`
	DnatExternalPort supertypes.StringValue `tfsdk:"dnat_external_port"`
	EdgeGatewayID    supertypes.StringValue `tfsdk:"edge_gateway_id"`
	EdgeGatewayName  supertypes.StringValue `tfsdk:"edge_gateway_name"`
	Enabled          supertypes.BoolValue   `tfsdk:"enabled"`
	ExternalAddress  supertypes.StringValue `tfsdk:"external_address"`
	FirewallMatch    supertypes.StringValue `tfsdk:"firewall_match"`
	ID               supertypes.StringValue `tfsdk:"id"`
	InternalAddress  supertypes.StringValue `tfsdk:"internal_address"`
	// Option not available in CloudAvenue
	// Logging                supertypes.BoolValue   `tfsdk:"logging"`
	Name                   supertypes.StringValue `tfsdk:"name"`
	Priority               supertypes.Int64Value  `tfsdk:"priority"`
	RuleType               supertypes.StringValue `tfsdk:"rule_type"`
	SnatDestinationAddress supertypes.StringValue `tfsdk:"snat_destination_address"`
}

func (rm *NATRuleModel) Copy() *NATRuleModel {
	x := &NATRuleModel{}
	utils.ModelCopy(rm, x)
	return x
}

func (rm *NATRuleModel) ToNsxtNATRule(_ context.Context) (values *govcdtypes.NsxtNatRule, err error) {
	values = &govcdtypes.NsxtNatRule{
		ApplicationPortProfile: func() *govcdtypes.OpenApiReference {
			if rm.AppPortProfileID.Get() != "" {
				return &govcdtypes.OpenApiReference{ID: rm.AppPortProfileID.Get()}
			}
			return nil
		}(),
		Name:                     rm.Name.Get(),
		Description:              rm.Description.Get(),
		Enabled:                  rm.Enabled.Get(),
		ExternalAddresses:        rm.ExternalAddress.Get(),
		InternalAddresses:        rm.InternalAddress.Get(),
		SnatDestinationAddresses: rm.SnatDestinationAddress.Get(),
		DnatExternalPort:         rm.DnatExternalPort.Get(),
		Type:                     rm.RuleType.Get(),
		FirewallMatch:            rm.FirewallMatch.Get(),
		Priority:                 rm.Priority.GetIntPtr(),
	}

	return values, err
}
