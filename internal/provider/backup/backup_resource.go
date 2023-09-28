package netbackup

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &BackupResource{}
	_ resource.ResourceWithConfigure   = &BackupResource{}
	_ resource.ResourceWithImportState = &BackupResource{}
	// _ resource.ResourceWithModifyPlan     = &BackupResource{}
	// _ resource.ResourceWithUpgradeState   = &BackupResource{}
	// _ resource.ResourceWithValidateConfig = &BackupResource{}.
)

// NewBackupResource is a helper function to simplify the provider implementation.
func NewBackupResource() resource.Resource {
	return &BackupResource{}
}

// BackupResource is the resource implementation.
type BackupResource struct {
	client *client.CloudAvenue

	// Uncomment the following lines if you need to access the resource's.
	// org    org.Org
	// vdc    vdc.VDC
	// vapp   vapp.VAPP
}

// If the resource don't have same schema/structure as the data source, you can use the following code:
// type BackupResourceModel struct {
// 	ID types.String `tfsdk:"id"`
// }

// Init Initializes the resource.
func (r *BackupResource) Init(ctx context.Context, rm *BackupModel) (diags diag.Diagnostics) {
	// Uncomment the following lines if you need to access to the Org
	// r.org, diags = org.Init(r.client)
	// if diags.HasError() {
	// 	return
	// }

	// Uncomment the following lines if you need to access to the VDC
	// r.vdc, diags = vdc.Init(r.client, rm.VDC)
	// if diags.HasError() {
	// 	return
	// }

	// Uncomment the following lines if you need to access to the VAPP
	// r.vapp, diags = vapp.Init(r.client, r.vdc, rm.VAppID, rm.VAppName)

	return
}

// Metadata returns the resource type name.
func (r *BackupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName
}

// Schema defines the schema for the resource.
func (r *BackupResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = backupSchema(ctx).GetResource(ctx)
}

func (r *BackupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *BackupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &BackupModel{}

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
func (r *BackupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &BackupModel{}

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
func (r *BackupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var (
		plan  = &BackupModel{}
		state = &BackupModel{}
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
func (r *BackupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	state := &BackupModel{}

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

func (r *BackupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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
func (r *BackupResource) read(ctx context.Context, planOrState *BackupModel) (stateRefreshed *BackupModel, found bool, diags diag.Diagnostics) {
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
