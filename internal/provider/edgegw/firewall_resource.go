// Package edgegw provides a Terraform resource.
package edgegw

import (
	"context"
	"fmt"

	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/edgegw"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/mutex"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/org"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/uuid"
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

type firewallModel struct {
	ID              types.String `tfsdk:"id"`
	EdgeGatewayID   types.String `tfsdk:"edge_gateway_id"`
	EdgeGatewayName types.String `tfsdk:"edge_gateway_name"`
	Rules           types.List   `tfsdk:"rules"`
}

type firewallModelRules []firewallModelRule

type firewallModelRule struct {
	ID                types.String `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	Enabled           types.Bool   `tfsdk:"enabled"`
	Direction         types.String `tfsdk:"direction"`
	IPProtocol        types.String `tfsdk:"ip_protocol"`
	Action            types.String `tfsdk:"action"`
	Logging           types.Bool   `tfsdk:"logging"`
	SourceIDs         types.Set    `tfsdk:"source_ids"`
	DestinationIDs    types.Set    `tfsdk:"destination_ids"`
	AppPortProfileIDs types.Set    `tfsdk:"app_port_profile_ids"`
}

// RulesFromPlan.
func (rm *firewallModel) RulesFromPlan(ctx context.Context) (rules firewallModelRules, diags diag.Diagnostics) {
	tflog.Info(ctx, "firewallModel.RulesFromPlan")
	rules = make(firewallModelRules, 0)
	diags.Append(rm.Rules.ElementsAs(ctx, &rules, false)...)
	return
}

// rulesToNsxtFirewallRule.
func (rules *firewallModelRules) rulesToNsxtFirewallRule(ctx context.Context) (nsxtFirewallRules []*govcdtypes.NsxtFirewallRule, diags diag.Diagnostics) {
	nsxtFirewallRules = make([]*govcdtypes.NsxtFirewallRule, len(*rules))
	for i, rule := range *rules {
		nsxtFirewallRules[i] = &govcdtypes.NsxtFirewallRule{
			Name:       rule.Name.ValueString(),
			Action:     rule.Action.ValueString(),
			Enabled:    rule.Enabled.ValueBool(),
			IpProtocol: rule.IPProtocol.ValueString(),
			Logging:    rule.Logging.ValueBool(),
			Direction:  rule.Direction.ValueString(),
			Version:    nil,
		}

		// ! If sourceIDs/destinationIDs is Null, it's an equivalent of any (source/destination)

		// * sourceIDs
		if !rule.SourceIDs.IsNull() {
			nsxtFirewallRules[i].SourceFirewallGroups = make([]govcdtypes.OpenApiReference, 0)
			sourcesIDs := make([]string, 0)
			if d := rule.SourceIDs.ElementsAs(ctx, &sourcesIDs, false); d.HasError() {
				diags.Append(d...)
				return
			}

			for index, sourceID := range sourcesIDs {
				nsxtFirewallRules[i].SourceFirewallGroups[index] = govcdtypes.OpenApiReference{
					ID: sourceID,
				}
			}
		}

		// * destinationIDs
		if !rule.DestinationIDs.IsNull() {
			nsxtFirewallRules[i].DestinationFirewallGroups = make([]govcdtypes.OpenApiReference, 0)
			destinationIDs := make([]string, 0)
			if d := rule.DestinationIDs.ElementsAs(ctx, &destinationIDs, false); d.HasError() {
				diags.Append(d...)
				return
			}

			for index, destinationID := range destinationIDs {
				nsxtFirewallRules[i].DestinationFirewallGroups[index] = govcdtypes.OpenApiReference{
					ID: destinationID,
				}
			}
		}

		// * appPortProfilesIDs
		if !rule.AppPortProfileIDs.IsNull() {
			nsxtFirewallRules[i].ApplicationPortProfiles = make([]govcdtypes.OpenApiReference, 0)
			appPortProfilesIDs := make([]string, 0)
			if d := rule.AppPortProfileIDs.ElementsAs(ctx, &appPortProfilesIDs, false); d.HasError() {
				diags.Append(d...)
				return
			}

			for index, appPortProfileID := range appPortProfilesIDs {
				nsxtFirewallRules[i].ApplicationPortProfiles[index] = govcdtypes.OpenApiReference{
					ID: appPortProfileID,
				}
			}
		}
	}

	return
}

// attrTypes returns the attribute types for the resource.
func (rules *firewallModelRules) AttrTypes(_ context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"id":                   types.StringType,
		"name":                 types.StringType,
		"enabled":              types.BoolType,
		"direction":            types.StringType,
		"ip_protocol":          types.StringType,
		"action":               types.StringType,
		"logging":              types.BoolType,
		"source_ids":           types.SetType{ElemType: types.StringType},
		"destination_ids":      types.SetType{ElemType: types.StringType},
		"app_port_profile_ids": types.SetType{ElemType: types.StringType},
	}
}

// objectTypes returns the object types for the resource.
func (rules *firewallModelRules) ObjectType(ctx context.Context) types.ObjectType {
	return types.ObjectType{
		AttrTypes: rules.AttrTypes(ctx),
	}
}

// ToPlan.
func (rules *firewallModelRules) ToPlan(ctx context.Context) (lV basetypes.ListValue, diags diag.Diagnostics) {
	if rules == nil {
		diags.AddError("Rules not initialized", "Failed to convert rules to plan")
		return types.ListNull(rules.ObjectType(ctx)), diags
	}

	return types.ListValueFrom(ctx, rules.ObjectType(ctx), rules)
}

// fwRulesRead.
func fwRulesRead(ctx context.Context, fwRules *govcd.NsxtFirewall) (rules firewallModelRules, diags diag.Diagnostics) {
	rules = make(firewallModelRules, 0)

	for _, rule := range fwRules.NsxtFirewallRuleContainer.UserDefinedRules {
		sourceIDs, d := types.SetValueFrom(ctx, types.StringType, rule.SourceFirewallGroups)
		diags.Append(d...)
		destinationIDs, d := types.SetValueFrom(ctx, types.StringType, rule.DestinationFirewallGroups)
		diags.Append(d...)
		appPortProfileIDs, d := types.SetValueFrom(ctx, types.StringType, rule.ApplicationPortProfiles)
		diags.Append(d...)
		if diags.HasError() {
			return
		}
		rules = append(rules, firewallModelRule{
			ID:                types.StringValue(rule.ID),
			Name:              types.StringValue(rule.Name),
			Enabled:           types.BoolValue(rule.Enabled),
			Direction:         types.StringValue(rule.Direction),
			IPProtocol:        types.StringValue(rule.IpProtocol),
			Action:            types.StringValue(rule.Action),
			Logging:           types.BoolValue(rule.Logging),
			SourceIDs:         sourceIDs,
			DestinationIDs:    destinationIDs,
			AppPortProfileIDs: appPortProfileIDs,
		})
	}

	return
}

// Init Initializes the resource.
func (r *firewallResource) Init(ctx context.Context, rm *firewallModel) (diags diag.Diagnostics) {
	var err error

	r.org, diags = org.Init(r.client)
	if diags.HasError() {
		return
	}

	r.edgegw, err = r.org.GetEdgeGateway(edgegw.BaseEdgeGW{
		ID:   rm.EdgeGatewayID,
		Name: rm.EdgeGatewayName,
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

	// Rules
	rules, d := plan.RulesFromPlan(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	/*
		Implement the resource creation logic here.
	*/

	// Lock object VDC or VDC Group
	vdcOrVDCGroup, err := r.edgegw.GetParent()
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving Edge Gateway parent", err.Error())
		return
	}

	mutex.GlobalMutex.KvLock(ctx, vdcOrVDCGroup.GetID())
	defer mutex.GlobalMutex.KvUnlock(ctx, vdcOrVDCGroup.GetID())

	// Set the rules
	vcdRules, d := rules.rulesToNsxtFirewallRule(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	if _, err := r.edgegw.UpdateNsxtFirewall(&govcdtypes.NsxtFirewallRuleContainer{
		UserDefinedRules: vcdRules,
	}); err != nil {
		resp.Diagnostics.AddError("Error to create Firewall", err.Error())
		return
	}

	state, d := r.read(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *firewallResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
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
		Implement the resource read here
	*/

	plan, d := r.read(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *firewallResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) { //nolint:dupl
	plan := &firewallModel{}

	// Get current state
	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Init the resource
	resp.Diagnostics.Append(r.Init(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Rules
	rules, d := plan.RulesFromPlan(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	/*
		Implement the resource update here
	*/

	vdcOrVDCGroup, err := r.edgegw.GetParent()
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving Edge Gateway parent", err.Error())
		return
	}

	mutex.GlobalMutex.KvLock(ctx, vdcOrVDCGroup.GetID())
	defer mutex.GlobalMutex.KvUnlock(ctx, vdcOrVDCGroup.GetID())

	vcdRules, d := rules.rulesToNsxtFirewallRule(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	if _, err := r.edgegw.UpdateNsxtFirewall(&govcdtypes.NsxtFirewallRuleContainer{
		UserDefinedRules: vcdRules,
	}); err != nil {
		resp.Diagnostics.AddError("Error to create Firewall", err.Error())
		return
	}

	plan, d = r.read(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *firewallResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
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

	mutex.GlobalMutex.KvLock(ctx, vdcOrVDCGroup.GetID())
	defer mutex.GlobalMutex.KvUnlock(ctx, vdcOrVDCGroup.GetID())

	fwRules, err := r.edgegw.GetNsxtFirewall()
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving Edge Gateway Firewall", err.Error())
		return
	}

	if err := fwRules.DeleteAllRules(); err != nil {
		resp.Diagnostics.AddError("Error deleting Edge Gateway Firewall", err.Error())
		return
	}
}

func (r *firewallResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var (
		edgegwID   string
		edgegwName string
		d          diag.Diagnostics
	)

	r.org, d = org.Init(r.client)
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	if uuid.IsValid(req.ID) {
		edgegwID = uuid.Normalize(uuid.Gateway, req.ID).String()
	} else {
		edgegwName = req.ID
	}

	edgegw, err := r.org.GetEdgeGateway(edgegw.BaseEdgeGW{
		ID:   types.StringValue(edgegwID),
		Name: types.StringValue(edgegwName),
	})
	if err != nil {
		resp.Diagnostics.AddError("Failed to import firewall.", err.Error())
		return
	}

	// ID
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), edgegw.GetID())...)
	// edge_gateway_id
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("edge_gateway_id"), edgegw.GetID())...)
	// edge_gateway_name
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("edge_gateway_name"), edgegw.GetName())...)
}

func (r *firewallResource) read(ctx context.Context) (plan *firewallModel, diags diag.Diagnostics) {
	fwRules, err := r.edgegw.GetNsxtFirewall()
	if err != nil {
		diags.AddError("Error retrieving Edge Gateway Firewall", err.Error())
		return
	}

	rules, d := fwRulesRead(ctx, fwRules)
	diags.Append(d...)
	if diags.HasError() {
		return
	}

	plan = &firewallModel{
		// ID is stored as Edge Gateway ID - because this is a "container" for all firewall rules at once and each child
		// TypeSet element will have a computed ID field for each rule
		ID:              types.StringValue(r.edgegw.GetID()),
		EdgeGatewayID:   types.StringValue(r.edgegw.GetID()),
		EdgeGatewayName: types.StringValue(r.edgegw.GetName()),
	}
	plan.Rules, d = rules.ToPlan(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return
	}

	return
}
