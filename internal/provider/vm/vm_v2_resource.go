package vm

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/adminvdc"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vapp"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vdc"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vm"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &VMV2Resource{}
	_ resource.ResourceWithConfigure   = &VMV2Resource{}
	_ resource.ResourceWithImportState = &VMV2Resource{}
	// _ resource.ResourceWithModifyPlan     = &VMV2Resource{}
	// _ resource.ResourceWithUpgradeState   = &VMV2Resource{}
	// _ resource.ResourceWithValidateConfig = &VMV2Resource{}.
)

// NewVMV2Resource is a helper function to simplify the provider implementation.
func NewVMV2Resource() resource.Resource {
	return &VMV2Resource{}
}

// VMV2Resource is the resource implementation.
type VMV2Resource struct {
	client   *client.CloudAvenue
	vdc      vdc.VDC
	adminVDC adminvdc.AdminVDC
	vapp     vapp.VAPP
	vm       vm.VM
}

// Init Initializes the resource.
func (r *VMV2Resource) Init(ctx context.Context, rm *VMV2Model) (diags diag.Diagnostics) {
	var d diag.Diagnostics

	r.vdc, d = vdc.Init(r.client, rm.VDC.StringValue)
	diags.Append(d...)
	if diags.HasError() {
		return
	}

	r.adminVDC, d = adminvdc.Init(r.client, rm.VDC.StringValue)
	diags.Append(d...)
	if diags.HasError() {
		return
	}

	r.vapp, d = vapp.Init(r.client, r.vdc, rm.VappID.StringValue, rm.VappName.StringValue)
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
func (r *VMV2Resource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_v2"
}

// Schema defines the schema for the resource.
func (r *VMV2Resource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = vmV2Schema(ctx).GetResource(ctx)
}

func (r *VMV2Resource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *VMV2Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer metrics.New("cloudavenue_vm_v_2", r.client.GetOrgName(), metrics.Create)()

	plan := &VMV2Model{}

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

	// Lock the vapp to prevent concurrent operations
	resp.Diagnostics.Append(r.vapp.LockVAPP(ctx)...)
	if resp.Diagnostics.HasError() {
		return
	}
	defer r.vapp.UnlockVAPP(ctx)

	// Two cases of creation:
	// 1. Create a new VM with VAppTemplate
	// 2. Create a new VM with BootImage

	// Case 1 
	if plan.VAppTemplateID.IsKnown() {
		




	// Use generic read function to refresh the state
	state, _, d := r.read(ctx, plan)
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *VMV2Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer metrics.New("cloudavenue_vm_v_2", r.client.GetOrgName(), metrics.Read)()

	state := &VMV2Model{}

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

	// Refresh the state
	stateRefreshed, found, d := r.read(ctx, state)
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, stateRefreshed)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *VMV2Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer metrics.New("cloudavenue_vm_v_2", r.client.GetOrgName(), metrics.Update)()

	var (
		plan  = &VMV2Model{}
		state = &VMV2Model{}
	)

	// Get current plan and state
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

	// Use generic read function to refresh the state
	stateRefreshed, _, d := r.read(ctx, plan)
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, stateRefreshed)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *VMV2Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer metrics.New("cloudavenue_vm_v_2", r.client.GetOrgName(), metrics.Delete)()

	state := &VMV2Model{}

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
}

func (r *VMV2Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	defer metrics.New("cloudavenue_vm_v_2", r.client.GetOrgName(), metrics.Import)()

	// * Import basic
	// resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	// * Import with custom logic
	// idParts := strings.Split(req.ID, ".")

	// if len(idParts) != 2 {
	// 	resp.Diagnostics.AddError(
	// 		"Unexpected Import Identifier",
	// 		fmt.Sprintf("Expected import identifier with format: xx.xx. Got: %q", req.ID),
	// 	)
	// 	return
	// }

	// resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("xx"), var1)...)
	// resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("xx"), var2)...)
}

// * CustomFuncs

// read is a generic read function that can be used by the resource Create, Read and Update functions.
func (r *VMV2Resource) read(ctx context.Context, planOrState *VMV2Model) (stateRefreshed *VMV2Model, found bool, diags diag.Diagnostics) {
	// TODO : Remove the comment line after you have run the types generator
	// stateRefreshed is commented because the Copy function is not before run the types generator
	// stateRefreshed = planOrState.Copy()

	/*
		Implement the resource read here
	*/

	/* Example

	data, err := r.foo.GetData()
	if err != nil {
		if govcd.ContainsNotFound(err) {
			return nil, false, nil
		}
		diags.AddError("Error retrieving foo", err.Error())
		return nil, true, diags
	}

	if !stateRefreshed.ID.IsKnown() {
		stateRefreshed.ID.Set(r.foo.GetID())
	}
	*/

	return stateRefreshed, true, nil
}
