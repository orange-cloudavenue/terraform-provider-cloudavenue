package edgegw

import (
	"context"
	"fmt"

	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"

	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

type IPSetModel struct {
	Description     supertypes.StringValue `tfsdk:"description"`
	EdgeGatewayID   supertypes.StringValue `tfsdk:"edge_gateway_id"`
	EdgeGatewayName supertypes.StringValue `tfsdk:"edge_gateway_name"`
	ID              supertypes.StringValue `tfsdk:"id"`
	IPAddresses     supertypes.SetValue    `tfsdk:"ip_addresses"`
	Name            supertypes.StringValue `tfsdk:"name"`
}

type IPSetModelIPAddresses []supertypes.StringValue

func NewIPSet(t any) *IPSetModel {
	switch x := t.(type) {
	case tfsdk.State:
		return &IPSetModel{
			Description:     supertypes.NewStringNull(),
			EdgeGatewayID:   supertypes.NewStringUnknown(),
			EdgeGatewayName: supertypes.NewStringUnknown(),
			ID:              supertypes.NewStringUnknown(),
			IPAddresses:     supertypes.NewSetNull(x.Schema.GetAttributes()["ip_addresses"].GetType().(supertypes.SetType).ElementType()),
			Name:            supertypes.NewStringNull(),
		}
	case tfsdk.Plan:
		return &IPSetModel{
			Description:     supertypes.NewStringNull(),
			EdgeGatewayID:   supertypes.NewStringUnknown(),
			EdgeGatewayName: supertypes.NewStringUnknown(),
			ID:              supertypes.NewStringUnknown(),
			IPAddresses:     supertypes.NewSetNull(x.Schema.GetAttributes()["ip_addresses"].GetType().(supertypes.SetType).ElementType()),
			Name:            supertypes.NewStringNull(),
		}
	case tfsdk.Config:
		return &IPSetModel{
			Description:     supertypes.NewStringNull(),
			EdgeGatewayID:   supertypes.NewStringUnknown(),
			EdgeGatewayName: supertypes.NewStringUnknown(),
			ID:              supertypes.NewStringUnknown(),
			IPAddresses:     supertypes.NewSetNull(x.Schema.GetAttributes()["ip_addresses"].GetType().(supertypes.SetType).ElementType()),
			Name:            supertypes.NewStringNull(),
		}
	default:
		panic(fmt.Sprintf("unexpected type %T", t))
	}
}

func (rm *IPSetModel) Copy() *IPSetModel {
	x := &IPSetModel{}
	utils.ModelCopy(rm, x)
	return x
}

// GetIpAddresses returns the value of the IpAddresses field.
func (rm *IPSetModel) GetIPAddresses(ctx context.Context) (values IPSetModelIPAddresses, diags diag.Diagnostics) {
	values = make(IPSetModelIPAddresses, 0)
	d := rm.IPAddresses.Get(ctx, &values, false)
	return values, d
}

// ToNsxtFirewallGroup transform the IPSetModel to a govcdtypes.NsxtFirewallGroup.
func (rm *IPSetModel) ToNsxtFirewallGroup(ctx context.Context, ownerID string) (values *govcdtypes.NsxtFirewallGroup, diags diag.Diagnostics) {
	values = &govcdtypes.NsxtFirewallGroup{
		Name:        rm.Name.Get(),
		Description: rm.Description.Get(),
		OwnerRef:    &govcdtypes.OpenApiReference{ID: ownerID},
		Type:        govcdtypes.FirewallGroupTypeIpSet,
	}

	if rm.ID.IsKnown() {
		values.ID = rm.ID.Get()
	}

	if rm.IPAddresses.IsKnown() {
		ipAddrs, d := rm.GetIPAddresses(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		values.IpAddresses = utils.SuperSliceTypesStringToSliceString(ipAddrs)
	}

	return values, diags
}
