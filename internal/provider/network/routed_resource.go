// Package network provides a Terraform resource.
package network

import (
	"context"
	"fmt"
	"strings"

	fstringvalidator "github.com/FrangipaneTeam/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/stringpm"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &networkRoutedResource{}
	_ resource.ResourceWithConfigure   = &networkRoutedResource{}
	_ resource.ResourceWithImportState = &networkRoutedResource{}
)

// NewNetworkRoutedResource is a helper function to simplify the provider implementation.
func NewNetworkRoutedResource() resource.Resource {
	return &networkRoutedResource{}
}

// networkRoutedResource is the resource implementation.
type networkRoutedResource struct {
	client *client.CloudAvenue
}

type networkRoutedResourceModel struct {
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

type staticIPPool struct {
	StartAddress types.String `tfsdk:"start_address"`
	EndAddress   types.String `tfsdk:"end_address"`
}

var staticIPPoolAttrTypes = map[string]attr.Type{
	"start_address": types.StringType,
	"end_address":   types.StringType,
}

// Metadata returns the resource type name.
func (r *networkRoutedResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_" + "routed"
}

// Schema defines the schema for the resource.
func (r *networkRoutedResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Provides a CloudAvenue Org VDC routed Network. This can be used to create, modify, and delete routed VDC networks.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID of the routed network.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Network name.",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Network description.",
			},
			"edge_gateway_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				MarkdownDescription: "Edge gateway ID in which Routed network should be located.",
			},
			"interface_type": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Optional interface type (only for NSX-V networks). One of `INTERNAL` (default), `DISTRIBUTED`, `SUBINTERFACE`",
				PlanModifiers: []planmodifier.String{
					stringpm.SetDefault("INTERNAL"),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("INTERNAL", "SUBINTERFACE", "DISTRIBUTED"),
				},
			},
			"gateway": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					fstringvalidator.IsValidIP(),
				},
				MarkdownDescription: "Gateway IP address.",
			},
			"prefix_length": schema.Int64Attribute{
				Required: true,
				Validators: []validator.Int64{
					int64validator.Between(1, 32),
				},
				MarkdownDescription: "Network prefix length.",
			},
			"dns1": schema.StringAttribute{
				Optional:            true,
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
				Optional:            true,
				MarkdownDescription: "IP ranges used for static pool allocation in the network.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"start_address": schema.StringAttribute{
							Required: true,
							Validators: []validator.String{
								fstringvalidator.IsValidIP(),
							},
							MarkdownDescription: "Start address of the IP range.",
						},
						"end_address": schema.StringAttribute{
							Required: true,
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

func (r *networkRoutedResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.CloudAvenue)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.CloudAvenue, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

// Create creates the resource and sets the initial Terraform state.
func (r *networkRoutedResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var (
		plan *networkRoutedResourceModel
	)

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	org, _, err := r.client.GetOrgAndVDC(r.client.GetOrg(), r.client.GetDefaultVDC())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving VDC", err.Error())
		return
	}

	parentEdgeGatewayOwnerID, errGet := getParentEdgeGatewayID(org, plan.EdgeGatewayID.ValueString())
	if errGet != nil {
		resp.Diagnostics.AddError(errGet.Summary, errGet.Detail)
		return
	}

	if parentEdgeGatewayOwnerID == nil {
		resp.Diagnostics.AddError("Error retrieving Edge Gateway", "parentEdgeGatewayOwnerID is nil")
		return
	}

	if govcd.OwnerIsVdcGroup(*parentEdgeGatewayOwnerID) {
		networkMutexKV.KvLock(ctx, *parentEdgeGatewayOwnerID)
		defer networkMutexKV.KvUnlock(ctx, *parentEdgeGatewayOwnerID)
	} else {
		networkMutexKV.KvLock(ctx, plan.EdgeGatewayID.ValueString())
		defer networkMutexKV.KvUnlock(ctx, plan.EdgeGatewayID.ValueString())
	}

	ipPool := []staticIPPool{}
	resp.Diagnostics.Append(plan.StaticIPPool.ElementsAs(ctx, &ipPool, true)...)

	orgVDCNetworkConfig := &govcdtypes.OpenApiOrgVdcNetwork{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		OwnerRef:    &govcdtypes.OpenApiReference{ID: *parentEdgeGatewayOwnerID},

		NetworkType: govcdtypes.OrgVdcNetworkTypeRouted,

		// Connection is used for "routed" network
		Connection: &govcdtypes.Connection{
			RouterRef: govcdtypes.OpenApiReference{
				ID: plan.EdgeGatewayID.ValueString(),
			},
			// API requires interface type in upper case, but we accept any case
			ConnectionType: plan.InterfaceType.ValueString(),
		},
		Subnets: govcdtypes.OrgVdcNetworkSubnets{
			Values: []govcdtypes.OrgVdcNetworkSubnetValues{
				{
					Gateway:      plan.Gateway.ValueString(),
					PrefixLength: int(plan.PrefixLength.ValueInt64()),
					DNSServer1:   plan.DNS1.ValueString(),
					DNSServer2:   plan.DNS2.ValueString(),
					DNSSuffix:    plan.DNSSuffix.ValueString(),
					IPRanges: govcdtypes.OrgVdcNetworkSubnetIPRanges{
						Values: processIPRanges(ipPool),
					},
				},
			},
		},
	}

	orgNetwork, err := org.CreateOpenApiOrgVdcNetwork(orgVDCNetworkConfig)
	if err != nil {
		resp.Diagnostics.AddError("Error creating routing network", err.Error())
		return
	}

	plan.ID = types.StringValue(orgNetwork.OpenApiOrgVdcNetwork.ID)

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *networkRoutedResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *networkRoutedResourceModel
	// Get current state
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	org, _, err := r.client.GetOrgAndVDC(r.client.GetOrg(), r.client.GetDefaultVDC())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving VDC", err.Error())
		return
	}

	orgNetwork, err := org.GetOpenApiOrgVdcNetworkById(state.ID.ValueString())
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

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *networkRoutedResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan *networkRoutedResourceModel

	// Get current state
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	org, _, err := r.client.GetOrgAndVDC(r.client.GetOrg(), r.client.GetDefaultVDC())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving VDC", err.Error())
		return
	}

	parentEdgeGatewayOwnerID, errGet := getParentEdgeGatewayID(org, plan.EdgeGatewayID.ValueString())
	if errGet != nil {
		resp.Diagnostics.AddError(errGet.Summary, errGet.Detail)
		return
	}

	if govcd.OwnerIsVdcGroup(*parentEdgeGatewayOwnerID) {
		networkMutexKV.KvLock(ctx, *parentEdgeGatewayOwnerID)
		defer networkMutexKV.KvUnlock(ctx, *parentEdgeGatewayOwnerID)
	} else {
		networkMutexKV.KvLock(ctx, plan.EdgeGatewayID.ValueString())
		defer networkMutexKV.KvUnlock(ctx, plan.EdgeGatewayID.ValueString())
	}

	orgNetwork, err := org.GetOpenApiOrgVdcNetworkById(plan.ID.ValueString())
	if err != nil {
		if govcd.ContainsNotFound(err) {
			tflog.Debug(ctx, "Network not found, removing resource from state")
			resp.State.RemoveResource(ctx)
		}
		resp.Diagnostics.AddError("Error retrieving routing network", err.Error())
		return
	}

	ipPool := []staticIPPool{}
	diags := plan.StaticIPPool.ElementsAs(ctx, &ipPool, true)
	resp.Diagnostics.Append(diags...)

	newOrgNetwork := &govcdtypes.OpenApiOrgVdcNetwork{
		ID:          orgNetwork.OpenApiOrgVdcNetwork.ID,
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		OwnerRef:    &govcdtypes.OpenApiReference{ID: *parentEdgeGatewayOwnerID},

		NetworkType: govcdtypes.OrgVdcNetworkTypeRouted,

		// Connection is used for "routed" network
		Connection: &govcdtypes.Connection{
			RouterRef: govcdtypes.OpenApiReference{
				ID: plan.EdgeGatewayID.ValueString(),
			},
			// API requires interface type in upper case, but we accept any case
			ConnectionType: plan.InterfaceType.ValueString(),
		},
		Subnets: govcdtypes.OrgVdcNetworkSubnets{
			Values: []govcdtypes.OrgVdcNetworkSubnetValues{
				{
					Gateway:      plan.Gateway.ValueString(),
					PrefixLength: int(plan.PrefixLength.ValueInt64()),
					DNSServer1:   plan.DNS1.ValueString(),
					DNSServer2:   plan.DNS2.ValueString(),
					DNSSuffix:    plan.DNSSuffix.ValueString(),
					IPRanges: govcdtypes.OrgVdcNetworkSubnetIPRanges{
						Values: processIPRanges(ipPool),
					},
				},
			},
		},
	}
	_, err = orgNetwork.Update(newOrgNetwork)
	if err != nil {
		resp.Diagnostics.AddError("Error updating routing network", err.Error())
		return
	}

	plan = &networkRoutedResourceModel{
		ID:            types.StringValue(orgNetwork.OpenApiOrgVdcNetwork.ID),
		Name:          plan.Name,
		Description:   plan.Description,
		EdgeGatewayID: plan.EdgeGatewayID,
		InterfaceType: plan.InterfaceType,
		Gateway:       plan.Gateway,
		PrefixLength:  plan.PrefixLength,
		DNS1:          plan.DNS1,
		DNS2:          plan.DNS2,
		DNSSuffix:     plan.DNSSuffix,
	}

	plan.StaticIPPool, diags = types.SetValueFrom(ctx, types.ObjectType{AttrTypes: staticIPPoolAttrTypes}, ipPool)
	resp.Diagnostics.Append(diags...)

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *networkRoutedResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *networkRoutedResourceModel

	// Get current state
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	org, _, err := r.client.GetOrgAndVDC(r.client.GetOrg(), r.client.GetDefaultVDC())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving VDC", err.Error())
		return
	}

	parentEdgeGatewayOwnerID, errGet := getParentEdgeGatewayID(org, state.EdgeGatewayID.ValueString())
	if errGet != nil {
		resp.Diagnostics.AddError(errGet.Summary, errGet.Detail)
		return
	}

	if govcd.OwnerIsVdcGroup(*parentEdgeGatewayOwnerID) {
		networkMutexKV.KvLock(ctx, *parentEdgeGatewayOwnerID)
		defer networkMutexKV.KvUnlock(ctx, *parentEdgeGatewayOwnerID)
	} else {
		networkMutexKV.KvLock(ctx, state.EdgeGatewayID.ValueString())
		defer networkMutexKV.KvUnlock(ctx, state.EdgeGatewayID.ValueString())
	}

	orgNetwork, err := org.GetOpenApiOrgVdcNetworkById(state.ID.ValueString())
	if err != nil {
		if govcd.ContainsNotFound(err) {
			tflog.Debug(ctx, "Network not found, removing resource from state")
			resp.State.RemoveResource(ctx)
		}
		resp.Diagnostics.AddError("Error retrieving routing network", err.Error())
		return
	}
	err = orgNetwork.Delete()
	if err != nil {
		resp.Diagnostics.AddError("Error deleting routing network", err.Error())
		return
	}
}

