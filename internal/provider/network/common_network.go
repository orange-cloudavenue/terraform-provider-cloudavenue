package network

import (
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
