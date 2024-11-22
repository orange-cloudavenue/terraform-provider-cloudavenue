package client

import (
	"fmt"

	v1 "github.com/orange-cloudavenue/cloudavenue-sdk-go/v1"
)

// GetVDC
// return the vdc using the name provided in the argument.
// If the name is empty, it will try to use the default vdc provided in the provider.
func (c *CloudAvenue) GetVDC(vdcNamestring string) (vdc *v1.VDC, err error) {
	if vdcNamestring == "" {
		if c.DefaultVDCExist() {
			vdcNamestring = c.GetDefaultVDC()
		} else {
			return nil, fmt.Errorf("%w", ErrEmptyVDCNameProvided)
		}
	}

	vdc, err = c.CAVSDK.V1.VDC().GetVDC(vdcNamestring)
	if err != nil {
		return nil, fmt.Errorf("%w: %s %w", ErrRetrievingVDC, vdcNamestring, err)
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
