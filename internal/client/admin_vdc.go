package client

import (
	"fmt"

	"github.com/vmware/go-vcloud-director/v2/govcd"
)

type AdminVDC struct {
	name string
	*govcd.AdminVdc
}

func (c *CloudAvenue) DefaultAdminVDCExist() bool {
	return c.VDC != ""
}

type GetAdminVDCOpts func(*AdminVDC)

func WithAdminVDCName(name string) GetAdminVDCOpts {
	return func(AdminVdc *AdminVDC) {
		AdminVdc.name = name
	}
}

func (v AdminVDC) GetName() string {
	return v.AdminVdc.AdminVdc.Name
}

func (v AdminVDC) GetID() string {
	return v.AdminVdc.AdminVdc.ID
}

// GetAdminVdc return the admin vdc using the name provided in the provider.
func (c CloudAvenue) GetAdminVDC(opts ...GetAdminVDCOpts) (*AdminVDC, error) {
	v := &AdminVDC{}

	for _, opt := range opts {
		opt(v)
	}

	if v.name == "" {
		if c.DefaultVDCExist() {
			v.name = c.GetDefaultVDC()
		} else {
			return nil, fmt.Errorf("%w", ErrEmptyVDCNameProvided)
		}
	}

	org, err := c.GetAdminOrg()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrRetrievingOrgAdmin, err)
	}

	v.AdminVdc, err = org.GetAdminVDCByName(v.name, false)
	if err != nil {
		return nil, fmt.Errorf("%w: %s %w", ErrRetrievingAdminVDC, v.name, err)
	}

	return v, nil
}
