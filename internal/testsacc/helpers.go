/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package testsacc

import (
	"fmt"
	"os"

	"github.com/orange-cloudavenue/common-go/validators"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/api/edgegateway/v1"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/api/organization/v1"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/api/vdc/v1"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/api/vdcgroup/v1"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
)

// ToValidate returns a CheckResourceAttrWithFunc that validates a string using the given validator name.
// If the value is empty, an error is returned indicating the value is not valid for the validator.
// Otherwise, the value is validated using the provided validator name via validators.New().Var().
func ToValidate(validatorToCheck string) resource.CheckResourceAttrWithFunc {
	return func(value string) error {
		if value == "" {
			return fmt.Errorf("empty value, is not a valid %s", validatorToCheck)
		}

		return validators.New().Var(value, validatorToCheck)
	}
}

// * This part is used to map resource names to their corresponding client functions for testing purposes.
type resourceClient struct {
	ClientCav cav.Client
}

func initClientCav() (rc resourceClient, err error) {
	url := os.Getenv("CLOUDAVENUE_URL")
	username := os.Getenv("CLOUDAVENUE_USERNAME")
	password := os.Getenv("CLOUDAVENUE_PASSWORD")
	orgName := os.Getenv("CLOUDAVENUE_ORG")
	if url == "" || username == "" || password == "" || orgName == "" {
		panic("CLOUDAVENUE_URL, CLOUDAVENUE_USERNAME, CLOUDAVENUE_PASSWORD or CLOUDAVENUE_ORG not set in environment")
	}

	// Create the CloudAvenue client
	rc.ClientCav, err = cav.NewClient(orgName, cav.WithCloudAvenueCredential(username, password))
	if err != nil {
		return rc, fmt.Errorf("failed to create CloudAvenue client: %w", err)
	}

	// Initialize the CloudAvenue client
	return
}

// return an organization.Client to interact with organization resources.
func (rc *resourceClient) GetOrganizationClient() (*organization.Client, error) {
	if rc.ClientCav == nil {
		return nil, fmt.Errorf("CloudAvenue client is not initialized")
	}
	return organization.New(rc.ClientCav)
}

// return a vdc.Client to interact with vdc resources.
func (rc *resourceClient) GetVDCClient() (*vdc.Client, error) {
	if rc.ClientCav == nil {
		return nil, fmt.Errorf("CloudAvenue client is not initialized")
	}
	return vdc.New(rc.ClientCav)
}

// return a vdcgroup.Client to interact with vdcgroup resources.
func (rc *resourceClient) GetVDCGroupClient() (*vdcgroup.Client, error) {
	if rc.ClientCav == nil {
		return nil, fmt.Errorf("CloudAvenue client is not initialized")
	}
	return vdcgroup.New(rc.ClientCav)
}

// return an edgegateway.Client to interact with edge resources.
func (rc *resourceClient) GetEdgeClient() (*edgegateway.Client, error) {
	if rc.ClientCav == nil {
		return nil, fmt.Errorf("CloudAvenue client is not initialized")
	}
	return edgegateway.New(rc.ClientCav)
}
