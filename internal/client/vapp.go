package client

import "github.com/vmware/go-vcloud-director/v2/govcd"

type VAPP struct {
	*govcd.VApp
}

// GetName give you the name of the vApp.
func (v VAPP) GetName() string {
	return v.VApp.VApp.Name
}

// GetID give you the ID of the vApp.
func (v VAPP) GetID() string {
	return v.VApp.VApp.ID
}

// GetStatusCode give you the status code of the vApp.
func (v VAPP) GetStatusCode() int {
	return v.VApp.VApp.Status
}

// GetHREF give you the HREF of the vApp.
func (v VAPP) GetHREF() string {
	return v.VApp.VApp.HREF
}

// GetDescription give you the status code of the vApp.
func (v VAPP) GetDescription() string {
	return v.VApp.VApp.Description
}
