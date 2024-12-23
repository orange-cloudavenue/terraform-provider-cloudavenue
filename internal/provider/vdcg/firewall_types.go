package vdcg

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/helpers"
	v1 "github.com/orange-cloudavenue/cloudavenue-sdk-go/v1"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

type (
	FirewallModel struct {
		ID           supertypes.StringValue                                `tfsdk:"id"`
		VDCGroupID   supertypes.StringValue                                `tfsdk:"vdc_group_id"`
		VDCGroupName supertypes.StringValue                                `tfsdk:"vdc_group_name"`
		Enabled      supertypes.BoolValue                                  `tfsdk:"enabled"`
		Rules        supertypes.ListNestedObjectValueOf[FirewallModelRule] `tfsdk:"rules"`
	}

	FirewallModelRule struct {
		ID                  supertypes.StringValue        `tfsdk:"id"`
		Name                supertypes.StringValue        `tfsdk:"name"`
		Enabled             supertypes.BoolValue          `tfsdk:"enabled"`
		Direction           supertypes.StringValue        `tfsdk:"direction"`
		IPProtocol          supertypes.StringValue        `tfsdk:"ip_protocol"`
		Action              supertypes.StringValue        `tfsdk:"action"`
		Logging             supertypes.BoolValue          `tfsdk:"logging"`
		SourceIDs           supertypes.SetValueOf[string] `tfsdk:"source_ids"`
		DestinationIDs      supertypes.SetValueOf[string] `tfsdk:"destination_ids"`
		AppPortProfileIDs   supertypes.SetValueOf[string] `tfsdk:"app_port_profile_ids"`
		SourcesExcluded     supertypes.BoolValue          `tfsdk:"source_groups_excluded"`
		DestinationExcluded supertypes.BoolValue          `tfsdk:"destination_groups_excluded"`
	}
)

func (rm *FirewallModel) Copy() *FirewallModel {
	x := &FirewallModel{}
	utils.ModelCopy(rm, x)
	return x
}

func (rm *FirewallModel) rulesToSDKRules(ctx context.Context) (sdkrules v1.VDCGroupFirewallType, diags diag.Diagnostics) {
	sdkrules.Enabled = rm.Enabled.Get()

	rules, d := rm.Rules.Get(ctx)
	if d.HasError() {
		diags.Append(d...)
		return
	}

	sdkrules.Rules = make(v1.VDCGroupFirewallTypeRules, len(rules))
	for i, r := range rules {
		sdkrules.Rules[i] = &v1.VDCGroupFirewallTypeRule{
			Name:    r.Name.Get(),
			Enabled: r.Enabled.Get(),
			Direction: func() v1.VDCGroupFirewallTypeRuleDirection {
				s, err := helpers.ParseVDCGroupFirewallRuleDirection(r.Direction.Get())
				if err != nil {
					diags.AddError("Error parsing Direction", err.Error())
				}
				return s
			}(),
			IPProtocol: func() v1.VDCGroupFirewallTypeRuleIPProtocol {
				s, err := helpers.ParseVDCGroupFirewallRuleIPProtocol(r.IPProtocol.Get())
				if err != nil {
					diags.AddError("Error parsing IPProtocol", err.Error())
				}
				return s
			}(),
			Action: func() v1.VDCGroupFirewallTypeRuleAction {
				s, err := helpers.ParseVDCGroupFirewallRuleAction(r.Action.Get())
				if err != nil {
					diags.AddError("Error parsing Action", err.Error())
				}
				return s
			}(),
			Logging: r.Logging.Get(),

			SourceGroupsExcluded:      r.SourcesExcluded.GetPtr(),
			DestinationGroupsExcluded: r.DestinationExcluded.GetPtr(),
		}

		sdkrules.Rules[i].SourceFirewallGroups, d = common.ToOpenAPIReferenceID(ctx, r.SourceIDs)
		diags.Append(d...)
		sdkrules.Rules[i].DestinationFirewallGroups, d = common.ToOpenAPIReferenceID(ctx, r.DestinationIDs)
		diags.Append(d...)
		sdkrules.Rules[i].ApplicationPortProfiles, d = common.ToOpenAPIReferenceID(ctx, r.AppPortProfileIDs)
		diags.Append(d...)

		if diags.HasError() {
			return
		}
	}

	return
}

// sdkRulesToRules.
func (rm *FirewallModel) sdkRulesToRules(ctx context.Context, rules v1.VDCGroupFirewallTypeRules) (terraformRules []*FirewallModelRule, diags diag.Diagnostics) {
	terraformRules = make([]*FirewallModelRule, len(rules))
	for i, r := range rules {
		terraformRules[i] = &FirewallModelRule{
			ID:                  supertypes.NewStringNull(),
			Name:                supertypes.NewStringNull(),
			Enabled:             supertypes.NewBoolNull(),
			Direction:           supertypes.NewStringNull(),
			IPProtocol:          supertypes.NewStringNull(),
			Action:              supertypes.NewStringNull(),
			Logging:             supertypes.NewBoolNull(),
			SourceIDs:           supertypes.NewSetValueOfNull[string](ctx),
			DestinationIDs:      supertypes.NewSetValueOfNull[string](ctx),
			AppPortProfileIDs:   supertypes.NewSetValueOfNull[string](ctx),
			SourcesExcluded:     supertypes.NewBoolNull(),
			DestinationExcluded: supertypes.NewBoolNull(),
		}

		terraformRules[i].ID.Set(r.ID)
		terraformRules[i].Name.Set(r.Name)
		terraformRules[i].Enabled.Set(r.Enabled)
		terraformRules[i].Direction.Set(string(r.Direction))
		terraformRules[i].IPProtocol.Set(string(r.IPProtocol))
		terraformRules[i].Action.Set(string(r.Action))
		terraformRules[i].Logging.Set(r.Logging)
		terraformRules[i].SourcesExcluded.SetPtr(r.SourceGroupsExcluded)
		terraformRules[i].DestinationExcluded.SetPtr(r.DestinationGroupsExcluded)
		diags.Append(terraformRules[i].SourceIDs.Set(ctx, common.FromOpenAPIReferenceID(ctx, r.SourceFirewallGroups))...)
		diags.Append(terraformRules[i].DestinationIDs.Set(ctx, common.FromOpenAPIReferenceID(ctx, r.DestinationFirewallGroups))...)
		diags.Append(terraformRules[i].AppPortProfileIDs.Set(ctx, common.FromOpenAPIReferenceID(ctx, r.ApplicationPortProfiles))...)
	}
	return
}
