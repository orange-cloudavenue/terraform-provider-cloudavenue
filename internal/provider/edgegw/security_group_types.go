package edgegw

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// * Security Group (singular) model.
type securityGroupModelMemberOrgNetworkIDs []string

type securityGroupModel struct {
	ID                  types.String `tfsdk:"id"`
	EdgeGatewayID       types.String `tfsdk:"edge_gateway_id"`
	EdgeGatewayName     types.String `tfsdk:"edge_gateway_name"`
	Name                types.String `tfsdk:"name"`
	Description         types.String `tfsdk:"description"`
	MemberOrgNetworkIDs types.Set    `tfsdk:"member_org_network_ids"`
}

// GetIDOrName returns the ID or the name of the security group.
func (rm *securityGroupModel) GetIDOrName() types.String {
	if rm.ID.IsNull() || rm.ID.IsUnknown() {
		return rm.Name
	}

	return rm.ID
}

// MemberOrgNetworkIDsFromPlan returns the member_org_network_ids from the plan.
func (rm *securityGroupModel) MemberOrgNetworkIDsFromPlan(ctx context.Context) (securityGroupModelMemberOrgNetworkIDs, diag.Diagnostics) {
	ids := securityGroupModelMemberOrgNetworkIDs{}
	return ids, rm.MemberOrgNetworkIDs.ElementsAs(ctx, &ids, false)
}
