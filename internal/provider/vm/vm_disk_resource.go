// Package vm provides a Terraform resource.
package vm

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/org"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vapp"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vdc"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vm"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &diskResource{}
	_ resource.ResourceWithConfigure = &diskResource{}
)

// NewDiskResource is a helper function to simplify the provider implementation.
func NewDiskResource() resource.Resource {
	return &diskResource{}
}

// diskResource is the resource implementation.
type diskResource struct {
	client *client.CloudAvenue
	vapp   vapp.VAPP
	vdc    vdc.VDC
	org    org.Org
}

type diskResourceModel vm.Disk

// Metadata returns the resource type name.
func (r *diskResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_" + "disk"
}

// Schema defines the schema for the resource.
func (r *diskResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The disk resource allows you to manage a disk in the vDC.",
		Attributes:          vm.DiskSchema(),
	}
}

func (r *diskResource) Init(ctx context.Context, rm *vm.Disk) (diags diag.Diagnostics) {
	r.org, diags = org.Init(r.client)
	if diags.HasError() {
		return
	}

	r.vdc, diags = vdc.Init(r.client, rm.VDC)
	if diags.HasError() {
		return
	}

	r.vapp, diags = vapp.Init(r.client, r.vdc, rm.VAppID, rm.VAppName)

	return
}

func (r *diskResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// ModifyPlan is called before Create, Update, and Delete to modify the plan.
func (r *diskResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var d diag.Diagnostics

	diskPlan := &diskResourceModel{}
	diskState := &diskResourceModel{}

	d = req.State.Get(ctx, diskState)
	if d.HasError() {
		// State is not available, so we can't validate the plan.
		return
	}

	d = req.Plan.Get(ctx, diskPlan)
	if d.HasError() {
		// Plan is not available, so we can't validate the plan.
		return
	}

	if diskPlan.IsDetachable.ValueBool() {
		if !diskPlan.SizeInMb.Equal(diskState.SizeInMb) {
			resp.Diagnostics.AddWarning(
				"Warning detach/attach disk is required",
				"Disk size cannot be changed when disk is detachable. Detach/attach disk is required. \n"+
					"If you apply this change, the disk will be detached and attached again.",
			)
		}
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *diskResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var (
		plan *vm.Disk

		myVM *govcd.VM
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

	plan.VAppName = types.StringValue(r.vapp.GetName())
	plan.VAppID = types.StringValue(r.vapp.GetID())

	// VMName or VMID was emptyString by default, so we need to check if it is emptyString or not
	if plan.VMName.ValueString() == "" && plan.VMID.ValueString() == "" &&
		!plan.IsDetachable.ValueBool() {
		resp.Diagnostics.AddError("Missing VM", "VM is required when disk is not detachable")
		return
	}

	if plan.VMName.ValueString() != "" || plan.VMID.ValueString() != "" {
		// Get VM Object
		var vmByNameOrID types.String
		if plan.VMName.ValueString() != "" {
			vmByNameOrID = plan.VMName
		} else {
			vmByNameOrID = plan.VMID
		}
		myVM, err := r.vapp.GetVMByNameOrId(vmByNameOrID.ValueString(), false)
		if err != nil {
			resp.Diagnostics.AddError("Error retrieving VM", err.Error())
			return
		}
		if myVM.VM == nil {
			resp.Diagnostics.AddError("Error retrieving VM", "VM not found")
			return
		}

		plan.VMName = types.StringValue(myVM.VM.Name)
		plan.VMID = types.StringValue(myVM.VM.ID)
	}

	var newPlan *vm.Disk

	if plan.IsDetachable.ValueBool() {
		// Create a detachable disk
		disk, d := vm.DiskCreate(ctx, r.org, r.vdc, myVM, plan, r.vapp)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}

		newPlan = disk
	} else {
		// Create a disk attached to a VM
		internalDisk, d := vm.InternalDiskCreate(ctx, r.client, vm.InternalDisk{
			ID:             plan.ID,
			BusType:        plan.BusType,
			BusNumber:      plan.BusNumber,
			UnitNumber:     plan.UnitNumber,
			SizeInMb:       plan.SizeInMb,
			StorageProfile: plan.StorageProfile,
		}, plan.VAppName, plan.VMName, plan.VDC)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}

		newPlan = plan
		newPlan.ID = internalDisk.ID
		newPlan.BusType = internalDisk.BusType
		newPlan.BusNumber = internalDisk.BusNumber
		newPlan.UnitNumber = internalDisk.UnitNumber
		newPlan.SizeInMb = internalDisk.SizeInMb
		newPlan.StorageProfile = internalDisk.StorageProfile
	}

	if myVM != nil && myVM.VM != nil {
		newPlan.VMID = types.StringValue(myVM.VM.ID)
		newPlan.VMName = types.StringValue(myVM.VM.Name)
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, &newPlan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *diskResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *vm.Disk

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

	disk, d := vm.DiskRead(ctx, r.client, r.org, r.vdc, state, r.vapp)
	if disk == nil && d != nil {
		// Disk not found, remove from state
		resp.State.RemoveResource(ctx)
	}
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, &disk)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *diskResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state *vm.Disk

	// Get current state
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.VMName.ValueString() == "" && plan.VMID.ValueString() == "" && !plan.IsDetachable.ValueBool() {
		resp.Diagnostics.AddError("Missing VM", "VM is required when disk is not detachable")
		return
	}

	// Init resource
	resp.Diagnostics.Append(r.Init(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	disk, d := vm.DiskUpdate(ctx, r.client, plan, state, r.org, r.vdc, r.vapp)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, &disk)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *diskResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *vm.Disk
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

	resp.Diagnostics.Append(vm.DiskDelete(ctx, r.client, state, r.org, r.vdc, r.vapp)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// func (r *diskResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
// 	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

// }
