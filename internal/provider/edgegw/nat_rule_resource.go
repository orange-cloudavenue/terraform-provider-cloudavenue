// Package edgegw provides a Terraform resource.
package edgegw

import (
	"context"
	"fmt"
	"strings"

	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/edgegw"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/mutex"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/org"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &natRuleResource{}
	_ resource.ResourceWithConfigure   = &natRuleResource{}
	_ resource.ResourceWithImportState = &natRuleResource{}
)

// NewNATRuleResource is a helper function to simplify the provider implementation.
func NewNATRuleResource() resource.Resource {
	return &natRuleResource{}
}

// natRuleResource is the resource implementation.
type natRuleResource struct {
	client *client.CloudAvenue
	org    org.Org
	edgegw edgegw.EdgeGateway
}

// Init Initializes the resource.
func (r *natRuleResource) Init(ctx context.Context, rm *NATRuleModel) (diags diag.Diagnostics) {
	var err error

	r.org, diags = org.Init(r.client)
	if diags.HasError() {
		return
	}

	r.edgegw, err = r.org.GetEdgeGateway(edgegw.BaseEdgeGW{
		ID:   types.StringValue(rm.EdgeGatewayID.Get()),
		Name: types.StringValue(rm.EdgeGatewayName.Get()),
	})
	if err != nil {
		diags.AddError("Error retrieving Edge Gateway", err.Error())
		return
	}

	return
}

// Metadata returns the resource type name.
func (r *natRuleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_nat_rule"
}

// Schema defines the schema for the resource.
func (r *natRuleResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = natRuleSchema(ctx).GetResource(ctx)
}

func (r *natRuleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *natRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer metrics.New("cloudavenue_edgegateway_nat_rule", r.client.GetOrgName(), metrics.Create)()

	plan := &NATRuleModel{}

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

	// Lock object EdgeGateway
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

	// Get data from plan
	nsxtNATRule, err := plan.ToNsxtNATRule(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error getting NSX-T NAT rule: %s", err.Error())
		return
	}

	// Create NAT Rule
	rule, err := r.edgegw.CreateNatRule(nsxtNATRule)
	if err != nil {
		resp.Diagnostics.AddError("Error creating NSX-T NAT rule: ", err.Error())
		return
	}

	// Set ID
	plan.ID.Set(rule.NsxtNatRule.ID)

	// read NAT Rule and update state
	stateRefreshed, _, d := r.read(plan)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, stateRefreshed)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *natRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer metrics.New("cloudavenue_edgegateway_nat_rule", r.client.GetOrgName(), metrics.Read)()

	state := &NATRuleModel{}

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

	// read NAT Rule and update state
	stateRefreshed, found, d := r.read(state)
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
func (r *natRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer metrics.New("cloudavenue_edgegateway_nat_rule", r.client.GetOrgName(), metrics.Update)()

	var (
		plan  = &NATRuleModel{}
		state = &NATRuleModel{}
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

	// Lock object EdgeGateway
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

	// Get data to plan
	nsxtNATRule, err := plan.ToNsxtNATRule(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error getting NSX-T NAT rule: %s", err.Error())
		return
	}

	// Get Nat Rule
	existingRule, err := r.edgegw.GetNatRuleById(plan.ID.Get())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving NAT Rule ID", err.Error())
		return
	}

	// Inject ID for update
	nsxtNATRule.ID = existingRule.NsxtNatRule.ID
	if _, err = existingRule.Update(nsxtNATRule); err != nil {
		resp.Diagnostics.AddError("Error updating NSX-T NAT rule: ", err.Error())
		return
	}

	// read NAT Rule and refresh state
	stateRefreshed, _, d := r.read(plan)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, stateRefreshed)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *natRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer metrics.New("cloudavenue_edgegateway_nat_rule", r.client.GetOrgName(), metrics.Delete)()

	state := &NATRuleModel{}

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

	// Lock object EdgeGateway
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

	// Get NAT Rule
	existingRule, err := r.edgegw.GetNatRuleById(state.ID.Get())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving NAT Rule ID", err.Error())
		return
	}

	if err = existingRule.Delete(); err != nil {
		resp.Diagnostics.AddError("Error Deleting NAT Rule ID", err.Error())
		return
	}
}

