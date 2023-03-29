// Package catalog provides a Terraform datasource.
package catalog

import (
	"context"
	"fmt"

	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/adminorg"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

var (
	_ datasource.DataSource              = &catalogMediasDataSource{}
	_ datasource.DataSourceWithConfigure = &catalogMediasDataSource{}
	_ catalog                            = &vAppTemplateDataSource{}
)

func NewCatalogMediasDataSource() datasource.DataSource {
	return &catalogMediasDataSource{}
}

type catalogMediasDataSource struct {
	client   *client.CloudAvenue
	adminOrg adminorg.AdminOrg
	catalog  base
}

type catalogMediasDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Medias      types.Map    `tfsdk:"medias"`
	MediasName  types.List   `tfsdk:"medias_name"`
	CatalogName types.String `tfsdk:"catalog_name"`
	CatalogID   types.String `tfsdk:"catalog_id"`
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
					Attributes: mediaDatasourceAttributes(),
				},
			},
			"medias_name": schema.ListAttribute{
				MarkdownDescription: "The list of medias name.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			catalogName: mediaSchema().GetDataSource(ctx).Attributes[catalogName],
			catalogID:   mediaSchema().GetDataSource(ctx).Attributes[catalogID],
		},
	}
}

func (d *catalogMediasDataSource) Init(ctx context.Context, rm *catalogMediasDataSourceModel) (diags diag.Diagnostics) {
	d.catalog = base{
		name: rm.CatalogName.ValueString(),
		id:   rm.CatalogID.ValueString(),
	}

	d.adminOrg, diags = adminorg.Init(d.client)

	return
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
	state := &catalogMediasDataSourceModel{}

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

	var (
		medias     = make(map[string]catalogMediaDataSourceModel)
		mediasName = make([]string, 0)
	)

	// Get all medias
	mediaList, err := catalog.QueryMediaList()
	if err != nil {
		resp.Diagnostics.AddError("Unable to query media list", err.Error())
		return
	}

	for _, media := range mediaList {
		s := catalogMediaDataSourceModel{
			ID:             types.StringValue(media.ID),
			Name:           types.StringValue(media.Name),
			CatalogID:      state.CatalogID,
			CatalogName:    state.CatalogName,
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

	listMediasName, diag := types.ListValueFrom(ctx, types.StringType, mediasName)
	resp.Diagnostics.Append(diag...)
	if resp.Diagnostics.HasError() {
		return
	}

	listMedias, diag := types.MapValueFrom(ctx, types.ObjectType{AttrTypes: catalogMediaDataSourceModelType()}, medias)
	resp.Diagnostics.Append(diag...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateState := catalogMediasDataSourceModel{
		ID:          utils.GenerateUUID("catalog_medias"),
		Medias:      listMedias,
		MediasName:  listMediasName,
		CatalogName: state.CatalogName,
		CatalogID:   state.CatalogID,
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, updateState)...)
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
