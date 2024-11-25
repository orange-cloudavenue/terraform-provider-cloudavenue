package org

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"

	v1 "github.com/orange-cloudavenue/cloudavenue-sdk-go/v1"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
)

type Org struct {
	*v1.Org
	c *client.CloudAvenue
}

// Init.
func Init(c *client.CloudAvenue) (org Org, diags diag.Diagnostics) {
	o, err := c.CAVSDK.V1.Org()
	if err != nil {
		diags.AddError("Unable to get ORG", err.Error())
		return org, diags
	}

	return Org{
		Org: o,
		c:   c,
	}, nil
}
