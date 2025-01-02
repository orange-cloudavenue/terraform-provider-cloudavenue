package edgegw

import (
	"context"
	"fmt"

	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

type StaticRouteModel struct {
	Description     supertypes.StringValue    `tfsdk:"description"`
	EdgeGatewayID   supertypes.StringValue    `tfsdk:"edge_gateway_id"`
	EdgeGatewayName supertypes.StringValue    `tfsdk:"edge_gateway_name"`
	ID              supertypes.StringValue    `tfsdk:"id"`
	Name            supertypes.StringValue    `tfsdk:"name"`
	NetworkCidr     supertypes.StringValue    `tfsdk:"network_cidr"`
	NextHops        supertypes.SetNestedValue `tfsdk:"next_hops"`
}

// * NextHops.
type StaticRouteModelNextHops []StaticRouteModelNextHop

// * NextHop.
type StaticRouteModelNextHop struct {
	AdminDistance supertypes.Int64Value  `tfsdk:"admin_distance"`
	IPAddress     supertypes.StringValue `tfsdk:"ip_address"`
}

func NewStaticRoute(t any) *StaticRouteModel {
	switch x := t.(type) {
	case tfsdk.State:
		return &StaticRouteModel{
			Description:     supertypes.NewStringNull(),
			EdgeGatewayID:   supertypes.NewStringUnknown(),
			EdgeGatewayName: supertypes.NewStringUnknown(),
			ID:              supertypes.NewStringUnknown(),
			Name:            supertypes.NewStringNull(),
			NetworkCidr:     supertypes.NewStringNull(),
			NextHops:        supertypes.NewSetNestedNull(x.Schema.GetAttributes()["next_hops"].GetType().(supertypes.SetNestedType).ElementType()),
		}

	case tfsdk.Plan:
		return &StaticRouteModel{
			Description:     supertypes.NewStringNull(),
			EdgeGatewayID:   supertypes.NewStringUnknown(),
			EdgeGatewayName: supertypes.NewStringUnknown(),
			ID:              supertypes.NewStringUnknown(),
			Name:            supertypes.NewStringNull(),
			NetworkCidr:     supertypes.NewStringNull(),
			NextHops:        supertypes.NewSetNestedNull(x.Schema.GetAttributes()["next_hops"].GetType().(supertypes.SetNestedType).ElementType()),
		}

	case tfsdk.Config:
		return &StaticRouteModel{
			Description:     supertypes.NewStringNull(),
			EdgeGatewayID:   supertypes.NewStringUnknown(),
			EdgeGatewayName: supertypes.NewStringUnknown(),
			ID:              supertypes.NewStringUnknown(),
			Name:            supertypes.NewStringNull(),
			NetworkCidr:     supertypes.NewStringNull(),
			NextHops:        supertypes.NewSetNestedNull(x.Schema.GetAttributes()["next_hops"].GetType().(supertypes.SetNestedType).ElementType()),
		}

	default:
		panic(fmt.Sprintf("unexpected type %T", t))
	}
}

func (rm *StaticRouteModel) Copy() *StaticRouteModel {
	x := &StaticRouteModel{}
	utils.ModelCopy(rm, x)
	return x
}

// GetNextHops returns the value of the NextHops field.
func (rm *StaticRouteModel) GetNextHops(ctx context.Context) (values StaticRouteModelNextHops, diags diag.Diagnostics) {
	values = make(StaticRouteModelNextHops, 0)
	d := rm.NextHops.Get(ctx, &values, false)
	return values, d
}

// * CustomFuncs

func (rm *StaticRouteModel) ToNsxtEdgeGatewayStaticRoute(ctx context.Context) (*govcdtypes.NsxtEdgeGatewayStaticRoute, diag.Diagnostics) {
	staticRouteConfig := &govcdtypes.NsxtEdgeGatewayStaticRoute{
		Name:        rm.Name.Get(),
		Description: rm.Description.Get(),
		NetworkCidr: rm.NetworkCidr.Get(),
		NextHops:    make([]govcdtypes.NsxtEdgeGatewayStaticRouteNextHops, 0),
	}

	nextHops, d := rm.GetNextHops(ctx)
	if d.HasError() {
		return nil, d
	}

	for _, nextHop := range nextHops {
		nH := govcdtypes.NsxtEdgeGatewayStaticRouteNextHops{
			IPAddress:     nextHop.IPAddress.Get(),
			AdminDistance: nextHop.AdminDistance.GetInt(),
		}

		staticRouteConfig.NextHops = append(staticRouteConfig.NextHops, nH)
	}

	return staticRouteConfig, d
}
