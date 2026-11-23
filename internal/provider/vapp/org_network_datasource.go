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

// Package vapp provides a Terraform datasource.
package vapp

import (
	"context"
	"fmt"

	"github.com/orange-cloudavenue/common-go/urn"

	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	sdktypes "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	cerrs "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/errors"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/network"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

var (
	_ datasource.DataSource              = &orgNetworkDataSource{}
	_ datasource.DataSourceWithConfigure = &orgNetworkDataSource{}
)

func NewOrgNetworkDataSource() datasource.DataSource {
	return &orgNetworkDataSource{}
}

type orgNetworkDataSource struct {
	client       *client.CloudAvenue
	vappClientV2 vAppClientV2
}

func (d *orgNetworkDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_org_network"
}

func (d *orgNetworkDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = network.GetSchema(network.SetRoutedVapp()).GetDataSource(ctx)
}

func (d *orgNetworkDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *orgNetworkDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	defer metrics.New("data.cloudavenue_vapp_org_network", d.client.GetOrgName(), metrics.Read)()

	var data *orgNetworkModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Resolve VDC and parent vApp
	vdc, diags := d.vappClientV2.resolveVDC(ctx, data.VDC.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	app, diags := d.vappClientV2.resolveVApp(ctx, vdc, data.VAppID.ValueString(), data.VAppName.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get vApp Network information
	section, err := d.vappClientV2.vappClient.GetVAppNetworkConfig(ctx, sdktypes.ParamsGetVAppNetworkConfig{ID: app.ID})
	if err != nil {
		cerrs.AddError(&resp.Diagnostics, cerrs.ActionRead, fmt.Sprintf("vApp %s(%s) network config", app.Name, app.ID), err)
		return
	}

	// Get Network Config (returns an error diagnostic if the network is not found)
	vAppNetwork, diags := data.findOrgNetwork(section.NetworkConfig)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set data
	state := &orgNetworkModel{
		ID:          types.StringValue(urn.Normalize(urn.Network, urn.ExtractUUID(vAppNetwork.ID)).String()),
		VAppName:    utils.StringValueOrNull(app.Name),
		VAppID:      utils.StringValueOrNull(app.ID),
		VDC:         types.StringValue(vdc.Name),
		NetworkName: data.NetworkName,
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
