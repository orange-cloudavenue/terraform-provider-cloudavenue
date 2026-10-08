// SPDX-FileCopyrightText: Copyright (c) 2026 Orange
// SPDX-License-Identifier: MPL-2.0
//
// This software is distributed under the MPL-2.0 license.
// the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
// or see the "LICENSE" file for more details.

// SPDX-FileCopyrightText: Copyright (c) 2023-2026 Orange
// SPDX-License-Identifier: MPL-2.0
//
// This software is distributed under the MPL-2.0 license.
// the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
// or see the "LICENSE" file for more details.

// Package client is the main client for the CloudAvenue provider.
package client

import (
	"fmt"
	"log/slog"

	"github.com/vmware/go-vcloud-director/v2/govcd"

	clientca "github.com/orange-cloudavenue/cloudavenue-sdk-go"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	v2consoles "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/consoles"
)

// CloudAvenue is the main struct for the CloudAvenue client.
type CloudAvenue struct {
	// API VMWARE
	//
	// Deprecated: use CAVSDK and CloudAvenue SDK accessors instead.
	// Vmware is kept only for compatibility with older callers.
	Vmware *govcd.VCDClient

	// SDK CLOUDAVENUE
	CAVSDK     *clientca.Client
	CAVSDKOpts *clientca.ClientOpts

	CAVSDKV2         cav.Client
	defaultVDC       string
	organizationName string
	username         string
	url              string
}

// New creates a new CloudAvenue client.
func (c *CloudAvenue) New() (*CloudAvenue, error) {
	var err error

	if c.CAVSDKOpts == nil {
		c.CAVSDKOpts = new(clientca.ClientOpts)
	}

	// Keep VMware compatibility through v1 SDK.
	c.CAVSDK, err = clientca.New(c.CAVSDKOpts)
	if err != nil {
		return nil, err
	}

	c.Vmware, err = c.CAVSDK.V1.Vmware()
	if err != nil {
		return nil, err
	}

	if err := c.newV2(); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *CloudAvenue) newV2() error {
	if c.CAVSDKOpts == nil || c.CAVSDKOpts.CloudAvenue == nil {
		return fmt.Errorf("initialize v2 client: missing cloudavenue options")
	}

	console, ok := v2consoles.FindByOrganizationName(c.CAVSDKOpts.CloudAvenue.Org)
	if !ok {
		return fmt.Errorf("initialize v2 client: console not found")
	}

	services := console.Services()
	if endpoint := c.CAVSDKOpts.CloudAvenue.URL; endpoint != "" {
		services.APIVCD.Endpoint = endpoint
	}
	if endpoint := c.CAVSDKOpts.CloudAvenue.CoreAPI; endpoint != "" {
		services.APICerberus.Endpoint = endpoint
	}

	client, err := cav.NewClient(
		c.CAVSDKOpts.CloudAvenue.Org,
		cav.WithCustomEndpoints(services),
		cav.WithCloudAvenueCredential(c.CAVSDKOpts.CloudAvenue.Username, c.CAVSDKOpts.CloudAvenue.Password),
		cav.WithLogger(slog.Default()),
	)
	if err != nil {
		return fmt.Errorf("initialize v2 client: %w", err)
	}

	c.CAVSDKV2 = client
	c.defaultVDC = c.CAVSDKOpts.CloudAvenue.VDC
	c.organizationName = c.CAVSDKOpts.CloudAvenue.Org
	c.username = c.CAVSDKOpts.CloudAvenue.Username
	c.url = services.APIVCD.Endpoint

	return nil
}

// DefaultVDCExist returns true if the default VDC exists.
func (c *CloudAvenue) DefaultVDCExist() bool {
	return c.defaultVDC != ""
}

// GetDefaultVDC returns the default VDC.
func (c *CloudAvenue) GetDefaultVDC() string {
	return c.defaultVDC
}

// GetURL returns the base path of the API.
func (c *CloudAvenue) GetURL() string {
	return c.url
}

// GetOrgName returns the name of the organization.
func (c *CloudAvenue) GetOrgName() string {
	return c.organizationName
}

// GetUserName returns the name of the user.
func (c *CloudAvenue) GetUserName() string {
	return c.username
}
