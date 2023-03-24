package network

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"
)

type Common struct {
	TypeOfNetwork Type
}

type CommonResourceModel struct {
	VDCOrParentEdgeGatewayID types.String

	// BASE
	ID           types.String
	Name         types.String
	Description  types.String
	Gateway      types.String
	PrefixLength types.Int64
	DNS1         types.String
	DNS2         types.String
	DNSSuffix    types.String
	StaticIPPool types.Set

	// ISOLATED

	// ROUTED
}

type Network interface {
	ConstructNetworkAPIObject(ctx context.Context, plan any) (*govcdtypes.OpenApiOrgVdcNetwork, diag.Diagnostics)
}

type StaticIPPool struct {
	StartAddress types.String `tfsdk:"start_address"`
	EndAddress   types.String `tfsdk:"end_address"`
}

func (c Common) ConstructNetworkAPIObject(ctx context.Context, data CommonResourceModel) (*govcdtypes.OpenApiOrgVdcNetwork, diag.Diagnostics) {
	apiObject, d := data.constructBaseNetworkAPIObject(context.Background())
	if d.HasError() {
		return nil, d
	}

	switch c.TypeOfNetwork {
	case ISOLATED:
		apiObject.NetworkType = govcdtypes.OrgVdcNetworkTypeIsolated
		//
	case NAT_ROUTED:
		apiObject.NetworkType = govcdtypes.OrgVdcNetworkTypeRouted
	}

	return apiObject, nil
}

// constructBaseNetworkAPIObject constructs the base network object with common attributes.
func (r CommonResourceModel) constructBaseNetworkAPIObject(ctx context.Context) (*govcdtypes.OpenApiOrgVdcNetwork, diag.Diagnostics) {
	ipRanges, d := r.processIPRanges(ctx)
	if d.HasError() {
		return nil, d
	}

	return &govcdtypes.OpenApiOrgVdcNetwork{
		Name:        r.Name.ValueString(),
		Description: r.Description.ValueString(),
		OwnerRef:    &govcdtypes.OpenApiReference{ID: r.VDCOrParentEdgeGatewayID.ValueString()},
		Subnets: govcdtypes.OrgVdcNetworkSubnets{
			Values: []govcdtypes.OrgVdcNetworkSubnetValues{
				{
					Gateway:      r.Gateway.ValueString(),
					PrefixLength: int(r.PrefixLength.ValueInt64()),
					IPRanges: govcdtypes.OrgVdcNetworkSubnetIPRanges{
						Values: ipRanges,
					},
					DNSServer1: r.DNS1.ValueString(),
					DNSServer2: r.DNS2.ValueString(),
					DNSSuffix:  r.DNSSuffix.ValueString(),
				},
			},
		},
	}, d
}

func (r CommonResourceModel) processIPRanges(ctx context.Context) ([]govcdtypes.ExternalNetworkV2IPRange, diag.Diagnostics) {
	d := diag.Diagnostics{}

	var (
		ipPool    = []StaticIPPool{}
		subnetRng = make([]govcdtypes.ExternalNetworkV2IPRange, len(ipPool))
	)
	d.Append(r.StaticIPPool.ElementsAs(ctx, &ipPool, true)...)
	if d.HasError() {
		return nil, d
	}

	for rangeIndex, subnetRange := range ipPool {
		oneRange := govcdtypes.ExternalNetworkV2IPRange{
			StartAddress: subnetRange.StartAddress.ValueString(),
			EndAddress:   subnetRange.EndAddress.ValueString(),
		}
		subnetRng[rangeIndex] = oneRange
	}
	return subnetRng, d
}
