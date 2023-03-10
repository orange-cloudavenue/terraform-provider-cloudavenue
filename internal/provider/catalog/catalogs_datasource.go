// Package catalog provides a Terraform datasource.
package catalog

import (
	"context"
	"fmt"
	"net/url"
	"sort"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

var (
	_ datasource.DataSource              = &catalogsDataSource{}
	_ datasource.DataSourceWithConfigure = &catalogsDataSource{}
)

func NewCatalogsDataSource() datasource.DataSource {
	return &catalogsDataSource{}
}

type catalogsDataSource struct {
	client *client.CloudAvenue
}

type catalogsDataSourceModel struct {
	ID           types.String                       `tfsdk:"id"`
	Catalogs     map[string]catalogDataSourceStruct `tfsdk:"catalogs"`
	CatalogsName types.List                         `tfsdk:"catalogs_name"`
}

func (d *catalogsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "s"
}

func (d *catalogsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The catalogs datasource show the details of all the catalogs.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"catalogs_name": schema.ListAttribute{
				MarkdownDescription: "List of catalogs name.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"catalogs": schema.MapNestedAttribute{
				MarkdownDescription: "Map of catalogs.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: schemaDataSource(),
				},
			},
		},
	}
}

func (d *catalogsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *catalogsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data catalogsDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	adminOrg, err := d.client.Vmware.GetAdminOrgByNameOrId(d.client.GetOrg())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving Org", err.Error())
		return
	}

	catalogs := make(map[string]catalogDataSourceStruct)
	catalogsName := make([]string, 0)

	for _, x := range adminOrg.AdminOrg.Catalogs.Catalog {
		catalog, err := adminOrg.GetAdminCatalogByNameOrId(x.Name, false)
		if err != nil {
			if govcd.ContainsNotFound(err) {
				return
			}
			resp.Diagnostics.AddError("Unable to get catalog", err.Error())
			continue
		} else {
			s := catalogDataSourceStruct{
				ID:          types.StringValue(catalog.AdminCatalog.ID),
				CatalogName: types.StringValue(catalog.AdminCatalog.Name),
				CreatedAt:   types.StringValue(catalog.AdminCatalog.DateCreated),
				Description: types.StringValue(catalog.AdminCatalog.Description),
				Href:        types.StringValue(catalog.AdminCatalog.HREF),
				IsPublished: types.BoolValue(catalog.AdminCatalog.IsPublished),
				IsLocal:     types.BoolValue(!catalog.AdminCatalog.IsPublished),
			}

			catalogsName = append(catalogsName, catalog.AdminCatalog.Name)

			if catalog.AdminCatalog.Owner != nil && catalog.AdminCatalog.Owner.User != nil {
				s.OwnerName = types.StringValue(catalog.AdminCatalog.Owner.User.Name)
			} else {
				s.OwnerName = types.StringValue("")
			}

			if catalog.AdminCatalog.PublishExternalCatalogParams != nil {
				s.IsCached = types.BoolValue(*catalog.AdminCatalog.PublishExternalCatalogParams.IsCachedEnabled)
				s.IsShared = types.BoolValue(*catalog.AdminCatalog.PublishExternalCatalogParams.IsPublishedExternally)
				s.PreserveIdentityInformation = types.BoolValue(*catalog.AdminCatalog.PublishExternalCatalogParams.PreserveIdentityInfoFlag)
			}

			var rawMediaItemsList []attr.Value
			var mediaItemList []string

			filter := fmt.Sprintf("catalog==%s", url.QueryEscape(catalog.AdminCatalog.HREF))
			mediaResults, err := d.client.Vmware.QueryWithNotEncodedParams(nil, map[string]string{"type": "media", "filter": filter, "filterEncoded": "true"})
			if err != nil {
				resp.Diagnostics.AddWarning(
					"Unable to get media records",
					fmt.Sprintf("Unable to get media records: %s", err),
				)
				continue
			} else {
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
			}

			catalogs[catalog.AdminCatalog.Name] = s
		}
	}

	sort.Strings(catalogsName)

	cn, diag := types.ListValueFrom(ctx, types.StringType, catalogsName)
	resp.Diagnostics.Append(diag...)
	if resp.Diagnostics.HasError() {
		return
	}

	state := &catalogsDataSourceModel{
		ID:           utils.GenerateUUID("catalogs"),
		Catalogs:     catalogs,
		CatalogsName: cn,
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
