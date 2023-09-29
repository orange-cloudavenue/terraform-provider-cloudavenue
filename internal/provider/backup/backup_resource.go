package backup

import (
	"context"
	"fmt"

	"github.com/orange-cloudavenue/netbackup-sdk-go/netbackupclient"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &backupResource{}
	_ resource.ResourceWithConfigure   = &backupResource{}
	_ resource.ResourceWithImportState = &backupResource{}
	// _ resource.ResourceWithModifyPlan     = &backupResource{}
	// _ resource.ResourceWithUpgradeState   = &backupResource{}
	// _ resource.ResourceWithValidateConfig = &backupResource{}.
)

// NewbackupResource is a helper function to simplify the provider implementation.
func NewBackupResource() resource.Resource {
	return &backupResource{}
}

// backupResource is the resource implementation.
type backupResource struct {
	client *client.CloudAvenue
	// Uncomment the following lines if you need to access the resource's.
	// org    org.Org
	// vdc vdc.VDC
	// vapp   vapp.VAPP
}

// If the resource don't have same schema/structure as the data source, you can use the following code:
// type backupResourceModel struct {
// 	ID types.String `tfsdk:"id"`
// }

// Init Initializes the resource.
func (r *backupResource) Init(ctx context.Context, rm *backupModel) (diags diag.Diagnostics) {
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
func (r *backupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName
}

// Schema defines the schema for the resource.
func (r *backupResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = backupSchema(ctx).GetResource(ctx)
}

func (r *backupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *backupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &backupModel{}

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

	// Define if the target is a VDC ID or a VDC Name
	var nameOrID string
	if plan.TargetID.IsNull() {
		nameOrID = plan.TargetName.Get()
	} else {
		nameOrID = plan.TargetID.Get()
	}

	// for policies extract values from the plan
	policies, d := plan.GetPolicies(ctx)
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	// switch on the type of the backup
	switch plan.Type.Get() {
	case "vdc":
		vdc, err := r.client.NetBackupClient.VCloud.GetVdcByNameOrIdentifier(nameOrID)
		if err != nil {
			resp.Diagnostics.AddError("Error getting vCloud Director Virtual Data Center", err.Error())
			return
		}
		// for each policy, protect the VDC
		for _, policy := range policies {
			if policy.Enabled.Get() {
				job, err := vdc.Protect(netbackupclient.ProtectUnprotectRequest{
					ProtectionLevelName: policy.PolicyName.Get(),
					ProtectionLevelID:   policy.PolicyID.GetIntPtr(),
				})
				if err != nil {
					resp.Diagnostics.AddError("Error protecting vCloud Director Virtual Data Center", err.Error())
					return
				}
				if err := job.Wait(1, 15); err != nil {
					resp.Diagnostics.AddError("Error waiting for job", err.Error())
				}
			}
		}
	case "vapp":
		// TODO: Create a backup for a VAPP
	case "vm":
		// TODO: Create a backup for a VM
	}

	// Use generic read function to refresh the state
	// state, _, d := r.read(ctx, plan)
	// if d.HasError() {
	// 	resp.Diagnostics.Append(d...)
	// 	return
	// }

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *backupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &backupModel{}

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
func (r *backupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var (
		plan  = &backupModel{}
		state = &backupModel{}
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
func (r *backupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	state := &backupModel{}

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

func (r *backupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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
func (r *backupResource) read(ctx context.Context, planOrState *backupModel) (stateRefreshed *backupModel, found bool, diags diag.Diagnostics) {
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
