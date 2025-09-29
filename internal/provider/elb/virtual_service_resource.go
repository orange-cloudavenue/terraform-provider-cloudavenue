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
	"strings"

	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
	v1 "github.com/orange-cloudavenue/cloudavenue-sdk-go/v1"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go/v1/edgeloadbalancer"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/mutex"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &VirtualServiceResource{}
	_ resource.ResourceWithConfigure   = &VirtualServiceResource{}
	_ resource.ResourceWithImportState = &VirtualServiceResource{}
	// _ resource.ResourceWithModifyPlan     = &VirtualServiceResource{}
	// _ resource.ResourceWithUpgradeState   = &VirtualServiceResource{}
	// _ resource.ResourceWithValidateConfig = &VirtualServiceResource{}.
)

// NewVirtualServiceResource is a helper function to simplify the provider implementation.
func NewVirtualServiceResource() resource.Resource {
	return &VirtualServiceResource{}
}

// VirtualServiceResource is the resource implementation.
type VirtualServiceResource struct {
	client *client.CloudAvenue
	elb    edgeloadbalancer.Client
	edge   *v1.EdgeClient
}

// Init Initializes the resource.
func (r *VirtualServiceResource) Init(_ context.Context, rm *VirtualServiceModel) (diags diag.Diagnostics) {
	var err error

	r.elb, err = edgeloadbalancer.NewClient()
	if err != nil {
		diags.AddError("Error creating elb client", err.Error())
	}

	eIDOrName := rm.EdgeGatewayID.Get()
	if eIDOrName == "" {
		eIDOrName = rm.EdgeGatewayName.Get()
	}
	r.edge, err = r.client.CAVSDK.V1.EdgeGateway.Get(eIDOrName)
	if err != nil {
		diags.AddError("Error creating edge client", err.Error())
	}

	rm.EdgeGatewayID.Set(urn.Normalize(urn.Gateway, r.edge.GetID()).String())
	rm.EdgeGatewayName.Set(r.edge.GetName())

	return diags
}

// Metadata returns the resource type name.
func (r *VirtualServiceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_virtual_service"
}

// Schema defines the schema for the resource.
func (r *VirtualServiceResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = virtualServiceSchema(ctx).GetResource(ctx)
}

