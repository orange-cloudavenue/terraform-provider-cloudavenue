// Package catalog provides a Terraform datasource.
package catalog

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

var (
	_ datasource.DataSource              = &catalogMediasDataSource{}
	_ datasource.DataSourceWithConfigure = &catalogMediasDataSource{}
)

func NewCatalogMediasDataSource() datasource.DataSource {
	return &catalogMediasDataSource{}
}

type catalogMediasDataSource struct {
	client *client.CloudAvenue
}

type catalogMediasDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Medias      types.Map    `tfsdk:"medias"`
	MediasName  types.List   `tfsdk:"medias_name"`
	CatalogName types.String `tfsdk:"catalog_name"`
	CatalogID   types.String `tfsdk:"catalog_id"`
}

func catalogMediasAttrType() map[string]attr.Type {
	return map[string]attr.Type{
		"id":              types.StringType,
		"name":            types.StringType,
		"description":     types.StringType,
		"is_iso":          types.BoolType,
		"owner_name":      types.StringType,
		"is_published":    types.BoolType,
		"created_at":      types.StringType,
		"size":            types.Int64Type,
		"status":          types.StringType,
		"storage_profile": types.StringType,
	}
}

func (d *catalogMediasDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_" + "medias"
}

func (d *catalogMediasDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The `catalog_medias` datasource allows you to list all the medias of a catalog.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"medias": schema.MapNestedAttribute{
				MarkdownDescription: "The map of medias.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: schemaCatalogDataSource(),
				},
			},
			"medias_name": schema.ListAttribute{
				MarkdownDescription: "The list of medias name.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"catalog_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The ID of the catalog to which media file belongs. Required if `catalog_name` is not set. ",
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot("catalog_name")),
				},
			},
			"catalog_name": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The name of the catalog to which media file belongs. Required if `catalog_id` is not set.",
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot("catalog_id")),
				},
			},
		},
	}
}

func (d *catalogMediasDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
	var data catalogMediasDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var (
		catalog *govcd.Catalog
		err     error
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

	medias := make(map[string]catalogMediaDataStruct)
	mediasName := make([]string, 0)
	x, err := catalog.QueryMediaList()
	// Get all medias
	for _, media := range x {
		s := catalogMediaDataStruct{
			ID:             types.StringValue(media.ID),
			Name:           types.StringValue(media.Name),
			IsISO:          types.BoolValue(media.IsIso),
			OwnerName:      types.StringValue(media.OwnerName),
			IsPublished:    types.BoolValue(media.IsPublished),
			CreatedAt:      types.StringValue(media.CreationDate),
			Status:         types.StringValue(media.Status),
			Size:           types.Int64Value(media.StorageB),
			StorageProfile: types.StringValue(media.StorageProfileName),
		}
		mediasName = append(mediasName, media.Name)
		medias[media.Name] = s
	}

	cn, diag := types.ListValueFrom(ctx, types.StringType, mediasName)
	resp.Diagnostics.Append(diag...)
	if resp.Diagnostics.HasError() {
		return
	}

	z, diag := types.MapValueFrom(ctx, types.ObjectType{AttrTypes: catalogMediasAttrType()}, medias)
	resp.Diagnostics.Append(diag...)
	if resp.Diagnostics.HasError() {
		return
	}

	state := &catalogMediasDataSourceModel{
		ID:          utils.GenerateUUID("catalog_medias"),
		Medias:      z,
		MediasName:  cn,
		CatalogName: types.StringValue(catalog.Catalog.Name),
		CatalogID:   types.StringValue(catalog.Catalog.ID),
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
