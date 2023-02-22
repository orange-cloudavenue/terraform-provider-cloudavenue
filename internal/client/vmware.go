// Package client is the main client for the CloudAvenue provider.
package client

import (
	"fmt"

	"github.com/vmware/go-vcloud-director/v2/govcd"
)

// GetOrgAndVDC finds a pair of org and vdc using the names provided
// in the args. If the names are empty, it will use the default
// org and vdc names from the provider.
func (c *CloudAvenue) GetOrgAndVDC(orgName, vdcName string) (org *govcd.Org, vdc *govcd.Vdc, err error) {
	if orgName == "" {
		return nil, nil, fmt.Errorf("empty Org name provided")
	}
	if vdcName == "" {
		return nil, nil, fmt.Errorf("empty VDC name provided")
	}

	org, err = c.Vmware.GetOrgByName(orgName)
	if err != nil {
		return nil, nil, fmt.Errorf("error retrieving Org %s: %w", orgName, err)
	}

	if org.Org.Name == "" || org.Org.HREF == "" || org.Org.ID == "" {
		return nil, nil, fmt.Errorf("empty Org %s found ", orgName)
	}

	vdc, err = org.GetVDCByName(vdcName, false)
	if err != nil {
		return nil, nil, fmt.Errorf("error retrieving VDC %s: %w", vdcName, err)
	}

	if vdc == nil || vdc.Vdc.ID == "" || vdc.Vdc.HREF == "" || vdc.Vdc.Name == "" {
		return nil, nil, fmt.Errorf("error retrieving VDC %s: not found", vdcName)
	}

	return org, vdc, err
}
