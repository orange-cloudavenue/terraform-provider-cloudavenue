// Package network provides a Terraform datasource.
package network

import (
	"context"
	"fmt"

	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
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

	var orgNetwork *govcd.OpenApiOrgVdcNetwork
	if data.EdgeGatewayID.IsNull() {
		orgNetwork, err = d.org.GetOpenApiOrgVdcNetworkByNameAndOwnerId(data.Name.ValueString(), vdc.GetID())
	} else {
		parentID, diag := GetParentEdgeGatewayID(d.org, data.EdgeGatewayID.ValueString())
		if diag != nil {
			resp.Diagnostics.Append(diag)
			return
		}
		orgNetwork, err = d.org.GetOpenApiOrgVdcNetworkByNameAndOwnerId(data.Name.ValueString(), *parentID)
	}

	if err != nil {
		if govcd.ContainsNotFound(err) {
			tflog.Debug(ctx, "Network not found, removing resource from state")
			resp.State.RemoveResource(ctx)
		}
		resp.Diagnostics.AddError("Error retrieving routing network", err.Error())
		return
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

	ipPools := []staticIPPool{}

	if len(orgNetwork.OpenApiOrgVdcNetwork.Subnets.Values[0].IPRanges.Values) > 0 {
		for _, ipRange := range orgNetwork.OpenApiOrgVdcNetwork.Subnets.Values[0].IPRanges.Values {
			ipPool := staticIPPool{
				StartAddress: types.StringValue(ipRange.StartAddress),
				EndAddress:   types.StringValue(ipRange.EndAddress),
			}
			ipPools = append(ipPools, ipPool)
		}
	}
	var diags diag.Diagnostics
	plan.StaticIPPool, diags = types.SetValueFrom(ctx, types.ObjectType{AttrTypes: staticIPPoolAttrTypes}, ipPools)
	resp.Diagnostics.Append(diags...)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}
