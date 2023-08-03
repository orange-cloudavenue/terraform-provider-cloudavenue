// Package edgegw provides a Terraform datasource.
package edgegw

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
)

var (
	_ datasource.DataSource              = &vpnIpsecDataSource{}
	_ datasource.DataSourceWithConfigure = &vpnIpsecDataSource{}
)

func NewVpnIpsecDataSource() datasource.DataSource {
	return &vpnIpsecDataSource{}
}

type vpnIpsecDataSource struct {
	client *client.CloudAvenue

	// Uncomment the following lines if you need to access the resource's.
	// org    org.Org
	// vdc    vdc.VDC
	// vapp   vapp.VAPP
}

// If the data source don't have same schema/structure as the resource, you can use the following code:
// type vpnIpsecDataSourceModel struct {
// 	ID types.String `tfsdk:"id"`
// }

// Init Initializes the data source.
func (d *vpnIpsecDataSource) Init(ctx context.Context, dm *vpnIpsecModel) (diags diag.Diagnostics) {
	
	// Uncomment the following lines if you need to access to the Org
	// d.org, diags = org.Init(d.client)
	// if diags.HasError() {
	// 	return
	// }

	// Uncomment the following lines if you need to access to the VDC
	// d.vdc, diags = vdc.Init(d.client, dm.VDC)
	// if diags.HasError() {
	// 	return
	// }

	// Uncomment the following lines if you need to access to the VAPP
	// d.vapp, diags = vapp.Init(d.client, d.vdc, dm.VAppID, dm.VAppName)

	return
}

func (d *vpnIpsecDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName  + "_vpn_ipsec"
}

func (d *vpnIpsecDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = vpnIpsecSchema(ctx).GetDataSource(ctx)
}

func (d *vpnIpsecDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *vpnIpsecDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	config := &vpnIpsecModel{}

	// If the data source don't have same schema/structure as the resource, you can use the following code:
	// config := &vpnIpsecDataSourceModel{}

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Init the resource
	resp.Diagnostics.Append(d.Init(ctx, config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	/*
		Implement the data source read logic here.
	*/

	// If read function is identical to the resource, you can use the following code:
	/* 
	s := &vpnIpsecResource{
		client: d.client,
		// org:    d.org,
		// vdc:    d.vdc,
		// vapp:   d.vapp,
	}

	// Read data from the API
	data, _, diags := s.read(ctx, config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	*/

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
