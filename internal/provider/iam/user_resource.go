// Package org provides a Terraform resource to manage org users.
package iam

import (
	"context"
	"fmt"

	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/adminorg"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &userResource{}
	_ resource.ResourceWithConfigure   = &userResource{}
	_ resource.ResourceWithImportState = &userResource{}
	_ user                             = &userResource{}
)

// NewuserResource is a helper function to simplify the provider implementation.
func NewIAMUserResource() resource.Resource {
	return &userResource{}
}

// userResource is the resource implementation.
type userResource struct {
	client   *client.CloudAvenue
	adminOrg adminorg.AdminOrg
	user     commonUser
}

// Metadata returns the resource type name.
func (r *userResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_user"
}

// Schema defines the schema for the resource.
func (r *userResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = userSchema().GetResource(ctx)
}

func (r *userResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *userResource) Init(_ context.Context, rm *userResourceModel) (diags diag.Diagnostics) {
	r.user = commonUser{
		ID:   rm.ID,
		Name: rm.Name,
	}

	r.adminOrg, diags = adminorg.Init(r.client)

	return
}

// Create creates the resource and sets the initial Terraform state.
func (r *userResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &userResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(r.Init(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	userData := govcd.OrgUserConfiguration{
		Name:            plan.Name.ValueString(),
		RoleName:        plan.RoleName.ValueString(),
		FullName:        plan.FullName.ValueString(),
		EmailAddress:    plan.Email.ValueString(),
		Telephone:       plan.Telephone.ValueString(),
		IsEnabled:       plan.Enabled.ValueBool(),
		Password:        plan.Password.ValueString(),
		DeployedVmQuota: int(plan.DeployedVMQuota.ValueInt64()),
		StoredVmQuota:   int(plan.StoredVMQuota.ValueInt64()),
	}

	user, err := r.adminOrg.CreateUserSimple(userData)
	if err != nil {
		resp.Diagnostics.AddError("Error creating user", err.Error())
		return
	}

	state := *plan
	state.ID = types.StringValue(user.User.ID)

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *userResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &userResourceModel{}

	// Get current state
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(r.Init(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	user, err := r.GetUser(true)
	if err != nil {
		if govcd.ContainsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error retrieving user", err.Error())
		return
	}

	plan := &userResourceModel{
		ID:              types.StringValue(user.User.ID),
		Name:            types.StringValue(user.User.Name),
		RoleName:        types.StringValue(user.User.Role.Name),
		FullName:        utils.StringValueOrNull(user.User.FullName),
		Email:           utils.StringValueOrNull(user.User.EmailAddress),
		Telephone:       utils.StringValueOrNull(user.User.Telephone),
		Enabled:         types.BoolValue(user.User.IsEnabled),
		DeployedVMQuota: types.Int64Value(int64(user.User.DeployedVmQuota)),
		StoredVMQuota:   types.Int64Value(int64(user.User.StoredVmQuota)),
		TakeOwnership:   state.TakeOwnership,
		Password:        state.Password,
	}

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *userResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	plan := &userResourceModel{}

	// Get current state
	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(r.Init(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	user, err := r.GetUser(false)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving user", err.Error())
		return
	}

	userData := govcd.OrgUserConfiguration{
		Name:            plan.Name.ValueString(),
		RoleName:        plan.RoleName.ValueString(),
		FullName:        plan.FullName.ValueString(),
		EmailAddress:    plan.Email.ValueString(),
		Telephone:       plan.Telephone.ValueString(),
		IsEnabled:       plan.Enabled.ValueBool(),
		Password:        plan.Password.ValueString(),
		DeployedVmQuota: int(plan.DeployedVMQuota.ValueInt64()),
		StoredVmQuota:   int(plan.StoredVMQuota.ValueInt64()),
	}

	if err = user.UpdateSimple(userData); err != nil {
		resp.Diagnostics.AddError("Error updating user", err.Error())
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *userResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	state := &userResourceModel{}

	// Get current state
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(r.Init(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	user, err := r.GetUser(false)
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

func (r *userResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

func (r *userResource) GetUser(refresh bool) (*govcd.OrgUser, error) {
	return r.user.GetUser(r.adminOrg, refresh)
}
