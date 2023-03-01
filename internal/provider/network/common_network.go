package network

import (
	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/mutex"
)

var networkMutexKV = mutex.NewKV()

// TODO: refactor -> go to common
type diagnosticError struct {
	Summary string
	Detail  string
}

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

func getParentEdgeGatewayID(org *govcd.Org, edgeGatewayID string) (*string, *diagnosticError) {
	anyEdgeGateway, err := org.GetAnyTypeEdgeGatewayById(edgeGatewayID)
	if err != nil {
		return nil, &diagnosticError{Summary: "Error retrieving edge gateway", Detail: err.Error()}
	}
	if anyEdgeGateway == nil {
		return nil, &diagnosticError{Summary: "Edge gateway not found", Detail: "anyEdgeGateway object is nil"}
	}
	id := anyEdgeGateway.EdgeGateway.OwnerRef.ID

	return &id, nil
}
