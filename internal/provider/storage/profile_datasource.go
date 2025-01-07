/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

// Package storage provides a Terraform datasource.
package storage

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vdc"
)

var (
	_ datasource.DataSource              = &profileDataSource{}
	_ datasource.DataSourceWithConfigure = &profileDataSource{}
)

func NewProfileDataSource() datasource.DataSource {
	return &profileDataSource{}
}

type profileDataSource struct {
	client *client.CloudAvenue
	vdc    vdc.VDC
}

// Init Initializes the data source.
func (d *profileDataSource) Init(ctx context.Context, dm *profileDataSourceModel) (diags diag.Diagnostics) {
	d.vdc, diags = vdc.Init(d.client, dm.VDC)
	if diags.HasError() {
		return
	}

	return
}

func (d *profileDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_profile"
}

func (d *profileDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = d.superSchema(ctx).GetDataSource(ctx)
}

func (d *profileDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *profileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	defer metrics.New("data.cloudavenue_storage_profile", d.client.GetOrgName(), metrics.Read)()

	config := &profileDataSourceModel{}

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

	storageProfileID, err := d.vdc.FindStorageProfileName(config.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Storage Profile (ID) not found",
			fmt.Sprintf("Storage profile %s not found in VDC %s", config.Name.ValueString(), d.vdc.GetName()),
		)
		return
	}

	storageProfileRef, err := d.vdc.GetStorageProfileReference(storageProfileID, true)
	if err != nil {
		resp.Diagnostics.AddError(
			"Storage Profile (Reference) not found",
			fmt.Sprintf("Storage profile %s not found in VDC %s", config.Name.ValueString(), d.vdc.GetName()),
		)
		return
	}

	storageProfile, err := d.client.Vmware.GetStorageProfileByHref(storageProfileRef.HREF)
	if err != nil {
		resp.Diagnostics.AddError(
			"Storage Profile (Reference) not found",
			fmt.Sprintf("Storage profile %s not found in VDC %s", config.Name.ValueString(), d.vdc.GetName()),
		)
		return
	}

	config.ID = types.StringValue(storageProfileID)
	config.VDC = types.StringValue(d.vdc.GetName())
	config.Limit = types.Int64Value(storageProfile.Limit)
	config.UsedStorage = types.Int64Value(storageProfile.StorageUsedMB)
	config.Default = types.BoolValue(storageProfile.Default)
	config.Enabled = types.BoolValue(*storageProfile.Enabled)
	config.IopsAllocated = types.Int64Value(storageProfile.IopsAllocated)
	config.Units = types.StringValue(storageProfile.Units)
	config.IopsLimitingEnabled = types.BoolNull()
	config.MaximumDiskIops = types.Int64Null()
	config.DefaultDiskIops = types.Int64Null()
	config.DiskIopsPerGbMax = types.Int64Null()
	config.IopsLimit = types.Int64Null()
	if storageProfile.IopsSettings != nil {
		config.IopsLimitingEnabled = types.BoolValue(storageProfile.IopsSettings.Enabled)
		config.MaximumDiskIops = types.Int64Value(storageProfile.IopsSettings.DiskIopsMax)
		config.DefaultDiskIops = types.Int64Value(storageProfile.IopsSettings.DiskIopsDefault)
		config.DiskIopsPerGbMax = types.Int64Value(storageProfile.IopsSettings.DiskIopsPerGbMax)
		config.IopsLimit = types.Int64Value(storageProfile.IopsSettings.StorageProfileIopsLimit)
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, config)...)
}
