// Package edgegw provides a Terraform resource to manage edge gateways.
package edgegw

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/uuid"
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

func (d *edgeGatewayDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName
}

func (d *edgeGatewayDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = edgegwSchema().GetDataSource(ctx)
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
	defer metrics.New("data.cloudavenue_edgegateway", d.client.GetOrgName(), metrics.Read)()

	config := &edgeGatewayDatasourceModel{}

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data := config.Copy()

	// Read data from the API
	edgegw, err := d.client.CAVSDK.V1.EdgeGateway.GetByName(config.Name.Get())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving edge gateway", err.Error())
		return
	}

	data.ID.Set(uuid.Normalize(uuid.Gateway, edgegw.GetID()).String())
	data.Tier0VrfID.Set(edgegw.GetTier0VrfID())
	data.OwnerName.Set(edgegw.GetOwnerName())
	data.OwnerType.Set(string(edgegw.GetOwnerType()))
	data.Description.Set(edgegw.GetDescription())
	data.Bandwidth.SetInt(int(edgegw.GetBandwidth()))

	// EnableLoadBalancing is now deprecated, but we still need to set it to false if it is unknown
	if !data.EnableLoadBalancing.IsKnown() {
		data.EnableLoadBalancing.Set(false)
	}

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
