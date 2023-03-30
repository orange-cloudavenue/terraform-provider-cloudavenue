// Package publicip provides a Terraform resource to manage public IPs.
package publicip

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers"
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

type publicIPDataSourceModel struct {
	ID        types.String                 `tfsdk:"id"`
	PublicIPs []publicIPNetworkConfigModel `tfsdk:"public_ips"`
}

type publicIPNetworkConfigModel struct {
	ID              types.String `tfsdk:"id"`
	PublicIP        types.String `tfsdk:"public_ip"`
	EdgeGatewayName types.String `tfsdk:"edge_gateway_name"`
	EdgeGatewayID   types.String `tfsdk:"edge_gateway_id"`
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
	resp.Schema = schema.Schema{
		Description: "The public IP data source show the list of public IP addresses.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"public_ips": schema.ListNestedAttribute{
				MarkdownDescription: "A list of public IPs.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":                publicIPSchema().GetDataSource(ctx).Attributes["id"],
						"public_ip":         publicIPSchema().GetDataSource(ctx).Attributes["public_ip"],
						"edge_gateway_name": publicIPSchema().GetDataSource(ctx).Attributes["edge_gateway_name"],
						"edge_gateway_id":   publicIPSchema().GetDataSource(ctx).Attributes["edge_gateway_id"],
					},
				},
			},
		},
	}
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
