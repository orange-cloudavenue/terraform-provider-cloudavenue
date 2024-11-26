package edgegw

import (
	"context"

	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"

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
		return
	}

	nsxtFirewallRules = make([]*govcdtypes.NsxtFirewallRule, len(rules))
	for i, rule := range rules {
		nsxtFirewallRules[i] = &govcdtypes.NsxtFirewallRule{
			Name:                      rule.Name.Get(),
			Action:                    rule.Action.Get(),
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
			return
		}

		nsxtFirewallRules[i].DestinationFirewallGroups, d = common.ToOpenAPIReferenceID(ctx, rule.DestinationIDs)
		if d.HasError() {
			diags.Append(d...)
			return
		}

		nsxtFirewallRules[i].ApplicationPortProfiles, d = common.ToOpenAPIReferenceID(ctx, rule.AppPortProfileIDs)
		if d.HasError() {
			diags.Append(d...)
			return
		}
	}

	return
}
