// Package vapp provides a Terraform datasource.
package vapp

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/org"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vapp"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vdc"
)

var (
	_ datasource.DataSource              = &orgNetworkDataSource{}
	_ datasource.DataSourceWithConfigure = &orgNetworkDataSource{}
)

func NewOrgNetworkDataSource() datasource.DataSource {
	return &orgNetworkDataSource{}
}

type orgNetworkDataSource struct {
	client *client.CloudAvenue

	// Uncomment the following lines if you need to access the resource's.
	org  org.Org
	vdc  vdc.VDC
	vapp vapp.VAPP
}

type orgNetworkDataSourceModel struct {
	ID                 types.String `tfsdk:"id"`
	VAppName           types.String `tfsdk:"vapp_name"`
	VAppID             types.String `tfsdk:"vapp_id"`
	VDC                types.String `tfsdk:"vdc"`
	NetworkName        types.String `tfsdk:"network_name"`
	IsFenced           types.Bool   `tfsdk:"is_fenced"`
	RetainIPMacEnabled types.Bool   `tfsdk:"retain_ip_mac_enabled"`
}

// Init Initializes the data source.
func (d *orgNetworkDataSource) Init(ctx context.Context, dm *orgNetworkDataSourceModel) (diags diag.Diagnostics) {

	// Uncomment the following lines if you need to access to the Org
	d.org, diags = org.Init(d.client)
	if diags.HasError() {
		return
	}

	// Uncomment the following lines if you need to access to the VDC
	d.vdc, diags = vdc.Init(d.client, dm.VDC)
	if diags.HasError() {
		return
	}

	// Uncomment the following lines if you need to access to the VAPP
	d.vapp, diags = vapp.Init(d.client, d.vdc, dm.VAppID, dm.VAppName)

	return
}

func (d *orgNetworkDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_"
}

func (d *orgNetworkDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The  datasource allows you to read a ...",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (d *orgNetworkDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *orgNetworkDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	config := &orgNetworkDataSourceModel{}

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

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
