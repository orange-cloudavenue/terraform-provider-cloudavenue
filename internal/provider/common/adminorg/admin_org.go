package adminorg

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
)

type AdminOrg struct {
	*client.AdminOrg
}

// Init.
func Init(c *client.CloudAvenue) (AdminOrg, diag.Diagnostics) {
	var (
		d   = diag.Diagnostics{}
		o   = AdminOrg{}
		err error
	)

	o.AdminOrg, err = c.GetAdminOrg()
	if err != nil {
		d.AddError("Unable to get ORG", err.Error())
		return o, d
	}

	return o, nil
}
