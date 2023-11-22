package backup

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	v1common "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/common/netbackup"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go/v1/netbackup"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
)

var (
	_ resource.Resource                = &backupResource{}
	_ resource.ResourceWithConfigure   = &backupResource{}
	_ resource.ResourceWithImportState = &backupResource{}
)

const (
	vdc  string = "vdc"
	vapp string = "vapp"
	vm   string = "vm"
)

// NewbackupResource is a helper function to simplify the provider implementation.
func NewBackupResource() resource.Resource {
	return &backupResource{}
}

// backupResource is the resource implementation.
type backupResource struct {
	client *client.CloudAvenue
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
	defer metrics.New("cloudavenue_backup", r.client.GetOrgName(), metrics.Create)()

	plan := &backupModel{}

	// Retrieve values from plan
	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Extract policies from plan
	policies, d := plan.getPolicies(ctx)
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	// Refresh data NetBackup from the API
	job, err := r.client.CAVSDK.V1.Netbackup.Inventory.Refresh()
	if err != nil {
		resp.Diagnostics.AddError("Error refreshing NetBackup inventory", err.Error())
		return
	}
	if err := job.Wait(1, 45); err != nil {
		resp.Diagnostics.AddError("Error waiting for NetBackup inventory refresh", err.Error())
		return
	}

	// Get the type target object
	typeTarget, d := r.getTarget(plan)
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	// // Set target name
	// if plan.TargetName.IsNull() {
	// 	plan.TargetName.Set(typeTarget.GetName())
	// }

	// Apply the protection levels policies for each policy
	if err := applyPolicies(typeTarget, policies); err != nil {
		resp.Diagnostics.AddError("Error applying protection levels", err.Error())
		return
	}

	d = plan.Policies.Set(ctx, policies)
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	plan.ID.SetInt(typeTarget.GetID())

	// Use generic read function to refresh the state
	stateRefreshed, _, d := r.read(ctx, plan)
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, stateRefreshed)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *backupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer metrics.New("cloudavenue_backup", r.client.GetOrgName(), metrics.Read)()

	state := &backupModel{}

	// Get current state
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
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
	defer metrics.New("cloudavenue_backup", r.client.GetOrgName(), metrics.Update)()

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

	// Extract policies values from the plan
	policiesPlan, d := plan.getPolicies(ctx)
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	// Extract policies values from the state
	policiesState, d := state.getPolicies(ctx)
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	// Get the type target object
	typeTarget, d := r.getTarget(plan)
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	// For each Policy in plan not found in state => apply the protection levels
	// For each Policy in state not found in plan => apply the unprotection levels
	// For each Policy found in plan and state => nothing to do
	toProtect := make(backupModelPolicies, 0)
	toUnprotect := make(backupModelPolicies, 0)
	for _, policyPlan := range *policiesPlan {
		var found bool
		for _, policyState := range *policiesState {
			if policyPlan.PolicyName.Get() == policyState.PolicyName.Get() {
				found = true
			}
		}
		if !found {
			toProtect = append(toProtect, policyPlan)
		}
	}
	for _, policyState := range *policiesState {
		var found bool
		for _, policyPlan := range *policiesPlan {
			if policyState.PolicyName.Get() == policyPlan.PolicyName.Get() {
				found = true
			}
		}
		if !found {
			toUnprotect = append(toUnprotect, policyState)
		}
	}

	// Apply the protection levels policies for each policy
	if err := applyPolicies(typeTarget, &toProtect); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error applying %s protection levels", plan.Type.Get()), err.Error())
		return
	}

	// Unapply the protection levels policies for each policy
	if err := unApplyPolicies(typeTarget, &toUnprotect); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error unprotecting %s", plan.Type.Get()), err.Error())
		return
	}

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
	defer metrics.New("cloudavenue_backup", r.client.GetOrgName(), metrics.Delete)()

	state := &backupModel{}

	// Get current state
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Extract policies from plan
	policies, d := state.getPolicies(ctx)
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	// Get the target type object (vdc, vapp or vm)
	typeTarget, d := r.getTarget(state)
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	// apply the unprotection levels
	for _, policy := range *policies {
		if err := unApplyPolicy(typeTarget, policy); err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Error unprotecting %s", state.Type.Get()), err.Error())
			return
		}
	}
}

// ImportState imports the resource into the Terraform state.
func (r *backupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	defer metrics.New("cloudavenue_backup", r.client.GetOrgName(), metrics.Import)()

	// * Import with custom logic
	idParts := strings.Split(req.ID, ".")

	if len(idParts) != 2 {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: Type.TargetName Got: %q", req.ID),
		)
		return
	}

	// Refresh data NetBackup from the API
	job, err := r.client.CAVSDK.V1.Netbackup.Inventory.Refresh()
	if err != nil {
		resp.Diagnostics.AddError("Error refreshing NetBackup inventory", err.Error())
		return
	}
	if err := job.Wait(1, 45); err != nil {
		resp.Diagnostics.AddError("Error waiting for NetBackup inventory refresh", err.Error())
		return
	}

	data := NewBackup()
	data.Type.Set(idParts[0])
	data.TargetName.Set(idParts[1])

	// Use generic read function to refresh the state
	stateRefreshed, _, d := r.read(ctx, data)
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, stateRefreshed)...)
}

// * CustomFuncs

