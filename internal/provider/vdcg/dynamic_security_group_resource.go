/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package vdcg

import (
	"context"
	"fmt"
	"strings"

	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
	v1 "github.com/orange-cloudavenue/cloudavenue-sdk-go/v1"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/mutex"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &DynamicSecurityGroupResource{}
	_ resource.ResourceWithConfigure   = &DynamicSecurityGroupResource{}
	_ resource.ResourceWithImportState = &DynamicSecurityGroupResource{}
)

// NewDynamicSecurityGroupResource is a helper function to simplify the provider implementation.
func NewDynamicSecurityGroupResource() resource.Resource {
	return &DynamicSecurityGroupResource{}
}

// DynamicSecurityGroupResource is the resource implementation.
type DynamicSecurityGroupResource struct {
	client   *client.CloudAvenue
	vdcGroup *v1.VDCGroup
}

// Init Initializes the resource.
func (r *DynamicSecurityGroupResource) Init(_ context.Context, rm *DynamicSecurityGroupModel) (diags diag.Diagnostics) {
	var err error

	idOrName := rm.VDCGroupName.Get()
	if rm.VDCGroupID.IsKnown() && urn.IsVDCGroup(rm.VDCGroupID.Get()) {
		// Use the ID
		idOrName = rm.VDCGroupID.Get()
	}

	r.vdcGroup, err = r.client.CAVSDK.V1.VDC().GetVDCGroup(idOrName)
	if err != nil {
		diags.AddError("Error retrieving VDC Group", err.Error())
		return diags
	}
	return diags
}

// Metadata returns the resource type name.
func (r *DynamicSecurityGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_dynamic_security_group"
}

// Schema defines the schema for the resource.
func (r *DynamicSecurityGroupResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = dynamicSecurityGroupSchema(ctx).GetResource(ctx)
}

