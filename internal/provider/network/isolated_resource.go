// Package network provides a Terraform resource.
package network

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	fstringvalidator "github.com/FrangipaneTeam/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vmware/go-vcloud-director/v2/govcd"
	govdctypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/mutex"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &networkIsolatedResource{}
	_ resource.ResourceWithConfigure   = &networkIsolatedResource{}
	_ resource.ResourceWithImportState = &networkIsolatedResource{}
)

// NewNetworkIsolatedResource is a helper function to simplify the provider implementation.
func NewNetworkIsolatedResource() resource.Resource {
	return &networkIsolatedResource{}
}

// networkIsolatedResource is the resource implementation.
type networkIsolatedResource struct {
	client *client.CloudAvenue
}

type networkIsolatedResourceModel struct {
	ID           types.String `tfsdk:"id"`
	OwnerID      types.String `tfsdk:"owner_id"`
	Name         types.String `tfsdk:"name"`
	Description  types.String `tfsdk:"description"`
	Gateway      types.String `tfsdk:"gateway"`
	PrefixLength types.Int64  `tfsdk:"prefix_length"`
	PrimaryDNS   types.String `tfsdk:"dns1"`
	SecondaryDNS types.String `tfsdk:"dns2"`
	SuffixDNS    types.String `tfsdk:"dns_suffix"`
	StaticIPPool types.Set    `tfsdk:"static_ip_pool"`
}

type staticIPPoolResourceModel struct {
	StartAddress types.String `tfsdk:"start_address"`
	EndAddress   types.String `tfsdk:"end_address"`
}

var staticIPPoolResourceModelAttrTypes = map[string]attr.Type{
	"start_address": types.StringType,
	"end_address":   types.StringType,
}

// Metadata returns the resource type name.
func (r *networkIsolatedResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_" + "isolated"
}

// Schema defines the schema for the resource.
func (r *networkIsolatedResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Provides a Cloud Avenue VDC isolated Network. This can be used to create, modify, and delete isolated VDC networks. An network isolated is not connected to any other network.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID of the network. This is a generated value and cannot be specified during creation. This value is used to identify the network in other resources.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"owner_id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "(ForceNew) The Uuid of the `VDC` or `VDC Group` that owns the network. If not specified, it use the vdc at provider level.",
				Validators: []validator.String{
					fstringvalidator.IsValidUUID(),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the network. This value must be unique within the `VDC` or `VDC Group` that owns the network.",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "A description of the network. This value is optional.",
			},
			"gateway": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The gateway IP address for the network. This value define also the network IP range with the prefix length.",
				Validators: []validator.String{
					fstringvalidator.IsValidIP(),
				},
			},
			"prefix_length": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "The prefix length for the network. This value must be a valid prefix length for the network IP range.(e.g. 24 for netmask 255.255.255.0)",
				Validators: []validator.Int64{
					int64validator.Between(1, 32),
				},
			},
			"dns1": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The primary DNS server IP address for the network.",
				Validators: []validator.String{
					fstringvalidator.IsValidIP(),
				},
			},
			"dns2": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The secondary DNS server IP address for the network.",
				Validators: []validator.String{
					fstringvalidator.IsValidIP(),
				},
			},
			"dns_suffix": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The DNS suffix for the network.",
			},
			"static_ip_pool": schema.SetNestedAttribute{
				Optional:            true,
				MarkdownDescription: "A set of static IP pools to be used for this network.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"start_address": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "The start address of the IP pool. This value must be a valid IP address in the network IP range.",
							Validators: []validator.String{
								fstringvalidator.IsValidIP(),
							},
						},
						"end_address": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "The end address of the IP pool. This value must be a valid IP address in the network IP range.",
							Validators: []validator.String{
								fstringvalidator.IsValidIP(),
							},
						},
					},
				},
			},
		},
	}
}

