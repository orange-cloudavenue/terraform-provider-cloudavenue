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

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/adminorg"
)

var (
	_ datasource.DataSource              = &vdcgDataSource{}
	_ datasource.DataSourceWithConfigure = &vdcgDataSource{}
)

func NewVDCGDataSource() datasource.DataSource {
	return &vdcgDataSource{}
}

type vdcgDataSource struct {
	client   *client.CloudAvenue
	adminOrg adminorg.AdminOrg
}

func (d *vdcgDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName
}

// Init Initializes the resource.
func (d *vdcgDataSource) Init(_ context.Context, _ *vdcgModel) (diags diag.Diagnostics) {
	d.adminOrg, diags = adminorg.Init(d.client)
	return
}

func (d *vdcgDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = vdcgSchema(ctx).GetDataSource(ctx)
}

func (d *vdcgDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *vdcgDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	defer metrics.New("data.cloudavenue_vdcg", d.client.GetOrgName(), metrics.Read)()

	config := &vdcgModel{}

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Init the resource
	resp.Diagnostics.Append(d.Init(ctx, config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	/*
		Implement the data source read logic here.
	*/

	s := &vdcgResource{
		client:   d.client,
		adminOrg: d.adminOrg,
	}

	// Read data from the API
	data, found, diags := s.read(ctx, config)
	if !found {
		resp.Diagnostics.AddError("Resource not found", fmt.Sprintf("The VDC Group %s(%s) was not found", config.Name.Get(), config.ID.Get()))
		return
	}
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