func (r *natRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var (
		edgegwID, edgegwName string
		d                    diag.Diagnostics
		err                  error
		natRule              *govcd.NsxtNatRule
	)

	defer metrics.New("cloudavenue_edgegateway_nat_rule", r.client.GetOrgName(), metrics.Import)()

	// Split req.ID with dot. ID format is EdgeGatewayIDOrName.NATRuleIDOrName
	idParts := strings.Split(req.ID, ".")

	if len(idParts) != 2 {
		resp.Diagnostics.AddError("Invalid ID format", "ID format is EdgeGatewayIDOrName.NATRuleIDOrName")
		return
	}

	// Get ORG
	r.org, d = org.Init(r.client)
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	// Get EdgeGW is ID or Name
	if urn.IsEdgeGateway(idParts[0]) {
		edgegwID = idParts[0]
	} else {
		edgegwName = idParts[0]
	}

	r.edgegw, err = r.org.GetEdgeGateway(edgegw.BaseEdgeGW{
		ID:   types.StringValue(edgegwID),
		Name: types.StringValue(edgegwName),
	})
	if err != nil {
		resp.Diagnostics.AddError("Failed to import NAT Rule.", err.Error())
		return
	}

	// Check if NATRule is ID or a Name
	if urn.IsUUIDV4(idParts[1]) {
		natRule, err = r.edgegw.GetNatRuleById(idParts[1])
		if err != nil {
			resp.Diagnostics.AddError("Failed to Get NAT Rule.", err.Error())
			return
		}
	} else {
		// Get all NAT Rules
		allRules, err := r.edgegw.GetAllNatRules(nil)
		if err != nil {
			resp.Diagnostics.AddError("Failed to Get NAT Rule.", err.Error())
			return
		}
		// Filter NAT Rule by Name
		allResults := make([]*govcd.NsxtNatRule, 0)
		for _, rule := range allRules {
			if rule.NsxtNatRule.Name == idParts[1] {
				allResults = append(allResults, rule)
			}
		}
		// Check if multiple NAT Rules found with the same name
		if len(allResults) > 1 {
			resp.Diagnostics.AddError("Failed to Get NAT Rule.", fmt.Sprintf("Multiple NAT Rules found with the same name: %s", idParts[1]))
			return
		}
		if len(allResults) == 0 {
			resp.Diagnostics.AddError("Failed to Get NAT Rule.", fmt.Sprintf("No NAT Rule found with the name: %s", idParts[1]))
			return
		}
		natRule = allResults[0]
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), natRule.NsxtNatRule.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), natRule.NsxtNatRule.Name)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("edge_gateway_id"), r.edgegw.GetID())...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("edge_gateway_name"), r.edgegw.GetName())...)
}

func (r *natRuleResource) read(planOrState *NATRuleModel) (stateRefreshed *NATRuleModel, found bool, diags diag.Diagnostics) {
	stateRefreshed = planOrState.Copy()

	var (
		rule *govcd.NsxtNatRule
		err  error
	)

	// Get Nat Rule by Name or ID
	if stateRefreshed.ID.IsKnown() {
		rule, err = r.edgegw.GetNatRuleById(stateRefreshed.ID.Get())
	} else {
		rule, err = r.edgegw.GetNatRuleByName(stateRefreshed.Name.Get())
	}
	if err != nil {
		if govcd.ContainsNotFound(err) {
			return stateRefreshed, false, nil
		}
		diags.AddError("Error retrieving NAT Rule ID", err.Error())
		return
	}

	// Check if ApplicationPortProfileID is known
	var appPortProfile *govcd.NsxtAppPortProfile
	if stateRefreshed.AppPortProfileID.IsKnown() {
		appPortProfile, err = r.org.GetNsxtAppPortProfileById(stateRefreshed.AppPortProfileID.Get())
		if err != nil {
			diags.AddError("Error retrieving NSX-T Application Port Profile", err.Error())
			return
		}
		stateRefreshed.AppPortProfileID.Set(appPortProfile.NsxtAppPortProfile.ID)
	}

	stateRefreshed.Description.Set(rule.NsxtNatRule.Description)
	stateRefreshed.DnatExternalPort.Set(rule.NsxtNatRule.DnatExternalPort)
	stateRefreshed.EdgeGatewayID.Set(r.edgegw.GetID())
	stateRefreshed.EdgeGatewayName.Set(r.edgegw.GetName())
	stateRefreshed.Enabled.Set(rule.NsxtNatRule.Enabled)
	stateRefreshed.ExternalAddress.Set(rule.NsxtNatRule.ExternalAddresses)
	stateRefreshed.FirewallMatch.Set(rule.NsxtNatRule.FirewallMatch)
	stateRefreshed.ID.Set(rule.NsxtNatRule.ID)
	stateRefreshed.InternalAddress.Set(rule.NsxtNatRule.InternalAddresses)
	stateRefreshed.Name.Set(rule.NsxtNatRule.Name)
	stateRefreshed.Priority.SetIntPtr(rule.NsxtNatRule.Priority)
	stateRefreshed.RuleType.Set(rule.NsxtNatRule.Type)
	stateRefreshed.SnatDestinationAddress.Set(rule.NsxtNatRule.SnatDestinationAddresses)

	return stateRefreshed, true, nil
}
