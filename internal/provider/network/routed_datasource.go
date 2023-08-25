// Package network provides a Terraform datasource.
package network

import (
	"context"
	"fmt"

	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/edgegw"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/network"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/org"
)

var (
	_ datasource.DataSource              = &networkRoutedDataSource{}
	_ datasource.DataSourceWithConfigure = &networkRoutedDataSource{}
)

func NewNetworkRoutedDataSource() datasource.DataSource {
	return &networkRoutedDataSource{}
}

type networkRoutedDataSource struct {
	client  *client.CloudAvenue
	org     org.Org
	network network.Kind
}

func (d *networkRoutedDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_routed"
}

func (d *networkRoutedDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = network.GetSchema(network.SetRouted()).GetDataSource(ctx)
}

// Init resource used to initialize the resource.
func (d *networkRoutedDataSource) Init(_ context.Context, rm *networkRoutedModel) (diags diag.Diagnostics) {
	// Init Network
	d.network.TypeOfNetwork = network.NAT_ROUTED
	// Init Org
	d.org, diags = org.Init(d.client)
	return
}

func (d *networkRoutedDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *networkRoutedDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	defer metrics.New("data.cloudavenue_network_routed", d.client.GetOrgName(), metrics.Read)()

	var data networkRoutedModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Init resource
	resp.Diagnostics.Append(d.Init(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get Network from Parent VDC or Edge Gateway
	var orgNetwork *govcd.OpenApiOrgVdcNetwork
	// Get Edge Gateway
	egw, err := d.org.GetEdgeGateway(edgegw.BaseEdgeGW{
		Name: data.EdgeGatewayName,
		ID:   data.EdgeGatewayID,
	})
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving Edge Gateway", err.Error())
		return
	}
	// Get Parent Edge Gateway
	parent, err := egw.GetParent()
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving parent Edge Gateway", err.Error())
		return
	}
	orgNetwork, err = d.org.GetOpenApiOrgVdcNetworkByNameAndOwnerId(data.Name.ValueString(), parent.GetID())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving Network of parent Edge Gateway", err.Error())
		return
	}

	// Set data into the model
	data = SetDataToNetworkRoutedModel(orgNetwork)

	// Set static IP pool
	var diags diag.Diagnostics
	data.StaticIPPool, diags = types.SetValueFrom(ctx, types.ObjectType{AttrTypes: staticIPPoolAttrTypes}, GetIPRanges(orgNetwork))
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
