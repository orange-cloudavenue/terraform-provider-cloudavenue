package edgegw

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers"
)

var (
	_ datasource.DataSource              = &edgeGatewaysDataSource{}
	_ datasource.DataSourceWithConfigure = &edgeGatewaysDataSource{}
)

func NewEdgeGatewaysDataSource() datasource.DataSource {
	return &edgeGatewaysDataSource{}
}

type edgeGatewaysDataSource struct {
	client *client.CloudAvenue
}

type edgeGatewaysDataSourceModel struct {
	ID           types.String `tfsdk:"id"`
	EdgeGateways []gateway    `tfsdk:"edge_gateways"`
}

type gateway struct {
	Tier0VrfID  types.String `tfsdk:"tier0_vrf_id"`
	EdgeName    types.String `tfsdk:"edge_name"`
	EdgeID      types.String `tfsdk:"edge_id"`
	OwnerType   types.String `tfsdk:"owner_type"`
	OwnerName   types.String `tfsdk:"owner_name"`
	Description types.String `tfsdk:"description"`
}

func (d *edgeGatewaysDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_edge_gateways"
}

func (d *edgeGatewaysDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The edge gateways data source show the list of edge gateways of an organization.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"edge_gateways": schema.ListNestedAttribute{
				Computed:    true,
				Description: "A list of Edge Gateways.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"tier0_vrf_id": schema.StringAttribute{
							Description: "The ID of the Tier-0 VRF.",
							Computed:    true,
						},
						"edge_name": schema.StringAttribute{
							Description: "The name of the Edge Gateway.",
							Computed:    true,
						},
						"edge_id": schema.StringAttribute{
							Description: "The ID of the Edge Gateway.",
							Computed:    true,
						},
						"owner_type": schema.StringAttribute{
							Description: "The type of the owner of the Edge Gateway.",
							Computed:    true,
						},
						"owner_name": schema.StringAttribute{
							Description: "The name of the owner of the Edge Gateway.",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "The description of the Edge Gateway.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *edgeGatewaysDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *edgeGatewaysDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data edgeGatewaysDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	gateways, httpR, err := d.client.APIClient.EdgeGatewaysApi.GetEdges(d.client.Auth)
	if x := helpers.CheckAPIError(err, httpR); x != nil {
		resp.Diagnostics.Append(x.GetTerraformDiagnostic())
		if resp.Diagnostics.HasError() {
			return
		}

		// Is Not Found
		data.EdgeGateways = []gateway{}
		data.ID = types.StringValue("")
	} else {
		for _, gw := range gateways {
			d := gateway{
				Tier0VrfID:  types.StringValue(gw.Tier0VrfId),
				EdgeName:    types.StringValue(gw.EdgeName),
				EdgeID:      types.StringValue(gw.EdgeId),
				OwnerType:   types.StringValue(gw.OwnerType),
				OwnerName:   types.StringValue(gw.OwnerName),
				Description: types.StringValue(gw.Description),
			}
			data.EdgeGateways = append(data.EdgeGateways, d)
		}

		data.ID = types.StringValue("frangipane")
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
