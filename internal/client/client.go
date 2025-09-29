/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

// Package client is the main client for the CloudAvenue provider.
package client

import (
	"context"
	"log/slog"

	"github.com/vmware/go-vcloud-director/v2/govcd"

	clientca "github.com/orange-cloudavenue/cloudavenue-sdk-go"
	sdkv2 "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
)

// CloudAvenue is the main struct for the CloudAvenue client.
type CloudAvenue struct {
	// API VMWARE
	Vmware *govcd.VCDClient // Deprecated

	// SDK CLOUDAVENUE
	CAVSDK     *clientca.Client     // Deprecated
	CAVSDKOpts *clientca.ClientOpts // Deprecated

	// SDK V2
	V2 sdkv2.Client
}

// New creates a new CloudAvenue client.
func (c *CloudAvenue) New(ctx context.Context) (*CloudAvenue, error) {
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

	c.V2, err = sdkv2.NewClient(
		c.CAVSDKOpts.CloudAvenue.Org,
		sdkv2.WithCloudAvenueCredential(c.CAVSDKOpts.CloudAvenue.Username, c.CAVSDKOpts.CloudAvenue.Password),
		sdkv2.WithLogger(slog.New(newTflogHandler())),
	)
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