type target interface {
	GetProtectionLevelAvailableByName(string) (*netbackup.ProtectionLevel, error)
	Protect(netbackup.ProtectUnprotectRequest) (*v1common.JobAPIResponse, error)
	Unprotect(netbackup.ProtectUnprotectRequest) (*v1common.JobAPIResponse, error)
	GetID() int
	GetName() string
	ListProtectionLevels() (*netbackup.ProtectionLevels, error)
}

// Apply the protection level for a policy to the target.
// A target can be a vdc, vapp or vm.
// Return a policy with the protection level ID.
// Return an error if any.
func applyPolicy[T target](t T, policy backupModelPolicy) (backupModelPolicy, error) {
	// apply the protection levels
	job, err := t.Protect(netbackup.ProtectUnprotectRequest{
		ProtectionLevelID:   policy.PolicyID.GetIntPtr(),
		ProtectionLevelName: policy.PolicyName.Get(),
	})
	if err != nil {
		return backupModelPolicy{}, err
	}
	if err := job.Wait(1, 30); err != nil {
		return backupModelPolicy{}, err
	}

	// get the protection level ID
	if policy.PolicyID.Get() == 0 {
		pl, err := t.GetProtectionLevelAvailableByName(policy.PolicyName.Get())
		if err != nil {
			return backupModelPolicy{}, err
		}
		policy.PolicyID.SetInt(pl.ID)
	}
	return policy, nil
}

// Apply policies to the target.
// A target can be a vdc, vapp or vm.
// Return an error if any.
func applyPolicies[T target](t T, policies *backupModelPolicies) (err error) {
	for i, policy := range *policies {
		p, err := applyPolicy(t, policy)
		if err != nil {
			return err
		}
		(*policies)[i] = p
	}
	return nil
}

// Unapply policies to the target.
// A target can be a vdc, vapp or vm.
// Return an error if any.
func unApplyPolicies[T target](t T, policies *backupModelPolicies) (err error) {
	for _, policy := range *policies {
		if err := unApplyPolicy(t, policy); err != nil {
			return err
		}
	}
	return nil
}

// Unapply the protection level for a policy to the target.
// A target can be a VDC, VAPP or VM.
// Return an error if any.
func unApplyPolicy[T target](t T, policy backupModelPolicy) error {
	// apply the protection levels
	job, err := t.Unprotect(netbackup.ProtectUnprotectRequest{
		ProtectionLevelID:   policy.PolicyID.GetIntPtr(),
		ProtectionLevelName: policy.PolicyName.Get(),
	})
	if err != nil {
		return err
	}
	if err := job.Wait(1, 15); err != nil {
		return err
	}

	return nil
}

// read is a generic read function that can be used by the resource Create, Read and Update functions.
func (r *backupResource) read(ctx context.Context, planOrState *backupModel) (stateRefreshed *backupModel, found bool, diags diag.Diagnostics) {
	stateRefreshed = planOrState.Copy()

	// Get the type target object
	typeTarget, d := r.getTarget(planOrState)
	if d.HasError() {
		diags.Append(d...)
		return
	}

	// Set ID
	if !planOrState.ID.IsKnown() {
		stateRefreshed.ID.SetInt(typeTarget.GetID())
	}

	// get VDC Protection Levels from NetBackup
	policiesFromAPI, err := typeTarget.ListProtectionLevels()
	if err != nil {
		diags.AddError("Error listing protection levels", err.Error())
		return nil, true, diags
	}

	// Set target name
	if !planOrState.TargetName.IsKnown() {
		stateRefreshed.TargetName.Set(typeTarget.GetName())
	}

	// Add policies from API
	policies := backupModelPolicies{}
	for _, policyFromAPI := range *policiesFromAPI {
		x := backupModelPolicy{}
		x.PolicyName.Set(policyFromAPI.Name)
		x.PolicyID.SetInt(policyFromAPI.ID)
		policies = append(policies, x)
	}

	// set policies to the state
	d = stateRefreshed.Policies.Set(ctx, policies)
	if d.HasError() {
		diags.Append(d...)
		return
	}

	return stateRefreshed, true, nil
}

// getTarget returns the target object from the plan or state.
// A target can be a vdc, vapp or vm netbackup object.
func (r *backupResource) getTarget(data *backupModel) (typeTarget target, d diag.Diagnostics) {
	var err error
	switch data.Type.Get() {
	case vdc:
		typeTarget, err = r.client.CAVSDK.V1.Netbackup.VCloud.GetVDCByNameOrIdentifier(data.getTargetIDOrName())
	case vapp:
		typeTarget, err = r.client.CAVSDK.V1.Netbackup.VCloud.GetVAppByNameOrIdentifier(data.getTargetIDOrName())
	case vm:
		typeTarget, err = r.client.CAVSDK.V1.Netbackup.Machines.GetMachineByNameOrIdentifier(data.getTargetIDOrName())
	}
	if err != nil {
		d.AddError(fmt.Sprintf("Error getting vCloud Director %s", data.Type.Get()), err.Error())
		return nil, d
	}

	if data.ID.IsKnown() && typeTarget.GetID() != data.ID.GetInt() {
		d.AddError(fmt.Sprintf("Error getting vCloud Director %s", data.Type.Get()), "ID not match")
		return nil, d
	}

	return typeTarget, d
}
