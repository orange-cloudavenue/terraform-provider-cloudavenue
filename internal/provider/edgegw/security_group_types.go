package edgegw

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"

	v1 "github.com/orange-cloudavenue/cloudavenue-sdk-go/v1"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

type SecurityGroupModel struct {
	ID              supertypes.StringValue        `tfsdk:"id"`
	Name            supertypes.StringValue        `tfsdk:"name"`
	Description     supertypes.StringValue        `tfsdk:"description"`
	EdgeGatewayName supertypes.StringValue        `tfsdk:"edge_gateway_name"`
	EdgeGatewayID   supertypes.StringValue        `tfsdk:"edge_gateway_id"`
	Members         supertypes.SetValueOf[string] `tfsdk:"member_org_network_ids"`
}

func (rm *SecurityGroupModel) ToSDKSecurityGroupModel(ctx context.Context) (*v1.FirewallGroupSecurityGroupModel, diag.Diagnostics) {
	members, d := rm.Members.Get(ctx)
	if d.HasError() {
		return nil, d
	}

	return &v1.FirewallGroupSecurityGroupModel{
		FirewallGroupModel: v1.FirewallGroupModel{
			ID:          rm.ID.Get(),
			Name:        rm.Name.Get(),
			Description: rm.Description.Get(),
		},
		Members: utils.SliceIDToOpenAPIReference(members),
	}, nil
}

func (rm *SecurityGroupModel) Copy() *SecurityGroupModel {
	x := &SecurityGroupModel{}
	utils.ModelCopy(rm, x)
	return x
}
