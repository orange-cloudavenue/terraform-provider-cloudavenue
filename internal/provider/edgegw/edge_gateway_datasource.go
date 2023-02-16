// Package edgegw provides a Terraform resource to manage edge gateways.
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

// NewEdgeGatewayDataSource returns a new datasource implementing the edge_gateway data source.
func NewEdgeGatewayDataSource() datasource.DataSource {
	return &edgeGatewayDataSource{}
}

type edgeGatewayDataSource struct {
	client *client.CloudAvenue
}

type edgeGatewayDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Tier0VrfID  types.String `tfsdk:"tier0_vrf_id"`
	EdgeName    types.String `tfsdk:"edge_name"`
	EdgeID      types.String `tfsdk:"edge_id"`
	OwnerType   types.String `tfsdk:"owner_type"`
	OwnerName   types.String `tfsdk:"owner_name"`
	Description types.String `tfsdk:"description"`
}

func (d *edgeGatewayDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_edge_gateway"
}

func (d *edgeGatewayDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The edge gateway data source show the details of the edge gateway.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
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
				Required:    true,
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
	}
}

func (d *edgeGatewayDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *edgeGatewayDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data edgeGatewayDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	gateway, httpR, err := d.client.APIClient.EdgeGatewaysApi.GetEdgeById(d.client.Auth, data.EdgeID.ValueString())
	if x := helpers.CheckAPIError(err, httpR); x != nil {
		defer httpR.Body.Close()
		resp.Diagnostics.Append(x.GetTerraformDiagnostic())
		if resp.Diagnostics.HasError() {
			return
		}

		// Is Not Found
		data.ID = types.StringValue("")
	} else {
		data = edgeGatewayDataSourceModel{
			Tier0VrfID:  types.StringValue(gateway.Tier0VrfId),
			EdgeName:    types.StringValue(gateway.EdgeName),
			EdgeID:      types.StringValue(gateway.EdgeId),
			OwnerType:   types.StringValue(gateway.OwnerType),
			OwnerName:   types.StringValue(gateway.OwnerName),
			Description: types.StringValue(gateway.Description),
		}

		data.ID = types.StringValue("frangipane")
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
