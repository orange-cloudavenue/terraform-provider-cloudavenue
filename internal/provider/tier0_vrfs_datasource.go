package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

var (
	_ datasource.DataSource              = &tier0VrfsDataSource{}
	_ datasource.DataSourceWithConfigure = &tier0VrfsDataSource{}
)

func NewTier0VrfsDataSource() datasource.DataSource {
	return &tier0VrfsDataSource{}
}

type tier0VrfsDataSource struct {
	client *CloudAvenueClient
}

type tier0VrfsDataSourceModel struct {
	ID        types.String     `tfsdk:"id"`
	Tier0Vrfs []tier0VrfsModel `tfsdk:"tier0_vrfs"`
}

type tier0VrfsModel struct {
	Name types.String `tfsdk:"name"`
}

func (d *tier0VrfsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tier0_vrfs"
}

func (d *tier0VrfsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The Tier-0 VRFs data source allow access to a list of Tier-0 that can be accessed by the user.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"tier0_vrfs": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "A list of Tier-0 VRFs.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "The name of the Tier-0 VRF.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *tier0VrfsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*CloudAvenueClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *CloudAvenueClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *tier0VrfsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data tier0VrfsDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tier0vrfs, _, err := d.client.Tier0Api.ApiCustomersV20Tier0VrfsGet(d.client.auth)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read example, got error: %s", err))
		return
	}

	names := make([]string, 0, len(tier0vrfs))
	for _, tier0vrf := range tier0vrfs {
		name := tier0VrfsModel{
			Name: types.StringValue(tier0vrf),
		}
		data.Tier0Vrfs = append(data.Tier0Vrfs, name)
		names = append(names, tier0vrf)
	}

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read a data source")

	// Generate a UUID from the list of names
	data.ID = utils.GenerateUUIDFromList(names)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
