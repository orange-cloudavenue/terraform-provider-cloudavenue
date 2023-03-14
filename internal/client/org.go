package client

import (
	"fmt"

	"github.com/vmware/go-vcloud-director/v2/govcd"
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

// GetOrg return the org using the name provided in the provider.
func (c *CloudAvenue) GetOrg() (*Org, error) {
	x, err := c.Vmware.GetOrgByName(c.Org)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrRetrievingOrg, err)
	}

	return &Org{x}, nil
}
