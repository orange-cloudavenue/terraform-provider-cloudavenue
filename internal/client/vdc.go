/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package client

import (
	"fmt"

	v1 "github.com/orange-cloudavenue/cloudavenue-sdk-go/v1"
)

// GetVDC
// return the vdc using the name provided in the argument.
// If the name is empty, it will try to use the default vdc provided in the provider.
func (c *CloudAvenue) GetVDC(vdcName string) (vdc *v1.VDC, err error) {
	if vdcName == "" {
		if c.DefaultVDCExist() {
			vdcName = c.GetDefaultVDC()
		} else {
			return nil, fmt.Errorf("%w", ErrEmptyVDCNameProvided)
		}
	}

	vdc, err = c.CAVSDK.V1.VDC().GetVDC(vdcName)
	if err != nil {
		return nil, fmt.Errorf("%w: %s %w", ErrRetrievingVDC, vdcName, err)
	}

	return vdc, nil
}

// GetVDCGroup return the vdc group using the name provided in the argument.
func (c *CloudAvenue) GetVDCGroup(vdcGroupName string) (vdcGroup *v1.VDCGroup, err error) {
	if vdcGroupName == "" {
		return nil, fmt.Errorf("%w", ErrEmptyVDCNameProvided)
	}

	vdcGroup, err = c.CAVSDK.V1.VDC().GetVDCGroup(vdcGroupName)
	if err != nil {
		return nil, fmt.Errorf("%w: %s %w", ErrRetrievingVDCGroup, vdcGroupName, err)
	}

	return vdcGroup, nil
}

// GetVDCOrVDCGroup return the vdc or vdc group using the name provided in the argument.
func (c *CloudAvenue) GetVDCOrVDCGroup(vdcOrVDCGroupName string) (v1.VDCOrVDCGroupInterface, error) {
	return c.CAVSDK.V1.VDC().GetVDCOrVDCGroup(vdcOrVDCGroupName)
}
