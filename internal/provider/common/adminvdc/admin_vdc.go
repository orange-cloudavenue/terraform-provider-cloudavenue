package adminvdc

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
)

type AdminVDC struct {
	*client.AdminVDC
	adminOrg *client.AdminOrg
}

func Init(c *client.CloudAvenue, adminvdc types.String) (AdminVDC, diag.Diagnostics) {
	var (
		d    = diag.Diagnostics{}
		opts = make([]client.GetAdminVDCOpts, 0)
		v    = AdminVDC{}

		err error
	)

	v.adminOrg, err = c.GetAdminOrg()
	if err != nil {
		d.AddError("Unable to get AdminORG", err.Error())
		return AdminVDC{}, d
	}

	if !adminvdc.IsNull() && !adminvdc.IsUnknown() {
		opts = append(opts, client.WithAdminVDCName(adminvdc.ValueString()))
	}

	v.AdminVDC, err = c.GetAdminVDC(opts...)
	if err != nil {
		d.AddError("Unable to get AdminVDC", err.Error())
		return AdminVDC{}, d
	}

	return v, nil
}
