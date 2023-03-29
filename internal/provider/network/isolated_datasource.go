// Package network provides a Terraform datasource.
package network

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
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

type networkIsolatedDataSourceModel struct {
	ID           types.String `tfsdk:"id"`
	VDC          types.String `tfsdk:"vdc"`
	Name         types.String `tfsdk:"name"`
	Description  types.String `tfsdk:"description"`
	Gateway      types.String `tfsdk:"gateway"`
	PrefixLength types.Int64  `tfsdk:"prefix_length"`
	PrimaryDNS   types.String `tfsdk:"dns1"`
	SecondaryDNS types.String `tfsdk:"dns2"`
	SuffixDNS    types.String `tfsdk:"dns_suffix"`
	StaticIPPool types.Set    `tfsdk:"static_ip_pool"`
}

type staticIPPoolDataSourceModel struct {
	StartAddress types.String `tfsdk:"start_address"`
	EndAddress   types.String `tfsdk:"end_address"`
}

func (d *networkIsolatedDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_isolated"
}

func (d *networkIsolatedDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The network_isolated datasource allows you to manage a ...",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"vdc": schema.StringAttribute{
				Computed:            true,
				Optional:            true,
				MarkdownDescription: "The name of VDC to use, optional if defined at provider level",
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "A unique name for this network",
			},
			"description": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "A description of the network.",
			},
			"gateway": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "(Force replacement) The gateway IP address for the network. This value define also the network IP range with the prefix length.",
			},
			"prefix_length": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "(Force replacement) The prefix length for the network. This value must be a valid prefix length for the network IP range.(e.g. 24 for netmask 255.255.255.0)",
			},
			"dns1": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The primary DNS server IP address for the network.",
			},
			"dns2": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The secondary DNS server IP address for the network.",
			},
			"dns_suffix": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The DNS suffix for the network.",
			},
			"static_ip_pool": schema.SetNestedAttribute{
				Computed:            true,
				MarkdownDescription: "A set of static IP pools to be used for this network.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"start_address": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The start address of the IP pool. This value must be a valid IP address in the network IP range.",
						},
						"end_address": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The end address of the IP pool. This value must be a valid IP address in the network IP range.",
						},
					},
				},
			},
		},
	}
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
	var data networkIsolatedDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the organization
	// if _, err := d.client.GetOrg(); err != nil {
	//	resp.Diagnostics.AddError("Unable to get organization", err.Error())
	//	return
	//}

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
	ipPools := []staticIPPoolDataSourceModel{}
	if len(network.OpenApiOrgVdcNetwork.Subnets.Values[0].IPRanges.Values) > 0 {
		for _, ipRange := range network.OpenApiOrgVdcNetwork.Subnets.Values[0].IPRanges.Values {
			ipPools = append(ipPools, staticIPPoolDataSourceModel{
				StartAddress: types.StringValue(ipRange.StartAddress),
				EndAddress:   types.StringValue(ipRange.EndAddress),
			})
		}
	}

	// Set Plan updated
	data = networkIsolatedDataSourceModel{
		ID:           types.StringValue(network.OpenApiOrgVdcNetwork.ID),
		Name:         types.StringValue(network.OpenApiOrgVdcNetwork.Name),
		Description:  types.StringValue(network.OpenApiOrgVdcNetwork.Description),
		VDC:          types.StringValue(vdcOrVDCGroup.GetName()),
		Gateway:      types.StringValue(network.OpenApiOrgVdcNetwork.Subnets.Values[0].Gateway),
		PrefixLength: types.Int64Value(int64(network.OpenApiOrgVdcNetwork.Subnets.Values[0].PrefixLength)),
		PrimaryDNS:   types.StringValue(network.OpenApiOrgVdcNetwork.Subnets.Values[0].DNSServer1),
		SecondaryDNS: types.StringValue(network.OpenApiOrgVdcNetwork.Subnets.Values[0].DNSServer2),
		SuffixDNS:    types.StringValue(network.OpenApiOrgVdcNetwork.Subnets.Values[0].DNSSuffix),
	}

	// Set static IP pools
	var diags diag.Diagnostics
	data.StaticIPPool, diags = types.SetValueFrom(ctx, types.ObjectType{AttrTypes: staticIPPoolResourceModelAttrTypes}, ipPools)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
