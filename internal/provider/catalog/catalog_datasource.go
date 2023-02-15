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
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
)

var (
	_ datasource.DataSource              = &catalogDataSource{}
	_ datasource.DataSourceWithConfigure = &catalogDataSource{}
)

func NewCatalogDataSource() datasource.DataSource {
	return &catalogDataSource{}
}

type catalogDataSource struct {
	client *client.CloudAvenue
}

type catalogDataSourceModel struct {
	ID                          types.String `tfsdk:"id"`
	Name                        types.String `tfsdk:"name"`
	Created                     types.String `tfsdk:"created"`
	Description                 types.String `tfsdk:"description"`
	PublishEnabled              types.Bool   `tfsdk:"publish_enabled"`
	CacheEnabled                types.Bool   `tfsdk:"cache_enabled"`
	PreserveIdentityInformation types.Bool   `tfsdk:"preserve_identity_information"`
	Href                        types.String `tfsdk:"href"`
	CatalogVersion              types.Int64  `tfsdk:"catalog_version"`
	OwnerName                   types.String `tfsdk:"owner_name"`
	NumberOfMedia               types.Int64  `tfsdk:"number_of_media"`
	MediaItemList               types.List   `tfsdk:"media_item_list"`
	IsShared                    types.Bool   `tfsdk:"is_shared"`
	IsPublished                 types.Bool   `tfsdk:"is_published"`
	IsLocal                     types.Bool   `tfsdk:"is_local"`
	PublishSubscriptionType     types.String `tfsdk:"publish_subscription_type"`
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
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Name of the catalog.",
			},
			"created": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Time stamp of when the catalog was created",
			},
			"description": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"publish_enabled": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "True allows to publish a catalog externally to make its vApp templates and media files available for subscription by organizations outside the Cloud Director installation. Default is `false`.",
			},
			"cache_enabled": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "True enables early catalog export to optimize synchronization",
			},
			"preserve_identity_information": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Include BIOS UUIDs and MAC addresses in the downloaded OVF package. Preserving the identity information limits the portability of the package and you should use it only when necessary.",
			},
			"href": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Catalog HREF",
			},
			"catalog_version": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Catalog version number.",
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
			"publish_subscription_type": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "PUBLISHED if published externally, SUBSCRIBED if subscribed to an external catalog, UNPUBLISHED otherwise.",
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

	catalogRecords, err := d.client.Vmware.Client.QueryCatalogRecords(data.Name.String(), govcd.TenantContext{})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to query catalog records",
			fmt.Sprintf("Unable to query catalog records: %s", err),
		)
		return
	}

	resp.Diagnostics.AddWarning("Catalog records", fmt.Sprintf("Catalog records: %v", catalogRecords))

	if len(catalogRecords) == 0 {
		resp.Diagnostics.AddError(
			"Unable to find catalog",
			fmt.Sprintf("Unable to find catalog: %s", data.Name.String()),
		)
		return
	}

	var catalogRecord *govcdtypes.CatalogRecord

	for _, cr := range catalogRecords {
		// if d.client.GetDefaultOrg() == cr.OrgName {
		// 	catalogRecord = cr
		// 	break
		// }
		if cr.Name == data.Name.String() {
			catalogRecord = cr
			break
		}
	}

	if catalogRecord == nil {
		resp.Diagnostics.AddError(
			"Unable to get catalog record",
			fmt.Sprintf("Unable to get catalog record: %s", err),
		)
		return
	}

	record, err := d.client.Vmware.Client.GetAdminCatalogByHref(catalogRecord.HREF)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to get catalog record",
			fmt.Sprintf("Unable to get catalog record: %s", err),
		)
		return
	}

	data.ID = types.StringValue(record.AdminCatalog.ID)
	data.Created = types.StringValue(record.AdminCatalog.DateCreated)
	data.Description = types.StringValue(record.AdminCatalog.Description)
	data.Href = types.StringValue(record.AdminCatalog.HREF)

	if record.AdminCatalog.PublishExternalCatalogParams != nil {
		data.CacheEnabled = types.BoolValue(*record.AdminCatalog.PublishExternalCatalogParams.IsCachedEnabled)
		data.PublishEnabled = types.BoolValue(*record.AdminCatalog.PublishExternalCatalogParams.IsPublishedExternally)
		data.PreserveIdentityInformation = types.BoolValue(*record.AdminCatalog.PublishExternalCatalogParams.PreserveIdentityInfoFlag)
	}

	data.CatalogVersion = types.Int64Value(catalogRecord.Version)
	data.OwnerName = types.StringValue(catalogRecord.OwnerName)
	data.IsPublished = types.BoolValue(catalogRecord.IsPublished)
	data.IsShared = types.BoolValue(catalogRecord.IsShared)
	data.IsLocal = types.BoolValue(catalogRecord.IsLocal)
	data.PublishSubscriptionType = types.StringValue(catalogRecord.PublishSubscriptionType)

	var rawMediaItemsList []attr.Value
	var mediaItemList []string

	filter := fmt.Sprintf("catalog==%s", url.QueryEscape(data.Href.String()))
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

	data.MediaItemList = basetypes.NewListValueMust(types.StringType, rawMediaItemsList)
	data.NumberOfMedia = types.Int64Value(int64(len(mediaItemList)))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
