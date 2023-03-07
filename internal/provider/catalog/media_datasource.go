// Package catalog provides a Terraform datasource.
package catalog

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
)

var (
	_ datasource.DataSource              = &catalogMediaDataSource{}
	_ datasource.DataSourceWithConfigure = &catalogMediaDataSource{}
)

func NewCatalogMediaDataSource() datasource.DataSource {
	return &catalogMediaDataSource{}
}

type catalogMediaDataSource struct {
	client *client.CloudAvenue
}

type catalogMediaDataSourceModel struct {
	ID             types.String `tfsdk:"id"`
	CatalogID      types.String `tfsdk:"catalog_id"`
	CatalogName    types.String `tfsdk:"catalog_name"`
	Name           types.String `tfsdk:"name"`
	Description    types.String `tfsdk:"description"`
	IsISO          types.Bool   `tfsdk:"is_iso"`
	OwnerName      types.String `tfsdk:"owner_name"`
	IsPublished    types.Bool   `tfsdk:"is_published"`
	CreatedAt      types.String `tfsdk:"created_at"`
	Size           types.Int64  `tfsdk:"size"`
	Status         types.String `tfsdk:"status"`
	StorageProfile types.String `tfsdk:"storage_profile"`
}

func (d *catalogMediaDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_" + "media"
}

func (d *catalogMediaDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The `catalog_media` datasource provides a CloudAvenue Catalog media data source.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID of the catalog media.",
			},
			"catalog_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The ID of the catalog to which media file belongs.",
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot("catalog_name")),
				},
			},
			"catalog_name": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The name of the catalog to which media file belongs.",
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot("catalog_id")),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the media.",
			},
			"description": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The description of the media.",
			},
			"is_iso": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "True if this media file is an Iso.",
			},
			"owner_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The name of the owner.",
			},
			"is_published": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "True if this media file is in a published catalog.",
			},
			"created_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The creation date of the media.",
			},
			"size": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The size of the media in bytes.",
			},
			"status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The media status.",
			},
			"storage_profile": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The name of the storage profile.",
			},
		},
	}
}

func (d *catalogMediaDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
	var data catalogMediaDataSourceModel
	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var (
		catalog *govcd.Catalog
		err     error
		media   *govcd.Media
	)

	// Check if catalog_id is set
	if data.CatalogID.IsNull() || data.CatalogID.IsUnknown() {
		// If not, try to find it using catalog_name
		catalog, err = d.client.Vmware.Client.GetCatalogByName(d.client.Org, data.CatalogName.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Unable to find catalog", err.Error())
			return
		}
		data.CatalogID = types.StringValue(catalog.Catalog.ID)
	} else {
		catalog, err = d.client.Vmware.Client.GetCatalogById(data.CatalogID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Unable to find catalog", err.Error())
			return
		}
		data.CatalogName = types.StringValue(catalog.Catalog.Name)
	}

	media, err = catalog.GetMediaByNameOrId(data.Name.ValueString(), false)
	if err != nil {
		resp.Diagnostics.AddError("Unable to find media", err.Error())
		return
	}

	data.ID = types.StringValue(media.Media.ID)

	mediaRecord, err := catalog.QueryMedia(media.Media.Name)
	if err != nil {
		resp.Diagnostics.AddError("Unable to query media", err.Error())
		return
	}
	// Check if mediaRecord.MediaRecord is nil
	if mediaRecord.MediaRecord == nil {
		resp.Diagnostics.AddError("Unable to find media record", "mediaRecord.MediaRecord is nil")
		return
	}
	// Check if catalog.Catalog is nil
	if catalog.Catalog == nil {
		resp.Diagnostics.AddError("Unable to find catalog ID", "catalog.Catalog.ID is nil")
		return
	}
	// Check if media.Media is nil
	if media.Media == nil {
		resp.Diagnostics.AddError("Unable to find media ID", "media.Media.ID is nil")
		return
	}

	data.CatalogName = types.StringValue(catalog.Catalog.Name)
	data.CatalogID = types.StringValue(catalog.Catalog.ID)
	data.Name = types.StringValue(media.Media.Name)
	data.Description = types.StringValue(media.Media.Description)
	data.IsISO = types.BoolValue(mediaRecord.MediaRecord.IsIso)
	data.OwnerName = types.StringValue(mediaRecord.MediaRecord.OwnerName)
	data.IsPublished = types.BoolValue(mediaRecord.MediaRecord.IsPublished)
	data.CreatedAt = types.StringValue(mediaRecord.MediaRecord.CreationDate)
	data.Size = types.Int64Value(mediaRecord.MediaRecord.StorageB)
	data.Status = types.StringValue(mediaRecord.MediaRecord.Status)
	data.StorageProfile = types.StringValue(mediaRecord.MediaRecord.StorageProfileName)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
