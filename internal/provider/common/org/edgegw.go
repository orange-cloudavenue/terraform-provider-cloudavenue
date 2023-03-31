package org

import (
	"context"

	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/edgegw"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/mutex"
)

var _ edgegw.Handler = (*Org)(nil)

var gwMutexKV = mutex.NewKV()

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

// GetEdgeGatewayID returns the edge gateway ID.
func (o *Org) GetEdgeGatewayID(egw edgegw.BaseEdgeGW) (string, error) {
	if egw.GetIDOrName() == "" {
		return "", edgegw.ErrEdgeGatewayIDOrNameIsEmpty
	}

	if egw.GetID() != "" {
		return egw.GetID(), nil
	}

	edgeGw, err := o.GetNsxtEdgeGatewayByName(egw.GetName())
	if err != nil {
		return "", err
	}

	return edgeGw.EdgeGateway.ID, nil
}

// LockParentEdgeGW locks the parent Edge Gateway.
func (o *Org) LockParentEdgeGW(ctx context.Context, egw edgegw.BaseEdgeGW) (d diag.Diagnostics) {
	key, err := o.GetEdgeGatewayID(egw)
	if err != nil {
		d.AddError("Unable to find Edge Gateway", err.Error())
		return
	}

	gwMutexKV.KvLock(ctx, key)
	return
}

// UnlockParentEdgeGW unlocks the Edge Gateway.
func (o *Org) UnlockParentEdgeGW(ctx context.Context, egw edgegw.BaseEdgeGW) (d diag.Diagnostics) {
	key, err := o.GetEdgeGatewayID(egw)
	if err != nil {
		d.AddError("Unable to find Edge Gateway", err.Error())
		return
	}

	gwMutexKV.KvUnlock(ctx, key)
	return
}
