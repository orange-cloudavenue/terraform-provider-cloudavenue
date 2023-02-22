// Package vm provides a Terraform resource.
package vm

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	myvapp "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vapp"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &vmInsertedMediaResource{}
	_ resource.ResourceWithConfigure = &vmInsertedMediaResource{}
)

// NewVMInsertedMediaResource is a helper function to simplify the provider implementation.
func NewVMInsertedMediaResource() resource.Resource {
	return &vmInsertedMediaResource{}
}

// vmInsertedMediaResource is the resource implementation.
type vmInsertedMediaResource struct {
	client *client.CloudAvenue
}

type vmInsertedMediaResourceModel struct {
	ID       types.String `tfsdk:"id"`
	VDC      types.String `tfsdk:"vdc"`
	Catalog  types.String `tfsdk:"catalog"`
	Name     types.String `tfsdk:"name"`
	VappName types.String `tfsdk:"vapp_name"`
	VMName   types.String `tfsdk:"vm_name"`
	// EjectForce types.Bool   `tfsdk:"eject_force"` - Disable attributes - Issue referrer: vmware/go-vcloud-director#552
}

// Metadata returns the resource type name.
func (r *vmInsertedMediaResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_" + "inserted_media"
}

// Schema defines the schema for the resource.
func (r *vmInsertedMediaResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The inserted_media resource resource for inserting or ejecting media (ISO) file for the VM. Create this resource for inserting the media, and destroy it for ejecting.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID of the inserted media. This is the vm Id where the media is inserted.",
			},
			"vdc": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The name of VDC to use, optional if defined at provider level",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"catalog": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the catalog where to find media file",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Media file name in catalog which will be inserted to VM",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"vapp_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "vApp name where VM is located",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"vm_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "VM name where media will be inserted or ejected",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			// "eject_force": schema.BoolAttribute{ - Disable attributes - Issue referrer: vmware/go-vcloud-director#552
			//	Optional:            true,
			//	MarkdownDescription: "Allows to pass answer to question in vCD when ejecting from a VM which is powered on. True means 'Yes' as answer to question. Default is true",
			//	PlanModifiers: []planmodifier.Bool{
			//		boolpm.SetDefault(true),
			//	},
			// },
		},
	}
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

	// If VDC is not defined at data source level, use the one defined at provider level
	if plan.VDC.IsNull() || plan.VDC.IsUnknown() {
		if r.client.DefaultVDCExist() {
			plan.VDC = types.StringValue(r.client.GetDefaultVDC())
		} else {
			resp.Diagnostics.AddError("Missing VDC", "VDC is required when not defined at provider level")
			return
		}
	}

	// Get Org & VDC
	org, vdc, err := r.client.GetOrgAndVDC(r.client.GetOrg(), plan.VDC.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving VDC", err.Error())
		return
	}

	// Check if vApp exists
	vapp, err := vdc.GetVAppByName(plan.VappName.ValueString(), true)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving vApp", err.Error())
		return
	}

	// Lock vApp
	lockvapp := myvapp.Ref{
		Name:  plan.VappName.ValueString(),
		Org:   org.Org.Name,
		VDC:   plan.VDC.ValueString(),
		TFCtx: ctx,
	}
	if errLock := lockvapp.LockParentVApp(); errors.Is(errLock, myvapp.ErrVAppRefEmpty) {
		resp.Diagnostics.AddError("Error locking vApp", errLock.Error())
		return
	}
	defer func() {
		if errUnlock := lockvapp.UnLockParentVApp(); errUnlock != nil {
			resp.Diagnostics.AddError("Error unlocking vApp", errUnlock.Error())
			return
		}
	}()

	// Check if VM exists
	vm, err := vapp.GetVMByName(plan.VMName.ValueString(), true)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving VM", err.Error())
		return
	}

	// Insert media
	task, err := vm.HandleInsertMedia(org, plan.Catalog.ValueString(), plan.Name.ValueString())
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
		VDC:      plan.VDC,
		Catalog:  plan.Catalog,
		Name:     plan.Name,
		VappName: plan.VappName,
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

	// If VDC is not defined at data source level, use the one defined at provider level
	if state.VDC.IsNull() || state.VDC.IsUnknown() {
		if r.client.DefaultVDCExist() {
			state.VDC = types.StringValue(r.client.GetDefaultVDC())
		} else {
			resp.Diagnostics.AddError("Missing VDC", "VDC is required when not defined at provider level")
			return
		}
	}

	// Get Org & VDC
	_, vdc, err := r.client.GetOrgAndVDC(r.client.GetOrg(), state.VDC.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving VDC or Org", err.Error())
		return
	}

	// Check if vApp exists
	vapp, err := vdc.GetVAppByName(state.VappName.ValueString(), true)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving vApp", err.Error())
		return
	}

	// Check if VM exists
	vm, err := vapp.GetVMByName(state.VMName.ValueString(), true)
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
		VDC:      state.VDC,
		Catalog:  state.Catalog,
		Name:     state.Name,
		VappName: state.VappName,
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

	// If VDC is not defined at data source level, use the one defined at provider level
	if plan.VDC.IsNull() || plan.VDC.IsUnknown() {
		if r.client.DefaultVDCExist() {
			plan.VDC = types.StringValue(r.client.GetDefaultVDC())
		} else {
			resp.Diagnostics.AddError("Missing VDC", "VDC is required when not defined at provider level")
			return
		}
	}

	// Get Org & VDC
	org, vdc, err := r.client.GetOrgAndVDC(r.client.GetOrg(), plan.VDC.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving VDC or Org", err.Error())
		return
	}

	// Check if vApp exists
	vapp, err := vdc.GetVAppByName(plan.VappName.ValueString(), true)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving vApp", err.Error())
		return
	}

	// Lock vApp
	lockvapp := myvapp.Ref{
		Name:  state.VappName.ValueString(),
		Org:   org.Org.Name,
		VDC:   state.VDC.ValueString(),
		TFCtx: ctx,
	}
	if errLock := lockvapp.LockParentVApp(); errors.Is(errLock, myvapp.ErrVAppRefEmpty) {
		resp.Diagnostics.AddError("Error locking vApp", errLock.Error())
		return
	}
	defer func() {
		if errUnlock := lockvapp.UnLockParentVApp(); errUnlock != nil {
			resp.Diagnostics.AddError("Error unlocking vApp", errUnlock.Error())
			return
		}
	}()

	// Check if VM exists
	vm, err := vapp.GetVMByName(state.VMName.ValueString(), true)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving VM", err.Error())
		return
	}

	plan = &vmInsertedMediaResourceModel{
		ID:       types.StringValue(vm.VM.ID),
		VDC:      plan.VDC,
		Catalog:  plan.Catalog,
		Name:     plan.Name,
		VappName: plan.VappName,
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

	// If VDC is not defined at data source level, use the one defined at provider level
	if state.VDC.IsNull() || state.VDC.IsUnknown() {
		if r.client.DefaultVDCExist() {
			state.VDC = types.StringValue(r.client.GetDefaultVDC())
		} else {
			resp.Diagnostics.AddError("Missing VDC", "VDC is required when not defined at provider level")
			return
		}
	}

	// Get Org & VDC
	org, vdc, err := r.client.GetOrgAndVDC(r.client.GetOrg(), state.VDC.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving VDC or Org", err.Error())
		return
	}

	// Check if vApp exists
	vapp, err := vdc.GetVAppByName(state.VappName.ValueString(), true)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving vApp", err.Error())
		return
	}

	// Lock vApp
	lockvapp := myvapp.Ref{
		Name:  state.VappName.ValueString(),
		Org:   org.Org.Name,
		VDC:   state.VDC.ValueString(),
		TFCtx: ctx,
	}
	if errLock := lockvapp.LockParentVApp(); errors.Is(errLock, myvapp.ErrVAppRefEmpty) {
		resp.Diagnostics.AddError("Error locking vApp", errLock.Error())
		return
	}
	defer func() {
		if errUnlock := lockvapp.UnLockParentVApp(); errUnlock != nil {
			resp.Diagnostics.AddError("Error unlocking vApp", errUnlock.Error())
			return
		}
	}()

	// Check if VM exists
	vm, err := vapp.GetVMByName(state.VMName.ValueString(), true)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving VM", err.Error())
		return
	}

	// Eject media
	_, err = vm.HandleEjectMediaAndAnswer(org, state.Catalog.ValueString(), state.Name.ValueString(), true)
	if err != nil {
		resp.Diagnostics.AddError("Error ejecting media", err.Error())
		return
	}
}