func (r *DynamicSecurityGroupResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *DynamicSecurityGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer metrics.New("cloudavenue_vdcg_dynamic_security_group", r.client.GetOrgName(), metrics.Create)()

	plan := &DynamicSecurityGroupModel{}

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

	values, d := plan.ToSDKDynamicSecurityGroupModel(ctx)
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	fwsg, err := r.vdcGroup.CreateFirewallDynamicSecurityGroup(values)
	if err != nil {
		resp.Diagnostics.AddError("Error creating dynamic security group", err.Error())
		return
	}

	// Set the ID
	plan.ID.Set(fwsg.ID)

	// Use generic read function to refresh the state
	state, found, d := r.read(ctx, plan)
	if !found {
		resp.Diagnostics.AddError("Resource not found", fmt.Sprintf("The dynamic security group '%s' was not found after creation.", plan.Name.Get()))
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
func (r *DynamicSecurityGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer metrics.New("cloudavenue_vdcg_dynamic_security_group", r.client.GetOrgName(), metrics.Read)()

	state := &DynamicSecurityGroupModel{}

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
func (r *DynamicSecurityGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer metrics.New("cloudavenue_vdcg_dynamic_security_group", r.client.GetOrgName(), metrics.Update)()

	var (
		plan  = &DynamicSecurityGroupModel{}
		state = &DynamicSecurityGroupModel{}
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

	fwsg, err := r.vdcGroup.GetFirewallDynamicSecurityGroup(state.ID.Get())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving dynamic security group", err.Error())
		return
	}

	values, d := plan.ToSDKDynamicSecurityGroupModel(ctx)
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	if err := fwsg.Update(values); err != nil {
		resp.Diagnostics.AddError("Error updating dynamic security group", err.Error())
		return
	}

	// Use generic read function to refresh the state
	stateRefreshed, found, d := r.read(ctx, plan)
	if !found {
		resp.Diagnostics.AddError("Resource not found", fmt.Sprintf("The dynamic security group '%s' was not found after update.", plan.Name.Get()))
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
func (r *DynamicSecurityGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer metrics.New("cloudavenue_vdcg_dynamic_security_group", r.client.GetOrgName(), metrics.Delete)()

	state := &DynamicSecurityGroupModel{}

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

	fwsg, err := r.vdcGroup.GetFirewallDynamicSecurityGroup(state.ID.Get())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving security group", err.Error())
		return
	}

	if err := fwsg.Delete(); err != nil {
		resp.Diagnostics.AddError("Error deleting security group", err.Error())
		return
	}
}

func (r *DynamicSecurityGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	defer metrics.New("cloudavenue_vdcg_dynamic_security_group", r.client.GetOrgName(), metrics.Import)()

	// * Import with custom logic
	idParts := strings.Split(req.ID, ".")

	if len(idParts) != 2 {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: VDCGroupNameOrID.DynamicSecurityGroupNameOrID Got: %q", req.ID),
		)
		return
	}
	vdcGroupNameOrID, dynamicSecurityGroupNameOrID := idParts[0], idParts[1]

	x := &DynamicSecurityGroupModel{
		ID:           supertypes.NewStringNull(),
		Name:         supertypes.NewStringNull(),
		VDCGroupName: supertypes.NewStringNull(),
		VDCGroupID:   supertypes.NewStringNull(),
		Description:  supertypes.NewStringNull(),
		Criteria:     supertypes.NewListNestedObjectValueOfNull[DynamicSecurityGroupModelCriteria](ctx),
	}

	if urn.IsVDCGroup(vdcGroupNameOrID) {
		x.VDCGroupID.Set(vdcGroupNameOrID)
	} else {
		x.VDCGroupName.Set(vdcGroupNameOrID)
	}

	if urn.IsSecurityGroup(dynamicSecurityGroupNameOrID) {
		x.ID.Set(dynamicSecurityGroupNameOrID)
	} else {
		x.Name.Set(dynamicSecurityGroupNameOrID)
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
func (r *DynamicSecurityGroupResource) read(ctx context.Context, planOrState *DynamicSecurityGroupModel) (stateRefreshed *DynamicSecurityGroupModel, found bool, diags diag.Diagnostics) {
	stateRefreshed = planOrState.Copy()

	/*
		Implement the resource read here
	*/

	idOrName := planOrState.Name.Get()
	if planOrState.ID.IsKnown() {
		idOrName = planOrState.ID.Get()
	}

	fwsg, err := r.vdcGroup.GetFirewallDynamicSecurityGroup(idOrName)
	if govcd.ContainsNotFound(err) {
		return nil, false, nil
	}
	if err != nil {
		diags.AddError("Error retrieving security group", err.Error())
		return nil, true, diags
	}

	stateRefreshed.ID.Set(fwsg.ID)
	stateRefreshed.Name.Set(fwsg.Name)
	stateRefreshed.Description.Set(fwsg.Description)
	stateRefreshed.VDCGroupName.Set(r.vdcGroup.GetName())
	stateRefreshed.VDCGroupID.Set(r.vdcGroup.GetID())

	if fwsg.Criteria != nil {
		criteria := make([]*DynamicSecurityGroupModelCriteria, 0)
		for _, c := range fwsg.Criteria {
			rules := make([]*DynamicSecurityGroupModelRule, 0)

			for _, r := range c.Rules {
				rule := &DynamicSecurityGroupModelRule{
					Type:     supertypes.NewStringNull(),
					Value:    supertypes.NewStringNull(),
					Operator: supertypes.NewStringNull(),
				}

				rule.Type.Set(string(r.RuleType))
				rule.Value.Set(r.Value)
				rule.Operator.Set(string(r.Operator))

				rules = append(rules, rule)
			}

			c := &DynamicSecurityGroupModelCriteria{
				Rules: supertypes.NewListNestedObjectValueOfNull[DynamicSecurityGroupModelRule](ctx),
			}
			diags.Append(c.Rules.Set(ctx, rules)...)

			criteria = append(criteria, c)
		}
		if len(criteria) == 0 {
			stateRefreshed.Criteria.SetNull(ctx)
		} else {
			diags.Append(stateRefreshed.Criteria.Set(ctx, criteria)...)
		}
	}

	return stateRefreshed, true, diags
}
