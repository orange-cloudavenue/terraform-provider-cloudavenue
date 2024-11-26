package org

import (
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/edgegw"
)

var _ edgegw.Handler = (*Org)(nil)

// GetEdgeGateway returns the edge gateway.
func (o *Org) GetEdgeGateway(egw edgegw.BaseEdgeGW) (edgegw.EdgeGateway, error) {
	if egw.GetIDOrName() == "" {
		return edgegw.EdgeGateway{}, edgegw.ErrEdgeGatewayIDOrNameIsEmpty
	}

	edge, err := o.c.CAVSDK.V1.EdgeGateway.Get(egw.GetIDOrName())
	if err != nil {
		return edgegw.EdgeGateway{}, err
	}

	vmwareEdgeGateway, err := edge.GetVmwareEdgeGateway()
	if err != nil {
		return edgegw.EdgeGateway{}, err
	}

	return edgegw.EdgeGateway{
		// Client is the CloudAvenue client.
		Client: o.c,

		// EdgeClient is the EdgeGateway client.
		EdgeClient: edge,

		// NsxtEdgeGateway is the NSX-T edge gateway.
		//
		// Deprecated: Use EdgeClient instead.
		NsxtEdgeGateway: vmwareEdgeGateway,
	}, err
}
