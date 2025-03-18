/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package elb

import (
	"context"
	"fmt"

	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/v1/edgeloadbalancer"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/mutex"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &PoliciesHTTPRequestResource{}
	_ resource.ResourceWithConfigure   = &PoliciesHTTPRequestResource{}
	_ resource.ResourceWithImportState = &PoliciesHTTPRequestResource{}
)

// NewPoliciesHTTPRequestResource is a helper function to simplify the provider implementation.
func NewPoliciesHTTPRequestResource() resource.Resource {
	return &PoliciesHTTPRequestResource{}
}

// PoliciesHTTPRequestResource is the resource implementation.
type PoliciesHTTPRequestResource struct {
	client *client.CloudAvenue
	elb    edgeloadbalancer.Client
}

// Init Initializes the resource.
func (r *PoliciesHTTPRequestResource) Init(ctx context.Context, rm *PoliciesHTTPRequestModel) (diags diag.Diagnostics) {
	var err error

	r.elb, err = edgeloadbalancer.NewClient()
	if err != nil {
		diags.AddError("Error creating elb client", err.Error())
	}

	return
}

// Metadata returns the resource type name.
func (r *PoliciesHTTPRequestResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_policies_http_request"
}

// Schema defines the schema for the resource.
func (r *PoliciesHTTPRequestResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = policiesHTTPRequestSchema(ctx).GetResource(ctx)
}

func (r *PoliciesHTTPRequestResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *PoliciesHTTPRequestResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer metrics.New("cloudavenue_elb_policies_http_request", r.client.GetOrgName(), metrics.Create)()

	plan := &PoliciesHTTPRequestModel{}

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

	private := &policiesHTTPPrivateModel{}
	resp.Diagnostics.Append(private.Get(ctx, plan.VirtualServiceID.Get(), resp.Private, r.getEdgeGateway)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Lock the EdgeGateway
	mutex.GlobalMutex.KvLock(ctx, private.EdgeGatewayID)
	defer mutex.GlobalMutex.KvUnlock(ctx, private.EdgeGatewayID)

	// Create the resource
	resp.Diagnostics.Append(r.createOrUpdate(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use generic read function to refresh the state
	state, found, d := r.read(ctx, plan)
	if !found {
		resp.Diagnostics.AddError("Resource not found", "The resource was not found after creation")
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
func (r *PoliciesHTTPRequestResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer metrics.New("cloudavenue_elb_policies_http_request", r.client.GetOrgName(), metrics.Read)()

	state := &PoliciesHTTPRequestModel{}

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
		resp.Diagnostics.AddError("Resource not found", "The resource was not found after refresh")
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
func (r *PoliciesHTTPRequestResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer metrics.New("cloudavenue_elb_policies_http_request", r.client.GetOrgName(), metrics.Update)()

	var (
		plan  = &PoliciesHTTPRequestModel{}
		state = &PoliciesHTTPRequestModel{}
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

	private := &policiesHTTPPrivateModel{}
	resp.Diagnostics.Append(private.Get(ctx, state.VirtualServiceID.Get(), req.Private, r.getEdgeGateway)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Lock the EdgeGateway
	mutex.GlobalMutex.KvLock(ctx, private.EdgeGatewayID)
	defer mutex.GlobalMutex.KvUnlock(ctx, private.EdgeGatewayID)

	// Update the resource
	resp.Diagnostics.Append(r.createOrUpdate(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use generic read function to refresh the state
	stateRefreshed, found, d := r.read(ctx, plan)
	if !found {
		resp.Diagnostics.AddError("Resource not found", "The resource was not found after update")
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
func (r *PoliciesHTTPRequestResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer metrics.New("cloudavenue_elb_policies_http_request", r.client.GetOrgName(), metrics.Delete)()

	state := &PoliciesHTTPRequestModel{}

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

	private := &policiesHTTPPrivateModel{}
	resp.Diagnostics.Append(private.Get(ctx, state.VirtualServiceID.Get(), req.Private, r.getEdgeGateway)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Lock the EdgeGateway
	mutex.GlobalMutex.KvLock(ctx, private.EdgeGatewayID)
	defer mutex.GlobalMutex.KvUnlock(ctx, private.EdgeGatewayID)

	// Delete the resource
	if err := r.elb.DeletePoliciesHTTPRequest(ctx, state.VirtualServiceID.Get()); err != nil {
		resp.Diagnostics.AddError("Error deleting policies http request", err.Error())
		return
	}
}

func (r *PoliciesHTTPRequestResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	defer metrics.New("cloudavenue_elb_policies_http_request", r.client.GetOrgName(), metrics.Import)()

	x := &PoliciesHTTPRequestModel{
		ID:               supertypes.NewStringNull(),
		VirtualServiceID: supertypes.NewStringNull(),
	}
	x.ID.Set(req.ID)
	x.VirtualServiceID.Set(req.ID)

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
func (r *PoliciesHTTPRequestResource) read(ctx context.Context, planOrState *PoliciesHTTPRequestModel) (stateRefreshed *PoliciesHTTPRequestModel, found bool, diags diag.Diagnostics) {
	stateRefreshed = planOrState.Copy()

	/*
		Implement the resource read here
	*/

	data, err := r.elb.GetPoliciesHTTPRequest(ctx, stateRefreshed.VirtualServiceID.Get())
	if err != nil {
		if govcd.ContainsNotFound(err) {
			return nil, false, nil
		}
		diags.AddError("Error retrieving policies http request", err.Error())
		return nil, true, diags
	}

	stateRefreshed = &PoliciesHTTPRequestModel{
		ID:               supertypes.NewStringValueOrNull(data.VirtualServiceID),
		VirtualServiceID: supertypes.NewStringValueOrNull(data.VirtualServiceID),
		Policies: func() supertypes.ListNestedObjectValueOf[PoliciesHTTPRequestModelPolicies] {
			policies := []*PoliciesHTTPRequestModelPolicies{}
			for _, v := range data.Policies {
				policy := &PoliciesHTTPRequestModelPolicies{
					Name:    supertypes.NewStringValueOrNull(v.Name),
					Active:  supertypes.NewBoolValue(v.Active),
					Logging: supertypes.NewBoolValue(v.Logging),
					Criteria: func() supertypes.SingleNestedObjectValueOf[PoliciesHTTPRequestMatchCriteria] {
						return supertypes.NewSingleNestedObjectValueOf(ctx, &PoliciesHTTPRequestMatchCriteria{
							Protocol:       supertypes.NewStringValueOrNull(v.MatchCriteria.Protocol),
							Query:          supertypes.NewSetValueOfSlice(ctx, v.MatchCriteria.QueryMatch),
							ClientIP:       policiesHTTPClientIPMatchFromSDK(ctx, v.MatchCriteria.ClientIPMatch),
							ServicePorts:   policiesHTTPServicePortMatchFromSDK(ctx, v.MatchCriteria.ServicePortMatch),
							HTTPMethods:    policiesHTTPMethodMatchFromSDK(ctx, v.MatchCriteria.MethodMatch),
							Path:           policiesHTTPPathMatchFromSDK(ctx, v.MatchCriteria.PathMatch),
							Cookie:         policiesHTTPCookieMatchFromSDK(ctx, v.MatchCriteria.CookieMatch),
							RequestHeaders: policiesHTTPHeadersMatchFromSDK(ctx, v.MatchCriteria.HeaderMatch),
						})
					}(),
					Actions: func() supertypes.SingleNestedObjectValueOf[PoliciesHTTPRequestActions] {
						return supertypes.NewSingleNestedObjectValueOf(ctx, &PoliciesHTTPRequestActions{
							Redirect:      policiesHTTPActionRedirectFromSDK(ctx, v.RedirectAction),
							RewriteURL:    policiesHTTPActionURLRewriteFromSDK(ctx, v.URLRewriteAction),
							ModifyHeaders: policiesHTTPActionHeadersRewriteFromSDK(ctx, v.HeaderRewriteActions),
						})
					}(),
				}
				policies = append(policies, policy)
			}
			if len(policies) == 0 {
				return supertypes.NewListNestedObjectValueOfNull[PoliciesHTTPRequestModelPolicies](ctx)
			}
			return supertypes.NewListNestedObjectValueOfSlice(ctx, policies)
		}(),
	}

	return stateRefreshed, true, nil
}

func (r *PoliciesHTTPRequestResource) createOrUpdate(ctx context.Context, goPlan *PoliciesHTTPRequestModel) (diags diag.Diagnostics) {
	model, d := goPlan.ToSDKPoliciesHTTPRequestModel(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return
	}

	_, err := r.elb.UpdatePoliciesHTTPRequest(ctx, model)
	if err != nil {
		diags.AddError("Error updating policies http request", err.Error())
	}
	return
}

func (r *PoliciesHTTPRequestResource) getEdgeGateway(ctx context.Context, virtualServiceID string) (string, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	vs, err := r.elb.GetVirtualService(ctx, "", virtualServiceID)
	if err != nil {
		if govcd.ContainsNotFound(err) {
			diags.AddError("virtual service not found", err.Error())
			return "", nil
		}
		diags.AddError("Error retrieving virtual service", err.Error())
		return "", diags
	}

	return vs.EdgeGatewayRef.ID, diags
}
