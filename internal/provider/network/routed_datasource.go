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
	client *client.CloudAvenue
	org    org.Org
}

type networkRoutedDataSourceModel struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	EdgeGatewayID types.String `tfsdk:"edge_gateway_id"`
	InterfaceType types.String `tfsdk:"interface_type"`
	Gateway       types.String `tfsdk:"gateway"`
	PrefixLength  types.Int64  `tfsdk:"prefix_length"`
	DNS1          types.String `tfsdk:"dns1"`
	DNS2          types.String `tfsdk:"dns2"`
	DNSSuffix     types.String `tfsdk:"dns_suffix"`
	StaticIPPool  types.Set    `tfsdk:"static_ip_pool"`
}

func (d *networkRoutedDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_routed"
}

func (d *networkRoutedDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = network.GetSchema(network.SetRouted()).GetDataSource(ctx)
}

// Init resource used to initialize the resource.
func (d *networkRoutedDataSource) Init(_ context.Context, rm *networkRoutedDataSourceModel) (diags diag.Diagnostics) {
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
	var data *networkRoutedDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Init resource
	resp.Diagnostics.Append(d.Init(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get VDC
	vdc, err := d.client.GetVDC()
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving VDC", err.Error())
		return
	}

	// Get Network from Parent VDC or Edge Gateway
	var orgNetwork *govcd.OpenApiOrgVdcNetwork
	if data.EdgeGatewayID.IsNull() { // Get Network from default VDC
		orgNetwork, err = d.org.GetOpenApiOrgVdcNetworkByNameAndOwnerId(data.Name.ValueString(), vdc.GetID())
		if err != nil {
			resp.Diagnostics.AddError("Error retrieving Network of VDC", err.Error())
			return
		}
	} else { // Get Network from Parent Edge Gateway
		// Get Edge Gateway
		egw, err := d.org.GetEdgeGateway(edgegw.BaseEdgeGW{ID: data.EdgeGatewayID})
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
	}

	plan := &networkRoutedDataSourceModel{
		ID:            types.StringValue(orgNetwork.OpenApiOrgVdcNetwork.ID),
		Name:          types.StringValue(data.Name.ValueString()),
		Description:   types.StringValue(orgNetwork.OpenApiOrgVdcNetwork.Description),
		EdgeGatewayID: types.StringValue(orgNetwork.OpenApiOrgVdcNetwork.Connection.RouterRef.ID),
		InterfaceType: types.StringValue(orgNetwork.OpenApiOrgVdcNetwork.Connection.ConnectionType),
		Gateway:       types.StringValue(orgNetwork.OpenApiOrgVdcNetwork.Subnets.Values[0].Gateway),
		PrefixLength:  types.Int64Value(int64(orgNetwork.OpenApiOrgVdcNetwork.Subnets.Values[0].PrefixLength)),
		DNS1:          types.StringValue(orgNetwork.OpenApiOrgVdcNetwork.Subnets.Values[0].DNSServer1),
		DNS2:          types.StringValue(orgNetwork.OpenApiOrgVdcNetwork.Subnets.Values[0].DNSServer2),
		DNSSuffix:     types.StringValue(orgNetwork.OpenApiOrgVdcNetwork.Subnets.Values[0].DNSSuffix),
	}

	var diags diag.Diagnostics
	plan.StaticIPPool, diags = types.SetValueFrom(ctx, types.ObjectType{AttrTypes: staticIPPoolAttrTypes}, GetIPRanges(orgNetwork))
	resp.Diagnostics.Append(diags...)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}
