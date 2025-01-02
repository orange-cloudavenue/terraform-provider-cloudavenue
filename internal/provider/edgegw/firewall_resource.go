// Package edgegw provides a Terraform resource.
package edgegw

import (
	"context"
	"fmt"

	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/edgegw"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/mutex"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/org"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &firewallResource{}
	_ resource.ResourceWithConfigure   = &firewallResource{}
	_ resource.ResourceWithImportState = &firewallResource{}
)

// NewFirewallResource is a helper function to simplify the provider implementation.
func NewFirewallResource() resource.Resource {
	return &firewallResource{}
}

// firewallResource is the resource implementation.
type firewallResource struct {
	client *client.CloudAvenue
	org    org.Org
	edgegw edgegw.EdgeGateway
}

// Init Initializes the resource.
func (r *firewallResource) Init(ctx context.Context, rm *firewallModel) (diags diag.Diagnostics) {
	var err error

	r.org, diags = org.Init(r.client)
	if diags.HasError() {
		return
	}

	r.edgegw, err = r.org.GetEdgeGateway(edgegw.BaseEdgeGW{
		ID:   rm.EdgeGatewayID.StringValue,
		Name: rm.EdgeGatewayName.StringValue,
	})
	if err != nil {
		diags.AddError("Error retrieving Edge Gateway", err.Error())
		return
	}

	return
}

// Metadata returns the resource type name.
func (r *firewallResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_firewall"
}

// Schema defines the schema for the resource.
func (r *firewallResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = firewallSchema(ctx).GetResource(ctx)
}

