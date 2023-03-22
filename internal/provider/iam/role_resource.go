package iam

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &iamRoleResource{}
	_ resource.ResourceWithConfigure   = &iamRoleResource{}
	_ resource.ResourceWithImportState = &iamRoleResource{}
)

// NewiamRoleResource is a helper function to simplify the provider implementation.
func NewIAMRoleResource() resource.Resource {
	return &iamRoleResource{}
}

// iamRoleResource is the resource implementation.
type iamRoleResource struct {
	client *client.CloudAvenue
}

// iamRoleResourceModel is the internal state representation of the resource.
type iamRoleResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	BundleKey   types.String `tfsdk:"bundle_key"`
	ReadOnly    types.Bool   `tfsdk:"read_only"`
	Rights      types.Set    `tfsdk:"rights"`
}

// Metadata returns the resource type name.
func (r *iamRoleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_" + "role"
}

// Schema defines the schema for the resource.
func (r *iamRoleResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The Role resource allows you to manage a role in CloudAvenue.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID is a unique identifier for the role",
			},
			"name": schema.StringAttribute{
				// Not ForceNew, to allow the resource name to be updated
				Required:            true,
				MarkdownDescription: "A name for the role",
			},
			"description": schema.StringAttribute{
				// Not ForceNew, to allow the resource name to be updated
				Required:            true,
				MarkdownDescription: "A description for the role",
			},
			// * Remove
			"bundle_key": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Key used for internationalization",
			},
			// * Remove in resource
			"read_only": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Indicates if the role is read only",
			},
			"rights": schema.SetAttribute{
				// Not ForceNew, to allow the resource name to be updated
				Required:            true,
				MarkdownDescription: "A list of rights for the role",
				ElementType:         types.StringType,
			},
		},
	}
}

