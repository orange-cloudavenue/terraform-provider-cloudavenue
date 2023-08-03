package edgegw

import (
	"context"
	"fmt"

	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"

	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

type NATRuleModel struct {
	// Option not implemented - see schema comment
	// AppPortProfileID supertypes.StringValue `tfsdk:"app_port_profile_id"`
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

func NewNATRule(t any) *NATRuleModel {
	switch x := t.(type) {
	case tfsdk.State:
		return &NATRuleModel{
			Description:            supertypes.NewStringNull(),
			DnatExternalPort:       supertypes.NewStringNull(),
			EdgeGatewayID:          supertypes.NewStringUnknown(),
			EdgeGatewayName:        supertypes.NewStringUnknown(),
			Enabled:                supertypes.NewBoolNull(),
			ExternalAddress:        supertypes.NewStringNull(),
			FirewallMatch:          supertypes.NewStringNull(),
			ID:                     supertypes.NewStringUnknown(),
			InternalAddress:        supertypes.NewStringNull(),
			Name:                   supertypes.NewStringUnknown(),
			Priority:               supertypes.NewInt64Unknown(),
			RuleType:               supertypes.NewStringNull(),
			SnatDestinationAddress: supertypes.NewStringNull(),
		}

	case tfsdk.Plan:
		return &NATRuleModel{
			Description:            supertypes.NewStringNull(),
			DnatExternalPort:       supertypes.NewStringNull(),
			EdgeGatewayID:          supertypes.NewStringUnknown(),
			EdgeGatewayName:        supertypes.NewStringUnknown(),
			Enabled:                supertypes.NewBoolNull(),
			ExternalAddress:        supertypes.NewStringNull(),
			FirewallMatch:          supertypes.NewStringNull(),
			ID:                     supertypes.NewStringUnknown(),
			InternalAddress:        supertypes.NewStringNull(),
			Name:                   supertypes.NewStringUnknown(),
			Priority:               supertypes.NewInt64Unknown(),
			RuleType:               supertypes.NewStringNull(),
			SnatDestinationAddress: supertypes.NewStringNull(),
		}

	case tfsdk.Config:
		return &NATRuleModel{
			Description:            supertypes.NewStringNull(),
			DnatExternalPort:       supertypes.NewStringNull(),
			EdgeGatewayID:          supertypes.NewStringUnknown(),
			EdgeGatewayName:        supertypes.NewStringUnknown(),
			Enabled:                supertypes.NewBoolNull(),
			ExternalAddress:        supertypes.NewStringNull(),
			FirewallMatch:          supertypes.NewStringNull(),
			ID:                     supertypes.NewStringUnknown(),
			InternalAddress:        supertypes.NewStringNull(),
			Name:                   supertypes.NewStringUnknown(),
			Priority:               supertypes.NewInt64Unknown(),
			RuleType:               supertypes.NewStringNull(),
			SnatDestinationAddress: supertypes.NewStringNull(),
		}

	default:
		panic(fmt.Sprintf("unexpected type %T %v", t, x))
	}
}

func (rm *NATRuleModel) Copy() *NATRuleModel {
	x := &NATRuleModel{}
	utils.ModelCopy(rm, x)
	return x
}

func (rm *NATRuleModel) ToNsxtNATRule(ctx context.Context) (values *govcdtypes.NsxtNatRule, err error) {
	values = &govcdtypes.NsxtNatRule{
		Name:                     rm.Name.Get(),
		Description:              rm.Description.Get(),
		Enabled:                  rm.Enabled.Get(),
		ExternalAddresses:        rm.ExternalAddress.Get(),
		InternalAddresses:        rm.InternalAddress.Get(),
		SnatDestinationAddresses: rm.SnatDestinationAddress.Get(),
		DnatExternalPort:         rm.DnatExternalPort.Get(),
		Type:                     rm.RuleType.Get(),
		FirewallMatch:            rm.FirewallMatch.Get(),
		Priority:                 utils.TakeIntPointer(int(rm.Priority.Get())),
	}

	return values, err
}
