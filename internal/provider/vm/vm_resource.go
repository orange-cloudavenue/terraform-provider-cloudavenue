// Package vm provides a Terraform resource.
package vm

import (
	"context"
	"fmt"
	"strings"

	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/adminvdc"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vapp"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vdc"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vm"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &vmResource{}
	_ resource.ResourceWithConfigure   = &vmResource{}
	_ resource.ResourceWithImportState = &vmResource{}
)

// NewVmResource is a helper function to simplify the provider implementation.
func NewVMResource() resource.Resource {
	return &vmResource{}
}

// vmResource is the resource implementation.
type vmResource struct {
	client *client.CloudAvenue

	// Uncomment the following lines if you need to access the resource's.
	// org    org.Org
	vdc      vdc.VDC
	adminVDC adminvdc.AdminVDC
	vapp     vapp.VAPP
	vm       vm.VM
}

// Init Initializes the resource.
func (r *vmResource) Init(ctx context.Context, rm *vm.VMResourceModel) (diags diag.Diagnostics) {
	var d diag.Diagnostics

	r.vdc, d = vdc.Init(r.client, rm.VDC)
	diags.Append(d...)
	if diags.HasError() {
		return
	}

	r.adminVDC, d = adminvdc.Init(r.client, rm.VDC)
	diags.Append(d...)
	if diags.HasError() {
		return
	}

	r.vapp, d = vapp.Init(r.client, r.vdc, rm.VappID, rm.VappName)
	diags.Append(d...)
	if diags.HasError() {
		return
	}

	if r.vapp.VAPP == nil {
		diags.AddError("Vapp not found", fmt.Sprintf("Vapp %s not found in VDC %s", rm.VappName, rm.VDC))
		return
	}

	// Vm is not initialized here because if VM is not found in read. Delete resource in state will be called.

	return
}

// Metadata returns the resource type name.
func (r *vmResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName
}

// Schema defines the schema for the resource.
func (r *vmResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = vmSuperSchema(ctx).GetResource(ctx)
}

