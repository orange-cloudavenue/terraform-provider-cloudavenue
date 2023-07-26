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

type DhcpForwardingModel struct {
	DhcpServers     supertypes.SetValue    `tfsdk:"dhcp_servers"`
	EdgeGatewayID   supertypes.StringValue `tfsdk:"edge_gateway_id"`
	EdgeGatewayName supertypes.StringValue `tfsdk:"edge_gateway_name"`
	Enabled         supertypes.BoolValue   `tfsdk:"enabled"`
	ID              supertypes.StringValue `tfsdk:"id"`
}

type DhcpForwardingModelDhcpServers []supertypes.StringValue

func NewDhcpForwarding(t any) *DhcpForwardingModel {
	switch x := t.(type) {
	case tfsdk.State:
		return &DhcpForwardingModel{
			DhcpServers:     supertypes.NewSetNull(x.Schema.GetAttributes()["dhcp_servers"].GetType().(supertypes.SetType).ElementType()),
			EdgeGatewayID:   supertypes.NewStringUnknown(),
			EdgeGatewayName: supertypes.NewStringUnknown(),
			Enabled:         supertypes.NewBoolUnknown(),
			ID:              supertypes.NewStringUnknown(),
		}

	case tfsdk.Plan:
		return &DhcpForwardingModel{
			DhcpServers:     supertypes.NewSetNull(x.Schema.GetAttributes()["dhcp_servers"].GetType().(supertypes.SetType).ElementType()),
			EdgeGatewayID:   supertypes.NewStringUnknown(),
			EdgeGatewayName: supertypes.NewStringUnknown(),
			Enabled:         supertypes.NewBoolUnknown(),
			ID:              supertypes.NewStringUnknown(),
		}

	case tfsdk.Config:
		return &DhcpForwardingModel{
			DhcpServers:     supertypes.NewSetNull(x.Schema.GetAttributes()["dhcp_servers"].GetType().(supertypes.SetType).ElementType()),
			EdgeGatewayID:   supertypes.NewStringUnknown(),
			EdgeGatewayName: supertypes.NewStringUnknown(),
			Enabled:         supertypes.NewBoolUnknown(),
			ID:              supertypes.NewStringUnknown(),
		}

	default:
		panic(fmt.Sprintf("unexpected type %T", t))
	}
}

func (rm *DhcpForwardingModel) Copy() *DhcpForwardingModel {
	x := &DhcpForwardingModel{}
	utils.ModelCopy(rm, x)
	return x
}

// GetDhcpServers returns the value of the DhcpServers field.
func (rm *DhcpForwardingModel) GetDhcpServers(ctx context.Context) (values DhcpForwardingModelDhcpServers, diags diag.Diagnostics) {
	values = make(DhcpForwardingModelDhcpServers, 0)
	d := rm.DhcpServers.Get(ctx, &values, false)
	return values, d
}

func (r *DhcpForwardingModelDhcpServers) Get() []string {
	return utils.SuperSliceTypesStringToSliceString(*r)
}

// ToNsxtEdgeGatewayDhcpForwarder returns the NSX-T Edge Gateway DHCP Forwarder representation of the model.
func (rm *DhcpForwardingModel) ToNsxtEdgeGatewayDhcpForwarder(ctx context.Context) (*govcdtypes.NsxtEdgeGatewayDhcpForwarder, diag.Diagnostics) {
	dhcpServers, d := rm.GetDhcpServers(ctx)
	if d.HasError() {
		return nil, d
	}

	return &govcdtypes.NsxtEdgeGatewayDhcpForwarder{
		Enabled:     rm.Enabled.Get(),
		DhcpServers: dhcpServers.Get(),
	}, nil
}
