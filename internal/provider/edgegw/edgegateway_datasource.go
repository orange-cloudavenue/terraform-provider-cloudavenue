// Package edgegw provides a Terraform resource to manage edge gateways.
package edgegw

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	apiclient "github.com/orange-cloudavenue/cloudavenue-sdk-go"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common"
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
	Name        types.String `tfsdk:"name"`
	OwnerType   types.String `tfsdk:"owner_type"`
	OwnerName   types.String `tfsdk:"owner_name"`
	Description types.String `tfsdk:"description"`
}

func (d *edgeGatewayDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName
}

func (d *edgeGatewayDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The edge gateway data source show the details of the edge gateway.",

		Attributes: map[string]schema.Attribute{
			"tier0_vrf_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the Tier-0 VRF.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the Edge Gateway.",
				Required:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the Edge Gateway.",
				Computed:            true,
			},
			"owner_type": schema.StringAttribute{
				MarkdownDescription: "The type of the owner of the Edge Gateway.",
				Computed:            true,
			},
			"owner_name": schema.StringAttribute{
				MarkdownDescription: "The name of the owner of the Edge Gateway.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the Edge Gateway.",
				Computed:            true,
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
	var (
		data    edgeGatewayDataSourceModel
		gateway apiclient.EdgeGateway
	)

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	gateways, httpR, err := d.client.APIClient.EdgeGatewaysApi.GetEdges(d.client.Auth)
	if apiErr := helpers.CheckAPIError(err, httpR); apiErr != nil {
		defer httpR.Body.Close()
		resp.Diagnostics.Append(apiErr.GetTerraformDiagnostic())
		if resp.Diagnostics.HasError() {
			return
		}
	}

	found := false
	for _, gateway = range gateways {
		if data.Name.Equal(types.StringValue(gateway.EdgeName)) {
			found = true
			break
		}
	}

	if !found {
		data.ID = types.StringValue("")
	} else {
		data = edgeGatewayDataSourceModel{
			Tier0VrfID:  types.StringValue(gateway.Tier0VrfId),
			Name:        types.StringValue(gateway.EdgeName),
			ID:          types.StringValue(common.NormalizeID("urn:vcloud:gateway:", gateway.EdgeId)),
			OwnerType:   types.StringValue(gateway.OwnerType),
			OwnerName:   types.StringValue(gateway.OwnerName),
			Description: types.StringValue(gateway.Description),
		}
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