func (r *VirtualServiceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *VirtualServiceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer metrics.New("cloudavenue_elb_virtual_service", r.client.GetOrgName(), metrics.Create)()

	plan := &VirtualServiceModel{}

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

	mutex.GlobalMutex.KvLock(ctx, plan.EdgeGatewayID.Get())
	defer mutex.GlobalMutex.KvUnlock(ctx, plan.EdgeGatewayID.Get())

	modelRequest, d := plan.ToSDKVirtualServiceModelRequest(ctx, r.elb)
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	vsCreated, err := r.elb.CreateVirtualService(ctx, *modelRequest)
	if err != nil {
		resp.Diagnostics.AddError("Error creating virtual service", err.Error())
		return
	}

	plan.ID.Set(vsCreated.ID)

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
func (r *VirtualServiceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer metrics.New("cloudavenue_elb_virtual_service", r.client.GetOrgName(), metrics.Read)()

	state := &VirtualServiceModel{}

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
func (r *VirtualServiceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer metrics.New("cloudavenue_elb_virtual_service", r.client.GetOrgName(), metrics.Update)()

	var (
		plan  = &VirtualServiceModel{}
		state = &VirtualServiceModel{}
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

	mutex.GlobalMutex.KvLock(ctx, plan.EdgeGatewayID.Get())
	defer mutex.GlobalMutex.KvUnlock(ctx, plan.EdgeGatewayID.Get())

	modelRequest, d := plan.ToSDKVirtualServiceModelRequest(ctx, r.elb)
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	_, err := r.elb.UpdateVirtualService(ctx, state.ID.Get(), *modelRequest)
	if err != nil {
		resp.Diagnostics.AddError("Error updating virtual service", err.Error())
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
func (r *VirtualServiceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer metrics.New("cloudavenue_elb_virtual_service", r.client.GetOrgName(), metrics.Delete)()

	state := &VirtualServiceModel{}

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

	mutex.GlobalMutex.KvLock(ctx, state.EdgeGatewayID.Get())
	defer mutex.GlobalMutex.KvUnlock(ctx, state.EdgeGatewayID.Get())

	if err := r.elb.DeleteVirtualService(ctx, state.ID.Get()); err != nil {
		resp.Diagnostics.AddError("Error deleting virtual service", err.Error())
	}
}

func (r *VirtualServiceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	defer metrics.New("cloudavenue_elb_virtual_service", r.client.GetOrgName(), metrics.Import)()

	// Import format is edgeGatewayIDOrName.virtualServiceIDOrName

	// * Import with custom logic
	idParts := strings.Split(req.ID, ".")

	if len(idParts) != 2 {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: edgeGatewayIDOrName.virtualServiceIDOrName. Got: %q", req.ID),
		)
		return
	}

	x := &VirtualServiceModel{
		ID:              supertypes.NewStringNull(),
		Name:            supertypes.NewStringNull(),
		EdgeGatewayID:   supertypes.NewStringNull(),
		EdgeGatewayName: supertypes.NewStringNull(),
	}

	if urn.IsEdgeGateway(idParts[0]) {
		x.EdgeGatewayID.Set(idParts[0])
	} else {
		edge, err := r.client.CAVSDK.V1.EdgeGateway.Get(idParts[0])
		if err != nil {
			resp.Diagnostics.AddError("Error retrieving Edge Gateway", err.Error())
			return
		}
		x.EdgeGatewayName.Set(idParts[0])
		x.EdgeGatewayID.Set(urn.Normalize(urn.Gateway, edge.GetID()).String())
	}

	if urn.IsLoadBalancerVirtualService(idParts[1]) {
		x.ID.Set(idParts[1])
	} else {
		x.Name.Set(idParts[1])
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
func (r *VirtualServiceResource) read(ctx context.Context, planOrState *VirtualServiceModel) (stateRefreshed *VirtualServiceModel, found bool, diags diag.Diagnostics) {
	stateRefreshed = planOrState.Copy()

	/*
		Implement the resource read here
	*/

	nameOrID := planOrState.Name.Get()
	if planOrState.ID.IsKnown() {
		nameOrID = planOrState.ID.Get()
	}

	data, err := r.elb.GetVirtualService(ctx, stateRefreshed.EdgeGatewayID.Get(), nameOrID)
	if err != nil {
		if govcd.ContainsNotFound(err) {
			return nil, false, nil
		}
		diags.AddError("Error retrieving virtual service", err.Error())
		return nil, true, diags
	}

	stateRefreshed.ID.Set(data.ID)
	stateRefreshed.Name.Set(data.Name)
	stateRefreshed.Description.Set(data.Description)
	stateRefreshed.Enabled.SetPtr(data.Enabled)
	stateRefreshed.EdgeGatewayID.Set(data.EdgeGatewayRef.ID)
	stateRefreshed.EdgeGatewayName.Set(data.EdgeGatewayRef.Name)
	stateRefreshed.PoolID.Set(data.PoolRef.ID)
	stateRefreshed.PoolName.Set(data.PoolRef.Name)
	stateRefreshed.ServiceEngineGroupName.Set(data.ServiceEngineGroupRef.Name)
	stateRefreshed.VirtualIP.Set(data.VirtualIPAddress)
	stateRefreshed.ServiceType.Set(string(data.ApplicationProfile))

	if data.CertificateRef != nil {
		stateRefreshed.CertificateID.SetPtr(&data.CertificateRef.ID)
	} else {
		stateRefreshed.CertificateID.SetNull()
	}

	servicePorts := make([]*VirtualServiceModelServicePort, 0)
	for _, port := range data.ServicePorts {
		sp := &VirtualServiceModelServicePort{
			Start: supertypes.NewInt64Null(),
			End:   supertypes.NewInt64Null(),
		}

		sp.Start.SetIntPtr(port.Start)
		sp.End.SetIntPtr(port.End)
		servicePorts = append(servicePorts, sp)
	}

	diags.Append(stateRefreshed.ServicePorts.Set(ctx, servicePorts)...)

	return stateRefreshed, true, diags
}
