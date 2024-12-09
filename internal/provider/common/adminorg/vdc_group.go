package adminorg

import (
	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
)

// GetVDCGroupByNameOrID returns the VDC group using the name or ID provided in the argument.
// Deprecated: Use GetVdcGroupByName or GetVdcGroupById instead.
func (ao *AdminOrg) GetVDCGroupByNameOrID(nameOrID string) (*govcd.VdcGroup, error) {
	if urn.IsVDCGroup(nameOrID) {
		return ao.GetVdcGroupById(nameOrID)
	}
	return ao.GetVdcGroupByName(nameOrID)
}
