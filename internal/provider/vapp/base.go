// SPDX-FileCopyrightText: Copyright (c) 2026 Orange
// SPDX-License-Identifier: MPL-2.0
//
// This software is distributed under the MPL-2.0 license.
// the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
// or see the "LICENSE" file for more details.

// SPDX-FileCopyrightText: Copyright (c) 2023-2026 Orange
// SPDX-License-Identifier: MPL-2.0
//
// This software is distributed under the MPL-2.0 license.
// the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
// or see the "LICENSE" file for more details.

package vapp

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"

	vapp "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/api/vapp/v1"
	vdc "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/api/vdc/v1"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
)

const (
	categoryName = "vapp"
)

// vAppClientV2 holds the SDK v2 vApp client for vapp operations.
type vAppClientV2 struct {
	client     *client.CloudAvenue
	vappClient *vapp.Client
	vdcClient  *vdc.Client
}

// InitVAppClientV2 initializes the v2 VAPP client for vapp operations.
func InitVAppClientV2(c *client.CloudAvenue, opts ...vapp.Option) (vAppClientV2, diag.Diagnostics) {
	var diags diag.Diagnostics
	v2, err := c.V2()
	if err != nil {
		diags.AddError("Error initializing vApp client", err.Error())
		return vAppClientV2{}, diags
	}
	vappClient, err := vapp.New(v2, opts...)
	if err != nil {
		diags.AddError("Error initializing vApp client", err.Error())
		return vAppClientV2{}, diags
	}
	vdcClient, err := vdc.New(v2)
	if err != nil {
		diags.AddError("Error initializing vApp client", err.Error())
		return vAppClientV2{}, diags
	}
	return vAppClientV2{
		client:     c,
		vappClient: vappClient,
		vdcClient:  vdcClient,
	}, diags
}
