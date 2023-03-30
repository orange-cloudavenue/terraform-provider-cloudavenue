package network

import (
	"context"

	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Kind struct {
	TypeOfNetwork Type
}

type GlobalResourceModel struct {
	// BASE
	ID                types.String
	Name              types.String
	Description       types.String
	Gateway           types.String
	PrefixLength      types.Int64
	DNS1              types.String
	DNS2              types.String
	DNSSuffix         types.String
	StaticIPPool      types.Set
	VDCIDOrVDCGroupID types.String

	// ISOLATED

	// ROUTED
	EdgeGatewayID types.String
	InterfaceType types.String
}

type Network interface {
	SetNetworkAPIObject(ctx context.Context, plan any) (*govcdtypes.OpenApiOrgVdcNetwork, diag.Diagnostics)
}

type StaticIPPool struct {
	StartAddress types.String `tfsdk:"start_address"`
	EndAddress   types.String `tfsdk:"end_address"`
}

// SetNetowrkAPIObject set the network object.
func (k Kind) SetNetworkAPIObject(ctx context.Context, data GlobalResourceModel) (*govcdtypes.OpenApiOrgVdcNetwork, diag.Diagnostics) {
	apiObject, d := data.setBaseNetworkAPIObject(context.Background())
	if d.HasError() {
		return nil, d
	}

	// Define the particular attributes according to the type of network
	switch k.TypeOfNetwork {
	case ISOLATED:
		apiObject.NetworkType = govcdtypes.OrgVdcNetworkTypeIsolated
		myshared := false // Cloudavenue does not support shared networks
		apiObject.Shared = &myshared

	case NAT_ROUTED:
		apiObject.NetworkType = govcdtypes.OrgVdcNetworkTypeRouted
		apiObject.Connection = &govcdtypes.Connection{
			RouterRef: govcdtypes.OpenApiReference{
				ID: data.EdgeGatewayID.ValueString(),
			},
			// API requires interface type in upper case, but we accept any case
			ConnectionType: data.InterfaceType.ValueString(),
		}
	}
	return apiObject, nil
}

// setBaseNetworkAPIObject set the base network object with common attributes.
func (g GlobalResourceModel) setBaseNetworkAPIObject(ctx context.Context) (*govcdtypes.OpenApiOrgVdcNetwork, diag.Diagnostics) {
	ipRanges, d := g.setIPRanges(ctx)
	if d.HasError() {
		return nil, d
	}

	return &govcdtypes.OpenApiOrgVdcNetwork{
		ID:          g.ID.ValueString(),
		Name:        g.Name.ValueString(),
		Description: g.Description.ValueString(),
		OwnerRef:    &govcdtypes.OpenApiReference{ID: g.VDCIDOrVDCGroupID.ValueString()},
		Subnets: govcdtypes.OrgVdcNetworkSubnets{
			Values: []govcdtypes.OrgVdcNetworkSubnetValues{
				{
					Gateway:      g.Gateway.ValueString(),
					PrefixLength: int(g.PrefixLength.ValueInt64()),
					IPRanges: govcdtypes.OrgVdcNetworkSubnetIPRanges{
						Values: ipRanges,
					},
					DNSServer1: g.DNS1.ValueString(),
					DNSServer2: g.DNS2.ValueString(),
					DNSSuffix:  g.DNSSuffix.ValueString(),
				},
			},
		},
	}, d
}

// Set Ip pool static for network.
func (g GlobalResourceModel) setIPRanges(ctx context.Context) ([]govcdtypes.ExternalNetworkV2IPRange, diag.Diagnostics) {
	var (
		d      = diag.Diagnostics{}
		ipPool = []StaticIPPool{}
	)

	d.Append(g.StaticIPPool.ElementsAs(ctx, &ipPool, true)...)
	if d.HasError() {
		return nil, d
	}
	subnetRng := make([]govcdtypes.ExternalNetworkV2IPRange, len(ipPool))

	for rangeIndex, subnetRange := range ipPool {
		oneRange := govcdtypes.ExternalNetworkV2IPRange{
			StartAddress: subnetRange.StartAddress.ValueString(),
			EndAddress:   subnetRange.EndAddress.ValueString(),
		}
		subnetRng[rangeIndex] = oneRange
	}
	return subnetRng, d
}