func (r *networkIsolatedResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *networkIsolatedResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var (
		plan     *networkIsolatedResourceModel
		err      error
		urn      string
		isGroup  bool
		vdcGroup *govcd.VdcGroup
		vdc      *govcd.Vdc
	)

	// Read the plan
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get Org and VDC
	org, err := r.client.GetOrg()
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving Org", err.Error())
		return
	}

	// If VDC is not defined at data source level, use the one defined at provider level
	if plan.OwnerID.IsNull() || plan.OwnerID.IsUnknown() {
		if r.client.DefaultVDCExist() {
			vdc, err = org.GetVDCByName(r.client.GetDefaultVDC(), true)
			if err != nil {
				resp.Diagnostics.AddError("Error retrieving VDC from Provider", err.Error())
				return
			}
			// Set the OwnerID with the VDC ID
			plan.OwnerID = types.StringValue(vdc.Vdc.ID)
		} else {
			resp.Diagnostics.AddError("[CREATE] Missing VDC", "OwnerID (VDC) is required when not defined at provider level")
			return
		}
	} else {
		// OwnerID is not null, so we need to retrieve the URN via Owner_ID (VDC or VDC Group)
		// First, we need to check if it is a VDC
		vdc, err = org.GetVDCById(plan.OwnerID.ValueString(), true)
		if err == nil {
			plan.OwnerID = types.StringValue(vdc.Vdc.ID)
		} else { // It is not a VDC, so it may be a VDC Group
			urn, err = govcd.BuildUrnWithUuid("urn:vcloud:vdcGroup:", plan.OwnerID.ValueString())
			if err == nil {
				vdcGroup, err = org.GetVdcGroupById(urn)
				if err != nil {
					resp.Diagnostics.AddError("[CREATE] Error retrieving VDC Group from OwnerID", err.Error())
					return
				}

				// OwnerId is VDCGroup ID
				isGroup = true
				plan.OwnerID = types.StringValue(urn)
			} else {
				resp.Diagnostics.AddError("[CREATE] Error retrieving VDC or VDC Group from OwnerID", err.Error())
				return
			}
		}
	}

	// Lock VDC or VDC Group to prevent concurrent access
	vcdMutexKV := mutex.NewKV()
	vcdMutexKV.KvLock(ctx, urn)
	defer vcdMutexKV.KvUnlock(ctx, urn)

	// Get network type
	ipPool := []staticIPPoolResourceModel{}
	resp.Diagnostics.Append(plan.StaticIPPool.ElementsAs(ctx, &ipPool, true)...)
	myshared := false // Cloudavenue does not support shared networks
	networkType := &govdctypes.OpenApiOrgVdcNetwork{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		Shared:      &myshared,
		NetworkType: govdctypes.OrgVdcNetworkTypeIsolated,
		Subnets: govdctypes.OrgVdcNetworkSubnets{
			Values: []govdctypes.OrgVdcNetworkSubnetValues{
				{
					Gateway:      plan.Gateway.ValueString(),
					PrefixLength: int(plan.PrefixLength.ValueInt64()),
					IPRanges: govdctypes.OrgVdcNetworkSubnetIPRanges{
						Values: myProcessIPRanges(ipPool),
					},
					DNSServer1: plan.PrimaryDNS.ValueString(),
					DNSServer2: plan.SecondaryDNS.ValueString(),
					DNSSuffix:  plan.SuffixDNS.ValueString(),
				},
			},
		},
	}

	// Set OwnerRef if VDC Group
	if isGroup {
		networkType.OwnerRef = &govdctypes.OpenApiReference{ID: vdcGroup.VdcGroup.Id}
	} else {
		networkType.OwnerRef = &govdctypes.OpenApiReference{ID: vdc.Vdc.ID}
	}

	// Create network
	orgNetwork, err := org.CreateOpenApiOrgVdcNetwork(networkType)
	if err != nil {
		resp.Diagnostics.AddError("[CREATE] Error creating isolated network", err.Error())
		return
	}

	// set Plan
	plan = &networkIsolatedResourceModel{
		ID:           types.StringValue(orgNetwork.OpenApiOrgVdcNetwork.ID),
		OwnerID:      types.StringValue(extractUUID(plan.OwnerID.ValueString())),
		Name:         plan.Name,
		Description:  plan.Description,
		Gateway:      plan.Gateway,
		PrefixLength: plan.PrefixLength,
		StaticIPPool: plan.StaticIPPool,
		PrimaryDNS:   plan.PrimaryDNS,
		SecondaryDNS: plan.SecondaryDNS,
		SuffixDNS:    plan.SuffixDNS,
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *networkIsolatedResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var (
		state *networkIsolatedResourceModel
		err   error
	)

	// Get current state
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get Org and VDC
	org, err := r.client.GetOrg()
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving Org", err.Error())
		return
	}

	// Build URN Network
	urn := state.ID.ValueString()
	// Get network
	orgNetwork, err := org.GetOpenApiOrgVdcNetworkById(urn)
	if err != nil {
		if govcd.ContainsNotFound(err) {
			// Network not found, so remove from state
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("[READ] Error retrieving isolated network", err.Error())
		return
	}

	plan := &networkIsolatedResourceModel{
		ID:           types.StringValue(orgNetwork.OpenApiOrgVdcNetwork.ID),
		OwnerID:      types.StringValue(extractUUID(orgNetwork.OpenApiOrgVdcNetwork.OwnerRef.ID)),
		Name:         types.StringValue(orgNetwork.OpenApiOrgVdcNetwork.Name),
		Description:  types.StringValue(orgNetwork.OpenApiOrgVdcNetwork.Description),
		Gateway:      types.StringValue(orgNetwork.OpenApiOrgVdcNetwork.Subnets.Values[0].Gateway),
		PrefixLength: types.Int64Value(int64(orgNetwork.OpenApiOrgVdcNetwork.Subnets.Values[0].PrefixLength)),
		PrimaryDNS:   types.StringValue(orgNetwork.OpenApiOrgVdcNetwork.Subnets.Values[0].DNSServer1),
		SecondaryDNS: types.StringValue(orgNetwork.OpenApiOrgVdcNetwork.Subnets.Values[0].DNSServer2),
		SuffixDNS:    types.StringValue(orgNetwork.OpenApiOrgVdcNetwork.Subnets.Values[0].DNSSuffix),
	}

	// Get network static IP pools
	ipPools := []staticIPPoolResourceModel{}
	if len(orgNetwork.OpenApiOrgVdcNetwork.Subnets.Values[0].IPRanges.Values) > 0 {
		for _, ipRange := range orgNetwork.OpenApiOrgVdcNetwork.Subnets.Values[0].IPRanges.Values {
			ipPools = append(ipPools, staticIPPoolResourceModel{
				StartAddress: types.StringValue(ipRange.StartAddress),
				EndAddress:   types.StringValue(ipRange.EndAddress),
			})
		}
	}
	var diags diag.Diagnostics
	plan.StaticIPPool, diags = types.SetValueFrom(ctx, types.ObjectType{AttrTypes: staticIPPoolResourceModelAttrTypes}, ipPools)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *networkIsolatedResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var (
		plan, state *networkIsolatedResourceModel
		err         error
		orgNetwork  *govcd.OpenApiOrgVdcNetwork
	)

	// Get current state
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get Org and VDC
	org, err := r.client.GetOrg()
	if err != nil {
		resp.Diagnostics.AddError("[UPDATE] Error retrieving Org", err.Error())
		return
	}

	// Get network

	orgNetwork, err = org.GetOpenApiOrgVdcNetworkById(state.ID.ValueString())
	if err != nil {
		if govcd.ContainsNotFound(err) {
			// Network not found, so remove from state
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("[UPDATE] Error retrieving isolated network", err.Error())
		return
	}

	// Set network type
	ipPools := []staticIPPoolResourceModel{}
	diags := plan.StaticIPPool.ElementsAs(ctx, &ipPools, true)
	resp.Diagnostics.Append(diags...)
	myshared := false // Cloudavenue does not support shared networks
	networkType := &govdctypes.OpenApiOrgVdcNetwork{
		ID:          orgNetwork.OpenApiOrgVdcNetwork.ID,
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		OwnerRef:    orgNetwork.OpenApiOrgVdcNetwork.OwnerRef,
		Shared:      &myshared,
		NetworkType: govdctypes.OrgVdcNetworkTypeIsolated,
		Subnets: govdctypes.OrgVdcNetworkSubnets{
			Values: []govdctypes.OrgVdcNetworkSubnetValues{
				{
					Gateway:      plan.Gateway.ValueString(),
					PrefixLength: int(plan.PrefixLength.ValueInt64()),
					IPRanges: govdctypes.OrgVdcNetworkSubnetIPRanges{
						Values: myProcessIPRanges(ipPools),
					},
					DNSServer1: plan.PrimaryDNS.ValueString(),
					DNSServer2: plan.SecondaryDNS.ValueString(),
					DNSSuffix:  plan.SuffixDNS.ValueString(),
				},
			},
		},
	}

	// Check if VdcGroup or Vdc
	var urn string
	vdc, err := org.GetVDCById(state.OwnerID.ValueString(), true)
	if err == nil { // It is a VDC
		urn, _ = govcd.BuildUrnWithUuid("urn:vcloud:vdc:", vdc.Vdc.ID)
	} else { // It is not a VDC, so it may be a VDC Group
		urn, err := govcd.BuildUrnWithUuid("urn:vcloud:vdcGroup:", state.OwnerID.ValueString())
		if err == nil {
			_, err := org.GetVdcGroupById(urn)
			if err != nil {
				resp.Diagnostics.AddError("[UPDATE] Error retrieving VDC Group from OwnerID", err.Error())
				return
			}
		} else {
			resp.Diagnostics.AddError("[UPDATE] Error retrieving VDC or VDC Group from OwnerID", err.Error())
			return
		}
	}

	// Lock VDC or VDC Group to prevent concurrent access
	vcdMutexKV := mutex.NewKV()
	vcdMutexKV.KvLock(ctx, urn)
	defer vcdMutexKV.KvUnlock(ctx, urn)

	// Update network
	_, err = orgNetwork.Update(networkType)
	if err != nil {
		resp.Diagnostics.AddError("Error updating isolated network", err.Error())
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *networkIsolatedResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var (
		state *networkIsolatedResourceModel
		err   error
	)

	// Get current state
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get Org and VDC
	org, err := r.client.GetOrg()
	if err != nil {
		resp.Diagnostics.AddError("[DELETE] Error retrieving Org", err.Error())
		return
	}

	// Get network
	orgNetwork, err := org.GetOpenApiOrgVdcNetworkById(state.ID.ValueString())
	if err != nil {
		if govcd.ContainsNotFound(err) {
			// Network not found, so remove from state
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("[DELETE] Error retrieving isolated network", err.Error())
		return
	}

	// Check if VdcGroup or Vdc
	var urn string
	vdc, err := org.GetVDCById(state.OwnerID.ValueString(), true)
	if err == nil { // It is a VDC
		urn, _ = govcd.BuildUrnWithUuid("urn:vcloud:vdc:", vdc.Vdc.ID)
	} else { // It is not a VDC, so it may be a VDC Group
		urn, err := govcd.BuildUrnWithUuid("urn:vcloud:vdcGroup:", state.OwnerID.ValueString())
		if err == nil {
			_, err := org.GetVdcGroupById(urn)
			if err != nil {
				resp.Diagnostics.AddError("[DELETE] Error retrieving VDC Group from OwnerID", err.Error())
				return
			}
		} else {
			resp.Diagnostics.AddError("[DELETE] Error retrieving VDC or VDC Group from OwnerID", err.Error())
			return
		}
	}

	// Lock VDC or VDC Group to prevent concurrent access
	vcdMutexKV := mutex.NewKV()
	vcdMutexKV.KvLock(ctx, urn)
	defer vcdMutexKV.KvUnlock(ctx, urn)

	// Delete network
	err = orgNetwork.Delete()
	if err != nil {
		resp.Diagnostics.AddError("Error deleting isolated network", err.Error())
		return
	}
}

func (r *networkIsolatedResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Get URI from import ID
	resourceURI := strings.Split(req.ID, ".")
	if len(resourceURI) != 2 {
		resp.Diagnostics.AddError("Error importing network_routed", "Resource name must be specified as vdc-name.network-name or vdc-group-name.network-name")
		return
	}
	vdcOrVDCGroupName, networkName := resourceURI[0], resourceURI[1]

	// Get VDC
	v, err := r.client.GetVDCOrVDCGroup(vdcOrVDCGroupName)
	if err != nil && !govcd.ContainsNotFound(err) {
		resp.Diagnostics.AddError("Error retrieving VDC", err.Error())
		return
	}

	// Get network
	orgNetwork, err := v.GetOpenApiOrgVdcNetworkByName(networkName)
	if err != nil && !govcd.ContainsNotFound(err) {
		resp.Diagnostics.AddError("Error retrieving org vdc network by name", err.Error())
		return
	}
	// If network is not found, return error
	if orgNetwork == nil {
		resp.Diagnostics.AddError("Error retrieving org network by name", err.Error())
		return
	}

	// Check if network is isolated
	if !orgNetwork.IsIsolated() {
		resp.Diagnostics.AddError("Error importing network_isolated", "Network is not isolated")
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), orgNetwork.OpenApiOrgVdcNetwork.ID)...)
}

// StaticIPPool is a helper function to get the static IP pool from the resource data.
func myProcessIPRanges(mystaticIPPool []staticIPPoolResourceModel) []govdctypes.ExternalNetworkV2IPRange {
	subnetRng := make([]govdctypes.ExternalNetworkV2IPRange, len(mystaticIPPool))
	for i, ipRange := range mystaticIPPool {
		subnetRng[i].StartAddress = ipRange.StartAddress.ValueString()
		subnetRng[i].EndAddress = ipRange.EndAddress.ValueString()
	}
	return subnetRng
}

// Extract Uuid from urn.
func extractUUID(input string) string {
	reGetID := regexp.MustCompile(`([a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12})`)
	matchListID := reGetID.FindAllStringSubmatch(input, -1)
	if len(matchListID) > 0 && len(matchListID[0]) > 0 {
		return matchListID[len(matchListID)-1][1]
	}
	return ""
}