func (r *networkRoutedResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceURI := strings.Split(req.ID, ".")

	if len(resourceURI) != 2 {
		resp.Diagnostics.AddError("Error importing network_routed", "Resource name must be specified as vdc-name.network-name or vdc-group-name.network-name")
		return
	}

	vdcOrVDCGroupName, networkName := resourceURI[0], resourceURI[1]

	org, vdc, err := r.client.GetOrgAndVDC(r.client.GetOrg(), r.client.GetDefaultVDC())
	if err != nil && govcd.ContainsNotFound(err) {
		resp.Diagnostics.AddError("Error retrieving VDC", err.Error())
		return
	}

	orgNetwork, err := vdc.GetOpenApiOrgVdcNetworkByName(networkName)
	if err != nil && !govcd.ContainsNotFound(err) {
		resp.Diagnostics.AddError("Error retrieving org vdc network by name", err.Error())
		return
	}

	if orgNetwork == nil {
		adminOrg, err := r.client.Vmware.GetAdminOrgByName(org.Org.Name)
		if err != nil {
			resp.Diagnostics.AddError("Error retrieving Admin Org", err.Error())
			return
		}
		vdcgroup, err := adminOrg.GetVdcGroupByName(vdcOrVDCGroupName)
		if err != nil {
			resp.Diagnostics.AddError("Error retrieving VDC group", err.Error())
			return
		}
		orgNetwork, err = vdcgroup.GetOpenApiOrgVdcNetworkByName(networkName)
		if err != nil && govcd.ContainsNotFound(err) {
			resp.Diagnostics.AddError("Error retrieving org vdc_group network by name", err.Error())
			return
		}
	}

	if !orgNetwork.IsRouted() {
		resp.Diagnostics.AddError("Error importing routed network", fmt.Sprintf("Org network with name '%s' found, but is not of type Routed (type is '%s')", networkName, orgNetwork.GetType()))
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), orgNetwork.OpenApiOrgVdcNetwork.ID)...)
}
