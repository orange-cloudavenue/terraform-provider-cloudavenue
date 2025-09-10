/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

// Package edgegw provides a Terraform resource to manage edge gateways.
package vdc

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/api/vdc/v1"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
	"github.com/orange-cloudavenue/common-go/utils"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	cavutils "github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
)

var (
	_ datasource.DataSource              = &storageProfilesDataSource{}
	_ datasource.DataSourceWithConfigure = &storageProfilesDataSource{}
)

// NewStorageProfilesDataSource returns a new resource implementing the storage_profiles data source.
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
	defer metrics.New("data.cloudavenue_storage_profiles", d.client.GetOrgName(), metrics.Read)()
	var (
		plan = new(storageProfilesDataSourceModel)
		data = new(storageProfilesDataSourceModel)
	)

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Call the API to get the list of storage profiles
	storageProfiles, err := d.eClient.ListStorageProfile(ctx, types.ParamsListStorageProfile{
		ID:   plan.VDCID.ValueString(),
		Name: plan.VDCName.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Unable to list vdc storage profiles", err.Error())
		return
	}

	// Set VDC ID and Name
	data.VDCID.Set(storageProfiles.VDCS[0].ID)
	data.VDCName.Set(storageProfiles.VDCS[0].Name)

	// Map each storage profile to the schema
	storageProfilesList := make([]storageProfilesDataSourceModelStorageProfile, 0, len(storageProfiles.VDCS[0].StorageProfiles))
	for _, sp := range storageProfiles.VDCS[0].StorageProfiles {
		storageProfilesList = append(storageProfilesList, storageProfilesDataSourceModelStorageProfile{
			ID:      supertypes.NewStringValue(sp.ID),
			Class:   supertypes.NewStringValue(sp.Class),
			Limit:   supertypes.NewInt64Value(int64(sp.Limit)),
			Used:    supertypes.NewInt64Value(int64(sp.Used)),
			Default: supertypes.NewBoolValue(sp.Default),
		})
	}

	data.StorageProfiles.Set(ctx, utils.ToPTRSlice(storageProfilesList))

	// Set the ID attribute to a static value as this data source does not have a unique identifier
	data.ID.Set(cavutils.GenerateUUID(data.VDCName.ValueString()).String())

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
