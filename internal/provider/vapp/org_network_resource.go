// Package vapp provides a Terraform resource.
package vapp

import (
	"context"
	"fmt"
	"strings"

	fboolplanmodifier "github.com/FrangipaneTeam/terraform-plugin-framework-planmodifiers/boolplanmodifier"
	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"

	fboolplanmodifier "github.com/FrangipaneTeam/terraform-plugin-framework-planmodifiers/boolplanmodifier"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/network"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vapp"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vdc"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &orgNetworkResource{}
	_ resource.ResourceWithConfigure   = &orgNetworkResource{}
	_ resource.ResourceWithImportState = &orgNetworkResource{}
)

// NewOrgNetworkResource is a helper function to simplify the provider implementation.
func NewOrgNetworkResource() resource.Resource {
	return &orgNetworkResource{}
}

// orgNetworkResource is the resource implementation.
type orgNetworkResource struct {
	client *client.CloudAvenue
	vdc    vdc.VDC
	vapp   vapp.VAPP
}

type orgNetworkResourceModel struct {
	ID                 types.String `tfsdk:"id"`
	VAppName           types.String `tfsdk:"vapp_name"`
	VAppID             types.String `tfsdk:"vapp_id"`
	VDC                types.String `tfsdk:"vdc"`
	NetworkName        types.String `tfsdk:"network_name"`
	IsFenced           types.Bool   `tfsdk:"is_fenced"`
	RetainIPMacEnabled types.Bool   `tfsdk:"retain_ip_mac_enabled"`
}

// Metadata returns the resource type name.
func (r *orgNetworkResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_org_network"
}

// Schema defines the schema for the resource.
func (r *orgNetworkResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	commonSchema := network.GetSchema(network.SetRoutedVapp()).GetResource()
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"vdc":       vdc.Schema(),
			"vapp_id":   vapp.Schema()["vapp_id"],
			"vapp_name": vapp.Schema()["vapp_name"],
			"is_fenced": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					fboolplanmodifier.SetDefault(false),
				},
				MarkdownDescription: "Fencing allows identical virtual machines in different vApp networks connect to organization VDC networks that are accessed in this vApp. Default is `false`.",
			},
			"retain_ip_mac_enabled": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					fboolplanmodifier.SetDefault(false),
				},
				MarkdownDescription: "Specifies whether the network resources such as IP/MAC of router will be retained across deployments. Default is `false`.",
			},
		},
	}
	// Add common attributes network
	for k, v := range commonSchema.Attributes {
		resp.Schema.Attributes[k] = v
	}
}

func (r *orgNetworkResource) Init(ctx context.Context, rm *orgNetworkResourceModel) (diags diag.Diagnostics) {
	r.vdc, diags = vdc.Init(r.client, rm.VDC)
	if diags.HasError() {
		return
	}

	r.vapp, diags = vapp.Init(r.client, r.vdc, rm.VAppID, rm.VAppName)

	return
}

