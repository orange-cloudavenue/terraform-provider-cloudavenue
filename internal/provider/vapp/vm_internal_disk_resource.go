// Package vapp provides a Terraform resource.
package vapp

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/boolpm"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vapp"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vm"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &internalDiskResource{}
	_ resource.ResourceWithConfigure   = &internalDiskResource{}
	_ resource.ResourceWithImportState = &internalDiskResource{}
)

// NewInternalDiskResource is a helper function to simplify the provider implementation.
func NewInternalDiskResource() resource.Resource {
	return &internalDiskResource{}
}

// internalDiskResource is the resource implementation.
type internalDiskResource struct {
	client *client.CloudAvenue
}

type internalDiskResourceModel struct {
	*vm.InternalDiskModel `tfsdk:"internal_disk"`
	ID                    types.String `tfsdk:"id"`
	VDC                   types.String `tfsdk:"vdc"`
	VAppName              types.String `tfsdk:"vapp_name"`
	VMName                types.String `tfsdk:"vm_name"`
	AllowVMReboot         types.Bool   `tfsdk:"allow_vm_reboot"`
}

// Metadata returns the resource type name.
func (r *internalDiskResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vapp_vm_internal_disk"
}

// Schema defines the schema for the resource.
func (r *internalDiskResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The vm_internal_disk resource allows you to manage an internal disk of a VM in a vApp.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID of the internal disk.",
			},
			"vdc": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The name of VDC to use, optional if defined at provider level.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"vapp_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				MarkdownDescription: "The vApp this VM internal disk belongs to.",
			},
			"vm_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				MarkdownDescription: "VM in vApp in which internal disk is created.",
			},
			"allow_vm_reboot": schema.BoolAttribute{
				Computed:            true,
				Optional:            true,
				PlanModifiers:       []planmodifier.Bool{boolpm.SetDefault(false)},
				MarkdownDescription: "Powers off VM when changing any attribute of an IDE disk or unit/bus number of other disk types, after the change is complete VM is powered back on. Without this setting enabled, such changes on a powered-on VM would fail.",
			},
		},
		Blocks: map[string]schema.Block{
			"internal_disk": schema.SingleNestedBlock{
				MarkdownDescription: "The internal disk configuration. See [internal_disk](#internal_disk) below for details.",
				Attributes:          vm.InternalDiskSchema(),
			},
		},
	}
}

