// Package catalog provides a Terraform resource.
package catalog

import (
	"context"
	"fmt"

	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/adminorg"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/uuid"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &aclResource{}
	_ resource.ResourceWithConfigure   = &aclResource{}
	_ resource.ResourceWithImportState = &aclResource{}
	// _ resource.ResourceWithModifyPlan     = &aclResource{}
	// _ resource.ResourceWithUpgradeState   = &aclResource{}
	// _ resource.ResourceWithValidateConfig = &aclResource{}.
)

// NewACLResource is a helper function to simplify the provider implementation.
func NewACLResource() resource.Resource {
	return &aclResource{}
}

// aclResource is the resource implementation.
type aclResource struct {
	client   *client.CloudAvenue
	adminOrg adminorg.AdminOrg
	catalog  base
}

// Init Initializes the resource.
func (r *aclResource) Init(ctx context.Context, rm *ACLModel) (diags diag.Diagnostics) {
	r.catalog = base{
		id:   rm.CatalogID.Get(),
		name: rm.CatalogName.Get(),
	}

	r.adminOrg, diags = adminorg.Init(r.client)
	return
}

// Metadata returns the resource type name.
func (r *aclResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_acl"
}

// Schema defines the schema for the resource.
func (r *aclResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = aclSchema(ctx).GetResource(ctx)
}

func (r *aclResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *aclResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &ACLModel{}

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

	resp.Diagnostics.Append(r.createOrUpdate(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, _, d := r.read(ctx, plan)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *aclResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &ACLModel{}

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
		Implement the resource read here
	*/

	stateRefreshed, found, d := r.read(ctx, state)
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, stateRefreshed)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *aclResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var (
		plan  = &ACLModel{}
		state = &ACLModel{}
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

	resp.Diagnostics.Append(r.createOrUpdate(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	stateRefreshed, _, d := r.read(ctx, plan)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, stateRefreshed)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *aclResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	state := &ACLModel{}

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

	catalog, err := r.adminOrg.GetCatalogByNameOrId(r.catalog.GetIDOrName(), false)
	if err != nil {
		resp.Diagnostics.AddError("error when getting catalog", err.Error())
		return
	}

	if err := catalog.RemoveAccessControl(true); err != nil {
		resp.Diagnostics.AddError("error when removing ACL", err.Error())
	}
}

func (r *aclResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// import format is catalogIDOrName

	var d diag.Diagnostics

	r.adminOrg, d = adminorg.Init(r.client)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	catalog, err := r.adminOrg.GetCatalogByNameOrId(req.ID, true)
	if err != nil {
		resp.Diagnostics.AddError("error when retrieving catalog", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), catalog.Catalog.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("catalog_id"), catalog.Catalog.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("catalog_name"), catalog.Catalog.Name)...)
}

// * Custom Funcs

// read the ACL from the API.
func (r *aclResource) read(ctx context.Context, planOrState *ACLModel) (stateRefreshed *ACLModel, found bool, diags diag.Diagnostics) {
	stateRefreshed = planOrState.Copy()

	catalog, err := r.adminOrg.GetCatalogByNameOrId(r.catalog.GetIDOrName(), true)
	if err != nil {
		if govcd.ContainsNotFound(err) {
			return stateRefreshed, false, diags
		}
		diags.AddError("error when retrieving catalog", err.Error())
		return stateRefreshed, true, diags
	}

	acl, err := catalog.GetAccessControl(true)
	if err != nil {
		if govcd.ContainsNotFound(err) {
			return stateRefreshed, false, diags
		}
		diags.AddError("error when retrieving access control", err.Error())
		return stateRefreshed, true, diags
	}

	stateRefreshed.ID.Set(catalog.Catalog.ID)
	stateRefreshed.CatalogID.Set(catalog.Catalog.ID)
	stateRefreshed.CatalogName.Set(catalog.Catalog.Name)
	stateRefreshed.SharedWithEveryone.Set(acl.IsSharedToEveryone)
	if acl.EveryoneAccessLevel != nil {
		stateRefreshed.EveryoneAccessLevel.Set(*acl.EveryoneAccessLevel)
	} else {
		stateRefreshed.EveryoneAccessLevel.SetNull()
	}

	sharedWithUsers := make(ACLModelSharedWithUsers, 0)
	if acl.AccessSettings != nil {
		for _, user := range acl.AccessSettings.AccessSetting {
			// Get the UUID from the HREF
			id, err := govcd.GetUuidFromHref(user.Subject.HREF, true)
			if err != nil {
				diags.AddError("unable to get UUID for user", err.Error())
				return stateRefreshed, true, diags
			}
			x := ACLModelSharedWithUser{}
			x.UserID.Set(uuid.Normalize(uuid.User, id).String())
			x.AccessLevel.Set(user.AccessLevel)

			sharedWithUsers = append(sharedWithUsers, x)
		}
		diags.Append(stateRefreshed.SharedWithUsers.Set(ctx, sharedWithUsers)...)
	} else {
		stateRefreshed.SharedWithUsers.SetNull(ctx)
	}

	return stateRefreshed, true, diags
}

func (r *aclResource) createOrUpdate(ctx context.Context, plan *ACLModel) (diags diag.Diagnostics) {
	catalog, err := r.adminOrg.GetCatalogByNameOrId(r.catalog.GetIDOrName(), true)
	if err != nil {
		diags.AddError("error when retrieving catalog", err.Error())
		return
	}

	aclConfig, d := plan.ToControlAccessParams(ctx, r.adminOrg)
	diags.Append(d...)
	if diags.HasError() {
		return
	}

	if err := catalog.SetAccessControl(&aclConfig, true); err != nil {
		diags.AddError("error when setting access control", err.Error())
		return
	}

	return
}
