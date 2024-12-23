package vdcg

import (
	"context"
	"fmt"
	"strings"

	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
	v1 "github.com/orange-cloudavenue/cloudavenue-sdk-go/v1"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/mutex"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &IPSetResource{}
	_ resource.ResourceWithConfigure   = &IPSetResource{}
	_ resource.ResourceWithImportState = &IPSetResource{}
)

// NewIPSetResource is a helper function to simplify the provider implementation.
func NewIPSetResource() resource.Resource {
	return &IPSetResource{}
}

// IPSetResource is the resource implementation.
type IPSetResource struct {
	client   *client.CloudAvenue
	vdcGroup *v1.VDCGroup
}

// Init Initializes the resource.
func (r *IPSetResource) Init(ctx context.Context, rm *IPSetModel) (diags diag.Diagnostics) {
	var err error

	idOrName := rm.VDCGroupName.Get()
	if rm.VDCGroupID.IsKnown() {
		// Use the ID
		idOrName = rm.VDCGroupID.Get()
	}

	r.vdcGroup, err = r.client.CAVSDK.V1.VDC().GetVDCGroup(idOrName)
	if err != nil {
		diags.AddError("Error retrieving VDC Group", err.Error())
		return
	}
	return
}

// Metadata returns the resource type name.
func (r *IPSetResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_ip_set"
}

// Schema defines the schema for the resource.
func (r *IPSetResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ipSetSchema(ctx).GetResource(ctx)
}

func (r *IPSetResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *IPSetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer metrics.New("cloudavenue_vdcg_ip_set", r.client.GetOrgName(), metrics.Create)()

	plan := &IPSetModel{}

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

	mutex.GlobalMutex.KvLock(ctx, r.vdcGroup.GetID())
	defer mutex.GlobalMutex.KvUnlock(ctx, r.vdcGroup.GetID())

	values, d := plan.ToSDKIPSetModel(ctx)
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	fwipset, err := r.vdcGroup.CreateFirewallIPSet(values)
	if err != nil {
		resp.Diagnostics.AddError("Error creating ip set", err.Error())
		return
	}

	// Set the ID
	plan.ID.Set(fwipset.ID)

	// Use generic read function to refresh the state
	state, found, d := r.read(ctx, plan)
	if !found {
		resp.Diagnostics.AddError("IP Set not found", "The ip set was not found after creation")
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
func (r *IPSetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer metrics.New("cloudavenue_vdcg_ip_set", r.client.GetOrgName(), metrics.Read)()

	state := &IPSetModel{}

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
func (r *IPSetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer metrics.New("cloudavenue_vdcg_ip_set", r.client.GetOrgName(), metrics.Update)()

	var (
		plan  = &IPSetModel{}
		state = &IPSetModel{}
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

	mutex.GlobalMutex.KvLock(ctx, r.vdcGroup.GetID())
	defer mutex.GlobalMutex.KvUnlock(ctx, r.vdcGroup.GetID())

	fwipset, err := r.vdcGroup.GetFirewallIPSet(state.ID.Get())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving ip set", err.Error())
		return
	}

	values, d := plan.ToSDKIPSetModel(ctx)
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	if err := fwipset.Update(values); err != nil {
		resp.Diagnostics.AddError("Error updating ip set", err.Error())
		return
	}

	// Use generic read function to refresh the state
	stateRefreshed, found, d := r.read(ctx, plan)
	if !found {
		resp.Diagnostics.AddError("IP Set not found", "The ip set was not found after update")
		return
	}
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, stateRefreshed)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *IPSetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer metrics.New("cloudavenue_vdcg_ip_set", r.client.GetOrgName(), metrics.Delete)()

	state := &IPSetModel{}

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

	mutex.GlobalMutex.KvLock(ctx, r.vdcGroup.GetID())
	defer mutex.GlobalMutex.KvUnlock(ctx, r.vdcGroup.GetID())

	fwipset, err := r.vdcGroup.GetFirewallIPSet(state.ID.Get())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving ip set", err.Error())
		return
	}

	if err := fwipset.Delete(); err != nil {
		resp.Diagnostics.AddError("Error deleting ip set", err.Error())
		return
	}
}

func (r *IPSetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	defer metrics.New("cloudavenue_vdcg_ip_set", r.client.GetOrgName(), metrics.Import)()

	// * Import with custom logic
	idParts := strings.Split(req.ID, ".")

	if len(idParts) != 2 {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: VDCGroupNameOrID.IPSetNameOrID Got: %q", req.ID),
		)
		return
	}
	vdcGroupNameOrID, ipSetNameOrID := idParts[0], idParts[1]

	x := &IPSetModel{
		ID:           supertypes.NewStringNull(),
		Name:         supertypes.NewStringNull(),
		VDCGroupName: supertypes.NewStringNull(),
		VDCGroupID:   supertypes.NewStringNull(),
		Description:  supertypes.NewStringNull(),
		IPAddresses:  supertypes.NewSetValueOfNull[string](ctx),
	}

	if urn.IsVDCGroup(vdcGroupNameOrID) {
		x.VDCGroupID.Set(vdcGroupNameOrID)
	} else {
		x.VDCGroupName.Set(vdcGroupNameOrID)
	}

	if urn.IsSecurityGroup(ipSetNameOrID) {
		x.ID.Set(ipSetNameOrID)
	} else {
		x.Name.Set(ipSetNameOrID)
	}

	resp.Diagnostics.Append(r.Init(ctx, x)...)
	if resp.Diagnostics.HasError() {
		return
	}

	stateRefreshed, found, d := r.read(ctx, x)
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

// * CustomFuncs

// read is a generic read function that can be used by the resource Create, Read and Update functions.
func (r *IPSetResource) read(ctx context.Context, planOrState *IPSetModel) (stateRefreshed *IPSetModel, found bool, diags diag.Diagnostics) {
	stateRefreshed = planOrState.Copy()

	/*
		Implement the resource read here
	*/

	idOrName := planOrState.Name.Get()
	if planOrState.ID.IsKnown() {
		idOrName = planOrState.ID.Get()
	}

	fwipset, err := r.vdcGroup.GetFirewallIPSet(idOrName)
	if govcd.ContainsNotFound(err) {
		return nil, false, nil
	}
	if err != nil {
		diags.AddError("Error retrieving ip set", err.Error())
		return nil, true, diags
	}

	stateRefreshed.ID.Set(fwipset.ID)
	stateRefreshed.Name.Set(fwipset.Name)
	stateRefreshed.Description.Set(fwipset.Description)
	stateRefreshed.VDCGroupName.Set(r.vdcGroup.GetName())
	stateRefreshed.VDCGroupID.Set(r.vdcGroup.GetID())

	if fwipset.IPAddresses != nil || len(fwipset.IPAddresses) > 0 {
		diags.Append(stateRefreshed.IPAddresses.Set(ctx, fwipset.IPAddresses)...)
	} else {
		stateRefreshed.IPAddresses.SetNull(ctx)
	}

	return stateRefreshed, true, diags
}
