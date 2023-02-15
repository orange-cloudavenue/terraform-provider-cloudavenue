// Package org provides a Terraform resource to manage org users.
package org

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &orgUserResource{}
	_ resource.ResourceWithConfigure   = &orgUserResource{}
	_ resource.ResourceWithImportState = &orgUserResource{}
)

// NewOrgUserResource is a helper function to simplify the provider implementation.
func NewOrgUserResource() resource.Resource {
	return &orgUserResource{}
}

// orgUserResource is the resource implementation.
type orgUserResource struct {
	client *client.CloudAvenue
}

type orgUserResourceModel struct {
	ID              types.String `tfsdk:"id"`
	UserName        types.String `tfsdk:"user_name"`
	FullName        types.String `tfsdk:"full_name"`
	Role            types.String `tfsdk:"role"`
	Password        types.String `tfsdk:"password"`
	Enabled         types.Bool   `tfsdk:"enabled"`
	Description     types.String `tfsdk:"description"`
	Email           types.String `tfsdk:"email"`
	Telephone       types.String `tfsdk:"telephone"`
	TakeOwnership   types.Bool   `tfsdk:"take_ownership"`
	DeployedVMQuota types.Int64  `tfsdk:"deployed_vm_quota"`
	StoredVMQuota   types.Int64  `tfsdk:"stored_vm_quota"`
}

// Metadata returns the resource type name.
func (r *orgUserResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_org_user"
}

