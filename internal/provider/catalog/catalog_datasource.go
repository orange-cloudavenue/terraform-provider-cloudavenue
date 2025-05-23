/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

// Package catalog provides a Terraform resource to manage catalogs.
package catalog

import (
	"context"
	"fmt"
	"net/url"
	"sort"

	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/adminorg"
)

var (
	_ datasource.DataSource              = &catalogDataSource{}
	_ datasource.DataSourceWithConfigure = &catalogDataSource{}
	_ catalog                            = &catalogDataSource{}
)

// NewCatalogDataSource returns a new resource implementing the catalog data source.
func NewCatalogDataSource() datasource.DataSource {
	return &catalogDataSource{}
}

type catalogDataSource struct {
	client   *client.CloudAvenue
	adminOrg adminorg.AdminOrg
	catalog  base
}

func (d *catalogDataSource) Init(_ context.Context, rm *catalogDataSourceModel) (diags diag.Diagnostics) {
	d.catalog = base{
		name: rm.Name.ValueString(),
		id:   rm.ID.ValueString(),
	}

	d.adminOrg, diags = adminorg.Init(d.client)

	return
}

func (d *catalogDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName
}

func (d *catalogDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = catalogSchema().GetDataSource(ctx)
}

func (d *catalogDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.CloudAvenue)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *CloudAvenueClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *catalogDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	defer metrics.New("data.cloudavenue_catalog", d.client.GetOrgName(), metrics.Read)()

	state := &catalogDataSourceModel{}

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(d.Init(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	catalog, err := d.GetCatalog()
	if err != nil {
		resp.Diagnostics.AddError("Unable to find catalog", err.Error())
		return
	}

	updatedState := state
	updatedState.ID = types.StringValue(catalog.AdminCatalog.ID)
	updatedState.CreatedAt = types.StringValue(catalog.AdminCatalog.DateCreated)
	updatedState.Description = types.StringValue(catalog.AdminCatalog.Description)
	updatedState.IsPublished = types.BoolValue(catalog.AdminCatalog.IsPublished)
	updatedState.IsLocal = types.BoolValue(!catalog.AdminCatalog.IsPublished)

	if catalog.AdminCatalog.Owner != nil && catalog.AdminCatalog.Owner.User != nil {
		updatedState.OwnerName = types.StringValue(catalog.AdminCatalog.Owner.User.Name)
	}

	if catalog.AdminCatalog.PublishExternalCatalogParams != nil {
		// Fx Issue #657 - The IsCachedEnabled flag is not always set. Now use BoolPointerValue to avoid panic
		updatedState.IsCached = types.BoolPointerValue(catalog.AdminCatalog.PublishExternalCatalogParams.IsCachedEnabled)
		updatedState.IsShared = types.BoolPointerValue(catalog.AdminCatalog.PublishExternalCatalogParams.IsPublishedExternally)
		updatedState.PreserveIdentityInformation = types.BoolPointerValue(catalog.AdminCatalog.PublishExternalCatalogParams.PreserveIdentityInfoFlag)
	}

	var (
		rawMediaItemsList = make([]attr.Value, 0)
		mediaItemList     = make([]string, 0)
	)

	filter := fmt.Sprintf("catalog==%s", url.QueryEscape(catalog.AdminCatalog.HREF))
	mediaResults, err := d.client.Vmware.QueryWithNotEncodedParams(nil, map[string]string{"type": "media", "filter": filter, "filterEncoded": "true"})
	if err != nil {
		resp.Diagnostics.AddWarning(
			"Unable to get media records",
			fmt.Sprintf("Unable to get media records: %s", err),
		)
		return
	}

	for _, media := range mediaResults.Results.MediaRecord {
		mediaItemList = append(mediaItemList, media.Name)
	}

	// Sort the lists, so that they will always match in state
	sort.Strings(mediaItemList)
	for _, mediaName := range mediaItemList {
		rawMediaItemsList = append(rawMediaItemsList, types.StringValue(mediaName))
	}

	updatedState.MediaItemList = basetypes.NewListValueMust(types.StringType, rawMediaItemsList)
	updatedState.NumberOfMedia = types.Int64Value(int64(len(mediaItemList)))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, updatedState)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// GetID returns the ID of the catalog.
func (d *catalogDataSource) GetID() string {
	return d.catalog.id
}

// GetName returns the name of the catalog.
func (d *catalogDataSource) GetName() string {
	return d.catalog.name
}

// GetIDOrName returns the ID if it is set, otherwise it returns the name.
func (d *catalogDataSource) GetIDOrName() string {
	if d.GetID() != "" {
		return d.GetID()
	}
	return d.GetName()
}

// GetCatalog returns the govcd.Catalog.
func (d *catalogDataSource) GetCatalog() (*govcd.AdminCatalog, error) {
	return d.adminOrg.GetAdminCatalogByNameOrId(d.GetIDOrName(), true)
}
