package iam

import (
	"context"
	"fmt"

	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/adminorg"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &UserSAMLResource{}
	_ resource.ResourceWithConfigure   = &UserSAMLResource{}
	_ resource.ResourceWithImportState = &UserSAMLResource{}
)

// NewUserSAMLResource is a helper function to simplify the provider implementation.
func NewUserSAMLResource() resource.Resource {
	return &UserSAMLResource{}
}

// UserSAMLResource is the resource implementation.
type UserSAMLResource struct {
	client   *client.CloudAvenue
	adminOrg adminorg.AdminOrg
}

// Init Initializes the resource.
func (r *UserSAMLResource) Init(ctx context.Context, rm *UserSAMLModel) (diags diag.Diagnostics) {
	r.adminOrg, diags = adminorg.Init(r.client)
	return
}

// Metadata returns the resource type name.
func (r *UserSAMLResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_user_saml"
}

// Schema defines the schema for the resource.
func (r *UserSAMLResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = userSAMLSchema(ctx).GetResource(ctx)
}

func (r *UserSAMLResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *UserSAMLResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer metrics.New("cloudavenue_iam_user_saml", r.client.GetOrgName(), metrics.Create)()

	plan := &UserSAMLModel{}

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

	userData := govcd.OrgUserConfiguration{
		ProviderType:    govcd.OrgUserProviderSAML,
		Name:            plan.UserName.Get(),
		RoleName:        plan.RoleName.Get(),
		IsEnabled:       plan.Enabled.Get(),
		DeployedVmQuota: plan.DeployedVMQuota.GetInt(),
		StoredVmQuota:   plan.StoredVMQuota.GetInt(),
		IsExternal:      true,
	}

	user, err := r.adminOrg.CreateUserSimple(userData)
	if err != nil {
		resp.Diagnostics.AddError("Error creating user", err.Error())
		return
	}

	plan.ID.Set(user.User.ID)

	// Use generic read function to refresh the state
	state, found, d := r.read(ctx, plan)
	if !found {
		resp.Diagnostics.AddError("User not found", fmt.Sprintf("User with name %s not found after import", plan.UserName.Get()))
		return
	}
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *UserSAMLResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer metrics.New("cloudavenue_iam_user_saml", r.client.GetOrgName(), metrics.Read)()

	state := &UserSAMLModel{}

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
func (r *UserSAMLResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer metrics.New("cloudavenue_iam_user_saml", r.client.GetOrgName(), metrics.Update)()

	var (
		plan  = &UserSAMLModel{}
		state = &UserSAMLModel{}
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

	user, err := r.adminOrg.GetUserById(state.ID.Get(), true)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving user", err.Error())
		return
	}

	userData := govcd.OrgUserConfiguration{
		ProviderType:    govcd.OrgUserProviderSAML,
		Name:            plan.UserName.Get(),
		RoleName:        plan.RoleName.Get(),
		IsEnabled:       plan.Enabled.Get(),
		DeployedVmQuota: plan.DeployedVMQuota.GetInt(),
		StoredVmQuota:   plan.StoredVMQuota.GetInt(),
		IsExternal:      true,
	}

	if err := user.UpdateSimple(userData); err != nil {
		resp.Diagnostics.AddError("Error updating user", err.Error())
		return
	}

	// Special case to inject TakeOwnership value
	state.TakeOwnership.Set(plan.TakeOwnership.Get())

	// Use generic read function to refresh the state
	stateRefreshed, _, d := r.read(ctx, state)
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, stateRefreshed)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *UserSAMLResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer metrics.New("cloudavenue_iam_user_saml", r.client.GetOrgName(), metrics.Delete)()

	state := &UserSAMLModel{}

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

	user, err := r.adminOrg.GetUserById(state.ID.Get(), true)
	if err != nil {
		if govcd.ContainsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error retrieving user", err.Error())
		return
	}

	if err = user.Delete(state.TakeOwnership.ValueBool()); err != nil {
		resp.Diagnostics.AddError("Error deleting user", err.Error())
	}
}

func (r *UserSAMLResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	defer metrics.New("cloudavenue_iam_user_saml", r.client.GetOrgName(), metrics.Import)()

	userData := &UserSAMLModel{}

	// Init the resource
	resp.Diagnostics.Append(r.Init(ctx, userData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	user, err := r.adminOrg.GetUserByName(req.ID, true)
	if err != nil {
		if govcd.ContainsNotFound(err) {
			resp.Diagnostics.AddError("User not found", fmt.Sprintf("User with name %s not found", req.ID))
			return
		}
		resp.Diagnostics.AddError("Error retrieving user", err.Error())
		return
	}

	if user.User.ProviderType != govcd.OrgUserProviderSAML {
		resp.Diagnostics.AddError("User is not SAML user", fmt.Sprintf("User with name %s is not a SAML user.", req.ID))
		return
	}

	userData.ID.Set(user.User.ID)
	userData.UserName.Set(user.User.Name)

	// Use generic read function to refresh the state
	stateRefreshed, found, d := r.read(ctx, userData)
	if !found {
		resp.Diagnostics.AddError("User not found", fmt.Sprintf("User with name %s not found", req.ID))
		return
	}
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, stateRefreshed)...)
}

// * CustomFuncs

// read is a generic read function that can be used by the resource Create, Read and Update functions.
func (r *UserSAMLResource) read(_ context.Context, planOrState *UserSAMLModel) (stateRefreshed *UserSAMLModel, found bool, diags diag.Diagnostics) {
	stateRefreshed = planOrState.Copy()

	user, err := r.adminOrg.GetUserById(planOrState.ID.Get(), true)
	if err != nil {
		if govcd.ContainsNotFound(err) {
			return nil, false, nil
		}
		diags.AddError("Error retrieving user", err.Error())
		return
	}

	stateRefreshed.ID.Set(user.User.ID)
	stateRefreshed.UserName.Set(user.User.Name)
	stateRefreshed.RoleName.Set(user.User.Role.Name)
	stateRefreshed.Enabled.Set(user.User.IsEnabled)
	stateRefreshed.DeployedVMQuota.SetInt(user.User.DeployedVmQuota)
	stateRefreshed.StoredVMQuota.SetInt(user.User.StoredVmQuota)

	return stateRefreshed, true, nil
}
