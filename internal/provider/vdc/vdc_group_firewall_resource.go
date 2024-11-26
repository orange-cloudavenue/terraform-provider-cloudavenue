package vdc

import (
	"context"
	"fmt"

	"github.com/k0kubun/pp"

	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"

	v1 "github.com/orange-cloudavenue/cloudavenue-sdk-go/v1"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &GroupFirewallResource{}
	_ resource.ResourceWithConfigure   = &GroupFirewallResource{}
	_ resource.ResourceWithImportState = &GroupFirewallResource{}
	// _ resource.ResourceWithModifyPlan     = &GroupFirewallResource{}
	// _ resource.ResourceWithUpgradeState   = &GroupFirewallResource{}
	// _ resource.ResourceWithValidateConfig = &GroupFirewallResource{}.
)

// NewGroupFirewallResource is a helper function to simplify the provider implementation.
func NewGroupFirewallResource() resource.Resource {
	return &GroupFirewallResource{}
}

// GroupFirewallResource is the resource implementation.
type GroupFirewallResource struct {
	client *client.CloudAvenue
	vdcg   *v1.VDCGroup
}

// Init Initializes the resource.
func (r *GroupFirewallResource) Init(ctx context.Context, rm *groupFirewallModel) (diags diag.Diagnostics) {
	var err error

	r.vdcg, err = r.client.CAVSDK.V1.VDC().GetVDCGroup(rm.VDCGroup.Get())
	if err != nil {
		diags.AddError("Error retrieving VDC Group", err.Error())
		return
	}

	return
}

// Metadata returns the resource type name.
func (r *GroupFirewallResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_group_firewall"
}

// Schema defines the schema for the resource.
func (r *GroupFirewallResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = groupFirewallSchema(ctx).GetResource(ctx)
}

func (r *GroupFirewallResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *GroupFirewallResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer metrics.New("cloudavenue_vdc_group_firewall", r.client.GetOrgName(), metrics.Create)()

	plan := &groupFirewallModel{}

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

	rules, d := plan.rulesToSDKRules(ctx)
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	if err := r.vdcg.Refresh(); err != nil {
		resp.Diagnostics.AddError("Error refreshing VDC Group", err.Error())
		return
	}

	_, err := r.vdcg.CreateFirewall(rules)
	if err != nil {
		resp.Diagnostics.AddError("Error creating VDC Group Firewall", err.Error())
		return
	}

	// Use generic read function to refresh the state
	state, found, d := r.read(ctx, plan)
	if !found {
		resp.Diagnostics.AddError("Fail to retrieve VDC Group Firewall after creation", "VDC Group Firewall not found")
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
func (r *GroupFirewallResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer metrics.New("cloudavenue_vdc_group_firewall", r.client.GetOrgName(), metrics.Read)()

	state := &groupFirewallModel{}

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
func (r *GroupFirewallResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer metrics.New("cloudavenue_vdc_group_firewall", r.client.GetOrgName(), metrics.Update)()

	var (
		plan  = &groupFirewallModel{}
		state = &groupFirewallModel{}
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

	vdcgfw, err := r.vdcg.GetFirewall()
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving VDC Group Firewall", err.Error())
		return
	}

	rules, d := plan.rulesToSDKRules(ctx)
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	if err := vdcgfw.UpdateFirewall(rules); err != nil {
		resp.Diagnostics.AddError("Error updating VDC Group Firewall rules", err.Error())
		return
	}

	// Use generic read function to refresh the state
	stateRefreshed, found, d := r.read(ctx, plan)
	if !found {
		resp.Diagnostics.AddError("Fail to retrieve VDC Group Firewall after update", "VDC Group Firewall not found")
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
func (r *GroupFirewallResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer metrics.New("cloudavenue_vdc_group_firewall", r.client.GetOrgName(), metrics.Delete)()

	state := &groupFirewallModel{}

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

	vdcgfw, err := r.vdcg.GetFirewall()
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving VDC Group Firewall", err.Error())
		return
	}

	if err := vdcgfw.Delete(); err != nil {
		resp.Diagnostics.AddError("Error deleting VDC Group Firewall", err.Error())
		return
	}
}

func (r *GroupFirewallResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	defer metrics.New("cloudavenue_vdc_group_firewall", r.client.GetOrgName(), metrics.Import)()

	state := &groupFirewallModel{
		VDCGroup: supertypes.NewStringNull(),
		Enabled:  supertypes.NewBoolNull(),
		Rules:    supertypes.NewListNestedObjectValueOfNull[groupFirewallModelRule](ctx),
	}

	state.VDCGroup.Set(req.ID)

	resp.Diagnostics.Append(r.Init(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use generic read function to refresh the state
	stateRefreshed, found, d := r.read(ctx, state)
	if !found {
		resp.Diagnostics.AddError("Fail to retrieve VDC Group Firewall after import", "VDC Group Firewall not found")
		return
	}
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, stateRefreshed)...)
}

// * CustomFuncs

// read is a generic read function that can be used by the resource Create, Read and Update functions.
func (r *GroupFirewallResource) read(ctx context.Context, planOrState *groupFirewallModel) (stateRefreshed *groupFirewallModel, found bool, diags diag.Diagnostics) {
	stateRefreshed = planOrState.Copy()

	vdcgfw, err := r.vdcg.GetFirewall()
	if err != nil {
		if govcd.ContainsNotFound(err) {
			return nil, false, nil
		}
		diags.AddError("Error retrieving VDC Group Firewall", err.Error())
		return
	}

	if !stateRefreshed.ID.IsKnown() {
		// firewall don't have an ID, use the VDC Group ID instead
		stateRefreshed.ID.Set(r.vdcg.GetID())
	}

	// * Enabled
	isEnabled, err := vdcgfw.IsEnabled()
	if err != nil {
		diags.AddError("Error retrieving VDC Group Firewall enabled status", err.Error())
		return
	}
	stateRefreshed.Enabled.Set(isEnabled)

	tflog.Debug(ctx, pp.Sprint(vdcgfw.GetRules()))
	// * Rules
	rules, d := stateRefreshed.sdkRulesToRules(ctx, vdcgfw.GetRules())
	if d.HasError() {
		diags.Append(d...)
		return stateRefreshed, true, diags
	}

	diags.Append(stateRefreshed.Rules.Set(ctx, rules)...)

	return stateRefreshed, true, diags
}
