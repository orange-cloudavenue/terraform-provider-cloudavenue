/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package vdc

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/api/vdc/v1"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
)

var (
	_ datasource.DataSource              = &storageProfileDataSource{}
	_ datasource.DataSourceWithConfigure = &storageProfileDataSource{}
)

// NewStorageProfileDataSource returns a new resource implementing the storage_profile data source.
func NewStorageProfileDataSource() datasource.DataSource {
	return &storageProfileDataSource{}
}

type storageProfileDataSource struct {
	// Client is a terraform Client
	client *client.CloudAvenue

	// eClient is the Storage Profile client from the SDK V2
	eClient *vdc.Client
}

func (d *storageProfileDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_storage_profile"
}

func (d *storageProfileDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = storageProfileSchema(ctx).GetDataSource(ctx)
}

func (d *storageProfileDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	vC, err := vdc.New(client.V2)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create VDC client, got error: %s", err),
		)
		return
	}

	d.client = client
	d.eClient = vC
}

func (d *storageProfileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	defer metrics.New("data.cloudavenue_vdc_storage_profile", d.client.GetOrgName(), metrics.Read)()
	var (
		plan = new(storageProfileDataSourceModel)
		data = new(storageProfileDataSourceModel)
	)

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Call the API to get the list of storage profile
	storageProfiles, err := d.eClient.ListStorageProfile(ctx, types.ParamsListStorageProfile{
		ID:      plan.ID.Get(),
		Class:   plan.Class.Get(),
		VdcID:   plan.VDCID.Get(),
		VdcName: plan.VDCName.Get(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Unable to get vdc storage profile", err.Error())
		return
	}

	// Check if we have at least one storage profile
	if len(storageProfiles.VDCS) == 0 || len(storageProfiles.VDCS[0].StorageProfiles) == 0 {
		resp.Diagnostics.AddError("Unable to find vdc storage profile", "No storage profile found with the given information")
		return
	}

	// Set VDC ID and Name
	data.VDCID.Set(storageProfiles.VDCS[0].ID)
	data.VDCName.Set(storageProfiles.VDCS[0].Name)

	data.ID.Set(storageProfiles.VDCS[0].StorageProfiles[0].ID)
	data.Class.Set(storageProfiles.VDCS[0].StorageProfiles[0].Class)
	data.Limit.SetInt(storageProfiles.VDCS[0].StorageProfiles[0].Limit)
	data.Used.SetInt(storageProfiles.VDCS[0].StorageProfiles[0].Used)
	data.Default.Set(storageProfiles.VDCS[0].StorageProfiles[0].Default)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}
