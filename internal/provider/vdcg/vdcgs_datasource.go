/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package vdcg

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/api/vdcgroup/v1"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
)

var (
	_ datasource.DataSource              = &vdcgsDataSource{}
	_ datasource.DataSourceWithConfigure = &vdcgsDataSource{}
)

func NewVDCGsDataSource() datasource.DataSource {
	return &vdcgsDataSource{}
}

type vdcgsDataSource struct {
	// Client is a terraform Client
	client *client.CloudAvenue

	// vgClient is the VDC Group client from the SDK V2
	vgClient *vdcgroup.Client
}

func (d *vdcgsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "s"
}

func (d *vdcgsDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = vdcgsSchema(ctx).GetDataSource(ctx)
}

func (d *vdcgsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	vgC, err := vdcgroup.New(client.V2)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create VDC Group client, got error: %s", err),
		)
		return
	}

	d.client = client
	d.vgClient = vgC
}

func (d *vdcgsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	defer metrics.New("data.cloudavenue_vdcgs", d.client.GetOrgName(), metrics.Read)()

	config := &vdcgsModel{}

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	vdcGroups, err := d.vgClient.ListVdcGroup(ctx, types.ParamsListVdcGroup{
		ID:   config.FilterID.Get(),
		Name: config.FilterID.Get(),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to read VDC Groups, got error: %s", err),
		)
		return
	}

	data := &vdcgsModel{}
	resp.Diagnostics.Append(data.fromSDK(ctx, vdcGroups)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.FilterID = config.FilterID
	data.FilterName = config.FilterName

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}
