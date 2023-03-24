package network

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/mutex"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/org"
)

type staticIPPool struct {
	StartAddress types.String `tfsdk:"start_address"`
	EndAddress   types.String `tfsdk:"end_address"`
}

var staticIPPoolAttrTypes = map[string]attr.Type{
	"start_address": types.StringType,
	"end_address":   types.StringType,
}

var networkMutexKV = mutex.NewKV()

func processIPRanges(pool []staticIPPool) []govcdtypes.ExternalNetworkV2IPRange {
	subnetRng := make([]govcdtypes.ExternalNetworkV2IPRange, len(pool))
	for rangeIndex, subnetRange := range pool {
		oneRange := govcdtypes.ExternalNetworkV2IPRange{
			StartAddress: subnetRange.StartAddress.ValueString(),
			EndAddress:   subnetRange.EndAddress.ValueString(),
		}
		subnetRng[rangeIndex] = oneRange
	}
	return subnetRng
}

func getParentEdgeGatewayID(org org.Org, edgeGatewayID string) (*string, diag.Diagnostic) {
	anyEdgeGateway, err := org.GetAnyTypeEdgeGatewayById(edgeGatewayID)
	if err != nil {
		return nil, diag.NewErrorDiagnostic("error retrieving edge gateway", err.Error())
	}
	if anyEdgeGateway == nil {
		return nil, diag.NewErrorDiagnostic("error retrieving edge gateway", "edge gateway is a nil object")
	}
	id := anyEdgeGateway.EdgeGateway.OwnerRef.ID

	return &id, nil
}

// Interface for network types.
type vcdNetworkIsolatedOrRouted interface {
	SetVCDNetwork(ctx context.Context, idVDC string, data any) (*govcdtypes.OpenApiOrgVdcNetwork, diag.Diagnostics)
	//GetNetworkType() string
	SetNetworkType(network *govcdtypes.OpenApiOrgVdcNetwork, data any) (*govcdtypes.OpenApiOrgVdcNetwork, diag.Diagnostics)
	//GetVDCNetwork(string, context.Context) (*govcdtypes.OpenApiOrgVdcNetwork, diag.Diagnostics)
}

type networkCommonResourceModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Description  types.String `tfsdk:"description"`
	Gateway      types.String `tfsdk:"gateway"`
	PrefixLength types.Int64  `tfsdk:"prefix_length"`
	DNS1         types.String `tfsdk:"dns1"`
	DNS2         types.String `tfsdk:"dns2"`
	DNSSuffix    types.String `tfsdk:"dns_suffix"`
	StaticIPPool types.Set    `tfsdk:"static_ip_pool"`
}

func (r *networkIsolatedResource) SetNetworkType(network *govcdtypes.OpenApiOrgVdcNetwork, _ any) (*govcdtypes.OpenApiOrgVdcNetwork, diag.Diagnostics) {
	network.NetworkType = govcdtypes.OrgVdcNetworkTypeIsolated
	return network, diag.Diagnostics{}
}

func (r *networkRoutedResource) SetNetworkType(network *govcdtypes.OpenApiOrgVdcNetwork, data any) (*govcdtypes.OpenApiOrgVdcNetwork, diag.Diagnostics) {
	x := diag.Diagnostics{}
	d, ok := data.(networkRoutedResourceModel)
	if !ok {
		x.AddError("Error in struct networkRoutedResourceModel", "error converting data to networkRoutedResourceModel")
		return nil, x
	}
	network.NetworkType = govcdtypes.OrgVdcNetworkTypeRouted
	network.Connection = &govcdtypes.Connection{
		RouterRef: govcdtypes.OpenApiReference{
			ID: d.EdgeGatewayID.ValueString(),
		},
		// API requires interface type in upper case, but we accept any case
		ConnectionType: d.InterfaceType.ValueString(),
	}
	return network, x
}

