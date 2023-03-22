// ! REMOVE THIS RESOURCE

package iam

import (
	"context"
	"fmt"

	fstringplanmodifier "github.com/FrangipaneTeam/terraform-plugin-framework-planmodifiers/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
)

var (
	_ resource.Resource                = &iamGroupResource{}
	_ resource.ResourceWithConfigure   = &iamGroupResource{}
	_ resource.ResourceWithImportState = &iamGroupResource{}
)

// NewiamGroupResource is a helper function to simplify the provider implementation.
func NewIAMGroupResource() resource.Resource {
	return &iamGroupResource{}
}

// iamGroupResource is the resource implementation.
type iamGroupResource struct {
	client *client.CloudAvenue
}

type iamGroupResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Role        types.String `tfsdk:"role"`
	UserNames   types.List   `tfsdk:"user_names"`
}

// Metadata returns the resource type name.
func (r *iamGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_" + "group"
}

// Schema defines the schema for the resource.
func (r *iamGroupResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Provides a Cloud Avenue IAM group. This can be used to create, update, and delete iam groups.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID is a unique identifier for the iam group",
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "A name for the iam group",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Computed:            true,
				Optional:            true,
				MarkdownDescription: "Description of the iam group",
				PlanModifiers: []planmodifier.String{
					fstringplanmodifier.SetDefaultEmptyString(),
				},
			},
			"role": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The role to assign to the iam group",
			},
			"user_names": schema.ListAttribute{
				ElementType:         types.StringType,
				Computed:            true,
				MarkdownDescription: "Set of user names that belong to the iam group",
			},
		},
	}
}

func (r *iamGroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *iamGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var (
		plan *iamGroupResourceModel
		err  error
	)

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// group creation is accessible only for administrator account
	adminOrg, err := r.client.Vmware.GetAdminOrgByNameOrId(r.client.GetOrgName())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving Org", err.Error())
		return
	}

	// Get role reference
	roleRef, err := adminOrg.GetRoleReference(plan.Role.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving role reference", err.Error())
		return
	}

	// Create the org group
	newGroup := govcd.NewGroup(&r.client.Vmware.Client, adminOrg)
	groupDefinition := govcdtypes.Group{
		Name:         plan.Name.ValueString(),
		Role:         roleRef,
		ProviderType: "SAML",
		Description:  plan.Description.ValueString(),
	}

	newGroup.Group = &groupDefinition

	createGroup, err := adminOrg.CreateGroup(newGroup.Group)
	if err != nil {
		resp.Diagnostics.AddError("Error creating Org group", err.Error())
		return
	}

	var userNames []attr.Value

	plan = &iamGroupResourceModel{
		ID:          types.StringValue(createGroup.Group.ID),
		Name:        plan.Name,
		Description: plan.Description,
		Role:        plan.Role,
		UserNames:   types.ListValueMust(types.StringType, userNames),
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *iamGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *iamGroupResourceModel

	// Get current state
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	adminOrg, err := r.client.Vmware.GetAdminOrgByNameOrId(r.client.GetOrgName())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving Org", err.Error())
		return
	}

	var groupID string
	if state.ID.IsNull() {
		groupID = state.Name.ValueString()
	} else {
		groupID = state.ID.ValueString()
	}

	group, err := adminOrg.GetGroupByNameOrId(groupID, false)
	if err != nil {
		if govcd.IsNotFound(err) {
			// Group not found, remove from state
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error retrieving Org group", err.Error())
		return
	}

	// users
	var userNames []attr.Value
	for _, user := range group.Group.UsersList.UserReference {
		userNames = append(userNames, types.StringValue(user.Name))
	}

	state = &iamGroupResourceModel{
		ID:          types.StringValue(group.Group.ID),
		Name:        types.StringValue(group.Group.Name),
		Description: types.StringValue(group.Group.Description),
		Role:        types.StringValue(group.Group.Role.Name),
		UserNames:   types.ListValueMust(types.StringType, userNames),
	}

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *iamGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state *iamGroupResourceModel

	// Get current state
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	adminOrg, err := r.client.Vmware.GetAdminOrgByNameOrId(r.client.GetOrgName())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving Org", err.Error())
		return
	}

	group, err := adminOrg.GetGroupById(state.ID.ValueString(), false)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving Org group", err.Error())
		return
	}

	if !plan.Role.Equal(state.Role) {
		// Get role reference
		roleRef, errGetRoleRef := adminOrg.GetRoleReference(plan.Role.ValueString())
		if errGetRoleRef != nil {
			resp.Diagnostics.AddError("Error retrieving role reference", errGetRoleRef.Error())
			return
		}

		group.Group.Role = roleRef
	}

	if !plan.Description.Equal(state.Description) {
		group.Group.Description = plan.Description.ValueString()
	}

	err = group.Update()
	if err != nil {
		resp.Diagnostics.AddError("Error updating Org group", err.Error())
		return
	}

	plan = &iamGroupResourceModel{
		ID:          state.ID,
		Name:        plan.Name,
		Description: plan.Description,
		Role:        plan.Role,
		UserNames:   state.UserNames,
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *iamGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *iamGroupResourceModel

	// Get current state
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	adminOrg, err := r.client.Vmware.GetAdminOrgByNameOrId(r.client.GetOrgName())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving Org", err.Error())
		return
	}

	group, err := adminOrg.GetGroupById(state.ID.ValueString(), false)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving Org group", err.Error())
		return
	}

	err = group.Delete()
	if err != nil {
		if govcd.IsNotFound(err) {
			return
		}
		resp.Diagnostics.AddError("Error deleting Org group", err.Error())
		return
	}
}

//go:generate go run github.com/FrangipaneTeam/tf-doc-extractor@latest -filename $GOFILE -example-dir ../../../examples -resource
func (r *iamGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
