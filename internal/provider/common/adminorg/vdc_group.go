package adminorg

import (
	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/uuid"
)

// GetVDCGroupByNameOrID returns the VDC group using the name or ID provided in the argument.
func (ao *AdminOrg) GetVDCGroupByNameOrID(nameOrID string) (*govcd.VdcGroup, error) {
	if uuid.IsVDCGroup(nameOrID) {
		return ao.GetVdcGroupById(nameOrID)
	}
	return ao.GetVdcGroupByName(nameOrID)
}
