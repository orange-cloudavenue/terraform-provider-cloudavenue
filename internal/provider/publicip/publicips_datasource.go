// Package publicip provides a Terraform resource to manage public IPs.
package publicip

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/adminorg"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/edgegw"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

var (
	_ datasource.DataSource              = &publicIPDataSource{}
	_ datasource.DataSourceWithConfigure = &publicIPDataSource{}
)

// NewPublicIPDataSource returns a new resource implementing the public IP data source.
func NewPublicIPDataSource() datasource.DataSource {
	return &publicIPDataSource{}
}

type publicIPDataSource struct {
	client   *client.CloudAvenue
	adminOrg adminorg.AdminOrg
}

// Init.
func (d *publicIPDataSource) Init(_ context.Context, rm *publicIPDataSourceModel) (diags diag.Diagnostics) {
	d.adminOrg, diags = adminorg.Init(d.client)

	return
}

func (d *publicIPDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "s"
}

func (d *publicIPDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = publicIPsSchema(ctx)
}

func (d *publicIPDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *publicIPDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	defer metrics.New("data.cloudavenue_publicips", d.client.GetOrgName(), metrics.Read)()

	data := &publicIPDataSourceModel{}

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(d.Init(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	publicIPs, httpR, err := d.client.APIClient.PublicIPApi.GetPublicIPs(d.client.Auth)
	if httpR != nil {
		defer func() {
			err = errors.Join(err, httpR.Body.Close())
		}()
	}

	if apiErr := helpers.CheckAPIError(err, httpR); apiErr != nil {
		resp.Diagnostics.Append(apiErr.GetTerraformDiagnostic())
		return
	}

	listOfIps := make([]string, 0)

	for _, cfg := range publicIPs.NetworkConfig {
		edgeGateway, err := d.adminOrg.GetEdgeGateway(edgegw.BaseEdgeGW{
			Name: types.StringValue(cfg.EdgeGatewayName),
		})
		if err != nil {
			resp.Diagnostics.AddError("Error while getting edge gateway", err.Error())
			return
		}

		x := publicIPNetworkConfigModel{
			ID:              types.StringValue(cfg.UplinkIp),
			EdgeGatewayName: types.StringValue(edgeGateway.GetName()),
			EdgeGatewayID:   types.StringValue(edgeGateway.GetID()),
			PublicIP:        types.StringValue(cfg.UplinkIp),
		}

		data.PublicIPs = append(data.PublicIPs, x)
		listOfIps = append(listOfIps, cfg.UplinkIp)
	}

	data.ID = utils.GenerateUUID(listOfIps)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
