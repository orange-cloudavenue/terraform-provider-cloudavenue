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
	_ resource.Resource                = &PoliciesHTTPResponseResource{}
	_ resource.ResourceWithConfigure   = &PoliciesHTTPResponseResource{}
	_ resource.ResourceWithImportState = &PoliciesHTTPResponseResource{}
)

// NewPoliciesHTTPResponseResource is a helper function to simplify the provider implementation.
func NewPoliciesHTTPResponseResource() resource.Resource {
	return &PoliciesHTTPResponseResource{}
}

// PoliciesHTTPResponseResource is the resource implementation.
type PoliciesHTTPResponseResource struct {
	client *client.CloudAvenue
	elb    edgeloadbalancer.Client
}

// Init Initializes the resource.
func (r *PoliciesHTTPResponseResource) Init(_ context.Context, _ *PoliciesHTTPResponseModel) (diags diag.Diagnostics) {
	var err error

	r.elb, err = edgeloadbalancer.NewClient()
	if err != nil {
		diags.AddError("Error creating elb client", err.Error())
	}

	return
}

// Metadata returns the resource type name.
func (r *PoliciesHTTPResponseResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_policies_http_response"
}

// Schema defines the schema for the resource.
func (r *PoliciesHTTPResponseResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = policiesHTTPResponseSchema(ctx).GetResource(ctx)
}

func (r *PoliciesHTTPResponseResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *PoliciesHTTPResponseResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer metrics.New("cloudavenue_elb_policies_http_response", r.client.GetOrgName(), metrics.Create)()

	plan := &PoliciesHTTPResponseModel{}

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
func (r *PoliciesHTTPResponseResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer metrics.New("cloudavenue_elb_policies_http_response", r.client.GetOrgName(), metrics.Read)()

	state := &PoliciesHTTPResponseModel{}

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
func (r *PoliciesHTTPResponseResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer metrics.New("cloudavenue_elb_policies_http_response", r.client.GetOrgName(), metrics.Update)()

	var (
		plan  = &PoliciesHTTPResponseModel{}
		state = &PoliciesHTTPResponseModel{}
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
func (r *PoliciesHTTPResponseResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer metrics.New("cloudavenue_elb_policies_http_response", r.client.GetOrgName(), metrics.Delete)()

	state := &PoliciesHTTPResponseModel{}

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
	if err := r.elb.DeletePoliciesHTTPResponse(ctx, state.VirtualServiceID.Get()); err != nil {
		resp.Diagnostics.AddError("Error deleting policies http request", err.Error())
		return
	}
}

func (r *PoliciesHTTPResponseResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	defer metrics.New("cloudavenue_elb_policies_http_response", r.client.GetOrgName(), metrics.Import)()

	x := &PoliciesHTTPResponseModel{
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
func (r *PoliciesHTTPResponseResource) read(ctx context.Context, planOrState *PoliciesHTTPResponseModel) (stateRefreshed *PoliciesHTTPResponseModel, found bool, diags diag.Diagnostics) {
	stateRefreshed = planOrState.Copy()

	/*
		Implement the resource read here
	*/

	data, err := r.elb.GetPoliciesHTTPResponse(ctx, stateRefreshed.VirtualServiceID.Get())
	if err != nil {
		if govcd.ContainsNotFound(err) {
			return nil, false, nil
		}
		diags.AddError("Error retrieving policies http request", err.Error())
		return nil, true, diags
	}

	stateRefreshed = &PoliciesHTTPResponseModel{
		ID:               supertypes.NewStringValueOrNull(data.VirtualServiceID),
		VirtualServiceID: supertypes.NewStringValueOrNull(data.VirtualServiceID),
		Policies: func() supertypes.ListNestedObjectValueOf[PoliciesHTTPResponseModelPolicies] {
			policies := []*PoliciesHTTPResponseModelPolicies{}
			for _, v := range data.Policies {
				policy := &PoliciesHTTPResponseModelPolicies{
					Name:    supertypes.NewStringValueOrNull(v.Name),
					Active:  supertypes.NewBoolValue(v.Active),
					Logging: supertypes.NewBoolValue(v.Logging),
					Criteria: func() supertypes.SingleNestedObjectValueOf[PoliciesHTTPResponseMatchCriteria] {
						return supertypes.NewSingleNestedObjectValueOf(ctx, &PoliciesHTTPResponseMatchCriteria{
							Protocol:        supertypes.NewStringValueOrNull(v.MatchCriteria.Protocol),
							ClientIP:        policiesHTTPClientIPMatchFromSDK(ctx, v.MatchCriteria.ClientIPMatch),
							ServicePorts:    policiesHTTPServicePortMatchFromSDK(ctx, v.MatchCriteria.ServicePortMatch),
							HTTPMethods:     policiesHTTPMethodMatchFromSDK(ctx, v.MatchCriteria.MethodMatch),
							Path:            policiesHTTPPathMatchFromSDK(ctx, v.MatchCriteria.PathMatch),
							Cookie:          policiesHTTPCookieMatchFromSDK(ctx, v.MatchCriteria.CookieMatch),
							Location:        policiesHTTPLocationMatchFromSDK(ctx, v.MatchCriteria.LocationMatch),
							RequestHeaders:  policiesHTTPHeadersMatchFromSDK(ctx, v.MatchCriteria.RequestHeaderMatch),
							ResponseHeaders: policiesHTTPHeadersMatchFromSDK(ctx, v.MatchCriteria.ResponseHeaderMatch),
							StatusCode:      policiesHTTPStatusCodeMatchFromSDK(ctx, v.MatchCriteria.StatusCodeMatch),
							Query:           supertypes.NewSetValueOfSlice(ctx, v.MatchCriteria.QueryMatch),
						})
					}(),
					Actions: func() supertypes.SingleNestedObjectValueOf[PoliciesHTTPResponseActions] {
						return supertypes.NewSingleNestedObjectValueOf(ctx, &PoliciesHTTPResponseActions{
							LocationRewrite: policiesHTTPActionLocationRewriteFromSDK(ctx, v.LocationRewriteAction),
							ModifyHeaders:   policiesHTTPActionHeadersRewriteFromSDK(ctx, v.HeaderRewriteActions),
						})
					}(),
				}
				policies = append(policies, policy)
			}
			if len(policies) == 0 {
				return supertypes.NewListNestedObjectValueOfNull[PoliciesHTTPResponseModelPolicies](ctx)
			}
			return supertypes.NewListNestedObjectValueOfSlice(ctx, policies)
		}(),
	}

	return stateRefreshed, true, nil
}

func (r *PoliciesHTTPResponseResource) createOrUpdate(ctx context.Context, goPlan *PoliciesHTTPResponseModel) (diags diag.Diagnostics) {
	model, d := goPlan.ToSDKPoliciesHTTPResponseModel(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return
	}

	_, err := r.elb.UpdatePoliciesHTTPResponse(ctx, model)
	if err != nil {
		diags.AddError("Error updating policies http request", err.Error())
	}
	return
}

func (r *PoliciesHTTPResponseResource) getEdgeGateway(ctx context.Context, virtualServiceID string) (string, diag.Diagnostics) {
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
