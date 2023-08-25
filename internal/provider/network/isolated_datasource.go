// Package network provides a Terraform datasource.
package network

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/network"
)

var (
	_ datasource.DataSource              = &networkIsolatedDataSource{}
	_ datasource.DataSourceWithConfigure = &networkIsolatedDataSource{}
)

func NewNetworkIsolatedDataSource() datasource.DataSource {
	return &networkIsolatedDataSource{}
}

type networkIsolatedDataSource struct {
	client *client.CloudAvenue
}

func (d *networkIsolatedDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_isolated"
}

func (d *networkIsolatedDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = network.GetSchema(network.SetIsolated()).GetDataSource(ctx)
}

func (d *networkIsolatedDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *networkIsolatedDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	defer metrics.New("data.cloudavenue_network_isolated", d.client.GetOrgName(), metrics.Read)()

	var data networkIsolatedModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Define VDC or VDCGroup
	vdcOrVDCGroup, err := d.client.GetVDCOrVDCGroup(data.VDC.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving VDC or VDCGroup", err.Error())
		return
	}

	// Get network
	network, err := vdcOrVDCGroup.GetOpenApiOrgVdcNetworkByName(data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("[READ] Error retrieving isolated network", err.Error())
		return
	}

	// Get network static IP pools
	ipPools := []staticIPPool{}
	if len(network.OpenApiOrgVdcNetwork.Subnets.Values[0].IPRanges.Values) > 0 {
		for _, ipRange := range network.OpenApiOrgVdcNetwork.Subnets.Values[0].IPRanges.Values {
			ipPools = append(ipPools, staticIPPool{
				StartAddress: types.StringValue(ipRange.StartAddress),
				EndAddress:   types.StringValue(ipRange.EndAddress),
			})
		}
	}

	// Set Plan updated
	data = networkIsolatedModel{
		ID:           types.StringValue(network.OpenApiOrgVdcNetwork.ID),
		Name:         types.StringValue(network.OpenApiOrgVdcNetwork.Name),
		Description:  types.StringValue(network.OpenApiOrgVdcNetwork.Description),
		VDC:          types.StringValue(vdcOrVDCGroup.GetName()),
		Gateway:      types.StringValue(network.OpenApiOrgVdcNetwork.Subnets.Values[0].Gateway),
		PrefixLength: types.Int64Value(int64(network.OpenApiOrgVdcNetwork.Subnets.Values[0].PrefixLength)),
		DNS1:         types.StringValue(network.OpenApiOrgVdcNetwork.Subnets.Values[0].DNSServer1),
		DNS2:         types.StringValue(network.OpenApiOrgVdcNetwork.Subnets.Values[0].DNSServer2),
		DNSSuffix:    types.StringValue(network.OpenApiOrgVdcNetwork.Subnets.Values[0].DNSSuffix),
	}

	// Set static IP pools
	var diags diag.Diagnostics
	data.StaticIPPool, diags = types.SetValueFrom(ctx, types.ObjectType{AttrTypes: staticIPPoolAttrTypes}, ipPools)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
