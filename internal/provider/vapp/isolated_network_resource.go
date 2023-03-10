// Package vapp provides a Terraform resource.
package vapp

import (
	"context"
	"fmt"
	"strings"

	"github.com/FrangipaneTeam/terraform-plugin-framework-validators/stringvalidator"
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
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/boolpm"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/stringpm"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vapp"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vdc"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &isolatedNetworkResource{}
	_ resource.ResourceWithConfigure   = &isolatedNetworkResource{}
	_ resource.ResourceWithImportState = &isolatedNetworkResource{}
)

// NewIsolatedNetworkResource is a helper function to simplify the provider implementation.
func NewIsolatedNetworkResource() resource.Resource {
	return &isolatedNetworkResource{}
}

// isolatedNetworkResource is the resource implementation.
type isolatedNetworkResource struct {
	client *client.CloudAvenue
	vdc    vdc.VDC
	vapp   vapp.VApp
}

type isolatedNetworkResourceModel struct {
	ID                 types.String `tfsdk:"id"`
	VDC                types.String `tfsdk:"vdc"`
	Name               types.String `tfsdk:"name"`
	Description        types.String `tfsdk:"description"`
	VAppName           types.String `tfsdk:"vapp_name"`
	VAppID             types.String `tfsdk:"vapp_id"`
	Netmask            types.String `tfsdk:"netmask"`
	Gateway            types.String `tfsdk:"gateway"`
	DNS1               types.String `tfsdk:"dns1"`
	DNS2               types.String `tfsdk:"dns2"`
	DNSSuffix          types.String `tfsdk:"dns_suffix"`
	GuestVLANAllowed   types.Bool   `tfsdk:"guest_vlan_allowed"`
	RetainIPMacEnabled types.Bool   `tfsdk:"retain_ip_mac_enabled"`
	StaticIPPool       types.Set    `tfsdk:"static_ip_pool"`
}

type staticIPPoolModel struct {
	StartAddress types.String `tfsdk:"start_address"`
	EndAddress   types.String `tfsdk:"end_address"`
}

var staticIPPoolModelAttrTypes = map[string]attr.Type{
	"start_address": types.StringType,
	"end_address":   types.StringType,
}

// Metadata returns the resource type name.
func (r *isolatedNetworkResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_" + "isolated_network"
}

// Schema defines the schema for the resource.
func (r *isolatedNetworkResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Provides a Cloud Avenue isolated vAPP Network resource. This can be used to create, modify, and delete isolated vAPP Network.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the vApp network.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"vdc":       vdc.Schema(),
			"vapp_id":   vapp.Schema()["vapp_id"],
			"vapp_name": vapp.Schema()["vapp_name"],
			"name": schema.StringAttribute{
				MarkdownDescription: "(ForceNew) The name of the vApp network.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the vApp network.",
				Optional:            true,
			},
			"netmask": schema.StringAttribute{
				MarkdownDescription: "(ForceNew) The netmask for the network. Default is `255.255.255.0`",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.IsValidNetmask(),
				},
				PlanModifiers: []planmodifier.String{
					stringpm.SetDefault("255.255.255.0"),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"gateway": schema.StringAttribute{
				MarkdownDescription: "The gateway of the network.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.IsValidIP(),
				},
			},
			"dns1": schema.StringAttribute{
				MarkdownDescription: "First DNS server.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.IsValidIP(),
				},
			},
			"dns2": schema.StringAttribute{
				MarkdownDescription: "Second DNS server.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.IsValidIP(),
				},
			},
			"dns_suffix": schema.StringAttribute{
				MarkdownDescription: "A FQDN for the virtual machines on this network.",
				Optional:            true,
			},
			"guest_vlan_allowed": schema.BoolAttribute{
				MarkdownDescription: "True if Network allows guest VLAN. Default to `false`.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolpm.SetDefault(false),
				},
			},
			"retain_ip_mac_enabled": schema.BoolAttribute{
				MarkdownDescription: "Specifies whether the network resources such as IP/MAC of router will be retained across deployments. Default to `false`.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolpm.SetDefault(false),
				},
			},
			"static_ip_pool": schema.SetNestedAttribute{
				MarkdownDescription: "Range(s) of IPs permitted to be used as static IPs for virtual machines",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"start_address": schema.StringAttribute{
							MarkdownDescription: "The first address in the IP Range.",
							Required:            true,
						},
						"end_address": schema.StringAttribute{
							MarkdownDescription: "The last address in the IP Range.",
							Required:            true,
						},
					},
				},
			},
		},
	}
}

func (r *isolatedNetworkResource) Init(ctx context.Context, rm *isolatedNetworkResourceModel) (diags diag.Diagnostics) {
	r.vdc, diags = vdc.Init(r.client, rm.VDC)
	if diags.HasError() {
		return
	}

	r.vapp, diags = vapp.Init(r.client, r.vdc, rm.VAppID, rm.VAppName)

	return
}

