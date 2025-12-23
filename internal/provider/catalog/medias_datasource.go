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
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

var (
	_ datasource.DataSource              = &catalogMediasDataSource{}
	_ datasource.DataSourceWithConfigure = &catalogMediasDataSource{}
	_ catalog                            = &catalogMediasDataSource{}
)

func NewCatalogMediasDataSource() datasource.DataSource {
	return &catalogMediasDataSource{}
}

type catalogMediasDataSource struct {
	client   *client.CloudAvenue
	adminOrg adminorg.AdminOrg
	catalog  base
}

func (d *catalogMediasDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_medias"
}

func (d *catalogMediasDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = mediasSchema().GetDataSource(ctx)
}

func (d *catalogMediasDataSource) Init(_ context.Context, rm *catalogMediasDataSourceModel) (diags diag.Diagnostics) {
	d.catalog = base{
		name: rm.CatalogName.ValueString(),
		id:   rm.CatalogID.ValueString(),
	}

	d.adminOrg, diags = adminorg.Init(d.client)
	return diags
}

func (d *catalogMediasDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *catalogMediasDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	defer metrics.New("data.cloudavenue_catalog_medias", d.client.GetOrgName(), metrics.Read)()

	config := &catalogMediasDataSourceModel{}

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(d.Init(ctx, config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	catalog, err := d.GetCatalog()
	if err != nil {
		resp.Diagnostics.AddError("Unable to find catalog", err.Error())
		return
	}

	var (
		medias     = make(catalogMediasDataSourceModelMedias)
		mediasName = make(catalogMediasDataSourceModelMediasName, 0)
	)

	// Get all medias
	mediaList, err := catalog.QueryMediaList()
	if err != nil {
		resp.Diagnostics.AddError("Unable to query media list", err.Error())
		return
	}

	for _, media := range mediaList {
		mediasName = append(mediasName, media.Name)
		medias[media.Name] = catalogMediaDataSourceModel{
			ID:             types.StringValue(media.ID),
			Name:           types.StringValue(media.Name),
			CatalogID:      types.StringValue(d.GetID()),
			CatalogName:    types.StringValue(d.GetName()),
			IsISO:          types.BoolValue(media.IsIso),
			OwnerName:      types.StringValue(media.OwnerName),
			IsPublished:    types.BoolValue(media.IsPublished),
			CreatedAt:      types.StringValue(media.CreationDate),
			Status:         types.StringValue(media.Status),
			Size:           types.Int64Value(media.StorageB),
			StorageProfile: types.StringValue(media.StorageProfileName),
			Description:    types.StringNull(),
		}
	}

	listMediasName, diag := types.ListValueFrom(ctx, types.StringType, mediasName)
	resp.Diagnostics.Append(diag...)
	if resp.Diagnostics.HasError() {
		return
	}

	listMedias, diag := types.MapValueFrom(ctx, config.Medias.ElementType(ctx), medias)
	resp.Diagnostics.Append(diag...)
	if resp.Diagnostics.HasError() {
		return
	}

	state := catalogMediasDataSourceModel{
		ID:          utils.GenerateUUID("catalog_medias"),
		Medias:      listMedias,
		MediasName:  listMediasName,
		CatalogName: config.CatalogName,
		CatalogID:   config.CatalogID,
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *catalogMediasDataSource) GetID() string {
	return d.catalog.id
}

// GetName returns the name of the catalog.
func (d *catalogMediasDataSource) GetName() string {
	return d.catalog.name
}

// GetIDOrName returns the ID if it is set, otherwise it returns the name.
func (d *catalogMediasDataSource) GetIDOrName() string {
	if d.GetID() != "" {
		return d.GetID()
	}
	return d.GetName()
}

// GetCatalog returns the govcd.Catalog.
func (d *catalogMediasDataSource) GetCatalog() (*govcd.AdminCatalog, error) {
	return d.adminOrg.GetAdminCatalogByNameOrId(d.GetIDOrName(), true)
}