func (r *firewallResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *firewallResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) { //nolint:dupl
	defer metrics.New("cloudavenue_edgegateway_firewall", r.client.GetOrgName(), metrics.Create)()

	plan := &firewallModel{}

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

	// Create or update the resource
	resp.Diagnostics.Append(r.createOrUpdate(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	stateRefreshed, found, d := r.read(ctx, plan)
	if !found {
		resp.State.RemoveResource(ctx)
		resp.Diagnostics.AddWarning("Resource not found", fmt.Sprintf("Unable to find firewall on edge %s", plan.EdgeGatewayName.Get()))
		return
	}
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated state into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, stateRefreshed)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *firewallResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer metrics.New("cloudavenue_edgegateway_firewall", r.client.GetOrgName(), metrics.Update)()

	var (
		plan  = new(firewallModel)
		state = new(firewallModel)
	)

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
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

	// Use generic createOrUpdate function to update the resource
	resp.Diagnostics.Append(r.createOrUpdate(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use generic read function to refresh the state
	stateRefreshed, _, d := r.read(ctx, plan)
	resp.Diagnostics.Append(d...)
	if d.HasError() {
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, stateRefreshed)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *firewallResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer metrics.New("cloudavenue_edgegateway_firewall", r.client.GetOrgName(), metrics.Delete)()

	state := &firewallModel{}

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

	vdcOrVDCGroup, err := r.edgegw.GetParent()
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving Edge Gateway parent", err.Error())
		return
	}

	if vdcOrVDCGroup.IsVDCGroup() {
		mutex.GlobalMutex.KvLock(ctx, vdcOrVDCGroup.GetID())
		defer mutex.GlobalMutex.KvUnlock(ctx, vdcOrVDCGroup.GetID())
	} else {
		mutex.GlobalMutex.KvLock(ctx, r.edgegw.GetID())
		defer mutex.GlobalMutex.KvUnlock(ctx, r.edgegw.GetID())
	}
	fwRules, err := r.edgegw.GetNsxtFirewall()
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving Edge Gateway Firewall", err.Error())
		return
	}

	if err := fwRules.DeleteAllRules(); err != nil {
		resp.Diagnostics.AddError("Error deleting Edge Gateway Firewall", err.Error())
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *firewallResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer metrics.New("cloudavenue_edgegateway_firewall", r.client.GetOrgName(), metrics.Read)()

	state := &firewallModel{}

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

func (r *firewallResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	defer metrics.New("cloudavenue_edgegateway_firewall", r.client.GetOrgName(), metrics.Import)()

	var (
		edgegwID   string
		edgegwName string
		d          diag.Diagnostics
		err        error
	)

	r.org, d = org.Init(r.client)
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	if urn.IsValid(req.ID) {
		edgegwID = urn.Normalize(urn.Gateway, req.ID).String()
	} else {
		edgegwName = req.ID
	}

	r.edgegw, err = r.org.GetEdgeGateway(edgegw.BaseEdgeGW{
		ID:   types.StringValue(edgegwID),
		Name: types.StringValue(edgegwName),
	})
	if err != nil {
		resp.Diagnostics.AddError("Failed to import firewall.", err.Error())
		return
	}

	state := &firewallModel{}
	state.ID.Set(r.edgegw.GetID())
	state.EdgeGatewayID.Set(r.edgegw.GetID())
	state.EdgeGatewayName.Set(r.edgegw.GetName())

	// Refresh the state
	stateRefreshed, found, d := r.read(ctx, state)
	if !found {
		resp.Diagnostics.AddError("Failed to import firewall.", fmt.Sprintf("Unable to find firewall on edge %s", r.edgegw.GetName()))
		return
	}
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, stateRefreshed)...)
}

// * custom functions

// createOrUpdate creates or updates the resource and sets the Terraform state.
func (r *firewallResource) createOrUpdate(ctx context.Context, plan *firewallModel) (diags diag.Diagnostics) {
	// Lock object VDC or VDC Group
	vdcOrVDCGroup, err := r.edgegw.GetParent()
	if err != nil {
		diags.AddError("Error retrieving Edge Gateway parent", err.Error())
		return
	}

	if vdcOrVDCGroup.IsVDCGroup() {
		mutex.GlobalMutex.KvLock(ctx, vdcOrVDCGroup.GetID())
		defer mutex.GlobalMutex.KvUnlock(ctx, vdcOrVDCGroup.GetID())
	} else {
		mutex.GlobalMutex.KvLock(ctx, r.edgegw.GetID())
		defer mutex.GlobalMutex.KvUnlock(ctx, r.edgegw.GetID())
	}

	// Set the rules
	fwRules, d := plan.rulesToNsxtFirewallRule(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return
	}

	if _, err := r.edgegw.UpdateNsxtFirewall(&govcdtypes.NsxtFirewallRuleContainer{
		UserDefinedRules: fwRules,
	}); err != nil {
		diags.AddError("Error to create Firewall", err.Error())
		return
	}

	return
}

// read is a generic read function for the resource.
func (r *firewallResource) read(ctx context.Context, planOrState *firewallModel) (stateRefreshed *firewallModel, found bool, diags diag.Diagnostics) {
	stateRefreshed = planOrState.Copy()

	fwRules, err := r.edgegw.GetNsxtFirewall()
	if err != nil {
		if govcd.IsNotFound(err) {
			return stateRefreshed, false, nil
		}
		diags.AddError("Error retrieving Edge Gateway Firewall", err.Error())
		return stateRefreshed, true, diags
	}

	stateRefreshed.ID.Set(r.edgegw.GetID())
	stateRefreshed.EdgeGatewayID.Set(r.edgegw.GetID())
	stateRefreshed.EdgeGatewayName.Set(r.edgegw.GetName())

	rules := make([]*firewallModelRule, 0)

	if fwRules.NsxtFirewallRuleContainer == nil {
		return stateRefreshed, true, nil
	}

	for _, rule := range fwRules.NsxtFirewallRuleContainer.UserDefinedRules {
		fwRule := &firewallModelRule{
			ID:                supertypes.NewStringNull(),
			Name:              supertypes.NewStringNull(),
			Enabled:           supertypes.NewBoolNull(),
			Direction:         supertypes.NewStringNull(),
			IPProtocol:        supertypes.NewStringNull(),
			Action:            supertypes.NewStringNull(),
			Logging:           supertypes.NewBoolNull(),
			SourceIDs:         supertypes.NewSetValueOfNull[string](ctx),
			DestinationIDs:    supertypes.NewSetValueOfNull[string](ctx),
			AppPortProfileIDs: supertypes.NewSetValueOfNull[string](ctx),
		}
		fwRule.ID.Set(rule.ID)
		fwRule.Name.Set(rule.Name)
		fwRule.Enabled.Set(rule.Enabled)
		fwRule.Direction.Set(rule.Direction)
		fwRule.IPProtocol.Set(rule.IpProtocol)
		fwRule.Action.Set(rule.Action)
		fwRule.Logging.Set(rule.Logging)
		fwRule.SourceIDs.Set(ctx, common.FromOpenAPIReferenceID(ctx, rule.SourceFirewallGroups))
		fwRule.DestinationIDs.Set(ctx, common.FromOpenAPIReferenceID(ctx, rule.DestinationFirewallGroups))
		fwRule.AppPortProfileIDs.Set(ctx, common.FromOpenAPIReferenceID(ctx, rule.ApplicationPortProfiles))
		rules = append(rules, fwRule)
	}

	diags.Append(stateRefreshed.Rules.Set(ctx, rules)...)
	return stateRefreshed, true, diags
}
