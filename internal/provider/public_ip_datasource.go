package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

var (
	_ datasource.DataSource              = &publicIPDataSource{}
	_ datasource.DataSourceWithConfigure = &publicIPDataSource{}
)

func NewPublicIPDataSource() datasource.DataSource {
	return &publicIPDataSource{}
}

type publicIPDataSource struct {
	client *CloudAvenueClient
}

type publicIPDataSourceModel struct {
	ID            types.String                 `tfsdk:"id"`
	InternalIP    types.String                 `tfsdk:"internal_ip"`
	NetworkConfig []publicIPNetworkConfigModel `tfsdk:"network_config"`
}

type publicIPNetworkConfigModel struct {
	UPLinkIP        types.String `tfsdk:"uplink_ip"`
	TranslatedIP    types.String `tfsdk:"translated_ip"`
	EdgeGatewayName types.String `tfsdk:"edge_gateway_name"`
}

func (d *publicIPDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_public_ip"
}

func (d *publicIPDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The public IP data source show the list of public IP addresses.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"internal_ip": schema.StringAttribute{
				Description: "Internal IP address.",
				Computed:    true,
			},
			"network_config": schema.ListNestedAttribute{
				Description: "List of networks.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"uplink_ip": schema.StringAttribute{
							Description: "Uplink IP address.",
							Computed:    true,
						},
						"translated_ip": schema.StringAttribute{
							Description: "Translated IP address.",
							Computed:    true,
						},
						"edge_gateway_name": schema.StringAttribute{
							Description: "The name of the edge gateway related to the public ip. This properties is only present for NGP.",
							Computed:    true,
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

func (d *publicIPDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data publicIPDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	publicIP, _, err := d.client.PublicIPApi.ApiCustomersV20IpGet(d.client.auth)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read example, got error: %s", err))
		return
	}

	data = publicIPDataSourceModel{
		InternalIP: types.StringValue(publicIP.InternalIp),
	}

	listOfIps := make([]string, 0)

	for _, cfg := range publicIP.NetworkConfig {
		data.NetworkConfig = append(data.NetworkConfig, publicIPNetworkConfigModel{
			UPLinkIP:        types.StringValue(cfg.UplinkIp),
			TranslatedIP:    types.StringValue(cfg.TranslatedIp),
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
