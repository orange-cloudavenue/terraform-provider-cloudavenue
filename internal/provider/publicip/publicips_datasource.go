// Package publicip provides a Terraform resource to manage public IPs.
package publicip

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers"
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
	client *client.CloudAvenue
}

type publicIPDataSourceModel struct {
	ID        types.String                 `tfsdk:"id"`
	PublicIPs []publicIPNetworkConfigModel `tfsdk:"public_ips"`
}

type publicIPNetworkConfigModel struct {
	ID              types.String `tfsdk:"id"`
	EdgeGatewayName types.String `tfsdk:"edge_gateway_name"`
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
						"id": schema.StringAttribute{
							MarkdownDescription: "The public IP address.",
							Computed:            true,
						},
						"edge_gateway_name": schema.StringAttribute{
							MarkdownDescription: "The name of the edge gateway related to the public ip. This properties is only present for NGP.",
							Computed:            true,
						},
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
	var data publicIPDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	publicIPs, httpR, err := d.client.APIClient.PublicIPApi.GetPublicIPs(d.client.Auth)
	if apiErr := helpers.CheckAPIError(err, httpR); apiErr != nil {
		defer func() {
			err = errors.Join(err, httpR.Body.Close())
		}()
		resp.Diagnostics.Append(apiErr.GetTerraformDiagnostic())
		if resp.Diagnostics.HasError() {
			return
		}
		// If 404 return
		return
	}

	listOfIps := make([]string, 0)

	for _, cfg := range publicIPs.NetworkConfig {
		data.PublicIPs = append(data.PublicIPs, publicIPNetworkConfigModel{
			ID:              types.StringValue(cfg.UplinkIp),
			EdgeGatewayName: types.StringValue(cfg.EdgeGatewayName),
		})

		listOfIps = append(listOfIps, cfg.UplinkIp)
	}

	data.ID = utils.GenerateUUID(listOfIps)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
