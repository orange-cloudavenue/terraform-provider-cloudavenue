// Package vm provides a Terraform resource.
package vm

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/org"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vapp"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vdc"
)

// Ensure the implementation satisfies the expected interfaces.VAppName.
var (
	_ resource.Resource              = &vmInsertedMediaResource{}
	_ resource.ResourceWithConfigure = &vmInsertedMediaResource{}
)

// NewVMInsertedMediaResource is a helper function to simplify the provider implementation.
func NewVMInsertedMediaResource() resource.Resource {
	return &vmInsertedMediaResource{}
}

// Metadata returns the resource type name.
func (r *vmInsertedMediaResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_" + "inserted_media"
}

// Schema defines the schema for the resource.
func (r *vmInsertedMediaResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = vmInsertedMediaSchema()
}

func (r *vmInsertedMediaResource) Init(ctx context.Context, rm *vmInsertedMediaResourceModel) (diags diag.Diagnostics) {
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

func (r *vmInsertedMediaResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *vmInsertedMediaResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var (
		plan *vmInsertedMediaResourceModel
		err  error
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
	resp.Diagnostics.Append(r.vapp.LockVAPP(ctx)...)
	if resp.Diagnostics.HasError() {
		return
	}
	defer r.vapp.UnlockVAPP(ctx)

	// Check if VM exists
	vm, err := r.vapp.GetVMByName(plan.VMName.ValueString(), true)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving VM", err.Error())
		return
	}

	// Insert media
	task, err := vm.HandleInsertMedia(r.org.Org.Org, plan.Catalog.ValueString(), plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error inserting media", err.Error())
		return
	}
	err = task.WaitTaskCompletion()
	if err != nil {
		resp.Diagnostics.AddError("Error during inserting media", err.Error())
		return
	}

	// Set Plan state
	plan = &vmInsertedMediaResourceModel{
		ID:       types.StringValue(vm.VM.ID),
		VDC:      types.StringValue(r.vdc.GetName()),
		Catalog:  plan.Catalog,
		Name:     plan.Name,
		VAppName: plan.VAppName,
		VAppID:   plan.VAppID,
		VMName:   plan.VMName,
		// EjectForce: plan.EjectForce,
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *vmInsertedMediaResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *vmInsertedMediaResourceModel

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

	// Check if VM exists
	vm, err := r.vapp.GetVMByName(state.VMName.ValueString(), true)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving VM", err.Error())
		return
	}

	// Check if media is mounted
	var isIsoMounted bool

	for _, hardwareItem := range vm.VM.VirtualHardwareSection.Item {
		if hardwareItem.ResourceType == int(15) { // 15 = CD/DVD Drive
			isIsoMounted = true
			break
		}
	}
	if !isIsoMounted {
		resp.Diagnostics.AddError("Media not mounted", "Media is not mounted on the VM")
		resp.State.RemoveResource(ctx)
		return
	}

	// Set Plan state
	plan := &vmInsertedMediaResourceModel{
		ID:       types.StringValue(vm.VM.ID),
		VDC:      types.StringValue(r.vdc.GetName()),
		Catalog:  state.Catalog,
		Name:     state.Name,
		VAppName: state.VAppName,
		VAppID:   state.VAppID,
		VMName:   state.VMName,
		// EjectForce: state.EjectForce,
	}

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *vmInsertedMediaResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	/* linked with issue - Disable attributes - Issue referrer: vmware/go-vcloud-director#552
	var plan, state *vmInsertedMediaResourceModel

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

	// Check if VM exists
	vm, err := r.vapp.GetVMByName(state.VMName.ValueString(), true)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving VM", err.Error())
		return
	}

	plan = &vmInsertedMediaResourceModel{
		ID:       types.StringValue(vm.VM.ID),
		VDC:      plan.VDC,
		Catalog:  plan.Catalog,
		Name:     plan.Name,
		VAppName: plan.VAppName,
		VMName:   plan.VMName,
		// EjectForce: plan.EjectForce,
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	*/
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *vmInsertedMediaResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *vmInsertedMediaResourceModel

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

	// Check if VM exists
	vm, err := r.vapp.GetVMByName(state.VMName.ValueString(), true)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving VM", err.Error())
		return
	}

	// Eject media
	_, err = vm.HandleEjectMediaAndAnswer(r.org.Org.Org, state.Catalog.ValueString(), state.Name.ValueString(), true)
	if err != nil {
		resp.Diagnostics.AddError("Error ejecting media", err.Error())
		return
	}
}
