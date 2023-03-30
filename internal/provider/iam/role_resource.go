package iam

import (
	"context"
	"fmt"
	"strings"

	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/adminorg"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &roleResource{}
	_ resource.ResourceWithConfigure   = &roleResource{}
	_ resource.ResourceWithImportState = &roleResource{}
	_ role                             = &roleResource{}
)

// NewroleResource is a helper function to simplify the provider implementation.
func NewRoleResource() resource.Resource {
	return &roleResource{}
}

// roleResource is the resource implementation.
type roleResource struct {
	client   *client.CloudAvenue
	adminOrg adminorg.AdminOrg
	role     commonRole
}

func (r *roleResource) Init(_ context.Context, rm *roleResourceModel) (diags diag.Diagnostics) {
	r.role = commonRole{
		ID:   rm.ID,
		Name: rm.Name,
	}
	r.adminOrg, diags = adminorg.Init(r.client)

	return
}

// Metadata returns the resource type name.
func (r *roleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_" + "role"
}

// Schema defines the schema for the resource.
func (r *roleResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = roleSchema().GetResource(ctx)
}

func (r *roleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *roleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	plan := &roleResourceModel{}

	// Read the plan
	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(r.Init(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check rights are valid
	rights := make([]govcdtypes.OpenApiReference, 0)
	for _, right := range plan.Rights.Elements() {
		rg := strings.Trim(right.String(), "\"")
		x, err := r.adminOrg.GetRightByName(rg)
		if err != nil {
			resp.Diagnostics.AddError("Error retrieving right", err.Error())
			return
		}
		rights = append(rights, govcdtypes.OpenApiReference{Name: rg, ID: x.ID})
	}

	// Add implied rights
	missingImpliedRights, err := govcd.FindMissingImpliedRights(&r.client.Vmware.Client, rights)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving implied rights", err.Error())
		return
	}

	// Print missing implied rights
	if len(missingImpliedRights) > 0 {
		message := "The rights set for this role require the following implied rights to be added:"
		rightsList := ""
		for _, right := range missingImpliedRights {
			rightsList += fmt.Sprintf("\"%s\",\n", right.Name)
		}
		resp.Diagnostics.AddError(message, rightsList)
		return
	}

	// Create the role
	role, err := r.adminOrg.CreateRole(&govcdtypes.Role{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		BundleKey:   govcdtypes.VcloudUndefinedKey,
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating role", err.Error())
		return
	}
	if len(rights) > 0 {
		err = role.AddRights(rights)
		if err != nil {
			resp.Diagnostics.AddError("Error adding rights to role", err.Error())
			return
		}
	}

	// Set Plan state
	plan.ID = types.StringValue(role.Role.ID)
	plan.Name = types.StringValue(role.Role.Name)
	plan.Description = types.StringValue(role.Role.Description)

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *roleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *roleResourceModel

	// Read state
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(r.Init(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get Role
	role, err := r.GetRole()
	if err != nil {
		if govcd.ContainsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error retrieving role", err.Error())
		return
	}

	// Get rights
	rights, err := role.GetRights(nil)
	if err != nil {
		return
	}
	assignedRights := []attr.Value{}
	for _, right := range rights {
		assignedRights = append(assignedRights, types.StringValue(right.Name))
	}

	plan := &roleResourceModel{
		ID:          types.StringValue(role.Role.ID),
		Name:        types.StringValue(role.Role.Name),
		Description: types.StringValue(role.Role.Description),
		Rights:      types.SetNull(types.StringType),
	}

	var y diag.Diagnostics
	if len(assignedRights) > 0 {
		plan.Rights, y = types.SetValue(types.StringType, assignedRights)
		if y.HasError() {
			return
		}
	}

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *roleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var (
		state *roleResourceModel
		err   error
		role  *govcd.Role
	)

	// Read state
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(r.Init(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the role
	role, err = r.GetRole()
	if err != nil {
		if govcd.ContainsNotFound(err) {
			tflog.Debug(ctx, "Unable to find role. Removing from tfstate")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error retrieving role", err.Error())
		return
	}
	err = role.Delete()
	if err != nil {
		resp.Diagnostics.AddError("Error deleting role", err.Error())
		return
	}
}

func (r *roleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var (
		plan, state *roleResourceModel
		err         error
		role        *govcd.Role
	)

	// Get state and plan
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(r.Init(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the role
	role, err = r.GetRole()
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving role", err.Error())
		return
	}

	// Update the role Name or Description
	if (plan.Name.Equal(state.Name)) || (plan.Description.Equal(state.Description)) {
		role.Role.Name = plan.Name.ValueString()
		role.Role.Description = plan.Description.ValueString()
		_, err = role.Update()
		if err != nil {
			resp.Diagnostics.AddError("Error updating role", err.Error())
			return
		}
	}

	// Check rights are valid
	rights := make([]govcdtypes.OpenApiReference, 0)
	for _, right := range plan.Rights.Elements() {
		rg := strings.Trim(right.String(), "\"")
		x, err := r.adminOrg.GetRightByName(rg)
		if err != nil {
			resp.Diagnostics.AddError("Error retrieving right", err.Error())
			return
		}
		rights = append(rights, govcdtypes.OpenApiReference{Name: rg, ID: x.ID})
	}
	// Add implied rights
	missingImpliedRights, err := govcd.FindMissingImpliedRights(&r.client.Vmware.Client, rights)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving implied rights", err.Error())
		return
	}
	// Print missing implied rights
	if len(missingImpliedRights) > 0 {
		message := "The rights set for this role require the following implied rights to be added:"
		rightsList := ""
		for _, right := range missingImpliedRights {
			rightsList += fmt.Sprintf("\"%s\",\n", right.Name)
		}
		resp.Diagnostics.AddError(message, rightsList)
		return
	}
	// Update the role rights
	if len(rights) > 0 {
		err = role.UpdateRights(rights)
		if err != nil {
			resp.Diagnostics.AddError("Error updating role rights", err.Error())
			return
		}
	} else {
		currentRights, err := role.GetRights(nil)
		if err != nil {
			resp.Diagnostics.AddError("Error retrieving role rights", err.Error())
			return
		}
		if len(currentRights) > 0 {
			err = role.RemoveAllRights()
			if err != nil {
				resp.Diagnostics.AddError("Error removing role rights", err.Error())
				return
			}
		}
	}

	// Set Plan state
	plan = &roleResourceModel{
		ID:          types.StringValue(role.Role.ID),
		Name:        types.StringValue(role.Role.Name),
		Description: types.StringValue(role.Role.Description),
		Rights:      plan.Rights,
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

//go:generate tf-doc-extractor -filename $GOFILE -example-dir ../../../examples -resource
func (r *roleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

func (r *roleResource) GetRole() (*govcd.Role, error) {
	return r.role.GetRole(r.adminOrg)
}