func (r *iamRoleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *iamRoleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var (
		plan *iamRoleResourceModel
		err  error
	)

	// Read the plan
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// role creation is accessible only in administrator API part
	// (only administrator, organization administrator and Catalog author are allowed)
	adminOrg, err := r.client.Vmware.GetAdminOrgByNameOrId(r.client.GetOrgName())
	if err != nil {
		resp.Diagnostics.AddError("[role create] Error retrieving Org", err.Error())
		return
	}

	// Check rights are valid
	rights := make([]govcdtypes.OpenApiReference, 0)
	for _, right := range plan.Rights.Elements() {
		rg := strings.Trim(right.String(), "\"")
		x, err := adminOrg.GetRightByName(rg)
		if err != nil {
			resp.Diagnostics.AddError("[role create] Error retrieving right", err.Error())
			return
		}
		rights = append(rights, govcdtypes.OpenApiReference{Name: rg, ID: x.ID})
	}

	// Add implied rights
	missingImpliedRights, err := govcd.FindMissingImpliedRights(&r.client.Vmware.Client, rights)
	if err != nil {
		resp.Diagnostics.AddError("[role create] Error retrieving implied rights", err.Error())
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
	role, err := adminOrg.CreateRole(&govcdtypes.Role{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		BundleKey:   govcdtypes.VcloudUndefinedKey,
	})
	if err != nil {
		resp.Diagnostics.AddError("[role create] Error creating role", err.Error())
		return
	}
	if len(rights) > 0 {
		err = role.AddRights(rights)
		if err != nil {
			resp.Diagnostics.AddError("[role create] Error adding rights to role", err.Error())
			return
		}
	}

	// Set Plan state
	plan = &iamRoleResourceModel{
		ID:          types.StringValue(role.Role.ID),
		Name:        types.StringValue(role.Role.Name),
		BundleKey:   types.StringValue(role.Role.BundleKey),
		ReadOnly:    types.BoolValue(role.Role.ReadOnly),
		Description: types.StringValue(role.Role.Description),
		Rights:      plan.Rights,
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *iamRoleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Retrieve values from state
	var (
		state *iamRoleResourceModel
		err   error
		role  *govcd.Role
	)

	// Read state
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// role read is accessible only in administrator
	adminOrg, err := r.client.Vmware.GetAdminOrgByNameOrId(r.client.GetOrgName())
	if err != nil {
		resp.Diagnostics.AddError("[role read] Error retrieving Org", err.Error())
		return
	}

	// Get the role
	if state.ID.IsNull() {
		role, err = adminOrg.GetRoleByName(state.Name.ValueString())
	} else {
		role, err = adminOrg.GetRoleById(state.ID.ValueString())
	}
	if err != nil {
		if govcd.ContainsNotFound(err) {
			tflog.Info(ctx, "[DEBUG] Unable to find role. Removing from tfstate")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("[role read] Error retrieving role", err.Error())
		return
	}

	// Set state to fully populated data
	plan := &iamRoleResourceModel{
		ID:          types.StringValue(role.Role.ID),
		Name:        types.StringValue(role.Role.Name),
		BundleKey:   types.StringValue(role.Role.BundleKey),
		ReadOnly:    types.BoolValue(role.Role.ReadOnly),
		Description: types.StringValue(role.Role.Description),
	}

	// Get rights
	rights, err := role.GetRights(nil)
	if err != nil {
		resp.Diagnostics.AddError("[role read] Error while querying role rights", err.Error())
		return
	}
	assignedRights := []attr.Value{}
	for _, right := range rights {
		assignedRights = append(assignedRights, types.StringValue(right.Name))
	}
	var y diag.Diagnostics
	if len(assignedRights) > 0 {
		plan.Rights, y = types.SetValue(types.StringType, assignedRights)
		resp.Diagnostics.Append(y...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *iamRoleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var (
		state *iamRoleResourceModel
		err   error
		role  *govcd.Role
	)

	// Read state
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get Org and role deletion is accessible only in administrator
	adminOrg, err := r.client.Vmware.GetAdminOrgByNameOrId(r.client.GetOrgName())
	if err != nil {
		resp.Diagnostics.AddError("[role update] Error retrieving Org", err.Error())
		return
	}

	// Get the role
	if state.ID.IsNull() {
		role, err = adminOrg.GetRoleByName(state.Name.ValueString())
	} else {
		role, err = adminOrg.GetRoleById(state.ID.ValueString())
	}
	if err != nil {
		if govcd.ContainsNotFound(err) {
			tflog.Info(ctx, "[DEBUG] Unable to find role. Removing from tfstate")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("[role delete] Error retrieving role", err.Error())
		return
	}
	err = role.Delete()
	if err != nil {
		resp.Diagnostics.AddError("[role delete] Error deleting role", err.Error())
		return
	}
}

func (r *iamRoleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var (
		plan, state *iamRoleResourceModel
		err         error
		role        *govcd.Role
	)

	// Get state and plan
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get Org and role update is accessible only in administrator
	adminOrg, err := r.client.Vmware.GetAdminOrgByNameOrId(r.client.GetOrgName())
	if err != nil {
		resp.Diagnostics.AddError("[role update] Error retrieving Org", err.Error())
		return
	}

	// Get the role
	role, err = adminOrg.GetRoleById(state.ID.ValueString())
	if err != nil {
		if govcd.ContainsNotFound(err) {
			tflog.Info(ctx, "[DEBUG] Unable to find role. Removing from tfstate")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("[role update] Error retrieving role", err.Error())
		return
	}

	// Update the role Name or Description
	if (plan.Name.Equal(state.Name)) || (plan.Description.Equal(state.Description)) {
		role.Role.Name = plan.Name.ValueString()
		role.Role.Description = plan.Description.ValueString()
		_, err = role.Update()
		if err != nil {
			resp.Diagnostics.AddError("[role update] Error updating role", err.Error())
			return
		}
	}

	// Check rights are valid
	rights := make([]govcdtypes.OpenApiReference, 0)
	for _, right := range plan.Rights.Elements() {
		rg := strings.Trim(right.String(), "\"")
		x, err := adminOrg.GetRightByName(rg)
		if err != nil {
			resp.Diagnostics.AddError("[role update] Error retrieving right", err.Error())
			return
		}
		rights = append(rights, govcdtypes.OpenApiReference{Name: rg, ID: x.ID})
	}
	// Add implied rights
	missingImpliedRights, err := govcd.FindMissingImpliedRights(&r.client.Vmware.Client, rights)
	if err != nil {
		resp.Diagnostics.AddError("[role create] Error retrieving implied rights", err.Error())
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
			resp.Diagnostics.AddError("[role update] Error updating role rights", err.Error())
			return
		}
	} else {
		currentRights, err := role.GetRights(nil)
		if err != nil {
			resp.Diagnostics.AddError("[role update] Error retrieving role rights", err.Error())
			return
		}
		if len(currentRights) > 0 {
			err = role.RemoveAllRights()
			if err != nil {
				resp.Diagnostics.AddError("[role update] Error removing role rights", err.Error())
				return
			}
		}
	}

	// Set Plan state
	plan = &iamRoleResourceModel{
		ID:          types.StringValue(role.Role.ID),
		Name:        types.StringValue(role.Role.Name),
		BundleKey:   types.StringValue(role.Role.BundleKey),
		ReadOnly:    types.BoolValue(role.Role.ReadOnly),
		Description: types.StringValue(role.Role.Description),
		Rights:      plan.Rights,
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

//go:generate go run github.com/FrangipaneTeam/tf-doc-extractor@latest -filename $GOFILE -example-dir ../../../examples -resource
func (r *iamRoleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
