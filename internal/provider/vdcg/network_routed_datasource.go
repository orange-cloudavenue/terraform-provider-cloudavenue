/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

// Package vdcg provides a Terraform datasource.
package vdcg

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	v1 "github.com/orange-cloudavenue/cloudavenue-sdk-go/v1"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
)

var (
	_ datasource.DataSource              = &NetworkRoutedDataSource{}
	_ datasource.DataSourceWithConfigure = &NetworkRoutedDataSource{}
)

func NewNetworkRoutedDataSource() datasource.DataSource {
	return &NetworkRoutedDataSource{}
}

type NetworkRoutedDataSource struct {
	client *client.CloudAvenue
	vdcg   *v1.VDCGroup
}

// Init Initializes the data source.
func (d *NetworkRoutedDataSource) Init(ctx context.Context, dm *NetworkRoutedModel) (diags diag.Diagnostics) {
	var err error

	idOrName := dm.VDCGroupID.Get()
	if idOrName == "" {
		idOrName = dm.VDCGroupName.Get()
	}

	d.vdcg, err = d.client.CAVSDK.V1.VDC().GetVDCGroup(idOrName)
	if err != nil {
		diags.AddError(
			"Error retrieving VDC Group",
			fmt.Sprintf("Error retrieving VDC Group %q: %s", idOrName, err),
		)
		return
	}

	dm.VDCGroupID.Set(d.vdcg.GetID())
	dm.VDCGroupName.Set(d.vdcg.GetName())

	return
}

func (d *NetworkRoutedDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_network_routed"
}

func (d *NetworkRoutedDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = networkRoutedSchema(ctx).GetDataSource(ctx)
}

func (d *NetworkRoutedDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *NetworkRoutedDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	defer metrics.New("data.cloudavenue_vdcg_network_routed", d.client.GetOrgName(), metrics.Read)()

	config := &NetworkRoutedModel{}

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

	s := &NetworkRoutedResource{
		client: d.client,
		vdcg:   d.vdcg,
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
