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
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

var (
	_ datasource.DataSource              = &storageProfilesDataSource{}
	_ datasource.DataSourceWithConfigure = &storageProfilesDataSource{}
)

// NewstorageProfilesDataSource returns a new resource implementing the storage_profiles data source.
func NewStorageProfilesDataSource() datasource.DataSource {
	return &storageProfilesDataSource{}
}

type storageProfilesDataSource struct {
	// Client is a terraform Client
	client *client.CloudAvenue

	// eClient is the Storage Profile client from the SDK V2
	eClient *vdc.Client
}

func (d *storageProfilesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_storage_profiles"
}

func (d *storageProfilesDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = storageProfilesSuperSchema(ctx).GetDataSource(ctx)
}

func (d *storageProfilesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *storageProfilesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	defer metrics.New("data.cloudavenue_vdc_storage_profiles", d.client.GetOrgName(), metrics.Read)()
	var (
		plan = new(storageProfilesDataSourceModel)
		data = new(storageProfilesDataSourceModel)
	)

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Call the API to get the list of storage profile
	storageProfiles, err := d.eClient.ListStorageProfile(ctx, types.ParamsListStorageProfile{
		VdcID:   plan.VDCID.Get(),
		VdcName: plan.VDCName.Get(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Unable to get list of vdc storage profiles", err.Error())
		return
	}

	// Check if we have a VDC
	if len(storageProfiles.VDCS) == 0 {
		resp.Diagnostics.AddError("No VDC found", "No VDC found with the provided information")
		return
	}

	// Set VDC ID and Name
	data.VDCID.Set(storageProfiles.VDCS[0].ID)
	data.VDCName.Set(storageProfiles.VDCS[0].Name)

	// Get storage profiles
	sps := make([]*storageProfileDataSourceModelStorageProfile, len(storageProfiles.VDCS[0].StorageProfiles))
	for i, sp := range storageProfiles.VDCS[0].StorageProfiles {
		sps[i] = &storageProfileDataSourceModelStorageProfile{}
		sps[i].ID.Set(sp.ID)
		sps[i].Class.Set(sp.Class)
		sps[i].Limit.SetInt(sp.Limit)
		sps[i].Used.SetInt(sp.Used)
		sps[i].Default.Set(sp.Default)
	}

	// Set storage profiles
	resp.Diagnostics.Append(data.StorageProfiles.Set(ctx, sps)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create a slice of storage profile Class names to generate a unique ID
	// This will help to know if the list of storage profile change
	x := []string{}
	for _, sp := range storageProfiles.VDCS[0].StorageProfiles {
		x = append(x, sp.Class)
	}

	// Set ID for the data source
	// We use a URN format to avoid conflicts with other resources
	// urn:vcloud:vdcstorageProfiles:<uuid>
	// The UUID is generated from the list of storage profile Class names
	data.ID.Set(fmt.Sprintf("urn:vcloud:vdcstorageProfiles=%s", utils.GenerateUUID(x).ValueString()))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}
