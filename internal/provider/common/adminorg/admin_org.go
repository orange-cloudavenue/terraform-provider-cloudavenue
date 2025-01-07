/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

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
