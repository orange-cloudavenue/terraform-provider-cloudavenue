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

// Package vapp provides a Terraform resource to manage vApps.
package vapp

import (
	"cmp"
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	cerrs "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/errors"
)

var (
	_ datasource.DataSource              = &vappDataSource{}
	_ datasource.DataSourceWithConfigure = &vappDataSource{}
)

// const vAppUnknownStatus = "-unknown-status-"

// NewVappDataSource returns a new vApp data source.
func NewVappDataSource() datasource.DataSource {
	return &vappDataSource{}
}

type vappDataSource struct {
	client       *client.CloudAvenue
	vappClientV2 vAppClientV2
}

func (d *vappDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName
}

func (d *vappDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = vappSchema().GetDataSource(ctx)
}

func (d *vappDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.CloudAvenue)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected data source configure type",
			fmt.Sprintf("Expected *client.CloudAvenue, got %T. Report this to provider maintainers.", req.ProviderData),
		)
		return
	}

	d.client = client

	v2, diags := InitVAppClientV2(client)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	d.vappClientV2 = v2
}

func (d *vappDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	defer metrics.New("data.cloudavenue_vapp", d.client.GetOrgName(), metrics.Read)()

	config := &vappResourceModel{}

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resolved, diags := d.vappClientV2.resolveVDC(ctx, config.VDC.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	p := types.ParamsGetVApp{ID: config.VAppID.ValueString(), Name: config.VAppName.ValueString()}
	app, err := d.vappClientV2.vappClient.GetVApp(ctx, p)
	if err != nil {
		cerrs.AddError(&resp.Diagnostics, cerrs.ActionRead, fmt.Sprintf("vApp %s", cmp.Or(p.Name, p.ID)), err)
		return
	}

	if !vappInVDC(app, resolved) {
		resp.Diagnostics.AddError("Vdc not found", "App vdc id different from resolved vdc id")
		return
	}

	gp, err := d.vappClientV2.vappClient.GetVAppGuestProperties(ctx, types.ParamsGetVAppGuestProperties{ID: app.ID})
	if err != nil {
		cerrs.AddError(&resp.Diagnostics, cerrs.ActionRead, fmt.Sprintf("vApp %s(%s) guest properties", app.Name, app.ID), err)
		return
	}

	diags = config.fromSDK(ctx, app, resolved.Name)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = config.guestPropertiesFromSDK(ctx, gp)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, config)...)
}
