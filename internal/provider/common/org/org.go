package org

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
)

type Org struct {
	*client.Org
	c *client.CloudAvenue
}

// Init.
func Init(c *client.CloudAvenue) (Org, diag.Diagnostics) {
	var (
		d = diag.Diagnostics{}
		o = Org{
			c: c,
		}
		err error
	)

	o.Org, err = c.GetOrg()
	if err != nil {
		d.AddError("Unable to get ORG", err.Error())
		return Org{}, d
	}

	return o, nil
}
