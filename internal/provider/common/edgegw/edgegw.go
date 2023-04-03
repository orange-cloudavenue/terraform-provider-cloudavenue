package edgegw

import (
	"context"
	"fmt"

	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/mutex"
)

var ErrEdgeGatewayIDOrNameIsEmpty = fmt.Errorf("edge gateway ID or name is empty")

var gwMutexKV = mutex.NewKV()

type Handler interface {
	// GetEdgeGateway allows retrieving NSX-T edge gateway by ID Or Name.
	GetEdgeGateway(egw BaseEdgeGW) (EdgeGateway, error)
}

type BaseEdgeGW struct {
	ID   types.String
	Name types.String
}

type EdgeGateway struct {
	Client *client.CloudAvenue
	*govcd.NsxtEdgeGateway
}

// GetID returns the ID of the Edge Gateway.
func (e BaseEdgeGW) GetID() string {
	return e.ID.ValueString()
}

// GetName returns the Edge Gateway Name.
func (e BaseEdgeGW) GetName() string {
	return e.Name.ValueString()
}

// GetIDOrName returns the ID or the name of the Edge Gateway.
func (e BaseEdgeGW) GetIDOrName() string {
	if e.GetID() != "" {
		return e.GetID()
	}
	return e.GetName()
}

// GetParent returns the parent of the Edge Gateway.
// The parent is VDC Or VDC Group object.
func (e EdgeGateway) GetParent() (client.VDCOrVDCGroupHandler, error) {
	return e.Client.GetVDCOrVDCGroup(e.EdgeGateway.OwnerRef.Name)
}

// GetName returns the name of the Edge Gateway.
func (e EdgeGateway) GetName() string {
	return e.EdgeGateway.Name
}

// GetID returns the ID of the Edge Gateway.
func (e EdgeGateway) GetID() string {
	return e.EdgeGateway.ID
}

// Lock locks the Edge Gateway.
func (e EdgeGateway) Lock(ctx context.Context) {
	gwMutexKV.KvLock(ctx, e.GetID())
}

// Unlock unlocks the Edge Gateway.
func (e EdgeGateway) Unlock(ctx context.Context) {
	gwMutexKV.KvUnlock(ctx, e.GetID())
}