// Schema defines the schema for the resource.
func (r *orgUserResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a user in an organization.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID is a unique identifier for the user.",
			},

			// Required attributes
			"user_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "User's name. Only lowercase letters allowed. Cannot be changed after creation",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(`^[a-z]+$`), "only lowercase letters allowed"),
				},
			},
			"role": schema.StringAttribute{
				MarkdownDescription: "The role of the user in the organization",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},

			// Optional attributes
			"enabled": schema.BoolAttribute{
				// TODO Add Planmodifier to set default to true
				// Actually default true is in Create Func
				MarkdownDescription: "`true` if the user is enabled and can log in. Default is `true`",
				Optional:            true,
				Computed:            true,
			},
			"description": schema.StringAttribute{
				// Not ForceNew, to allow the resource name to be updated
				MarkdownDescription: "Optional description of the catalog",
				Optional:            true,
				Computed:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "The user's password. This value is never returned on read. ",
				Optional:            true,
				Sensitive:           true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(8),
				},
			},
			"full_name": schema.StringAttribute{
				MarkdownDescription: "The user's full name",
				Optional:            true,
				Computed:            true,
			},
			"email": schema.StringAttribute{
				MarkdownDescription: "The user's email address",
				Optional:            true,
				Computed:            true,
			},
			"telephone": schema.StringAttribute{
				MarkdownDescription: "The user's telephone number",
				Optional:            true,
				Computed:            true,
			},
			"take_ownership": schema.BoolAttribute{
				MarkdownDescription: "Take ownership of user's objects on deletion.",
				Optional:            true,
			},
			"deployed_vm_quota": schema.Int64Attribute{
				MarkdownDescription: "Quota of vApps that this user can deploy. A value of `0` specifies an unlimited quota.",
				Optional:            true,
				Computed:            true,
			},
			"stored_vm_quota": schema.Int64Attribute{
				MarkdownDescription: "Quota of vApps that this user can store. A value of `0` specifies an unlimited quota.",
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

func (r *orgUserResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *orgUserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var (
		plan *orgUserResourceModel
		err  error
	)

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// user creation is accessible only for administator account
	adminOrg, err := r.client.Vmware.GetAdminOrgByNameOrId(r.client.GetOrg())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving Org", err.Error())
		return
	}

	var enabled bool
	if plan.Enabled.IsNull() || plan.Enabled.IsUnknown() {
		enabled = true
	} else {
		enabled = plan.Enabled.ValueBool()
	}

	var userData govcd.OrgUserConfiguration
	userData.RoleName = plan.Role.ValueString()
	userData.Name = plan.UserName.ValueString()
	userData.Description = plan.Description.ValueString()
	userData.FullName = plan.FullName.ValueString()
	userData.EmailAddress = plan.Email.ValueString()
	userData.Telephone = plan.Telephone.ValueString()
	userData.IsEnabled = enabled
	userData.Password = plan.Password.ValueString()
	userData.DeployedVmQuota = int(plan.DeployedVMQuota.ValueInt64())
	userData.StoredVmQuota = int(plan.StoredVMQuota.ValueInt64())

	user, err := adminOrg.CreateUserSimple(userData)
	if err != nil {
		resp.Diagnostics.AddError("Error creating user", err.Error())
		return
	}

	plan = &orgUserResourceModel{
		ID:              types.StringValue(user.User.ID),
		UserName:        types.StringValue(user.User.Name),
		FullName:        types.StringValue(user.User.FullName),
		Role:            types.StringValue(user.User.Role.Name),
		Email:           types.StringValue(user.User.EmailAddress),
		Telephone:       types.StringValue(user.User.Telephone),
		Enabled:         types.BoolValue(user.User.IsEnabled),
		Description:     types.StringValue(user.User.Description),
		DeployedVMQuota: types.Int64Value(int64(user.User.DeployedVmQuota)),
		StoredVMQuota:   types.Int64Value(int64(user.User.StoredVmQuota)),

		Password:      plan.Password,
		TakeOwnership: plan.TakeOwnership,
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *orgUserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *orgUserResourceModel

	// Get current state
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// user creation is accessible only for administator account
	adminOrg, err := r.client.Vmware.GetAdminOrgByNameOrId(r.client.GetOrg())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving Org", err.Error())
		return
	}

	user, err := adminOrg.GetUserByName(state.UserName.ValueString(), false)
	if err != nil {
		if govcd.ContainsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error retrieving user", err.Error())
		return
	}

	plan := &orgUserResourceModel{
		ID:              types.StringValue(user.User.ID),
		UserName:        types.StringValue(user.User.Name),
		FullName:        types.StringValue(user.User.FullName),
		Role:            types.StringValue(user.User.Role.Name),
		Email:           types.StringValue(user.User.EmailAddress),
		Telephone:       types.StringValue(user.User.Telephone),
		Enabled:         types.BoolValue(user.User.IsEnabled),
		Description:     types.StringValue(user.User.Description),
		DeployedVMQuota: types.Int64Value(int64(user.User.DeployedVmQuota)),
		StoredVMQuota:   types.Int64Value(int64(user.User.StoredVmQuota)),
		TakeOwnership:   state.TakeOwnership,
		Password:        state.Password,
	}

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *orgUserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state *orgUserResourceModel

	// Get current state
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// user update is accessible only for administator account
	adminOrg, err := r.client.Vmware.GetAdminOrgByNameOrId(r.client.GetOrg())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving Org", err.Error())
		return
	}

	user, err := adminOrg.GetUserByName(state.UserName.ValueString(), false)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving user", err.Error())
		return
	}

	var userData govcd.OrgUserConfiguration
	userData.RoleName = plan.Role.ValueString()
	userData.Name = plan.UserName.ValueString()
	userData.Description = plan.Description.ValueString()
	userData.FullName = plan.FullName.ValueString()
	userData.EmailAddress = plan.Email.ValueString()
	userData.Telephone = plan.Telephone.ValueString()
	userData.IsEnabled = plan.Enabled.ValueBool()
	userData.Password = plan.Password.ValueString()
	userData.DeployedVmQuota = int(plan.DeployedVMQuota.ValueInt64())
	userData.StoredVmQuota = int(plan.StoredVMQuota.ValueInt64())

	err = user.UpdateSimple(userData)
	if err != nil {
		resp.Diagnostics.AddError("Error updating user", err.Error())
	}

	userRefresh, err := adminOrg.GetUserByName(state.UserName.ValueString(), true)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving user", err.Error())
		return
	}

	plan = &orgUserResourceModel{
		ID:              types.StringValue(userRefresh.User.ID),
		UserName:        types.StringValue(userRefresh.User.Name),
		FullName:        types.StringValue(userRefresh.User.FullName),
		Role:            types.StringValue(userRefresh.User.Role.Name),
		Email:           types.StringValue(userRefresh.User.EmailAddress),
		Telephone:       types.StringValue(userRefresh.User.Telephone),
		Enabled:         types.BoolValue(userRefresh.User.IsEnabled),
		Description:     types.StringValue(userRefresh.User.Description),
		DeployedVMQuota: types.Int64Value(int64(userRefresh.User.DeployedVmQuota)),
		StoredVMQuota:   types.Int64Value(int64(userRefresh.User.StoredVmQuota)),
		TakeOwnership:   state.TakeOwnership,
		Password:        state.Password,
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *orgUserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *orgUserResourceModel

	// Get current state
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// user delete is accessible only for administator account
	adminOrg, err := r.client.Vmware.GetAdminOrgByNameOrId(r.client.GetOrg())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving Org", err.Error())
		return
	}

	user, err := adminOrg.GetUserByName(state.UserName.ValueString(), false)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving user", err.Error())
		return
	}

	err = user.Delete(state.TakeOwnership.ValueBool())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting user", err.Error())
		return
	}
}

func (r *orgUserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("user_name"), req, resp)
}
