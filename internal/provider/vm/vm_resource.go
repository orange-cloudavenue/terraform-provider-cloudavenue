// Package vm provides a Terraform resource.
package vm

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vm"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &vmResource{}
	_ resource.ResourceWithConfigure   = &vmResource{}
	_ resource.ResourceWithImportState = &vmResource{}
)

// NewVMResource is a helper function to simplify the provider implementation.
func NewVMResource() resource.Resource {
	return &vmResource{}
}

// vmResource is the resource implementation.
type vmResource struct {
	client *client.CloudAvenue
}

type vmResourceModel struct {
	ID  types.String `tfsdk:"id"`
	VDC types.String `tfsdk:"vdc"`

	VappName       types.String `tfsdk:"vapp_name"`
	VappTemplateID types.String `tfsdk:"vapp_template_id"`

	VMName           types.String `tfsdk:"vm_name"`
	VMNameInTemplate types.String `tfsdk:"vm_name_in_template"`
	ComputerName     types.String `tfsdk:"computer_name"`

	Resource vm.Resource `tfsdk:"resource"`

	Description    types.String `tfsdk:"description"`
	Href           types.String `tfsdk:"href"`
	AcceptAllEulas types.Bool   `tfsdk:"accept_all_eulas"`

	PowerON               types.Bool `tfsdk:"power_on"`
	PreventUpdatePowerOff types.Bool `tfsdk:"prevent_update_power_off"`

	Disks                 []vm.DiskModel         `tfsdk:"disks"`
	InternalDisks         []vm.InternalDiskModel `tfsdk:"internal_disks"`
	StorageProfile        types.String           `tfsdk:"storage_profile"`
	BootImageID           types.String           `tfsdk:"boot_image_id"`
	OverrideTemplateDisks []vm.TemplateDiskModel `tfsdk:"override_template_disks"`

	OsType types.String `tfsdk:"os_type"`

	Networks               []vm.Network `tfsdk:"networks"`
	NetworkDhcpWaitSeconds types.Int64  `tfsdk:"network_dhcp_wait_seconds"`

	ExposeHardwareVirtualization types.Bool         `tfsdk:"expose_hardware_virtualization"`
	GuestProperties              types.Map          `tfsdk:"guest_properties"`
	Customization                []vm.Customization `tfsdk:"customization"`
	SizingPolicyID               types.String       `tfsdk:"sizing_policy_id"`
	PlacementPolicyID            types.String       `tfsdk:"placement_policy_id"`
	StatusCode                   types.Int64        `tfsdk:"status_code"`
	StatusText                   types.String       `tfsdk:"status_text"`
}

// Metadata returns the resource type name.
func (r *vmResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName
}

// Schema defines the schema for the resource.
func (r *vmResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = vmSchema()
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
	// Retrieve values from plan
	var (
		plan  *vmResourceModel
		err   error
		vm    *govcd.VM
		state *vmResourceModel
	)

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	v := &VMClient{
		Client: r.client,
		Plan:   plan,
		State:  state,
	}

	vm, err = createVM(ctx, v)
	if err != nil {
		resp.Diagnostics.AddError("Error creating VM from template", err.Error())
		return
	}

	if err = vm.Refresh(); err != nil {
		resp.Diagnostics.AddError("Error refreshing VM", err.Error())
		return
	}

	// Handle Guest Properties
	// Such schema fields are processed:
	// * guest_properties
	err = addRemoveGuestProperties(v, vm)
	if err != nil {
		resp.Diagnostics.AddError("Error setting guest properties during creation", err.Error())
		return
	}

	// vm.VM structure contains ProductSection so it needs to be refreshed after
	// `addRemoveGuestProperties`
	if err = vm.Refresh(); err != nil {
		resp.Diagnostics.AddError("Error refreshing VM", err.Error())
		return
	}

	// Handle Guest Customization Section
	// Such schema fields are processed:
	// * customization
	// * computer_name
	// * name
	err = updateGuestCustomizationSetting(v, vm)
	if err != nil {
		resp.Diagnostics.AddError("Error setting guest customization during creation", err.Error())
		return
	}

	// vm.VM structure contains GuestCustomizationSection so it needs to be refreshed after
	// `updateGuestCustomizationSetting`
	if err = vm.Refresh(); err != nil {
		resp.Diagnostics.AddError("Error refreshing VM", err.Error())
		return
	}

	// Explicitly setting CPU and Memory Hot Add settings
	// Note. VM Creation bodies allow specifying these values, but they are ignored therefore using
	// an explicit "/vmCapabilities" API endpoint
	// Such schema fields are processed:
	// * cpu_hot_add_enabled
	// * memory_hot_add_enabled
	_, err = vm.UpdateVmCpuAndMemoryHotAdd(v.Plan.Resource.CPUHotAddEnabled.ValueBool(), v.Plan.Resource.MemoryHotAddEnabled.ValueBool())
	if err != nil {
		resp.Diagnostics.AddError("Error setting VM CPU/Memory HotAdd capabilities", err.Error())
		return
	}

	// Independent disk handling
	// Such schema fields are processed:
	// * disk
	_, vdc, err := r.client.GetOrgAndVDC(r.client.GetOrg(), v.Plan.VDC.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error getting VDC", err.Error())
		return
	}
	err = attachDetachIndependentDisks(v, *vm, vdc)
	if err != nil {
		resp.Diagnostics.AddError("Error attaching-detaching independent disks when creating VM", err.Error())
		return
	}

	////////////////////////////////////////////////////////////////////////////////////////////////
	// VM power on handling is the last step, no other VM adjustment operations should be performed
	// after this
	////////////////////////////////////////////////////////////////////////////////////////////////

	// By default, the VM is created in POWERED_OFF state
	if v.Plan.PowerON.ValueBool() {
		// When customization is requested VM must be un-deployed before starting it
		customizationNeeded := isForcedCustomization(v)
		if customizationNeeded {
			err := vm.PowerOnAndForceCustomization()
			if err != nil {
				resp.Diagnostics.AddError("Error powering on VM with customization", err.Error())
				return
			}
		} else {
			task, err := vm.PowerOn()
			if err != nil {
				resp.Diagnostics.AddError("Error powering on VM", err.Error())
				return
			}
			err = task.WaitTaskCompletion()
			if err != nil {
				resp.Diagnostics.AddError("error waiting for power on", err.Error())
				return
			}
		}
	}

	plan.ID = types.StringValue(vm.VM.ID)
	plan.VMName = types.StringValue(vm.VM.Name)
	plan.ComputerName = types.StringValue(vm.VM.GuestCustomizationSection.ComputerName)
	plan.VDC = types.StringValue(vdc.Vdc.Name)
	plan.Href = types.StringValue(vm.VM.HREF)
	plan.StatusCode = types.Int64Value(int64(vm.VM.Status))

	statusText, err := vm.GetStatus()
	if err != nil {
		statusText = vmUnknownStatus
	}

	plan.StatusText = types.StringValue(statusText)

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *vmResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *vmResourceModel

	// Get current state
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan := &vmResourceModel{}

	v := &VMClient{
		Client: r.client,
		Plan:   plan,
		State:  state,
	}
	_, _ = readVM(v)

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *vmResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state *vmResourceModel

	// Get current state
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	v := &VMClient{
		Client: r.client,
		Plan:   plan,
		State:  state,
	}
	_, _ = updateVM(v)

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *vmResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *vmResourceModel

	// Get current state
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *vmResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
