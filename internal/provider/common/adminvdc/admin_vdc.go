/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package adminvdc

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	v1 "github.com/orange-cloudavenue/cloudavenue-sdk-go/v1"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
)

type AdminVDC struct {
	*v1.AdminVDC
	adminOrg *v1.AdminOrg
}

func Init(c *client.CloudAvenue, adminvdc types.String) (avdc AdminVDC, diags diag.Diagnostics) {
	var err error

	avdc.adminOrg, err = c.CAVSDK.V1.AdminOrg()
	if err != nil {
		diags.AddError("Unable to get AdminOrg", err.Error())
		return AdminVDC{}, diags
	}

	vdcName := adminvdc.ValueString()
	if vdcName == "" {
		if c.DefaultVDCExist() {
			vdcName = c.GetDefaultVDC()
		} else {
			diags.AddError("Empty VDC name provided", client.ErrEmptyVDCNameProvided.Error())
			return avdc, diags
		}
	}

	avdc.AdminVDC, err = c.CAVSDK.V1.AdminVDC().Get(vdcName)
	if err != nil {
		diags.AddError("Unable to get AdminVDC", err.Error())
		return AdminVDC{}, diags
	}

	return avdc, nil
}
