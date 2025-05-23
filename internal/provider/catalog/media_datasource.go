/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

// Package catalog provides a Terraform datasource.
package catalog

import (
	"context"
	"fmt"

	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/adminorg"
)

var (
	_ datasource.DataSource              = &catalogMediaDataSource{}
	_ datasource.DataSourceWithConfigure = &catalogMediaDataSource{}
	_ catalog                            = &catalogMediaDataSource{}
)

func NewCatalogMediaDataSource() datasource.DataSource {
	return &catalogMediaDataSource{}
}

type catalogMediaDataSource struct {
	client   *client.CloudAvenue
	adminOrg adminorg.AdminOrg
	catalog  base
}

func (d *catalogMediaDataSource) Init(_ context.Context, rm *catalogMediaDataSourceModel) (diags diag.Diagnostics) {
	d.catalog = base{
		name: rm.CatalogName.ValueString(),
		id:   rm.CatalogID.ValueString(),
	}

	d.adminOrg, diags = adminorg.Init(d.client)

	return
}

func (d *catalogMediaDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_" + "media"
}

func (d *catalogMediaDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = mediaSchema().GetDataSource(ctx)
}

func (d *catalogMediaDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *catalogMediaDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	defer metrics.New("data.cloudavenue_catalog_media", d.client.GetOrgName(), metrics.Read)()

	state := &catalogMediaDataSourceModel{}

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(d.Init(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get catalog
	catalog, err := d.GetCatalog()
	if err != nil {
		resp.Diagnostics.AddError("Unable to find catalog", err.Error())
		return
	}

	// Get media by name (This is use only for getting the media Description)
	media, err := catalog.GetMediaByNameOrId(state.Name.ValueString(), false)
	if err != nil {
		resp.Diagnostics.AddError("Unable to find media", err.Error())
		return
	}

	// Query media
	mediaRecord, err := catalog.QueryMedia(state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to query media", err.Error())
		return
	}

	// Check if mediaRecord.MediaRecord is nil
	if mediaRecord.MediaRecord == nil {
		resp.Diagnostics.AddError("Unable to find media record", "mediaRecord.MediaRecord is nil")
		return
	}

	updatedState := &catalogMediaDataSourceModel{
		CatalogName: state.CatalogName,
		CatalogID:   state.CatalogID,

		ID:             types.StringValue(media.Media.ID),
		Name:           types.StringValue(media.Media.Name),
		Description:    types.StringValue(media.Media.Description),
		IsISO:          types.BoolValue(mediaRecord.MediaRecord.IsIso),
		OwnerName:      types.StringValue(mediaRecord.MediaRecord.OwnerName),
		IsPublished:    types.BoolValue(mediaRecord.MediaRecord.IsPublished),
		CreatedAt:      types.StringValue(mediaRecord.MediaRecord.CreationDate),
		Size:           types.Int64Value(mediaRecord.MediaRecord.StorageB),
		Status:         types.StringValue(mediaRecord.MediaRecord.Status),
		StorageProfile: types.StringValue(mediaRecord.MediaRecord.StorageProfileName),
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, updatedState)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *catalogMediaDataSource) GetID() string {
	return d.catalog.id
}

// GetName returns the name of the catalog.
func (d *catalogMediaDataSource) GetName() string {
	return d.catalog.name
}

// GetIDOrName returns the ID if it is set, otherwise it returns the name.
func (d *catalogMediaDataSource) GetIDOrName() string {
	if d.GetID() != "" {
		return d.GetID()
	}
	return d.GetName()
}

// GetCatalog returns the govcd.Catalog.
func (d *catalogMediaDataSource) GetCatalog() (*govcd.AdminCatalog, error) {
	return d.adminOrg.GetAdminCatalogByNameOrId(d.GetIDOrName(), true)
}
