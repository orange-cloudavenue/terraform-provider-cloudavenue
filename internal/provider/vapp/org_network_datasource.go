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

	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/network"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/org"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vapp"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vdc"
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
	client *client.CloudAvenue

	org  org.Org
	vdc  vdc.VDC
	vapp vapp.VAPP
}

// Init Initializes the data source.
func (d *orgNetworkDataSource) Init(ctx context.Context, dm *orgNetworkModel) (diags diag.Diagnostics) {
	d.org, diags = org.Init(d.client)
	if diags.HasError() {
		return
	}

	d.vdc, diags = vdc.Init(d.client, dm.VDC)
	if diags.HasError() {
		return
	}

	d.vapp, diags = vapp.Init(d.client, d.vdc, dm.VAppID, dm.VAppName)

	return
}

func (d *orgNetworkDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_org_network"
}

func (d *orgNetworkDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = network.GetSchema(network.SetRoutedVapp()).GetDataSource(ctx)
}

func (d *orgNetworkDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *orgNetworkDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	defer metrics.New("data.cloudavenue_vapp_org_network", d.client.GetOrgName(), metrics.Read)()

	var data *orgNetworkModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Init resource
	resp.Diagnostics.Append(d.Init(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get vApp Network information
	vAppNetworkConfig, err := d.vapp.GetNetworkConfig()
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving vApp network config", err.Error())
		return
	}

	// Get Network Config
	vAppNetwork, networkID, errFindNetwork := data.findOrgNetwork(vAppNetworkConfig)
	resp.Diagnostics.Append(errFindNetwork...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Remove resource from state if not found
	if vAppNetwork == (&govcdtypes.VAppNetworkConfiguration{}) {
		resp.State.RemoveResource(ctx)
		return
	}

	// Set data
	plan := &orgNetworkModel{
		ID:          types.StringValue(urn.Normalize(urn.Network, *networkID).String()),
		VAppName:    utils.StringValueOrNull(d.vapp.GetName()),
		VAppID:      utils.StringValueOrNull(d.vapp.GetID()),
		VDC:         types.StringValue(d.vdc.GetName()),
		NetworkName: data.NetworkName,
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}
