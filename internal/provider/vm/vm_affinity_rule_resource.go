// Package vm provides a Terraform resource.
package vm

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vdc"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/uuid"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &vmAffinityRuleResource{}
	_ resource.ResourceWithConfigure   = &vmAffinityRuleResource{}
	_ resource.ResourceWithImportState = &vmAffinityRuleResource{}
)

// NewVMAffinityRuleResource is a helper function to simplify the provider implementation.
func NewVMAffinityRuleResource() resource.Resource {
	return &vmAffinityRuleResource{}
}

// vmAffinityRuleResource is the resource implementation.
type vmAffinityRuleResource struct {
	client *client.CloudAvenue
	vdc    vdc.VDC
}

// Metadata returns the resource type name.
func (r *vmAffinityRuleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_affinity_rule"
}

// Schema defines the schema for the resource.
func (r *vmAffinityRuleResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = vmAffinityRuleSchema().GetResource(ctx)
}

func (r *vmAffinityRuleResource) Init(ctx context.Context, rm *vmAffinityRuleResourceModel) (diags diag.Diagnostics) {
	r.vdc, diags = vdc.Init(r.client, rm.VDC)

	return
}

func (r *vmAffinityRuleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *vmAffinityRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var (
		plan *vmAffinityRuleResourceModel
	)

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Init resource
	resp.Diagnostics.Append(r.Init(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	vmAffinityRuleDef, err := resourceToAffinityRule(r, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create affinity rule : resourceToAffinityRule() error", err.Error())
		return
	}

	vmAffinityRule, err := r.vdc.CreateVmAffinityRule(vmAffinityRuleDef)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create affinity rule", err.Error())
		return
	}

	plan.ID = types.StringValue(vmAffinityRule.VmAffinityRule.ID)
	plan.VDC = types.StringValue(r.vdc.GetName())

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *vmAffinityRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *vmAffinityRuleResourceModel

	// Get current state
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Init resource
	resp.Diagnostics.Append(r.Init(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan := &vmAffinityRuleResourceModel{}

	vmAffinityRule, err := getVMAffinityRule(r.vdc, state.Name.ValueString(), state.ID.ValueString())
	if err != nil {
		if govcd.IsNotFound(err) {
			tflog.Debug(ctx, fmt.Sprintf("Affinity rule not found with id %s and name %s", state.ID.ValueString(), state.Name.ValueString()))
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Failed to read affinity rule", err.Error())
		return
	}

	plan = &vmAffinityRuleResourceModel{
		ID:       types.StringValue(vmAffinityRule.VmAffinityRule.ID),
		VDC:      types.StringValue(r.vdc.GetName()),
		Name:     types.StringValue(vmAffinityRule.VmAffinityRule.Name),
		Required: types.BoolValue(*vmAffinityRule.VmAffinityRule.IsMandatory),
		Enabled:  types.BoolValue(*vmAffinityRule.VmAffinityRule.IsEnabled),
		Polarity: types.StringValue(vmAffinityRule.VmAffinityRule.Polarity),
	}

	endpointVMs := vmReferencesToListValue(vmAffinityRule.VmAffinityRule.VmReferences)
	plan.VMIDs = types.SetValueMust(types.StringType, endpointVMs)

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *vmAffinityRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan *vmAffinityRuleResourceModel

	// Get current state
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Init resource
	resp.Diagnostics.Append(r.Init(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	vmAffinityRuleDef, err := resourceToAffinityRule(r, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create affinity rule : resourceToAffinityRule() error", err.Error())
		return
	}

	vmAffinityRule, err := getVMAffinityRule(r.vdc, plan.Name.ValueString(), plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to read affinity rule", err.Error())
		return
	}

	vmAffinityRule.VmAffinityRule.Name = vmAffinityRuleDef.Name
	vmAffinityRule.VmAffinityRule.IsMandatory = vmAffinityRuleDef.IsMandatory
	vmAffinityRule.VmAffinityRule.IsEnabled = vmAffinityRuleDef.IsEnabled
	vmAffinityRule.VmAffinityRule.VmReferences = vmAffinityRuleDef.VmReferences

	err = vmAffinityRule.Update()
	if err != nil {
		resp.Diagnostics.AddError("Failed to update affinity rule", err.Error())
		return
	}

	plan = &vmAffinityRuleResourceModel{
		ID:       types.StringValue(vmAffinityRule.VmAffinityRule.ID),
		VDC:      types.StringValue(r.vdc.GetName()),
		Polarity: types.StringValue(vmAffinityRule.VmAffinityRule.Polarity),
		Name:     types.StringValue(vmAffinityRule.VmAffinityRule.Name),
		Required: types.BoolValue(*vmAffinityRule.VmAffinityRule.IsMandatory),
		Enabled:  types.BoolValue(*vmAffinityRule.VmAffinityRule.IsEnabled),
	}

	endpointVMs := vmReferencesToListValue(vmAffinityRule.VmAffinityRule.VmReferences)
	plan.VMIDs = types.SetValueMust(types.StringType, endpointVMs)

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *vmAffinityRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *vmAffinityRuleResourceModel

	// Get current state
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Init resource
	resp.Diagnostics.Append(r.Init(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	vmAffinityRule, err := getVMAffinityRule(r.vdc, state.Name.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to read affinity rule", err.Error())
		return
	}

	err = vmAffinityRule.Delete()
	if err != nil {
		resp.Diagnostics.AddError("Failed to delete affinity rule", err.Error())
	}
}

func (r *vmAffinityRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var state *vmAffinityRuleResourceModel

	resourceURI := strings.Split(req.ID, ".")
	if len(resourceURI) != 2 && len(resourceURI) != 1 {
		resp.Diagnostics.AddError("Failed to import resource", "resource URI must be specified as vdc-name.affinity-rule-name or just affinity-rule-name if VDC is set at provider level.")
	}

	// Init resource
	affinityRuleIdentifier := ""
	vdcName := ""

	if len(resourceURI) == 1 {
		affinityRuleIdentifier = resourceURI[0]
	}

	if len(resourceURI) == 2 {
		vdcName, affinityRuleIdentifier = resourceURI[0], resourceURI[1]
	}
	state = &vmAffinityRuleResourceModel{
		VDC: types.StringValue(vdcName),
	}

	// Init resource
	resp.Diagnostics.Append(r.Init(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	vdcName = r.vdc.GetName()

	if vdcName == "" {
		resp.Diagnostics.AddError("Failed to import resource", "VDC must be set at provider level or in resource URI.")
		return
	}

	if affinityRuleIdentifier == "" {
		resp.Diagnostics.AddError("Failed to import resource", "affinity rule name must be specified in resource URI.")
		return
	}

	lookingForID := govcd.IsUuid(affinityRuleIdentifier)

	ruleList, err := r.vdc.GetAllVmAffinityRuleList()
	if err != nil {
		resp.Diagnostics.AddError("Failed to import resource", err.Error())
	}
	var foundRules []*govcdtypes.VmAffinityRule

	for _, rule := range ruleList {
		if lookingForID && rule.ID == affinityRuleIdentifier {
			foundRules = append(foundRules, rule)
			break
		}
		if rule.Name == affinityRuleIdentifier {
			foundRules = append(foundRules, rule)
		}
	}

	if len(foundRules) == 0 {
		resp.Diagnostics.AddError("Failed to import resource", "no affinity rule found with name or ID "+affinityRuleIdentifier)
		return
	}
	if len(foundRules) > 1 {
		resp.Diagnostics.AddError("Failed to import resource", "more than one affinity rule found with name or ID "+affinityRuleIdentifier)
		return
	}
	vmAffinityRule := foundRules[0]
	state = &vmAffinityRuleResourceModel{
		ID:   types.StringValue(vmAffinityRule.ID),
		Name: types.StringValue(vmAffinityRule.Name),
		VDC:  types.StringValue(vdcName),
	}

	endpointVMs := vmReferencesToListValue(vmAffinityRule.VmReferences)
	state.VMIDs = types.SetValueMust(types.StringType, endpointVMs)

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// resourceToAffinityRule prepares a VM affinity rule definition from the data in the resource.
func resourceToAffinityRule(r *vmAffinityRuleResource, m *vmAffinityRuleResourceModel) (*govcdtypes.VmAffinityRule, error) {
	name := m.Name.ValueString()
	polarity := m.Polarity.ValueString()
	required := m.Required.ValueBool()
	enabled := m.Enabled.ValueBool()
	rawVms := m.VMIDs

	fullVMList, err := r.vdc.QueryVmList(govcdtypes.VmQueryFilterOnlyDeployed)
	if err != nil {
		return nil, err
	}

	var (
		vmReferences   []*govcdtypes.Reference
		invalidEntries = make(map[string]bool)
		foundEntries   = make(map[string]bool)
	)

	for _, vmID := range rawVms.Elements() {
		for _, vm := range fullVMList {
			uuid := common.ExtractUUID(vmID.String())
			if uuid != "" {
				if uuid == common.ExtractUUID(vm.HREF) {
					vmReferences = append(vmReferences, &govcdtypes.Reference{HREF: vm.HREF})
					foundEntries[vmID.String()] = true
				}
			} else {
				invalidEntries[vmID.String()] = true
			}
		}
	}
	if len(invalidEntries) > 0 {
		var invalidItems []string
		for k := range invalidEntries {
			invalidItems = append(invalidItems, k)
		}
		return nil, fmt.Errorf("invalid entries (not a VM ID) detected: %v", invalidItems)
	}
	if len(rawVms.Elements()) > len(foundEntries) {
		var notExistingVms []string
		for _, vmID := range rawVms.Elements() {
			_, exists := foundEntries[vmID.String()]
			if !exists {
				notExistingVms = append(notExistingVms, vmID.String())
			}
		}
		return nil, fmt.Errorf("not existing VMs detected: %v", notExistingVms)
	}

	vmAffinityRuleDef := &govcdtypes.VmAffinityRule{
		Name:        name,
		IsEnabled:   &enabled,
		IsMandatory: &required,
		Polarity:    polarity,
		VmReferences: []*govcdtypes.VMs{
			{
				VMReference: vmReferences,
			},
		},
	}
	return vmAffinityRuleDef, nil
}

func getVMAffinityRule(vdc vdc.VDC, name, id string) (*govcd.VmAffinityRule, error) {
	d := diag.Diagnostics{}
	// The last method of retrieval is by name
	if id == "" {
		if name == "" {
			d.AddError("Failed to read affinity rule", "no identifier or name provided")
			return nil, errors.New("failed to read affinity rule : no identifier or name provided")
		}
		id = name
	}

	vmAffinityRule, err := vdc.GetVmAffinityRuleByNameOrId(id)
	if err != nil {
		return nil, err
	}

	return vmAffinityRule, nil
}

func vmReferencesToListValue(refs []*govcdtypes.VMs) []attr.Value {
	var endpointVMs []attr.Value
	for _, vmr := range refs {
		for _, ref := range vmr.VMReference {
			endpointVMs = append(endpointVMs, types.StringValue(uuid.Normalize(uuid.VM, ref.ID).String()))
		}
	}
	return endpointVMs
}
