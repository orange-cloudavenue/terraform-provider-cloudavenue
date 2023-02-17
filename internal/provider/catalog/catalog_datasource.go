// Package catalog provides a Terraform resource to manage catalogs.
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
)

var (
	_ datasource.DataSource              = &catalogDataSource{}
	_ datasource.DataSourceWithConfigure = &catalogDataSource{}
)

// NewCatalogDataSource returns a new resource implementing the catalog data source.
func NewCatalogDataSource() datasource.DataSource {
	return &catalogDataSource{}
}

type catalogDataSource struct {
	client *client.CloudAvenue
}

type catalogDataSourceModel struct {
	ID                          types.String `tfsdk:"id"`
	CatalogName                 types.String `tfsdk:"catalog_name"`
	CreatedAt                   types.String `tfsdk:"created_at"`
	Description                 types.String `tfsdk:"description"`
	PreserveIdentityInformation types.Bool   `tfsdk:"preserve_identity_information"`
	Href                        types.String `tfsdk:"href"`
	OwnerName                   types.String `tfsdk:"owner_name"`
	NumberOfMedia               types.Int64  `tfsdk:"number_of_media"`
	MediaItemList               types.List   `tfsdk:"media_item_list"`
	IsShared                    types.Bool   `tfsdk:"is_shared"`
	IsPublished                 types.Bool   `tfsdk:"is_published"`
	IsLocal                     types.Bool   `tfsdk:"is_local"`
	IsCached                    types.Bool   `tfsdk:"is_cached"`
}

func (d *catalogDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_catalog"
}

func (d *catalogDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The catalog data source show the details of the catalog.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"catalog_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Name of the catalog.",
			},
			"created_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Time stamp of when the catalog was created",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Description of the catalog.",
			},
			"preserve_identity_information": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Include BIOS UUIDs and MAC addresses in the downloaded OVF package. Preserving the identity information limits the portability of the package and you should use it only when necessary.",
			},
			"href": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Catalog HREF",
			},
			"owner_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Owner name from the catalog.",
			},
			"number_of_media": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Number of Medias this catalog contains.",
			},
			"media_item_list": schema.ListAttribute{
				Computed:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "List of Media items in this catalog",
			},
			"is_shared": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "True if this catalog is shared.",
			},
			"is_local": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "True if this catalog belongs to the current organization.",
			},
			"is_published": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "True if this catalog is shared to all organizations.",
			},
			"is_cached": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "True if this catalog is cached.",
			},
		},
	}
}

func (d *catalogDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
	var data catalogDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// catalog creation is accessible only for administator account
	adminOrg, err := d.client.Vmware.GetAdminOrgByNameOrId(d.client.GetOrg())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving Org", err.Error())
		return
	}

	catalog, err := adminOrg.GetAdminCatalogByNameOrId(data.CatalogName.ValueString(), false)
	if err != nil {
		if govcd.ContainsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Unable to query catalog records", fmt.Sprintf("Unable to query catalog records: %s", err))
		return
	}

	state := catalogDataSourceModel{
		ID:          types.StringValue(catalog.AdminCatalog.ID),
		CatalogName: types.StringValue(catalog.AdminCatalog.Name),
		CreatedAt:   types.StringValue(catalog.AdminCatalog.DateCreated),
		Description: types.StringValue(catalog.AdminCatalog.Description),
		Href:        types.StringValue(catalog.AdminCatalog.HREF),
		OwnerName:   types.StringValue(catalog.AdminCatalog.Owner.User.Name),
		IsPublished: types.BoolValue(catalog.AdminCatalog.IsPublished),
		IsLocal:     types.BoolValue(!catalog.AdminCatalog.IsPublished),
	}

	if catalog.AdminCatalog.PublishExternalCatalogParams != nil {
		state.IsCached = types.BoolValue(*catalog.AdminCatalog.PublishExternalCatalogParams.IsCachedEnabled)
		state.IsShared = types.BoolValue(*catalog.AdminCatalog.PublishExternalCatalogParams.IsPublishedExternally)
		state.PreserveIdentityInformation = types.BoolValue(*catalog.AdminCatalog.PublishExternalCatalogParams.PreserveIdentityInfoFlag)
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

	state.MediaItemList = basetypes.NewListValueMust(types.StringType, rawMediaItemsList)
	state.NumberOfMedia = types.Int64Value(int64(len(mediaItemList)))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
