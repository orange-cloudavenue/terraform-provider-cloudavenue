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
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

var (
	_ datasource.DataSource              = &catalogsDataSource{}
	_ datasource.DataSourceWithConfigure = &catalogsDataSource{}
)

type catalogsDataSource struct {
	client   *client.CloudAvenue
	adminOrg adminorg.AdminOrg
}

func NewCatalogsDataSource() datasource.DataSource {
	return &catalogsDataSource{}
}

func (d *catalogsDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = catalogsSuperSchema(ctx).GetDataSource(ctx)
}

func (d *catalogsDataSource) Init(_ context.Context, _ *catalogsDataSourceModel) (diags diag.Diagnostics) {
	d.adminOrg, diags = adminorg.Init(d.client)
	return
}

func (d *catalogsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "s"
}

func (d *catalogsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *catalogsDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	defer metrics.New("data.cloudavenue_catalogs", d.client.GetOrgName(), metrics.Read)()

	state := &catalogsDataSourceModel{}
	resp.Diagnostics.Append(d.Init(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	catalogs := make(map[string]catalogDataSourceModel)
	catalogsName := make([]string, 0)

	for _, x := range d.adminOrg.ListCatalogs().Catalog {
		catalog, err := d.adminOrg.GetAdminCatalogByNameOrId(x.Name, false)
		if err != nil {
			if govcd.ContainsNotFound(err) {
				return
			}
			resp.Diagnostics.AddError("Unable to get catalog", err.Error())
			continue
		}

		s := catalogDataSourceModel{
			ID:                          types.StringValue(catalog.AdminCatalog.ID),
			Name:                        types.StringValue(catalog.AdminCatalog.Name),
			CreatedAt:                   types.StringValue(catalog.AdminCatalog.DateCreated),
			Description:                 utils.StringValueOrNull(catalog.AdminCatalog.Description),
			IsPublished:                 types.BoolValue(catalog.AdminCatalog.IsPublished),
			IsLocal:                     types.BoolValue(!catalog.AdminCatalog.IsPublished),
			IsCached:                    types.BoolNull(),
			IsShared:                    types.BoolNull(),
			PreserveIdentityInformation: types.BoolNull(),
			OwnerName:                   types.StringNull(),
			MediaItemList:               types.ListNull(types.StringType),
			NumberOfMedia:               types.Int64Null(),
		}

		catalogsName = append(catalogsName, catalog.AdminCatalog.Name)

		if catalog.AdminCatalog.Owner != nil && catalog.AdminCatalog.Owner.User != nil {
			s.OwnerName = utils.StringValueOrNull(catalog.AdminCatalog.Owner.User.Name)
		}

		if catalog.AdminCatalog.PublishExternalCatalogParams != nil {
			if catalog.AdminCatalog.PublishExternalCatalogParams.IsCachedEnabled != nil {
				s.IsCached = types.BoolValue(*catalog.AdminCatalog.PublishExternalCatalogParams.IsCachedEnabled)
			}
			if catalog.AdminCatalog.PublishExternalCatalogParams.IsPublishedExternally != nil {
				s.IsShared = types.BoolValue(*catalog.AdminCatalog.PublishExternalCatalogParams.IsPublishedExternally)
			}
			if catalog.AdminCatalog.PublishExternalCatalogParams.PreserveIdentityInfoFlag != nil {
				s.PreserveIdentityInformation = types.BoolValue(*catalog.AdminCatalog.PublishExternalCatalogParams.PreserveIdentityInfoFlag)
			}
		}

		var (
			rawMediaItemsList []attr.Value
			mediaItemList     []string
		)

		filter := fmt.Sprintf("catalog==%s", url.QueryEscape(catalog.AdminCatalog.HREF))
		mediaResults, err := d.client.Vmware.QueryWithNotEncodedParams(nil, map[string]string{"type": "media", "filter": filter, "filterEncoded": "true"})
		if err != nil {
			resp.Diagnostics.AddWarning(
				"Unable to get media records",
				fmt.Sprintf("Unable to get media records: %s", err),
			)
			continue
		}

		for _, media := range mediaResults.Results.MediaRecord {
			mediaItemList = append(mediaItemList, media.Name)
		}
		// Sort the lists, so that they will always match in state
		sort.Strings(mediaItemList)
		for _, mediaName := range mediaItemList {
			rawMediaItemsList = append(rawMediaItemsList, types.StringValue(mediaName))
		}
		s.MediaItemList = basetypes.NewListValueMust(types.StringType, rawMediaItemsList)
		s.NumberOfMedia = types.Int64Value(int64(len(mediaItemList)))

		catalogs[catalog.AdminCatalog.Name] = s
	}

	sort.Strings(catalogsName)
	cn, diag := types.ListValueFrom(ctx, types.StringType, catalogsName)
	resp.Diagnostics.Append(diag...)
	if resp.Diagnostics.HasError() {
		return
	}

	updatedState := catalogsDataSourceModel{
		ID:           utils.GenerateUUID("catalogs"),
		Catalogs:     catalogs,
		CatalogsName: cn,
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, updatedState)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
