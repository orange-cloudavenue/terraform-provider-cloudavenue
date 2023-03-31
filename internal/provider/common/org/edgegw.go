package org

import (
	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/edgegw"
)

var _ edgegw.Handler = (*Org)(nil)

// GetEdgeGateway returns the edge gateway.
func (o *Org) GetEdgeGateway(egw edgegw.BaseEdgeGW) (edgegw.EdgeGateway, error) {
	if egw.GetIDOrName() == "" {
		return edgegw.EdgeGateway{}, edgegw.ErrEdgeGatewayIDOrNameIsEmpty
	}

	var (
		govcdValues *govcd.NsxtEdgeGateway
		err         error
	)

	if egw.GetID() != "" {
		govcdValues, err = o.GetNsxtEdgeGatewayById(egw.GetID())
	} else {
		govcdValues, err = o.GetNsxtEdgeGatewayByName(egw.GetName())
	}
	if err != nil {
		return edgegw.EdgeGateway{}, err
	}

	return edgegw.EdgeGateway{
		Client:          o.c,
		NsxtEdgeGateway: govcdValues,
	}, err
}