func (r *internalDiskResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *internalDiskResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var (
		plan *internalDiskResourceModel
	)

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// If VDC is not defined at data source level, use the one defined at provider level
	if plan.VDC.IsNull() || plan.VDC.IsUnknown() {
		if r.client.DefaultVDCExist() {
			plan.VDC = types.StringValue(r.client.GetDefaultVDC())
		} else {
			resp.Diagnostics.AddError("Missing VDC", "VDC is required when not defined at provider level")
			return
		}
	}

	_, vdc, err := r.client.GetOrgAndVDC(r.client.GetOrg(), plan.VDC.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving VDC", err.Error())
		return
	}

	myVM, err := vm.GetVM(vdc, plan.VAppName.ValueString(), plan.VMName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving VM", err.Error())
		return
	}

	myVApp := vapp.Ref{
		Name:  plan.VAppName.ValueString(),
		Org:   r.client.GetOrg(),
		VDC:   plan.VDC.ValueString(),
		TFCtx: ctx,
	}

	if errLock := myVApp.LockParentVApp(); errors.Is(errLock, vapp.ErrVAppRefEmpty) {
		resp.Diagnostics.AddError("Error locking vApp", "Empty name, org or vdc in vapp.VAppRef")
		return
	}
	defer func() {
		if errUnlock := myVApp.UnLockParentVApp(); errUnlock != nil {
			// tflog print error is enough ?
			resp.Diagnostics.AddError("Error unlocking vApp", errUnlock.Error())
			return
		}
	}()

	// storage profile
	var storageProfilePrt *govcdtypes.Reference
	var overrideVMDefault bool
	if plan.StorageProfile.IsNull() || plan.StorageProfile.IsUnknown() {
		storageProfilePrt = myVM.VM.StorageProfile
		overrideVMDefault = false
	} else {
		storageProfile, errFindStorage := vdc.FindStorageProfileReference(plan.StorageProfile.ValueString())
		if errFindStorage != nil {
			resp.Diagnostics.AddError("Error retrieving storage profile", errFindStorage.Error())
			return
		}
		storageProfilePrt = &storageProfile
		overrideVMDefault = true
	}

	// value is required but not treated.
	isThinProvisioned := true

	diskSetting := &govcdtypes.DiskSettings{
		SizeMb:              plan.SizeInMb.ValueInt64(),
		UnitNumber:          int(plan.UnitNumber.ValueInt64()),
		BusNumber:           int(plan.BusNumber.ValueInt64()),
		AdapterType:         vm.InternalDiskBusTypes[plan.BusType.ValueString()],
		ThinProvisioned:     &isThinProvisioned,
		StorageProfile:      storageProfilePrt,
		VirtualQuantityUnit: "byte",
		OverrideVmDefault:   overrideVMDefault,
	}

	vmStatusBefore, err := vm.PowerOffIfNeeded(myVM, plan.BusType.ValueString(), plan.AllowVMReboot.ValueBool())
	if err != nil {
		resp.Diagnostics.AddError("Error powering off VM", err.Error())
	}

	diskID, err := myVM.AddInternalDisk(diskSetting)
	if err != nil {
		resp.Diagnostics.AddError("Error creating disk", err.Error())
		return
	}
	err = vm.PowerOnIfNeeded(myVM, plan.BusType.ValueString(), plan.AllowVMReboot.ValueBool(), vmStatusBefore)
	if err != nil {
		resp.Diagnostics.AddError("Error powering on VM", err.Error())
	}

	plan = &internalDiskResourceModel{
		InternalDiskModel: &vm.InternalDiskModel{
			ID:             types.StringValue(diskID),
			BusType:        plan.BusType,
			SizeInMb:       plan.SizeInMb,
			BusNumber:      plan.BusNumber,
			UnitNumber:     plan.UnitNumber,
			StorageProfile: types.StringValue(storageProfilePrt.Name),
		},
		ID:            types.StringValue(diskID),
		VDC:           plan.VDC,
		VAppName:      plan.VAppName,
		VMName:        plan.VMName,
		AllowVMReboot: plan.AllowVMReboot,
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *internalDiskResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *internalDiskResourceModel

	// Get current state
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If VDC is not defined at data source level, use the one defined at provider level
	if state.VDC.IsNull() || state.VDC.IsUnknown() {
		if r.client.DefaultVDCExist() {
			state.VDC = types.StringValue(r.client.GetDefaultVDC())
		} else {
			resp.Diagnostics.AddError("Missing VDC", "VDC is required when not defined at provider level")
			return
		}
	}

	_, vdc, err := r.client.GetOrgAndVDC(r.client.GetOrg(), state.VDC.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving VDC", err.Error())
		return
	}

	myVM, err := vm.GetVM(vdc, state.VAppName.ValueString(), state.VMName.ValueString())
	if err != nil {
		if govcd.ContainsNotFound(err) {
			tflog.Debug(ctx, "VM not found, removing resource from state")
			// VM not found, remove from state
			resp.State.RemoveResource(ctx)
		}
		resp.Diagnostics.AddError(fmt.Sprintf("Error retrieving VM. VM : %s, VApp : %s", state.VAppName.ValueString(), state.VMName.ValueString()), err.Error())
		return
	}

	diskSettings, err := myVM.GetInternalDiskById(state.ID.ValueString(), true)
	if err != nil {
		if err == govcd.ErrorEntityNotFound {
			tflog.Debug(ctx, "Disk not found, removing resource from state")
			resp.State.RemoveResource(ctx)
		}
		resp.Diagnostics.AddError("Error retrieving disk with id "+state.ID.ValueString(), err.Error())
		return
	}

	plan := &internalDiskResourceModel{
		InternalDiskModel: &vm.InternalDiskModel{
			ID:             types.StringValue(diskSettings.DiskId),
			BusType:        types.StringValue(vm.InternalDiskBusTypesFromValues[strings.ToLower(diskSettings.AdapterType)]),
			SizeInMb:       types.Int64Value(diskSettings.SizeMb),
			BusNumber:      types.Int64Value(int64(diskSettings.BusNumber)),
			UnitNumber:     types.Int64Value(int64(diskSettings.UnitNumber)),
			StorageProfile: types.StringValue(diskSettings.StorageProfile.Name),
		},
		ID:            types.StringValue(diskSettings.DiskId),
		VDC:           state.VDC,
		VAppName:      state.VAppName,
		VMName:        state.VMName,
		AllowVMReboot: state.AllowVMReboot,
	}

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *internalDiskResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state *internalDiskResourceModel

	// Get current state
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If VDC is not defined at data source level, use the one defined at provider level
	if plan.VDC.IsNull() || plan.VDC.IsUnknown() {
		if r.client.DefaultVDCExist() {
			plan.VDC = types.StringValue(r.client.GetDefaultVDC())
		} else {
			resp.Diagnostics.AddError("Missing VDC", "VDC is required when not defined at provider level")
			return
		}
	}

	_, vdc, err := r.client.GetOrgAndVDC(r.client.GetOrg(), plan.VDC.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving VDC", err.Error())
		return
	}

	myVM, err := vm.GetVM(vdc, plan.VAppName.ValueString(), plan.VMName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving VM", err.Error())
		return
	}

	myVApp := vapp.Ref{
		Name:  plan.VAppName.ValueString(),
		Org:   r.client.GetOrg(),
		VDC:   plan.VDC.ValueString(),
		TFCtx: ctx,
	}

	if errLock := myVApp.LockParentVApp(); errors.Is(errLock, vapp.ErrVAppRefEmpty) {
		resp.Diagnostics.AddError("Error locking vApp", "Empty name, org or vdc in vapp.VAppRef")
		return
	}
	defer func() {
		if errUnlock := myVApp.UnLockParentVApp(); errUnlock != nil {
			// tflog print error is enough ?
			resp.Diagnostics.AddError("Error unlocking vApp", errUnlock.Error())
			return
		}
	}()

	vmStatusBefore, err := vm.PowerOffIfNeeded(myVM, plan.BusType.ValueString(), plan.AllowVMReboot.ValueBool())
	if err != nil {
		resp.Diagnostics.AddError("Error powering off VM", err.Error())
	}

	diskSettingsToUpdate, err := myVM.GetInternalDiskById(state.ID.ValueString(), false)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving disk", err.Error())
		return
	}

	infos := make(map[string]interface{})
	infos["id"] = state.ID.ValueString()
	tflog.Trace(ctx, "Internal disk found", infos)

	diskSettingsToUpdate.SizeMb = plan.SizeInMb.ValueInt64()
	// Note can't change adapter type, bus number, unit number as vSphere changes diskId

	var storageProfilePrt *govcdtypes.Reference
	var overrideVMDefault bool

	storageProfileName := plan.StorageProfile.ValueString()
	if storageProfileName != "" {
		storageProfile, errFindStorage := vdc.FindStorageProfileReference(storageProfileName)
		if errFindStorage != nil {
			resp.Diagnostics.AddError("Error retrieving storage profile", errFindStorage.Error())
		}
		storageProfilePrt = &storageProfile
		overrideVMDefault = true
	} else {
		storageProfilePrt = myVM.VM.StorageProfile
		overrideVMDefault = false
	}

	diskSettingsToUpdate.StorageProfile = storageProfilePrt
	diskSettingsToUpdate.OverrideVmDefault = overrideVMDefault

	_, err = myVM.UpdateInternalDisks(myVM.VM.VmSpecSection)
	if err != nil {
		resp.Diagnostics.AddError("Error updating disk", err.Error())
		return
	}

	err = vm.PowerOnIfNeeded(myVM, plan.BusType.ValueString(), plan.AllowVMReboot.ValueBool(), vmStatusBefore)
	if err != nil {
		resp.Diagnostics.AddError("Error powering on VM", err.Error())
	}

	plan = &internalDiskResourceModel{
		InternalDiskModel: &vm.InternalDiskModel{
			ID:             state.ID,
			BusType:        plan.BusType,
			SizeInMb:       types.Int64Value(diskSettingsToUpdate.SizeMb),
			BusNumber:      plan.BusNumber,
			UnitNumber:     plan.UnitNumber,
			StorageProfile: types.StringValue(diskSettingsToUpdate.StorageProfile.Name),
		},
		ID:            state.ID,
		VDC:           plan.VDC,
		VAppName:      plan.VAppName,
		VMName:        plan.VMName,
		AllowVMReboot: plan.AllowVMReboot,
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *internalDiskResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *internalDiskResourceModel
	// Get current state
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// If VDC is not defined at data source level, use the one defined at provider level
	if state.VDC.IsNull() || state.VDC.IsUnknown() {
		if r.client.DefaultVDCExist() {
			state.VDC = types.StringValue(r.client.GetDefaultVDC())
		} else {
			resp.Diagnostics.AddError("Missing VDC", "VDC is required when not defined at provider level")
			return
		}
	}

	_, vdc, err := r.client.GetOrgAndVDC(r.client.GetOrg(), state.VDC.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving VDC", err.Error())
		return
	}

	myVM, err := vm.GetVM(vdc, state.VAppName.ValueString(), state.VMName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving VM", err.Error())
		return
	}

	myVApp := vapp.Ref{
		Name:  state.VAppName.ValueString(),
		Org:   r.client.GetOrg(),
		VDC:   state.VDC.ValueString(),
		TFCtx: ctx,
	}

	if errLock := myVApp.LockParentVApp(); errors.Is(errLock, vapp.ErrVAppRefEmpty) {
		resp.Diagnostics.AddError("Error locking vApp", "Empty name, org or vdc in vapp.VAppRef")
		return
	}
	defer func() {
		if errUnlock := myVApp.UnLockParentVApp(); errUnlock != nil {
			// tflog print error is enough ?
			resp.Diagnostics.AddError("Error unlocking vApp", errUnlock.Error())
			return
		}
	}()

	vmStatusBefore, err := vm.PowerOffIfNeeded(myVM, state.BusType.ValueString(), state.AllowVMReboot.ValueBool())
	if err != nil {
		resp.Diagnostics.AddError("Error powering off VM", err.Error())
	}

	errDelete := myVM.DeleteInternalDisk(state.ID.ValueString())
	if errDelete != nil {
		resp.Diagnostics.AddError("Error deleting disk", errDelete.Error())
		return
	}
	err = vm.PowerOnIfNeeded(myVM, state.BusType.ValueString(), state.AllowVMReboot.ValueBool(), vmStatusBefore)
	if err != nil {
		resp.Diagnostics.AddError("Error powering on VM", err.Error())
	}

	infos := make(map[string]interface{})
	infos["id"] = state.ID.ValueString()
	tflog.Trace(ctx, "Disk deleted", infos)
}

func (r *internalDiskResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var state *internalDiskResourceModel
	resourceURI := strings.Split(req.ID, ".")

	if len(resourceURI) != 4 && len(resourceURI) != 3 {
		resp.Diagnostics.AddError("Error importing disk", "Wrong resource URI format. Expected vdc.vapp.vm.disk_id or va.vapp.vm.disk_id")
		return
	}

	state = &internalDiskResourceModel{
		InternalDiskModel: &vm.InternalDiskModel{
			ID: types.StringValue(resourceURI[2]),
		},
		ID:       types.StringValue(resourceURI[2]),
		VAppName: types.StringValue(resourceURI[0]),
		VMName:   types.StringValue(resourceURI[1]),
	}

	if len(resourceURI) == 4 {
		state = &internalDiskResourceModel{
			InternalDiskModel: &vm.InternalDiskModel{
				ID: types.StringValue(resourceURI[3]),
			},
			ID:       types.StringValue(resourceURI[3]),
			VDC:      types.StringValue(resourceURI[0]),
			VAppName: types.StringValue(resourceURI[1]),
			VMName:   types.StringValue(resourceURI[2]),
		}
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
