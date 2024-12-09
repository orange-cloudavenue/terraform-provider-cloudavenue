package vdc

import (
	"context"
	"fmt"

	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/adminorg"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &groupResource{}
	_ resource.ResourceWithConfigure   = &groupResource{}
	_ resource.ResourceWithImportState = &groupResource{}
)

// NewGroupResource is a helper function to simplify the provider implementation.
func NewGroupResource() resource.Resource {
	return &groupResource{}
}

// groupResource is the resource implementation.
type groupResource struct {
	client   *client.CloudAvenue
	adminOrg adminorg.AdminOrg
}

// Init Initializes the resource.
func (r *groupResource) Init(ctx context.Context, rm *GroupModel) (diags diag.Diagnostics) {
	r.adminOrg, diags = adminorg.Init(r.client)
	return
}

// Metadata returns the resource type name.
func (r *groupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_group"
}

// Schema defines the schema for the resource.
func (r *groupResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = groupSchema().GetResource(ctx)
}

func (r *groupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *groupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer metrics.New("cloudavenue_vdc_group", r.client.GetOrgName(), metrics.Create)()

	plan := &GroupModel{}

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

	vdcIDs, d := plan.GetVDCIds(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	vdcGroup, err := r.adminOrg.CreateNsxtVdcGroup(plan.Name.Get(), plan.Description.Get(), vdcIDs.Get()[0], vdcIDs.Get())
	if err != nil {
		resp.Diagnostics.AddError("Error creating VDC Group", err.Error())
		return
	}

	plan.ID.Set(vdcGroup.VdcGroup.Id)
	state, _, d := r.read(ctx, plan)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *groupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer metrics.New("cloudavenue_vdc_group", r.client.GetOrgName(), metrics.Read)()

	state := &GroupModel{}

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

	statRefreshed, found, d := r.read(ctx, state)
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, statRefreshed)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *groupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer metrics.New("cloudavenue_vdc_group", r.client.GetOrgName(), metrics.Update)()

	var (
		plan  = &GroupModel{}
		state = &GroupModel{}
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

	vdcIDs, d := plan.GetVDCIds(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Here use GetVdcGroupById instead of GetVdcGroupByNameOrID because we want to update the name of VDC Group
	vdcGroup, err := r.adminOrg.GetVdcGroupById(state.ID.Get())
	if err != nil {
		resp.Diagnostics.AddError("Error reading VDC Group", err.Error())
		return
	}

	if _, err := vdcGroup.Update(plan.Name.Get(), plan.Description.Get(), vdcIDs.Get()); err != nil {
		resp.Diagnostics.AddError("Error updating VDC Group", err.Error())
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
func (r *groupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer metrics.New("cloudavenue_vdc_group", r.client.GetOrgName(), metrics.Delete)()

	state := &GroupModel{}

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

	vdcGroup, err := r.adminOrg.GetVDCGroupByNameOrID(state.GetNameOrID())
	if err != nil {
		resp.Diagnostics.AddError("Error reading VDC Group", err.Error())
		return
	}

	if err = vdcGroup.Delete(); err != nil {
		resp.Diagnostics.AddError("Error deleting VDC Group", err.Error())
		return
	}
}

func (r *groupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	defer metrics.New("cloudavenue_vdc_group", r.client.GetOrgName(), metrics.Import)()

	// id format is vdcGroupIDOrName

	var (
		d   diag.Diagnostics
		err error
	)

	r.adminOrg, d = adminorg.Init(r.client)
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	vdcGroup, err := r.adminOrg.GetVDCGroupByNameOrID(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving VDC Group", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), vdcGroup.VdcGroup.Id)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), vdcGroup.VdcGroup.Name)...)
}

// * Custom Functions.
// read is a generic function to read a resource.
func (r *groupResource) read(ctx context.Context, planOrState *GroupModel) (stateRefreshed *GroupModel, found bool, diags diag.Diagnostics) {
	stateRefreshed = planOrState.Copy()

	vdcGroup, err := r.adminOrg.GetVDCGroupByNameOrID(planOrState.GetNameOrID())
	if err != nil {
		if govcd.ContainsNotFound(err) {
			return nil, false, nil
		}
		diags.AddError("Error reading VDC Group", err.Error())
		return nil, true, diags
	}

	stateRefreshed.ID.Set(vdcGroup.VdcGroup.Id)
	stateRefreshed.Name.Set(vdcGroup.VdcGroup.Name)
	stateRefreshed.Description = utils.SuperStringValueOrNull(vdcGroup.VdcGroup.Description)
	stateRefreshed.Status.Set(vdcGroup.VdcGroup.Status)
	stateRefreshed.Type.Set(vdcGroup.VdcGroup.Type)

	var vdcIDs []string
	for _, vdc := range vdcGroup.VdcGroup.ParticipatingOrgVdcs {
		vdcIDs = append(vdcIDs, vdc.VdcRef.ID)
	}

	diags.Append(stateRefreshed.VDCIds.Set(ctx, vdcIDs)...)

	return stateRefreshed, true, diags
}
