// Package client is the main client for the CloudAvenue provider.
package client

import (
	"errors"
	"fmt"

	"github.com/vmware/go-vcloud-director/v2/govcd"
)

var (
	ErrEmptyOrgNameProvided = errors.New("empty Org name provided")
	ErrEmptyVDCNameProvided = errors.New("empty VDC name provided")
	ErrRetrievingOrg        = errors.New("error retrieving Org")
	ErrRetrievingOrgAdmin   = errors.New("error retrieving Org admin")
	ErrEmptyOrgFound        = errors.New("empty Org found")
	ErrRetrievingVDC        = errors.New("error retrieving VDC")
	ErrRetrievingVDCGroup   = errors.New("error retrieving VDC Group")
	ErrEmptyVDCFound        = errors.New("empty VDC found")
)

// VDCOrVDCGroupHandler is an interface to access some common methods on VDC or VDC Group without
// explicitly handling exact types.
type VDCOrVDCGroupHandler interface {
	GetOpenApiOrgVdcNetworkByName(string) (*govcd.OpenApiOrgVdcNetwork, error)
}

// GetOrgAndVDC finds a pair of org and vdc using the names provided
// in the args. If the names are empty, it will use the default
// org and vdc names from the provider.
func (c *CloudAvenue) GetOrgAndVDC(orgName, vdcName string) (org *govcd.Org, vdcOrVDCGroup VDCOrVDCGroupHandler, err error) {
	if orgName == "" {
		return nil, nil, fmt.Errorf("%w", ErrEmptyOrgNameProvided)
	}
	if vdcName == "" {
		return nil, nil, fmt.Errorf("%w", ErrEmptyVDCNameProvided)
	}

	org, err = c.Vmware.GetOrgByName(orgName)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %w", ErrRetrievingOrg, err)
	}

	if org.Org.Name == "" || org.Org.HREF == "" || org.Org.ID == "" {
		return nil, nil, fmt.Errorf("%w : %s", ErrEmptyOrgFound, orgName)
	}

	vdcOrVDCGroup, err = org.GetVDCByName(vdcName, false)
	if err != nil && !govcd.ContainsNotFound(err) {
		return nil, nil, fmt.Errorf("%w: %s %w", ErrRetrievingVDC, vdcName, err)
	}

	if govcd.ContainsNotFound(err) {
		var adminOrg *govcd.AdminOrg
		adminOrg, err = c.Vmware.GetAdminOrgByName(orgName)
		if err != nil {
			return nil, nil, fmt.Errorf("%w: %s %w", ErrRetrievingOrgAdmin, vdcName, err)
		}

		vdcOrVDCGroup, err = adminOrg.GetVdcGroupByName(vdcName)
		if err != nil {
			return nil, nil, fmt.Errorf("%w: %s %w", ErrRetrievingVDCGroup, vdcName, err)
		}
	}
	if vdcOrVDCGroup == nil {
		return nil, nil, fmt.Errorf("error retrieving VDC %s: not found", vdcName)
	}

	return org, vdcOrVDCGroup, err
}