func (r *vmResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *vmResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &vm.VMResourceModel{}

	var (
		vmCreated vm.VM
		d         diag.Diagnostics
	)

	// Retrieve values from plan
	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Init the resource
	resp.Diagnostics.Append(r.Init(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	/*
		Implement the resource creation logic here.
	*/

	// * DeployOS
	deployOS, d := plan.DeployOSFromPlan(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// * State
	state, d := plan.StateFromPlan(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// * Settings
	settingsConfig, d := plan.SettingsFromPlan(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// * Customization
	customizationConfig, d := settingsConfig.CustomizationFromPlan(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// * Create VM with Template
	if !deployOS.VappTemplateID.IsNull() {
		vmCreated, d = r.createVMWithTemplate(ctx, *plan)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
	} else {
		// Create VM with ISO
		vmCreated, d = r.createVMWithBootImage(ctx, *plan)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	r.vm, d = r.processAfterCreate(ctx, vmCreated, *plan)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(r.vmPowerOn(ctx, *plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.vm.Refresh(); err != nil {
		resp.Diagnostics.AddError(
			"Unable to refresh VM",
			fmt.Sprintf("Error: %s", err),
		)
		return
	}

	status, err := r.vm.GetStatus()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to get VM status",
			fmt.Sprintf("Error: %s", err),
		)
		return
	}
	state.Status = types.StringValue(status)

	settings, err := r.vm.SettingsRead(ctx, customizationConfig)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to get VM settings",
			fmt.Sprintf("Error: %s", err),
		)
		return
	}

	networks, err := r.vm.NetworksRead()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to get VM networks",
			fmt.Sprintf("Error: %s", err),
		)
		return
	}

	tfState := *plan
	tfState.ID = types.StringValue(r.vm.GetID())
	tfState.VappID = types.StringValue(r.vapp.GetID())
	tfState.VappName = types.StringValue(r.vapp.GetName())
	tfState.State = state.ToPlan(ctx)
	tfState.VDC = types.StringValue(r.vdc.GetName())
	tfState.Settings = settings.ToPlan(ctx)
	tfState.Resource = r.vm.ResourceRead(ctx).ToPlan(ctx, networks)

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, tfState)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *vmResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &vm.VMResourceModel{}

	// Get current state
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Init the resource
	resp.Diagnostics.Append(r.Init(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var d diag.Diagnostics

	r.vm, d = vm.Init(r.client, r.vapp, vm.GetVMOpts{
		ID:   state.ID,
		Name: types.StringNull(),
	})

	if d.HasError() {
		if d.Contains(diag.NewErrorDiagnostic("VM not found", govcd.ErrorEntityNotFound.Error())) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.Append(d...)
		return
	}

	plan, d := r.read(ctx, state)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *vmResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) { //nolint:gocyclo
	plan := &vm.VMResourceModel{}
	state := &vm.VMResourceModel{}

	// Get current state
	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Init the resource
	resp.Diagnostics.Append(r.Init(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	/*
		Implement the resource update here
	*/

	var d diag.Diagnostics

	r.vm, d = vm.Init(r.client, r.vapp, vm.GetVMOpts{
		ID:   state.ID,
		Name: types.StringNull(),
	})
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	/*
		Update VM was 2 major steps:
		1. Hot update (VM must be powered on)
		2. Cold update (VM must be powered off)
	*/

	needColdChange := struct {
		memory  bool
		cpu     bool
		network bool
	}{
		memory:  false,
		cpu:     false,
		network: false,
	}

	allStructsPlan, d := plan.AllStructsFromPlan(ctx)
	resp.Diagnostics.Append(d...)

	allStructsState, d := state.AllStructsFromPlan(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// ! Hot update

	// ? Resource
	if !allStructsPlan.Resource.Equal(allStructsState.Resource) {
		// * CPU and CPU cores
		if !allStructsPlan.Resource.CPUs.Equal(allStructsState.Resource.CPUs) || !allStructsPlan.Resource.CPUsCores.Equal(allStructsState.Resource.CPUsCores) {
			// Detected change on CPU or CPU cores
			if r.vm.GetCPUHotAddEnabled() {
				// CPU hot update is enabled
				if err := r.vm.ChangeCPUAndCoreCount(utils.TakeIntPointer(int(allStructsPlan.Resource.CPUs.ValueInt64())), utils.TakeIntPointer(int(allStructsPlan.Resource.CPUsCores.ValueInt64()))); err != nil {
					resp.Diagnostics.AddError(
						"Unable to change CPU and CPU Cores",
						fmt.Sprintf("Error: %s", err),
					)
					return
				}
			} else {
				needColdChange.cpu = true
			}
		}

		// * Memory
		if !allStructsPlan.Resource.Memory.Equal(allStructsState.Resource.Memory) {
			// Detected change on memory
			if r.vm.GetMemoryHotAddEnabled() {
				// Memory hot update is enabled
				if err := r.vm.ChangeMemory(allStructsPlan.Resource.Memory.ValueInt64()); err != nil {
					resp.Diagnostics.AddError(
						"Unable to change memory size",
						fmt.Sprintf("Error: %s", err),
					)
					return
				}
			} else {
				needColdChange.memory = true
			}
		}
	}

	// ? Resource -> Networks
	if !allStructsPlan.Resource.Networks.Equal(allStructsState.Resource.Networks) {
		networkPlan, d := allStructsPlan.Resource.NetworksFromPlan(ctx)
		resp.Diagnostics.Append(d...)
		networkState, d := allStructsState.Resource.NetworksFromPlan(ctx)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}

		var networkChanged bool

		for i, network := range *networkPlan {
			if *networkState != nil || len(*networkState) == 0 || (*networkState)[i] == (vm.VMResourceModelResourceNetwork{}) || !network.Equal((*networkState)[i]) {
				if network.IsPrimary.ValueBool() {
					needColdChange.network = true
					break
				}
			}
		}

		if networkChanged && !needColdChange.network {
			// Found network change, but not primary network
			networkConnection := []vm.NetworkConnection{}
			for _, n := range *networkPlan {
				networkConnection = append(networkConnection, n.ConvertToNetworkConnection())
			}
			networkConfig, err := r.vm.ConstructNetworksConnection(networkConnection)
			if err != nil {
				resp.Diagnostics.AddError("Error retrieving network config", err.Error())
				return
			}
			if err = r.vm.UpdateNetworkConnectionSection(&networkConfig); err != nil {
				resp.Diagnostics.AddError("Error updating network config", err.Error())
				return
			}
		}
	}

	// ? Settings
	if !allStructsPlan.Settings.Equal(allStructsState.Settings) {
		// * Guest properties
		if !allStructsPlan.Settings.GuestProperties.Equal(allStructsState.Settings.GuestProperties) {
			// Detected change on guest properties
			guestProperties := make(map[string]string, 0)
			for key, value := range allStructsPlan.Settings.GuestProperties.Elements() {
				guestProperties[key] = strings.Trim(value.String(), "\"")
			}
			if err := r.vm.SetGuestProperties(guestProperties); err != nil {
				resp.Diagnostics.AddError("Error updating guest properties", err.Error())
				return
			}
		}

		// * AffinityRule
		if !allStructsPlan.Settings.AffinityRuleID.Equal(allStructsState.Settings.AffinityRuleID) {
			// Detected change on affinity rule
			affinityRuleID := allStructsPlan.Settings.AffinityRuleID.ValueString()
			if affinityRuleID == "" {
				if r.vdc.Vdc.Vdc.DefaultComputePolicy == nil {
					resp.Diagnostics.AddError("Error updating affinity rule", "Default affinity rule is not set")
					return
				}
				affinityRuleID = r.vdc.Vdc.Vdc.DefaultComputePolicy.ID
			}

			if _, err := r.vm.UpdateComputePolicyV2("", affinityRuleID, ""); err != nil {
				resp.Diagnostics.AddError("Error updating affinity rule", err.Error())
				return
			}
		}

		// * StorageProfile
		if !allStructsPlan.Settings.StorageProfile.Equal(allStructsState.Settings.StorageProfile) {
			var (
				storageProfile *govcdtypes.Reference
				err            error
			)

			if allStructsPlan.Settings.StorageProfile.ValueString() == "" {
				// If storage profile is empty, we use the default one
				storageProfile, err = r.adminVDC.GetDefaultStorageProfileReference()
				if err != nil {
					resp.Diagnostics.AddError("Error retrieving default storage profile", err.Error())
					return
				}
			} else {
				storageProfile, err = r.vdc.GetStorageProfileReference(allStructsPlan.Settings.StorageProfile.ValueString(), false)
				if err != nil {
					resp.Diagnostics.AddError("Error retrieving storage profile", err.Error())
					return
				}
			}

			if _, err := r.vm.UpdateStorageProfile(storageProfile.HREF); err != nil {
				resp.Diagnostics.AddError("Error updating storage profile", err.Error())
				return
			}
		}
	}

	// ! Hot Update

	// ! Cold Update
	vmStatusBeforeUpdate, err := r.vm.GetStatus()
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving VM status", err.Error())
		return
	}

	if !allStructsPlan.State.PowerON.Equal(allStructsState.State.PowerON) ||
		!allStructsPlan.Settings.ExposeHardwareVirtualization.Equal(allStructsState.Settings.ExposeHardwareVirtualization) ||
		!allStructsPlan.Settings.OsType.Equal(allStructsState.Settings.OsType) ||
		!allStructsPlan.Resource.CPUHotAddEnabled.Equal(allStructsState.Resource.CPUHotAddEnabled) ||
		!allStructsPlan.Resource.MemoryHotAddEnabled.Equal(allStructsState.Resource.MemoryHotAddEnabled) ||
		!plan.Description.Equal(state.Description) ||
		needColdChange.cpu ||
		needColdChange.memory ||
		needColdChange.network {
		if vmStatusBeforeUpdate != poweredOFF {
			task, err := r.vm.Undeploy()
			if err != nil {
				resp.Diagnostics.AddError("Error undeploying VM", err.Error())
				return
			}

			err = task.WaitTaskCompletion()
			if err != nil {
				resp.Diagnostics.AddError("Error waiting undeploying VM", err.Error())
				return
			}
		}

		// * ExposeHardwareVirtualization
		if !allStructsPlan.Settings.ExposeHardwareVirtualization.Equal(allStructsState.Settings.ExposeHardwareVirtualization) {
			task, err := r.vm.ToggleHardwareVirtualization(allStructsPlan.Settings.ExposeHardwareVirtualization.ValueBool())
			if err != nil {
				resp.Diagnostics.AddError("Error updating ExposeHardwareVirtualization", err.Error())
				return
			}

			if err = task.WaitTaskCompletion(); err != nil {
				resp.Diagnostics.AddError("Error waiting ExposeHardwareVirtualization", err.Error())
				return
			}
		}

		// * OsType And Description
		var (
			vmSpecSectionUpdate = false
			vmSpecSection       = r.vm.VM.VM.VM.VmSpecSection
			description         = r.vm.GetDescription()
		)

		if !allStructsPlan.Settings.OsType.Equal(allStructsState.Settings.OsType) {
			vmSpecSectionUpdate = true
			vmSpecSection.OsType = allStructsPlan.Settings.OsType.ValueString()
		}

		if !plan.Description.Equal(state.Description) {
			vmSpecSectionUpdate = true
			description = plan.Description.ValueString()
		}

		if vmSpecSectionUpdate {
			task, err := r.vm.UpdateVmSpecSectionAsync(vmSpecSection, description)
			if err != nil {
				resp.Diagnostics.AddError("Error updating VM spec section", err.Error())
				return
			}

			err = task.WaitTaskCompletion()
			if err != nil {
				resp.Diagnostics.AddError("Error waiting VM spec section", err.Error())
				return
			}
		}

		// * CPUHotAddEnabled And MemoryHotAddEnabled
		if !allStructsPlan.Resource.CPUHotAddEnabled.Equal(allStructsState.Resource.CPUHotAddEnabled) ||
			!allStructsPlan.Resource.MemoryHotAddEnabled.Equal(allStructsState.Resource.MemoryHotAddEnabled) {
			if _, err := r.vm.UpdateVmCpuAndMemoryHotAdd(allStructsPlan.Resource.CPUHotAddEnabled.ValueBool(), allStructsState.Resource.MemoryHotAddEnabled.ValueBool()); err != nil {
				resp.Diagnostics.AddError("Error updating CPUHotAddEnabled", err.Error())
				return
			}
		}
		// * Cold CPU Change
		if needColdChange.cpu {
			if err := r.vm.ChangeCPUAndCoreCount(utils.TakeIntPointer(int(allStructsPlan.Resource.CPUs.ValueInt64())), utils.TakeIntPointer(int(allStructsPlan.Resource.CPUsCores.ValueInt64()))); err != nil {
				resp.Diagnostics.AddError(
					"Unable to change CPU and CPU Cores",
					fmt.Sprintf("Error: %s", err),
				)
				return
			}
		}

		// * Cold Memory Change
		if needColdChange.memory {
			if err := r.vm.ChangeMemory(allStructsPlan.Resource.Memory.ValueInt64()); err != nil {
				resp.Diagnostics.AddError(
					"Unable to change memory size",
					fmt.Sprintf("Error: %s", err),
				)
				return
			}
		}

		// * Cold Network Change
		if needColdChange.network {
			networkPlan, d := allStructsPlan.Resource.NetworksFromPlan(ctx)
			resp.Diagnostics.Append(d...)
			if resp.Diagnostics.HasError() {
				return
			}

			// Found network change, but not primary network
			networkConnection := []vm.NetworkConnection{}
			for _, n := range *networkPlan {
				networkConnection = append(networkConnection, n.ConvertToNetworkConnection())
			}
			networkConfig, err := r.vm.ConstructNetworksConnection(networkConnection)
			if err != nil {
				resp.Diagnostics.AddError("Error retrieving network config", err.Error())
				return
			}
			if err := r.vm.UpdateNetworkConnectionSection(&networkConfig); err != nil {
				resp.Diagnostics.AddError("Error updating network config", err.Error())
				return
			}
		}
	} // ! End of Cold update

	vmStatus, err := r.vm.GetStatus()
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving VM status", err.Error())
		return
	}

	if allStructsPlan.State.PowerON.ValueBool() {
		if !allStructsPlan.Settings.Customization.Attributes()["force"].(types.Bool).ValueBool() && vmStatus != poweredON {
			task, err := r.vm.PowerOn()
			if err != nil {
				resp.Diagnostics.AddError("Error powering on VM", err.Error())
				return
			}

			err = task.WaitTaskCompletion()
			if err != nil {
				resp.Diagnostics.AddError("Error waiting powering on VM", err.Error())
				return
			}
		}

		if allStructsPlan.Settings.Customization.Attributes()["force"].(types.Bool).ValueBool() {
			if vmStatus != poweredOFF {
				task, err := r.vm.Undeploy()
				if err != nil {
					resp.Diagnostics.AddError("Error undeploying VM", err.Error())
					return
				}

				err = task.WaitTaskCompletion()
				if err != nil {
					resp.Diagnostics.AddError("Error waiting undeploying VM", err.Error())
					return
				}
			}

			if err := r.vm.PowerOnAndForceCustomization(); err != nil {
				resp.Diagnostics.AddError("Error powering on VM", err.Error())
				return
			}
		}
	} else if !allStructsPlan.State.PowerON.ValueBool() && vmStatus != poweredOFF {
		task, err := r.vm.Undeploy()
		if err != nil {
			resp.Diagnostics.AddError("Error undeploying VM", err.Error())
			return
		}

		err = task.WaitTaskCompletion()
		if err != nil {
			resp.Diagnostics.AddError("Error waiting undeploying VM", err.Error())
			return
		}
	}

	newPlan, d := r.read(ctx, state)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, newPlan)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *vmResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	state := &vm.VMResourceModel{}

	// Get current state
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Init the resource
	resp.Diagnostics.Append(r.Init(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	/*
		Implement the resource deletion here
	*/

	var d diag.Diagnostics

	r.vm, d = vm.Init(r.client, r.vapp, vm.GetVMOpts{
		ID:   state.ID,
		Name: types.StringNull(),
	})
	if d.HasError() {
		if d.Contains(diag.NewErrorDiagnostic("VM not found", govcd.ErrorEntityNotFound.Error())) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.Append(d...)
		return
	}

	// Check if all disks are detached
	for _, disk := range r.vm.GetDiskSettings() {
		// Disk detachable
		if disk.Disk != nil && disk.Disk.Name != "" {
			resp.Diagnostics.AddError("One or more disks are not detached", "Detach all additional disks before deleting the VM")
			return
		}
	}

	deployed, err := r.vm.IsDeployed()
	if err != nil {
		resp.Diagnostics.AddError("Error getting VM deploy status", err.Error())
		return
	}

	if deployed {
		task, err := r.vm.Undeploy()
		if err != nil {
			resp.Diagnostics.AddError("Error undeploying VM", err.Error())
			return
		}

		err = task.WaitTaskCompletion()
		if err != nil {
			resp.Diagnostics.AddError("Error waiting for undeploy", err.Error())
			return
		}
	}

	err = r.vapp.RemoveVM(*r.vm.VM.VM)
	if err != nil {
		resp.Diagnostics.AddError("Error removing VM", err.Error())
	}
}

func (r *vmResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *vmResource) createVMWithTemplate(ctx context.Context, rm vm.VMResourceModel) (vmCreated vm.VM, diags diag.Diagnostics) {
	var (
		err             error
		vmComputePolicy *govcdtypes.ComputePolicy
	)

	// * Resource
	resource, d := rm.ResourceFromPlan(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return vm.VM{}, diags
	}

	// * Network
	resourceNetworks, d := resource.NetworksFromPlan(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return vm.VM{}, diags
	}

	networkConnection := []vm.NetworkConnection{}
	for _, n := range *resourceNetworks {
		networkConnection = append(networkConnection, n.ConvertToNetworkConnection())
	}

	networkConfig, err := vm.ConstructNetworksConnectionWithoutVM(r.vapp, networkConnection)
	if err != nil {
		diags.AddError("Error retrieving network config", err.Error())
		return vm.VM{}, diags
	}

	// * DeployOS
	deployOS, d := rm.DeployOSFromPlan(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return vm.VM{}, diags
	}

	var (
		vappTemplate   *govcd.VAppTemplate
		storageProfile *govcdtypes.Reference
	)

	if !deployOS.VMNameInTemplate.IsNull() {
		vappTemplate, err = r.client.GetTemplateWithVMName(deployOS.VappTemplateID.ValueString(), deployOS.VMNameInTemplate.ValueString())
		if err != nil {
			diags.AddError("Error retrieving vAppTemplate", err.Error())
			return vm.VM{}, diags
		}
	} else {
		vappTemplate, err = r.client.GetTemplate(deployOS.VappTemplateID.ValueString())
		if err != nil {
			diags.AddError("Error retrieving vAppTemplate", err.Error())
			return vm.VM{}, diags
		}
	}

	// * Settings
	settings, d := rm.SettingsFromPlan(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return vm.VM{}, diags
	}

	// * * Compute Policy
	if !settings.AffinityRuleID.IsUnknown() && !settings.AffinityRuleID.IsNull() {
		affinityRule, err := r.client.GetAffinityRule(settings.AffinityRuleID.ValueString())
		if err != nil {
			diags.AddError("Error retrieving affinity rule", err.Error())
			return vm.VM{}, diags
		}

		if affinityRule != nil {
			vmComputePolicy.VmPlacementPolicy = &govcdtypes.Reference{HREF: affinityRule.Href}
		}
	}

	// * * Compute StorageProfile
	if !settings.StorageProfile.IsUnknown() && !settings.StorageProfile.IsNull() {
		storageProfile, err = r.vdc.GetStorageProfileReference(settings.StorageProfile.ValueString(), false)
		if err != nil {
			diags.AddError("Error retrieving storage profile", err.Error())
			return vm.VM{}, diags
		}
	}

	// * Construct GoVDC VM Object
	vmFromTemplateParams := &govcdtypes.ReComposeVAppParams{
		Ovf:              govcdtypes.XMLNamespaceOVF,
		Xsi:              govcdtypes.XMLNamespaceXSI,
		Xmlns:            govcdtypes.XMLNamespaceVCloud,
		AllEULAsAccepted: deployOS.AcceptAllEulas.ValueBool(),
		Name:             r.vapp.GetName(),
		PowerOn:          false, // VM will be powered on after all configuration is done
		SourcedItem: &govcdtypes.SourcedCompositionItemParam{
			Source: &govcdtypes.Reference{
				HREF: vappTemplate.VAppTemplate.HREF,
				Name: rm.Name.ValueString(), // This VM name defines the VM name after creation
			},
			VMGeneralParams: &govcdtypes.VMGeneralParams{
				Description: rm.Description.ValueString(),
			},
			InstantiationParams: &govcdtypes.InstantiationParams{
				// If a MAC address is specified for NIC - it does not get set with this call,
				// therefore an additional `vm.UpdateNetworkConnectionSection` is required.
				NetworkConnectionSection: &networkConfig,
			},
			ComputePolicy:  vmComputePolicy,
			StorageProfile: storageProfile,
		},
	}

	// * Create VM
	x, err := r.vapp.AddRawVM(vmFromTemplateParams)
	if err != nil {
		diags.AddError("Error creating VM", err.Error())
		return vm.VM{}, diags
	}

	// * Get VM
	w := vm.ConstructObject(r.vapp, x)

	return w, diags
}

func (r *vmResource) createVMWithBootImage(ctx context.Context, rm vm.VMResourceModel) (vmCreated vm.VM, diags diag.Diagnostics) {
	var (
		err             error
		vmComputePolicy *govcdtypes.ComputePolicy
		storageProfile  *govcdtypes.Reference
	)

	// * DeployOS
	deployOS, d := rm.DeployOSFromPlan(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return vm.VM{}, diags
	}

	// * Resource
	resource, d := rm.ResourceFromPlan(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return vm.VM{}, diags
	}

	// * Settings
	settings, d := rm.SettingsFromPlan(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return vm.VM{}, diags
	}

	// * * Compute Policy
	if !settings.AffinityRuleID.IsUnknown() && !settings.AffinityRuleID.IsNull() {
		affinityRule, err := r.client.GetAffinityRule(settings.AffinityRuleID.ValueString())
		if err != nil {
			diags.AddError("Error retrieving affinity rule", err.Error())
			return vm.VM{}, diags
		}

		if affinityRule != nil {
			vmComputePolicy.VmPlacementPolicy = &govcdtypes.Reference{HREF: affinityRule.Href}
		}
	}

	// * * Compute StorageProfile
	if !settings.StorageProfile.IsUnknown() && !settings.StorageProfile.IsNull() {
		storageProfile, err = r.vdc.GetStorageProfileReference(settings.StorageProfile.ValueString(), false)
		if err != nil {
			diags.AddError("Error retrieving storage profile", err.Error())
			return vm.VM{}, diags
		}
	}

	// * VirtualCPU Type
	virtualCPUType := "VM64"
	if strings.Contains(settings.OsType.ValueString(), "32") {
		virtualCPUType = "VM32"
	}

	bootImage, err := r.client.GetBootImage(deployOS.BootImageID.ValueString())
	if err != nil {
		diags.AddError("Error retrieving boot image", err.Error())
		return vm.VM{}, diags
	}

	// * Customization
	customization, d := settings.CustomizationFromPlan(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return vm.VM{}, diags
	}

	// * Construct GoVDC VM Object
	vmParams := &govcdtypes.RecomposeVAppParamsForEmptyVm{
		XmlnsVcloud: govcdtypes.XMLNamespaceVCloud,
		XmlnsOvf:    govcdtypes.XMLNamespaceOVF,
		PowerOn:     false, // Power on is handled at the end of VM creation process
		CreateItem: &govcdtypes.CreateItem{
			Name: rm.Name.ValueString(),
			// bug in VCD, do not allow empty NetworkConnectionSection, so we pass simplest
			// network configuration and after VM created update with real config
			NetworkConnectionSection: &govcdtypes.NetworkConnectionSection{
				PrimaryNetworkConnectionIndex: 0,
				NetworkConnection: []*govcdtypes.NetworkConnection{
					{Network: "none", NetworkConnectionIndex: 0, IPAddress: "any", IsConnected: false, IPAddressAllocationMode: "NONE"},
				},
			},
			StorageProfile:            storageProfile,
			ComputePolicy:             vmComputePolicy,
			Description:               rm.Description.ValueString(),
			GuestCustomizationSection: customization.GetCustomizationSection(rm.Name.ValueString()),
			VmSpecSection: &govcdtypes.VmSpecSection{
				Modified:          utils.TakeBoolPointer(true),
				Info:              "Virtual Machine specification",
				OsType:            settings.OsType.ValueString(),
				CpuResourceMhz:    &govcdtypes.CpuResourceMhz{Configured: 0},
				NumCpus:           utils.TakeIntPointer(int(resource.CPUs.ValueInt64())),
				NumCoresPerSocket: utils.TakeIntPointer(int(resource.CPUsCores.ValueInt64())),
				MemoryResourceMb:  &govcdtypes.MemoryResourceMb{Configured: resource.Memory.ValueInt64()},
				// can be created with resource internal_disk
				DiskSection:     &govcdtypes.DiskSection{DiskSettings: []*govcdtypes.DiskSettings{}},
				VirtualCpuType:  virtualCPUType,
				HardwareVersion: &govcdtypes.HardwareVersion{Value: "vmx-19"},
			},
			BootImage: bootImage,
		},
	}

	x, err := r.vapp.AddEmptyVm(vmParams)
	if err != nil {
		diags.AddError("Error creating VM", err.Error())
		return vm.VM{}, diags
	}

	// * Get VM
	return vm.Get(r.vapp, vm.GetVMOpts{
		Name: rm.Name,
		ID:   types.StringValue(x.VM.ID),
	})
}

// processAfterCreate is a common function for VM creation. It is called after VM is created and
// it is responsible for:
// * Updating VM network configuration
// * Update OS Type
// * Update CPU/Memory Hot Add
// * Update CPU and Memory
// * Update Guest Properties
// * Update Customization.
func (r *vmResource) processAfterCreate(ctx context.Context, vmCreated vm.VM, rm vm.VMResourceModel) (vmUpdated vm.VM, diags diag.Diagnostics) {
	var err error

	// * Resource
	resource, d := rm.ResourceFromPlan(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return
	}

	// * Network
	resourceNetworks, d := resource.NetworksFromPlan(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return
	}

	networkConnection := []vm.NetworkConnection{}
	for _, n := range *resourceNetworks {
		networkConnection = append(networkConnection, n.ConvertToNetworkConnection())
	}

	networkConfig, err := vm.ConstructNetworksConnectionWithoutVM(r.vapp, networkConnection)
	if err != nil {
		diags.AddError("Error retrieving network config", err.Error())
		return
	}

	// * Settings
	settings, d := rm.SettingsFromPlan(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return
	}

	// * Expose Hardware Virtualization
	err = vmCreated.SetExposeHardwareVirtualization(settings.ExposeHardwareVirtualization.ValueBool())
	if err != nil {
		diags.AddError("Error updating expose hardware virtualization", err.Error())
		return
	}

	// * Guest Properties
	guestProperties := make(map[string]string, 0)
	for key, value := range settings.GuestProperties.Elements() {
		guestProperties[key] = strings.Trim(value.String(), "\"")
	}

	if err = vmCreated.SetGuestProperties(guestProperties); err != nil {
		diags.AddError("Error updating guest properties", err.Error())
		return
	}

	// * Os Type
	if err = vmCreated.SetOSType(settings.OsType.ValueString()); err != nil {
		diags.AddError("Error updating OS Type", err.Error())
		return
	}

	// * Update CPU and Memory
	// ? CPU
	if !resource.CPUs.IsNull() && !resource.Memory.IsNull() {
		if err = vmCreated.ChangeCPUAndCoreCount(
			utils.TakeIntPointer(int(resource.CPUs.ValueInt64())),
			utils.TakeIntPointer(int(resource.CPUsCores.ValueInt64())),
		); err != nil {
			diags.AddError("Error updating CPU and Memory", err.Error())
			return
		}
	}

	// ? Memory
	if !resource.Memory.IsNull() {
		if err = vmCreated.ChangeMemory(resource.Memory.ValueInt64()); err != nil {
			diags.AddError("Error updating Memory", err.Error())
			return
		}
	}

	// * Network
	if err = vmCreated.UpdateNetworkConnectionSection(&networkConfig); err != nil {
		diags.AddError("Error updating network config", err.Error())
		return
	}

	// * Customization
	customization, d := settings.CustomizationFromPlan(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return
	}

	if err = vmCreated.SetCustomization(customization.GetCustomizationSection(rm.Name.ValueString())); err != nil {
		diags.AddError("Error updating customization", err.Error())
		return
	}

	// * CPU/Memory Hot Add
	if _, err = vmCreated.UpdateVmCpuAndMemoryHotAdd(resource.CPUHotAddEnabled.ValueBool(), resource.MemoryHotAddEnabled.ValueBool()); err != nil {
		diags.AddError("Error updating CPU/Memory Hot Add", err.Error())
		return
	}

	err = vmCreated.Refresh()
	if err != nil {
		diags.AddError("Error refreshing VM", err.Error())
		return
	}

	return vmCreated, diags
}

// vmPowerOn is a common function for VM power ON. It is called after VM is created.
func (r *vmResource) vmPowerOn(ctx context.Context, rm vm.VMResourceModel) (diags diag.Diagnostics) {
	// * State
	state, d := rm.StateFromPlan(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return
	}

	// * Settings
	settings, d := rm.SettingsFromPlan(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return
	}

	// * Customization
	customization, d := settings.CustomizationFromPlan(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return
	}

	// * Power On
	if state.PowerON.ValueBool() {
		if customization.Force.ValueBool() {
			if err := r.vm.PowerOnAndForceCustomization(); err != nil {
				diags.AddError("Error powering on VM", err.Error())
				return
			}
		} else {
			task, err := r.vm.PowerOn()
			if err != nil {
				diags.AddError("Error powering on VM", err.Error())
				return
			}
			if err = task.WaitTaskCompletion(); err != nil {
				diags.AddError("error waiting for power on", err.Error())
				return
			}
		}
	}

	return diags
}

// read is a common function for VM read. It is called in Update and Read.
func (r *vmResource) read(ctx context.Context, rm *vm.VMResourceModel) (plan *vm.VMResourceModel, diags diag.Diagnostics) {
	if err := r.vm.Refresh(); err != nil {
		diags.AddError("Error refreshing VM", err.Error())
		return
	}

	// ? deployOS -> Use state for unknown value

	// ? State
	stateStruct, err := r.vm.StateRead(ctx)
	if err != nil {
		diags.AddError(
			"Unable to get VM state",
			fmt.Sprintf("Error: %s", err),
		)
		return
	}

	// ? Resource
	networks, err := r.vm.NetworksRead()
	if err != nil {
		diags.AddError(
			"Unable to get VM networks",
			fmt.Sprintf("Error: %s", err),
		)
		return
	}

	// ? Settings
	settings, err := r.vm.SettingsRead(ctx, rm.Settings.Attributes()["customization"])
	if err != nil {
		diags.AddError(
			"Unable to get VM settings",
			fmt.Sprintf("Error: %s", err),
		)
		return
	}

	return &vm.VMResourceModel{
		ID:          types.StringValue(r.vm.GetID()),
		VDC:         types.StringValue(r.vdc.GetName()),
		Name:        types.StringValue(r.vm.GetName()),
		VappID:      types.StringValue(r.vapp.GetID()),
		VappName:    types.StringValue(r.vapp.GetName()),
		Description: rm.Description,
		State:       stateStruct.ToPlan(ctx),
		Resource:    r.vm.ResourceRead(ctx).ToPlan(ctx, networks),
		Settings:    settings.ToPlan(ctx),
		DeployOS:    rm.DeployOS,
	}, nil
}
