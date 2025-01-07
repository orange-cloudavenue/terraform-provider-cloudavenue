/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

// Package vapp provides a Terraform datasource.
package vapp

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vapp"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vdc"
)

var (
	_ datasource.DataSource              = &isolatedNetworkDataSource{}
	_ datasource.DataSourceWithConfigure = &isolatedNetworkDataSource{}
)

func NewIsolatedNetworkDataSource() datasource.DataSource {
	return &isolatedNetworkDataSource{}
}

type isolatedNetworkDataSource struct {
	client *client.CloudAvenue
	vdc    vdc.VDC
	vapp   vapp.VAPP
}

// Init Initializes the data source.
func (d *isolatedNetworkDataSource) Init(ctx context.Context, dm *isolatedNetworkModel) (diags diag.Diagnostics) {
	d.vdc, diags = vdc.Init(d.client, dm.VDC.StringValue)
	if diags.HasError() {
		return
	}

	d.vapp, diags = vapp.Init(d.client, d.vdc, dm.VAppID.StringValue, dm.VAppName.StringValue)

	return
}

func (d *isolatedNetworkDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_" + "isolated_network"
}

func (d *isolatedNetworkDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = isolatedNetworkSchema().GetDataSource(ctx)
}

func (d *isolatedNetworkDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.CloudAvenue)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.CloudAvenue, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *isolatedNetworkDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	defer metrics.New("data.cloudavenue_vapp_isolated_network", d.client.GetOrgName(), metrics.Read)()
	config := &isolatedNetworkModel{}

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Init datasource
	resp.Diagnostics.Append(d.Init(ctx, config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	s := &isolatedNetworkResource{
		client: d.client,
		vdc:    d.vdc,
		vapp:   d.vapp,
	}

	// Read data
	data, _, diags := s.read(ctx, config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
