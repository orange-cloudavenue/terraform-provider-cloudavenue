package edgegw

import (
	"context"

	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/hashicorp/terraform-plugin-framework/diag"

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

func (rm *IPSetModel) Copy() *IPSetModel {
	x := &IPSetModel{}
	utils.ModelCopy(rm, x)
	return x
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
		ipAddrs := make([]string, 0)
		diags.Append(rm.IPAddresses.Get(ctx, &ipAddrs, false)...)
		if diags.HasError() {
			return nil, diags
		}

		values.IpAddresses = ipAddrs
	}

	return values, diags
}
