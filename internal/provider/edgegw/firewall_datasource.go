// Package edgegw provides a Terraform datasource.
package edgegw

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/edgegw"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/org"
)

var (
	_ datasource.DataSource              = &firewallDataSource{}
	_ datasource.DataSourceWithConfigure = &firewallDataSource{}
)

func NewFirewallDataSource() datasource.DataSource {
	return &firewallDataSource{}
}

type firewallDataSource struct {
	client *client.CloudAvenue
	org    org.Org
	edgegw edgegw.EdgeGateway
}

// Init Initializes the data source.
func (d *firewallDataSource) Init(ctx context.Context, dm *firewallModel) (diags diag.Diagnostics) {
	var err error

	d.org, diags = org.Init(d.client)
	if diags.HasError() {
		return
	}

	d.edgegw, err = d.org.GetEdgeGateway(edgegw.BaseEdgeGW{
		ID:   dm.EdgeGatewayID,
		Name: dm.EdgeGatewayName,
	})
	if err != nil {
		diags.AddError("Error retrieving Edge Gateway", err.Error())
		return
	}

	return
}

func (d *firewallDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_firewall"
}

func (d *firewallDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = firewallSchema(ctx).GetDataSource(ctx)
}

func (d *firewallDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *firewallDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	config := &firewallModel{}

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

	f := &firewallResource{
		client: d.client,
		org:    d.org,
		edgegw: d.edgegw,
	}

	data, diag := f.read(ctx)
	resp.Diagnostics.Append(diag...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
