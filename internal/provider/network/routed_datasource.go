// Package network provides a Terraform datasource.
package network

import (
	"context"
	"fmt"

	fstringvalidator "github.com/FrangipaneTeam/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
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

type staticIPPoolDataSourceModel struct {
	StartAddress types.String `tfsdk:"start_address"`
	EndAddress   types.String `tfsdk:"end_address"`
}

func (d *networkRoutedDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + "network_routed"
}

func (d *networkRoutedDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The network_routed datasource allows you to manage a ...",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID of the routed network.",

				//	PlanModifiers: []planmodifier.String{
				//		stringplanmodifier.UseStateForUnknown(),
				//	},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Network name.",
			},
			"description": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Network description.",
			},
			"edge_gateway_id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Edge gateway ID in which Routed network should be located.",
			},
			"interface_type": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Optional interface type (only for NSX-V networks). One of `INTERNAL` (default), `DISTRIBUTED`, `SUBINTERFACE`",
				Validators: []validator.String{
					stringvalidator.OneOf("INTERNAL", "SUBINTERFACE", "DISTRIBUTED"),
				},
			},
			"gateway": schema.StringAttribute{
				Computed: true,
				Validators: []validator.String{
					fstringvalidator.IsValidIP(),
				},
				MarkdownDescription: "Gateway IP address.",
			},
			"prefix_length": schema.Int64Attribute{
				Computed: true,
				Validators: []validator.Int64{
					int64validator.Between(1, 32),
				},
				MarkdownDescription: "Network prefix length.",
			},
			"dns1": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "DNS server 1.",
			},
			"dns2": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "DNS server 2.",
			},
			"dns_suffix": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "DNS suffix.",
			},
			"static_ip_pool": schema.SetNestedAttribute{
				Computed:            true,
				MarkdownDescription: "IP ranges used for static pool allocation in the network.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"start_address": schema.StringAttribute{
							Computed: true,
							Validators: []validator.String{
								fstringvalidator.IsValidIP(),
							},
							MarkdownDescription: "Start address of the IP range.",
						},
						"end_address": schema.StringAttribute{
							Computed: true,
							Validators: []validator.String{
								fstringvalidator.IsValidIP(),
							},
							MarkdownDescription: "End address of the IP range.",
						},
					},
				},
			},
		},
	}
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

	// Get Org
	org, err := d.client.GetOrg()
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving ORG", err.Error())
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
		orgNetwork, err = org.GetOpenApiOrgVdcNetworkByNameAndOwnerId(data.Name.ValueString(), vdc.GetID())
	} else {
		parentID, _ := getParentEdgeGatewayID(org.Org, data.EdgeGatewayID.ValueString())
		orgNetwork, err = org.GetOpenApiOrgVdcNetworkByNameAndOwnerId(data.Name.ValueString(), *parentID)
	}

	// Define VDC or VDCGroup
	// vdcOrVDCGroup, err := d.client.GetVDCOrVDCGroup(parentName)
	// if err != nil {
	//	resp.Diagnostics.AddError("Error retrieving VDC or VDCGroup", err.Error())
	//	return
	//}

	if err != nil {
		if govcd.ContainsNotFound(err) {
			tflog.Debug(ctx, "Network not found, removing resource from state")
			resp.State.RemoveResource(ctx)
		}
		resp.Diagnostics.AddError("Error retrieving routing network", err.Error())
		return
	}

	plan := &networkRoutedResourceModel{
		ID:            types.StringValue(orgNetwork.OpenApiOrgVdcNetwork.ID),
		Name:          types.StringValue(orgNetwork.OpenApiOrgVdcNetwork.Name),
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
	if resp.Diagnostics.HasError() {
		return
	}
}