func (r *orgNetworkResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *orgNetworkResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var (
		plan *orgNetworkResourceModel
	)

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
	resp.Diagnostics.Append(r.vapp.LockVAPP(ctx)...)
	if resp.Diagnostics.HasError() {
		return
	}
	defer r.vapp.UnlockVAPP(ctx)

	orgNetworkName := plan.NetworkName.ValueString()
	orgNetwork, err := r.vdc.GetOrgVdcNetworkByNameOrId(orgNetworkName, true)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving org network", err.Error())
		return
	}

	retainIPMac := plan.RetainIPMacEnabled.ValueBool()
	isFenced := plan.IsFenced.ValueBool()

	vappNetworkSettings := &govcd.VappNetworkSettings{RetainIpMacEnabled: &retainIPMac}

	vAppNetworkConfig, err := r.vapp.AddOrgNetwork(vappNetworkSettings, orgNetwork.OrgVDCNetwork, isFenced)
	if err != nil {
		resp.Diagnostics.AddError("Error creating vApp network", err.Error())
		return
	}

	vAppNetwork := govcdtypes.VAppNetworkConfiguration{}
	for _, networkConfig := range vAppNetworkConfig.NetworkConfig {
		if networkConfig.NetworkName == orgNetwork.OrgVDCNetwork.Name {
			vAppNetwork = networkConfig
		}
	}

	if vAppNetwork == (govcdtypes.VAppNetworkConfiguration{}) {
		resp.Diagnostics.AddError("Error creating vApp network", "vApp network not found in vApp network config")
		return
	}

	networkID, err := govcd.GetUuidFromHref(vAppNetwork.Link.HREF, false)
	if err != nil {
		resp.Diagnostics.AddError("Error creating vApp network uuid", err.Error())
		return
	}

	id := common.NormalizeID("urn:vcloud:network:", networkID)

	plan = &orgNetworkResourceModel{
		ID:                 types.StringValue(id),
		VAppName:           plan.VAppName,
		VDC:                types.StringValue(r.vdc.GetName()),
		NetworkName:        plan.NetworkName,
		IsFenced:           plan.IsFenced,
		RetainIPMacEnabled: plan.RetainIPMacEnabled,
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *orgNetworkResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *orgNetworkResourceModel

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
		resp.Diagnostics.AddError("Error retrieving vApp network config", err.Error())
		return
	}

	vAppNetwork, networkID, errFindNetwork := state.findOrgNetwork(vAppNetworkConfig)
	resp.Diagnostics.Append(errFindNetwork...)
	if resp.Diagnostics.HasError() {
		return
	}

	if vAppNetwork == (&govcdtypes.VAppNetworkConfiguration{}) {
		resp.State.RemoveResource(ctx)
		return
	}

	id := common.NormalizeID("urn:vcloud:network:", *networkID)
	isFenced := vAppNetwork.Configuration.FenceMode == govcdtypes.FenceModeNAT

	plan := &orgNetworkResourceModel{
		ID:                 types.StringValue(id),
		VAppName:           state.VAppName,
		VAppID:             state.VAppID,
		VDC:                types.StringValue(r.vdc.GetName()),
		NetworkName:        state.NetworkName,
		IsFenced:           types.BoolValue(isFenced),
		RetainIPMacEnabled: types.BoolValue(*vAppNetwork.Configuration.RetainNetInfoAcrossDeployments),
	}

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *orgNetworkResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state *orgNetworkResourceModel

	// Get current state
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Init resource
	resp.Diagnostics.Append(r.Init(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Lock vApp
	resp.Diagnostics.Append(r.vapp.LockVAPP(ctx)...)
	if resp.Diagnostics.HasError() {
		return
	}
	defer r.vapp.UnlockVAPP(ctx)

	vAppNetworkConfig, err := r.vapp.GetNetworkConfig()
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving vApp network config", err.Error())
		return
	}

	vAppNetwork, _, errFindNetwork := plan.findOrgNetwork(vAppNetworkConfig)
	resp.Diagnostics.Append(errFindNetwork...)
	if resp.Diagnostics.HasError() {
		return
	}

	if vAppNetwork == (&govcdtypes.VAppNetworkConfiguration{}) {
		resp.State.RemoveResource(ctx)
		return
	}

	isFenced := vAppNetwork.Configuration.FenceMode == govcdtypes.FenceModeNAT

	if plan.IsFenced.ValueBool() != isFenced || plan.RetainIPMacEnabled.ValueBool() != *vAppNetwork.Configuration.RetainNetInfoAcrossDeployments {
		tflog.Debug(ctx, "updating vApp network")
		retainIP := plan.RetainIPMacEnabled.ValueBool()
		vappNetworkSettings := &govcd.VappNetworkSettings{
			ID:                 state.ID.ValueString(),
			RetainIpMacEnabled: &retainIP,
		}
		_, err = r.vapp.UpdateOrgNetwork(vappNetworkSettings, plan.IsFenced.ValueBool())
		if err != nil {
			resp.Diagnostics.AddError("Error updating vApp network", err.Error())
			return
		}
	}

	plan = &orgNetworkResourceModel{
		ID:                 state.ID,
		VAppName:           state.VAppName,
		VAppID:             state.VAppID,
		VDC:                state.VDC,
		NetworkName:        state.NetworkName,
		IsFenced:           plan.IsFenced,
		RetainIPMacEnabled: plan.RetainIPMacEnabled,
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *orgNetworkResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *orgNetworkResourceModel

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
	resp.Diagnostics.Append(r.vapp.LockVAPP(ctx)...)
	if resp.Diagnostics.HasError() {
		return
	}
	defer r.vapp.UnlockVAPP(ctx)

	_, err := r.vapp.RemoveNetwork(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting vApp network", err.Error())
		return
	}
}

func (r *orgNetworkResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var state *orgNetworkResourceModel
	resourceURI := strings.Split(req.ID, ".")

	if len(resourceURI) != 3 && len(resourceURI) != 2 {
		resp.Diagnostics.AddError("Error importing org_network", "Wrong resource URI format. Expected vdc.vapp.org_network_name or vapp.org_network_name")
		return
	}

	state = &orgNetworkResourceModel{
		VAppName:    types.StringValue(resourceURI[0]),
		NetworkName: types.StringValue(resourceURI[1]),
	}

	if len(resourceURI) == 3 {
		state = &orgNetworkResourceModel{
			VDC:         types.StringValue(resourceURI[0]),
			VAppName:    types.StringValue(resourceURI[1]),
			NetworkName: types.StringValue(resourceURI[2]),
		}
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