func (r *isolatedNetworkResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *isolatedNetworkResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var (
		plan *isolatedNetworkResourceModel
	)

	// Read the plan
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Init resource
	resp.Diagnostics.Append(r.Init(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Lock vApp
	resp.Diagnostics.Append(r.vapp.LockParentVApp(ctx)...)
	if resp.Diagnostics.HasError() {
		return
	}
	defer r.vapp.UnlockParentVApp(ctx)

	// Configure network
	retainIPMac := plan.RetainIPMacEnabled.ValueBool()
	guestVLAN := plan.GuestVLANAllowed.ValueBool()
	vappNetworkName := plan.Name.ValueString()

	var staticIPPools []*staticIPPoolModel
	resp.Diagnostics.Append(plan.StaticIPPool.ElementsAs(ctx, &staticIPPools, true)...)

	staticIPRanges := make([]*govcdtypes.IPRange, 0)
	for _, staticIPPool := range staticIPPools {
		staticIPRanges = append(staticIPRanges, &govcdtypes.IPRange{
			StartAddress: staticIPPool.StartAddress.ValueString(),
			EndAddress:   staticIPPool.EndAddress.ValueString(),
		})
	}

	vappNetworkSettings := &govcd.VappNetworkSettings{
		Name:               vappNetworkName,
		Description:        plan.Description.ValueString(),
		Gateway:            plan.Gateway.ValueString(),
		NetMask:            plan.Netmask.ValueString(),
		DNS1:               plan.DNS1.ValueString(),
		DNS2:               plan.DNS2.ValueString(),
		DNSSuffix:          plan.DNSSuffix.ValueString(),
		StaticIPRanges:     staticIPRanges,
		RetainIpMacEnabled: &retainIPMac,
		GuestVLANAllowed:   &guestVLAN,
	}

	// Create network
	vAppNetworkConfig, err := r.vapp.CreateVappNetwork(vappNetworkSettings, nil)
	if err != nil {
		resp.Diagnostics.AddError("Error creating vApp network", err.Error())
		return
	}

	vAppNetwork := govcdtypes.VAppNetworkConfiguration{}
	for _, networkConfig := range vAppNetworkConfig.NetworkConfig {
		if networkConfig.NetworkName == vappNetworkName {
			vAppNetwork = networkConfig
		}
	}

	if vAppNetwork == (govcdtypes.VAppNetworkConfiguration{}) {
		resp.Diagnostics.AddError("didn't find vApp network: %s", vappNetworkName)
		return
	}

	// Get UUID.
	networkID, err := govcd.GetUuidFromHref(vAppNetwork.Link.HREF, false)
	if err != nil {
		resp.Diagnostics.AddError("Error creating vApp network uuid", err.Error())
		return
	}

	plan.ID = types.StringValue(common.NormalizeID("urn:vcloud:network:", networkID))
	plan.VDC = types.StringValue(r.vdc.GetName())

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *isolatedNetworkResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var (
		state *isolatedNetworkResourceModel
		diag  diag.Diagnostics
	)

	// Get current state
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Init resource
	resp.Diagnostics.Append(r.Init(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	vAppNetworkConfig, err := r.vapp.GetNetworkConfig()
	if err != nil {
		resp.Diagnostics.AddError("Error getting vApp networks", err.Error())
		return
	}

	vAppNetwork := govcdtypes.VAppNetworkConfiguration{}
	for _, networkConfig := range vAppNetworkConfig.NetworkConfig {
		if networkConfig.NetworkName == state.Name.ValueString() {
			vAppNetwork = networkConfig
		}
	}

	if vAppNetwork == (govcdtypes.VAppNetworkConfiguration{}) {
		resp.State.RemoveResource(ctx)
		return
	}

	// Get UUID.
	networkID, err := govcd.GetUuidFromHref(vAppNetwork.Link.HREF, false)
	if err != nil {
		resp.Diagnostics.AddError("Error creating vApp network uuid", err.Error())
		return
	}

	id := common.NormalizeID("urn:vcloud:network:", networkID)

	plan := &isolatedNetworkResourceModel{
		ID:                 types.StringValue(id),
		VDC:                types.StringValue(r.vdc.GetName()),
		Name:               types.StringValue(vAppNetwork.NetworkName),
		Description:        types.StringValue(vAppNetwork.Description),
		VAppName:           state.VAppName,
		Netmask:            types.StringValue(vAppNetwork.Configuration.IPScopes.IPScope[0].Netmask),
		Gateway:            types.StringValue(vAppNetwork.Configuration.IPScopes.IPScope[0].Gateway),
		DNS1:               types.StringValue(vAppNetwork.Configuration.IPScopes.IPScope[0].DNS1),
		DNS2:               types.StringValue(vAppNetwork.Configuration.IPScopes.IPScope[0].DNS2),
		DNSSuffix:          types.StringValue(vAppNetwork.Configuration.IPScopes.IPScope[0].DNSSuffix),
		GuestVLANAllowed:   types.BoolValue(*vAppNetwork.Configuration.GuestVlanAllowed),
		RetainIPMacEnabled: types.BoolValue(*vAppNetwork.Configuration.RetainNetInfoAcrossDeployments),
	}

	if len(vAppNetwork.Configuration.IPScopes.IPScope) == 0 {
		plan.Netmask = types.StringNull()
		plan.Gateway = types.StringNull()
		plan.DNS1 = types.StringNull()
		plan.DNS2 = types.StringNull()
		plan.DNSSuffix = types.StringNull()
	} else {
		plan.Netmask = types.StringValue(vAppNetwork.Configuration.IPScopes.IPScope[0].Netmask)
		plan.Gateway = types.StringValue(vAppNetwork.Configuration.IPScopes.IPScope[0].Gateway)
		plan.DNS1 = types.StringValue(vAppNetwork.Configuration.IPScopes.IPScope[0].DNS1)
		plan.DNS2 = types.StringValue(vAppNetwork.Configuration.IPScopes.IPScope[0].DNS2)
		plan.DNSSuffix = types.StringValue(vAppNetwork.Configuration.IPScopes.IPScope[0].DNSSuffix)
	}

	// Fix empty string as StringNull for optional attributes
	if plan.Description.ValueString() == "" {
		plan.Description = types.StringNull()
	}
	if plan.DNS1.ValueString() == "" {
		plan.DNS1 = types.StringNull()
	}
	if plan.DNS2.ValueString() == "" {
		plan.DNS2 = types.StringNull()
	}
	if plan.DNSSuffix.ValueString() == "" {
		plan.DNSSuffix = types.StringNull()
	}

	// Loop on static_ip_pool if it is not nil
	staticIPRanges := make([]staticIPPoolModel, 0)
	if vAppNetwork.Configuration.IPScopes.IPScope[0].IPRanges != nil {
		for _, staticIPRange := range vAppNetwork.Configuration.IPScopes.IPScope[0].IPRanges.IPRange {
			staticIPRanges = append(staticIPRanges, staticIPPoolModel{
				StartAddress: types.StringValue(staticIPRange.StartAddress),
				EndAddress:   types.StringValue(staticIPRange.EndAddress),
			})
		}
		plan.StaticIPPool, diag = types.SetValueFrom(ctx, types.ObjectType{AttrTypes: staticIPPoolModelAttrTypes}, staticIPRanges)

		resp.Diagnostics.Append(diag...)
		if resp.Diagnostics.HasError() {
			return
		}
	} else {
		plan.StaticIPPool = types.SetNull(types.ObjectType{AttrTypes: staticIPPoolModelAttrTypes})
	}

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *isolatedNetworkResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan *isolatedNetworkResourceModel

	// Get current state
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Init resource
	resp.Diagnostics.Append(r.Init(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Lock vApp
	resp.Diagnostics.Append(r.vapp.LockParentVApp(ctx)...)
	if resp.Diagnostics.HasError() {
		return
	}
	defer r.vapp.UnlockParentVApp(ctx)

	// Configure network
	retainIPMac := plan.RetainIPMacEnabled.ValueBool()
	guestVLAN := plan.GuestVLANAllowed.ValueBool()
	vappNetworkName := plan.Name.ValueString()

	var staticIPPools []*staticIPPoolModel
	resp.Diagnostics.Append(plan.StaticIPPool.ElementsAs(ctx, &staticIPPools, true)...)

	staticIPRanges := make([]*govcdtypes.IPRange, 0)
	for _, staticIPPool := range staticIPPools {
		staticIPRanges = append(staticIPRanges, &govcdtypes.IPRange{
			StartAddress: staticIPPool.StartAddress.ValueString(),
			EndAddress:   staticIPPool.EndAddress.ValueString(),
		})
	}

	vappNetworkSettings := &govcd.VappNetworkSettings{
		ID:                 common.ExtractUUID(plan.ID.ValueString()),
		Name:               vappNetworkName,
		Description:        plan.Description.ValueString(),
		Gateway:            plan.Gateway.ValueString(),
		NetMask:            plan.Netmask.ValueString(),
		DNS1:               plan.DNS1.ValueString(),
		DNS2:               plan.DNS2.ValueString(),
		DNSSuffix:          plan.DNSSuffix.ValueString(),
		StaticIPRanges:     staticIPRanges,
		RetainIpMacEnabled: &retainIPMac,
		GuestVLANAllowed:   &guestVLAN,
	}

	// Update network
	_, err := r.vapp.UpdateNetwork(vappNetworkSettings, nil)
	if err != nil {
		resp.Diagnostics.AddError("Error updating vApp network", err.Error())
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *isolatedNetworkResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *isolatedNetworkResourceModel

	// Get current state
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Init resource
	resp.Diagnostics.Append(r.Init(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Lock vApp
	resp.Diagnostics.Append(r.vapp.LockParentVApp(ctx)...)
	if resp.Diagnostics.HasError() {
		return
	}
	defer r.vapp.UnlockParentVApp(ctx)

	_, err := r.vapp.RemoveNetwork(state.ID.String())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting vApp network", err.Error())
		return
	}
}

func (r *isolatedNetworkResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ".")

	if len(idParts) != 3 && len(idParts) != 2 {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: vdc.vapp_name.network_name or vapp_name.network_name. Got: %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("vapp_name"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), idParts[1])...)

	if len(idParts) == 3 {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("vdc"), idParts[0])...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("vapp_name"), idParts[1])...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), idParts[2])...)
	}
}
