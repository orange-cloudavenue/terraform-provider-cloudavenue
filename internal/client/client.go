// Package client is the main client for the CloudAvenue provider.
package client

import (
	"github.com/vmware/go-vcloud-director/v2/govcd"

	clientca "github.com/orange-cloudavenue/cloudavenue-sdk-go"
)

// CloudAvenue is the main struct for the CloudAvenue client.
type CloudAvenue struct {
	// API VMWARE
	Vmware *govcd.VCDClient // Deprecated

	// SDK CLOUDAVENUE
	CAVSDK     *clientca.Client
	CAVSDKOpts *clientca.ClientOpts
}

// New creates a new CloudAvenue client.
func (c *CloudAvenue) New() (*CloudAvenue, error) {
	var err error

	if c.CAVSDKOpts == nil {
		c.CAVSDKOpts = new(clientca.ClientOpts)
	}

	// SDK CloudAvenue
	c.CAVSDK, err = clientca.New(c.CAVSDKOpts)
	if err != nil {
		return nil, err
	}

	c.Vmware, err = c.CAVSDK.V1.Vmware()
	if err != nil {
		return nil, err
	}

	return c, nil
}

// DefaultVDCExist returns true if the default VDC exists.
func (c *CloudAvenue) DefaultVDCExist() bool {
	return c.CAVSDKOpts.CloudAvenue.VDC != ""
}

// GetDefaultVDC returns the default VDC.
func (c *CloudAvenue) GetDefaultVDC() string {
	return c.CAVSDKOpts.CloudAvenue.VDC
}

// GetURL returns the base path of the API.
func (c *CloudAvenue) GetURL() string {
	// Error is not returned for maintein compatibility with the previous version
	v, e := c.CAVSDK.Config().GetURL()
	if e != nil {
		return ""
	}

	return v
}

// GetOrgName() returns the name of the organization.
func (c *CloudAvenue) GetOrgName() string {
	// Error is not returned for maintein compatibility with the previous version
	v, e := c.CAVSDK.Config().GetOrganization()
	if e != nil {
		return ""
	}

	return v
}

// GetUserName() returns the name of the user.
func (c *CloudAvenue) GetUserName() string {
	// Error is not returned for maintein compatibility with the previous version
	v, e := c.CAVSDK.Config().GetUsername()
	if e != nil {
		return ""
	}

	return v
}