/*
SetNetwork
This function is used to sets parameters to the network and returns the network object
networktype: Type of network to be created (e.g: govcdtypes.OrgVdcNetworkTypeIsolated or govcdtypes.OrgVdcNetworkTypeRouted)
*/
func (r *networkIsolatedResource) SetVCDNetwork(ctx context.Context, idVDC string, data any) (*govcdtypes.OpenApiOrgVdcNetwork, diag.Diagnostics) {
	x := diag.Diagnostics{}
	d, ok := data.(networkIsolatedResourceModel)
	if !ok {
		x.AddError("Error in struct networkIsolatedResourceModel", "error converting data to networkIsolatedResourceModel")
		return nil, x
	}
	// Set common parameters for a network
	commonData := networkCommonResourceModel{
		ID:           d.ID,
		Name:         d.Name,
		Description:  d.Description,
		Gateway:      d.Gateway,
		PrefixLength: d.PrefixLength,
		DNS1:         d.DNS1,
		DNS2:         d.DNS2,
		DNSSuffix:    d.DNSSuffix,
		StaticIPPool: d.StaticIPPool,
	}
	// Set these commom parameters to the object network
	networkType, diag := commonVCDNetwork(id, ctx, commonData)
	// Add specific parameters for an isolated network
	networkType = r.SetNetworkType(networkType, data)
	return networkType, diag
}

/*
SetNetwork
This function is used to sets parameters to the network and returns the network object
networktype: Type of network to be created (: govcdtypes.OrgVdcNetworkTypeRouted
*/
func (r *networkRoutedResource) SetVCDNetwork(id string, ctx context.Context, data any) (*govcdtypes.OpenApiOrgVdcNetwork, diag.Diagnostics) {
	x := diag.Diagnostics{}
	d, ok := data.(networkIsolatedResourceModel)
	if !ok {
		x.AddError("Error in struct networkRoutedResourceModel", "error converting data to networkRoutedResourceModel")
		return nil, x
	}
	// Set common parameters for a network
	commonData := networkCommonResourceModel{
		ID:           d.ID,
		Name:         d.Name,
		Description:  d.Description,
		Gateway:      d.Gateway,
		PrefixLength: d.PrefixLength,
		DNS1:         d.DNS1,
		DNS2:         d.DNS2,
		DNSSuffix:    d.DNSSuffix,
		StaticIPPool: d.StaticIPPool,
	}
	// Set these commom parameters to the object network
	networkType, diag := commonVCDNetwork(id, ctx, commonData)
	// Add specific parameters for an routed network
	networkType = r.SetNetworkType(networkType, data)
	return networkType, diag
}

func commonVCDNetwork(id string, ctx context.Context, data networkCommonResourceModel) (*govcdtypes.OpenApiOrgVdcNetwork, diag.Diagnostics) {
	ipPool := []staticIPPool{}
	diag := diag.Diagnostics{}
	diag.Append(data.StaticIPPool.ElementsAs(ctx, &ipPool, true)...)
	networkType := &govcdtypes.OpenApiOrgVdcNetwork{
		Name:        data.Name.ValueString(),
		Description: data.Description.ValueString(),
		OwnerRef:    &govcdtypes.OpenApiReference{ID: id},
		Subnets: govcdtypes.OrgVdcNetworkSubnets{
			Values: []govcdtypes.OrgVdcNetworkSubnetValues{
				{
					Gateway:      data.Gateway.ValueString(),
					PrefixLength: int(data.PrefixLength.ValueInt64()),
					IPRanges: govcdtypes.OrgVdcNetworkSubnetIPRanges{
						Values: processIPRanges(ipPool),
					},
					DNSServer1: data.DNS1.ValueString(),
					DNSServer2: data.DNS2.ValueString(),
					DNSSuffix:  data.DNSSuffix.ValueString(),
				},
			},
		},
	}
	return networkType, diag
}

//func optVCDNetwork(data networkCommonResourceModel) (network *govcdtypes.OpenApiOrgVdcNetwork) {
//	switch data[0].(type) {
//	case *networkIsolatedResourceModel:
//		myshared := false // Cloudavenue does not support shared networks
//		network = &govcdtypes.OpenApiOrgVdcNetwork{
//			NetworkType: govcdtypes.OrgVdcNetworkTypeIsolated,
//			Shared:      &myshared,
//		}
//	case *networkRoutedResourceModel:
//		network = &govcdtypes.OpenApiOrgVdcNetwork{
//			NetworkType: govcdtypes.OrgVdcNetworkTypeRouted,
//			Connection: &govcdtypes.Connection{
//				RouterRef: govcdtypes.OpenApiReference{
//					ID: data.EdgeGatewayID.ValueString(),
//				},
//				// API requires interface type in upper case, but we accept any case
//				ConnectionType: data.InterfaceType.ValueString(),
//			},
//		}
//	}
//	return network
//}
