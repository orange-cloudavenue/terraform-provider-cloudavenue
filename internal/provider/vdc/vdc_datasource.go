/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

// Package vdc provides a resource to manage VDCs.
package vdc

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
)

var (
	_ datasource.DataSource              = &vdcDataSource{}
	_ datasource.DataSourceWithConfigure = &vdcDataSource{}
)

// NewVDCDataSource returns a new resource implementing the vdcs data source.
func NewVDCDataSource() datasource.DataSource {
	return &vdcDataSource{}
}

type vdcDataSource struct {
	client *client.CloudAvenue
}

func (d *vdcDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName
}

func (d *vdcDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = vdcSchema().GetDataSource(ctx)
}

func (d *vdcDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *vdcDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	defer metrics.New("data.cloudavenue_vdc", d.client.GetOrgName(), metrics.Read)()

	data := new(vdcDataSourceModel)

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	s := &vdcResource{
		client: d.client,
	}

	dataResource := new(vdcResourceModel)
	dataResource.Name = data.Name

	// Read data from the API
	dataRefreshed, found, diags := s.read(ctx, dataResource)
	if !found {
		resp.Diagnostics.AddError("VDC not found", fmt.Sprintf("The VDC with the name %q was not found", data.Name.ValueString()))
		return
	}
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.BillingModel = dataRefreshed.BillingModel
	data.CPUAllocated = dataRefreshed.CPUAllocated
	data.Description = dataRefreshed.Description
	data.DisponibilityClass = dataRefreshed.DisponibilityClass
	data.ID = dataRefreshed.ID
	data.MemoryAllocated = dataRefreshed.MemoryAllocated
	data.Name = dataRefreshed.Name
	data.ServiceClass = dataRefreshed.ServiceClass
	data.StorageBillingModel = dataRefreshed.StorageBillingModel
	data.VCPUInMhz = dataRefreshed.VCPUInMhz
	data.StorageProfiles = dataRefreshed.StorageProfiles

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
