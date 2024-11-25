package adminorg

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"

	v1 "github.com/orange-cloudavenue/cloudavenue-sdk-go/v1"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
)

type AdminOrg struct {
	*v1.AdminOrg
	c *client.CloudAvenue
}

// Init.
func Init(c *client.CloudAvenue) (adminOrg AdminOrg, diags diag.Diagnostics) {
	o, err := c.CAVSDK.V1.AdminOrg()
	if err != nil {
		diags.AddError("Unable to get ORG", err.Error())
		return adminOrg, diags
	}

	return AdminOrg{
		AdminOrg: o,
		c:        c,
	}, nil
}
