/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

// Package org provides a Terraform datasource.
package org

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/v1/org"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
)

var (
	_ datasource.DataSource              = &OrgDataSource{}
	_ datasource.DataSourceWithConfigure = &OrgDataSource{}
)

func NewOrgDataSource() datasource.DataSource {
	return &OrgDataSource{}
}

type OrgDataSource struct { //nolint:revive
	client *client.CloudAvenue
	org    org.Client
}

// Init Initializes the data source.
func (d *OrgDataSource) Init(_ context.Context, _ *OrgModel) (diags diag.Diagnostics) {
	var err error

	d.org, err = org.NewClient()
	if err != nil {
		diags.AddError("Error creating org client", err.Error())
		return diags
	}
	return diags
}

func (d *OrgDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName
}

func (d *OrgDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = orgSchema(ctx).GetDataSource(ctx)
}

func (d *OrgDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *OrgDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	defer metrics.New("data.cloudavenue_org", d.client.GetOrgName(), metrics.Read)()

	config := &OrgModel{}

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

	// If read function is identical to the resource, you can use the following code:

	s := &OrgResource{
		client: d.client,
		org:    d.org,
	}

	// Read data from the API
	data, found, diags := s.read(ctx, config)
	if !found {
		resp.Diagnostics.AddError("Resource not found", "The resource was not found")
		return
	}
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
