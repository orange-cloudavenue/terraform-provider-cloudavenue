// Package vapp provides a Terraform resource.
package vapp

import (
	"context"
	"fmt"
	"strings"

	"golang.org/x/exp/slices"

	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/network"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vapp"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vdc"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/uuid"
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

// Metadata returns the resource type name.
func (r *orgNetworkResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_org_network"
}

// Schema defines the schema for the resource.
func (r *orgNetworkResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = network.GetSchema(network.SetRoutedVapp()).GetResource(ctx)
}

func (r *orgNetworkResource) Init(ctx context.Context, rm *orgNetworkModel) (diags diag.Diagnostics) {
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
	defer metrics.New("cloudavenue_vapp_org_network", r.client.GetOrgName(), metrics.Create)()

	// Retrieve values from plan
	var (
		plan *orgNetworkModel
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

	state := &orgNetworkModel{
		ID:                 types.StringValue(uuid.Normalize(uuid.Network, networkID).String()),
		VAppName:           plan.VAppName,
		VAppID:             plan.VAppID,
		VDC:                types.StringValue(r.vdc.GetName()),
		NetworkName:        plan.NetworkName,
		IsFenced:           plan.IsFenced,
		RetainIPMacEnabled: plan.RetainIPMacEnabled,
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *orgNetworkResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer metrics.New("cloudavenue_vapp_org_network", r.client.GetOrgName(), metrics.Read)()

	var state *orgNetworkModel

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

	// Delete resource require vApp is Powered Off
	// Lock
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

	vAppNetwork, networkID, errFindNetwork := state.findOrgNetwork(vAppNetworkConfig)
	resp.Diagnostics.Append(errFindNetwork...)
	if resp.Diagnostics.HasError() {
		return
	}

	if vAppNetwork == (&govcdtypes.VAppNetworkConfiguration{}) {
		resp.State.RemoveResource(ctx)
		return
	}

	isFenced := vAppNetwork.Configuration.FenceMode == govcdtypes.FenceModeNAT

	plan := &orgNetworkModel{
		ID:                 types.StringValue(uuid.Normalize(uuid.Network, *networkID).String()),
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
	defer metrics.New("cloudavenue_vapp_org_network", r.client.GetOrgName(), metrics.Update)()

	var plan, state *orgNetworkModel

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

	plan = &orgNetworkModel{
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
	defer metrics.New("cloudavenue_vapp_org_network", r.client.GetOrgName(), metrics.Delete)()

	var state *orgNetworkModel

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

	// Vapp Statuses
	// 1:  "RESOLVED",
	// 3:  "SUSPENDED",
	// 8:  "POWERED_OFF",

	var (
		vAppRequiredStatuses   = []int{1, 3, 8}
		vAppStatusBeforeAction = r.vapp.VApp.VApp.Status
	)

	// Suspended vApp
	if err := r.vapp.Refresh(); err != nil {
		resp.Diagnostics.AddError("Error retrieving vApp status", err.Error())
		return
	}

	// if vapp not contains VMs, is not possible to PowerOff or Undeploy vApp (error 400)
	if !slices.Contains(vAppRequiredStatuses, r.vapp.VApp.VApp.Status) {
		if r.vapp.VApp.VApp.Children == nil || len(r.vapp.VApp.VApp.Children.VM) == 0 {
			task, err := r.vapp.Undeploy()
			if err != nil {
				resp.Diagnostics.AddError("Error undeploying vApp", err.Error())
				return
			}

			if err = task.WaitTaskCompletion(); err != nil {
				resp.Diagnostics.AddError("Error undeploying vApp", err.Error())
				return
			}
		} else {
			task, err := r.vapp.Suspend()
			if err != nil {
				resp.Diagnostics.AddError("Error suspending vApp", err.Error())
				return
			}

			if err = task.WaitTaskCompletion(); err != nil {
				resp.Diagnostics.AddError("Error suspending vApp", err.Error())
				return
			}
		}
	}

	if _, err := r.vapp.RemoveNetwork(state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting vApp network", err.Error())
	}

	// Vapp Statuses
	// 4: "POWERED_ON"
	// 19: "VAPP_PARTIALLY_DEPLOYED"
	// 20: "PARTIALLY_POWERED_OFF"
	// 21: "PARTIALLY_SUSPENDED"
	if slices.Contains([]int{4, 19, 20, 21}, vAppStatusBeforeAction) {
		// Power On vApp
		task, err := r.vapp.PowerOn()
		if err != nil {
			resp.Diagnostics.AddError("Error powering on vApp", err.Error())
			return
		}

		if err := task.WaitTaskCompletion(); err != nil {
			resp.Diagnostics.AddError("Error powering on vApp", err.Error())
			return
		}
	}
}

func (r *orgNetworkResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	defer metrics.New("cloudavenue_vapp_org_network", r.client.GetOrgName(), metrics.Import)()

	var state *orgNetworkModel
	resourceURI := strings.Split(req.ID, ".")

	if len(resourceURI) != 3 && len(resourceURI) != 2 {
		resp.Diagnostics.AddError("Error importing org_network", "Wrong resource URI format. Expected vdc.vapp.org_network_name or vapp.org_network_name")
		return
	}

	state = &orgNetworkModel{
		VAppName:    types.StringValue(resourceURI[0]),
		NetworkName: types.StringValue(resourceURI[1]),
	}

	if len(resourceURI) == 3 {
		state = &orgNetworkModel{
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
