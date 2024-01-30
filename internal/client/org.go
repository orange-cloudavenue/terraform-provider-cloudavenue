package client

import (
	"fmt"

	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"
)

type Org struct {
	*govcd.Org
}

// Returns the name of the organization.
func (o *Org) GetName() string {
	return o.Org.Org.Name
}

// Returns the ID of the organization.
func (o *Org) GetID() string {
	return o.Org.Org.ID
}

// GetNetworkDHCP returns the DHCP object for the org network provided in parameter.
func (o *Org) GetNetworkDHCP(orgNetworkID string) (*govcd.OpenApiOrgVdcNetworkDhcp, error) {
	if err := o.Refresh(); err != nil {
		return nil, err
	}

	orgNetwork, err := o.GetOpenApiOrgVdcNetworkById(orgNetworkID)
	if err != nil {
		return nil, err
	}

	return orgNetwork.GetOpenApiOrgVdcNetworkDhcp()
}

// UpdateNetworkDHCP updates the DHCP object for the org network provided in parameter.
func (o *Org) UpdateNetworkDHCP(orgNetworkID string, dhcpParameters *govcdtypes.OpenApiOrgVdcNetworkDhcp) error {
	orgNetwork, err := o.GetOpenApiOrgVdcNetworkById(orgNetworkID)
	if err != nil {
		return err
	}

	_, err = orgNetwork.UpdateDhcp(dhcpParameters)
	return err
}

// DeleteNetworkDHCP deletes the DHCP object for the org network provided in parameter.
func (o *Org) DeleteNetworkDHCP(orgNetworkID string) error {
	orgNetwork, err := o.GetOpenApiOrgVdcNetworkById(orgNetworkID)
	if err != nil {
		return err
	}

	return orgNetwork.DeletNetworkDhcp()
}

// GetOrg return the org using the name provided in the provider.
func (c *CloudAvenue) GetOrg() (*Org, error) {
	x, err := c.Vmware.GetOrgByName(c.GetOrgName())
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrRetrievingOrg, err)
	}

	return &Org{x}, nil
}
