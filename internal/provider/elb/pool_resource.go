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
	_ resource.Resource                = &PoolResource{}
	_ resource.ResourceWithConfigure   = &PoolResource{}
	_ resource.ResourceWithImportState = &PoolResource{}
)

// NewPoolResource is a helper function to simplify the provider implementation.
func NewPoolResource() resource.Resource {
	return &PoolResource{}
}

// PoolResource is the resource implementation.
type PoolResource struct {
	client *client.CloudAvenue
	elb    edgeloadbalancer.Client
	edge   *v1.EdgeClient
}

// Init Initializes the resource.
func (r *PoolResource) Init(_ context.Context, rm *PoolModel) (diags diag.Diagnostics) {
	var err error

	r.elb, err = edgeloadbalancer.NewClient()
	if err != nil {
		diags.AddError("Error creating elb client", err.Error())
		return diags
	}

	eIDOrName := rm.EdgeGatewayID.Get()
	if eIDOrName == "" {
		eIDOrName = rm.EdgeGatewayName.Get()
	}
	r.edge, err = r.client.CAVSDK.V1.EdgeGateway.Get(eIDOrName)
	if err != nil {
		diags.AddError("Error creating edge client", err.Error())
		return diags
	}

	rm.EdgeGatewayID.Set(urn.Normalize(urn.Gateway, r.edge.GetID()).String())
	rm.EdgeGatewayName.Set(r.edge.GetName())

	return diags
}

// Metadata returns the resource type name.
func (r *PoolResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_pool"
}

// Schema defines the schema for the resource.
func (r *PoolResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = poolSchema(ctx).GetResource(ctx)
}

func (r *PoolResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *PoolResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer metrics.New("cloudavenue_elb_pool", r.client.GetOrgName(), metrics.Create)()

	plan := &PoolModel{}

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

	model, d := plan.ToSDKPoolModelRequest(ctx, r.client)
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	poolCreated, err := r.elb.CreatePool(ctx, *model)
	if err != nil {
		resp.Diagnostics.AddError("Error creating pool", err.Error())
		return
	}

	plan.ID.Set(poolCreated.ID)

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
func (r *PoolResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer metrics.New("cloudavenue_elb_pool", r.client.GetOrgName(), metrics.Read)()

	state := &PoolModel{}

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
func (r *PoolResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer metrics.New("cloudavenue_elb_pool", r.client.GetOrgName(), metrics.Update)()

	var (
		plan  = &PoolModel{}
		state = &PoolModel{}
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

	model, d := plan.ToSDKPoolModelRequest(ctx, r.client)
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	_, err := r.elb.UpdatePool(ctx, state.ID.Get(), *model)
	if err != nil {
		resp.Diagnostics.AddError("Error updating pool", err.Error())
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
func (r *PoolResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer metrics.New("cloudavenue_elb_pool", r.client.GetOrgName(), metrics.Delete)()

	state := &PoolModel{}

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

	if err := r.elb.DeletePool(ctx, state.ID.Get()); err != nil {
		resp.Diagnostics.AddError("Error deleting pool", err.Error())
		return
	}
}

func (r *PoolResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	defer metrics.New("cloudavenue_elb_pool", r.client.GetOrgName(), metrics.Import)()

	// Import format is edgeGatewayIDOrName.poolIDOrName

	// * Import with custom logic
	idParts := strings.Split(req.ID, ".")

	if len(idParts) != 2 {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: edgeGatewayIDOrName.poolIDOrName. Got: %q", req.ID),
		)
		return
	}

	x := &PoolModel{
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

	if urn.IsLoadBalancerPool(idParts[1]) {
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
func (r *PoolResource) read(ctx context.Context, planOrState *PoolModel) (stateRefreshed *PoolModel, found bool, diags diag.Diagnostics) {
	stateRefreshed = planOrState.Copy()

	/*
		Implement the resource read here
	*/

	idOrName := planOrState.ID.Get()
	if idOrName == "" {
		idOrName = planOrState.Name.Get()
	}

	data, err := r.elb.GetPool(ctx, planOrState.EdgeGatewayID.Get(), idOrName)
	if err != nil {
		if govcd.ContainsNotFound(err) {
			return nil, false, nil
		}
		diags.AddError("Error retrieving pool", err.Error())
		return nil, true, diags
	}

	stateRefreshed.ID.Set(data.ID)
	stateRefreshed.Name.Set(data.Name)
	stateRefreshed.Description.Set(data.Description)
	stateRefreshed.EdgeGatewayID.Set(data.GatewayRef.ID)
	stateRefreshed.EdgeGatewayName.Set(data.GatewayRef.Name)
	stateRefreshed.Enabled.SetPtr(data.Enabled)
	stateRefreshed.Algorithm.Set(string(data.Algorithm))
	stateRefreshed.DefaultPort.SetIntPtr(data.DefaultPort)

	// * Members
	members := &PoolModelMembers{
		GracefulTimeoutPeriod: supertypes.NewStringNull(),
		TargetGroup:           supertypes.NewStringNull(),
		Targets:               supertypes.NewListNestedObjectValueOfNull[PoolModelMembersIPAddress](ctx),
	}

	if data.MemberGroupRef != nil {
		members.TargetGroup.Set(data.MemberGroupRef.ID)
	}
	if data.GracefulTimeoutPeriod != nil {
		members.GracefulTimeoutPeriod.Set(fmt.Sprintf("%d", *data.GracefulTimeoutPeriod))
	}

	if len(data.Members) != 0 {
		ipAddress := make([]*PoolModelMembersIPAddress, 0)
		for _, m := range data.Members {
			ipa := &PoolModelMembersIPAddress{
				Enabled:   supertypes.NewBoolNull(),
				IPAddress: supertypes.NewStringNull(),
				Port:      supertypes.NewInt64Null(),
				Ratio:     supertypes.NewInt64Null(),
			}

			ipa.Enabled.Set(m.Enabled)
			ipa.IPAddress.Set(m.IPAddress)
			ipa.Port.SetInt(m.Port)
			ipa.Ratio.SetIntPtr(m.Ratio)

			ipAddress = append(ipAddress, ipa)
		}

		diags.Append(members.Targets.Set(ctx, ipAddress)...)
		if diags.HasError() {
			return nil, true, diags
		}
	}

	diags.Append(stateRefreshed.Members.Set(ctx, members)...)
	if diags.HasError() {
		return nil, true, diags
	}

	// * Health
	health := &PoolModelHealth{
		PassiveMonitoringEnabled: supertypes.NewBoolNull(),
		Monitors:                 supertypes.NewListValueOfNull[string](ctx),
	}

	health.PassiveMonitoringEnabled.SetPtr(data.PassiveMonitoringEnabled)

	// prevent unexpected new value: .health.monitors: was null, but now cty.ListValEmpty(cty.String).
	if len(data.HealthMonitors) != 0 {
		monitors := []string{}
		for _, m := range data.HealthMonitors {
			monitors = append(monitors, string(m.Type))
		}

		diags.Append(health.Monitors.Set(ctx, monitors)...)
		if diags.HasError() {
			return nil, true, diags
		}
	}

	diags.Append(stateRefreshed.Health.Set(ctx, health)...)
	if diags.HasError() {
		return nil, true, diags
	}

	// * TLS
	tls := &PoolModelTLS{
		Enabled:                supertypes.NewBoolNull(),
		DomainNames:            supertypes.NewListValueOfNull[string](ctx),
		CaCertificateRefs:      supertypes.NewListValueOfNull[string](ctx),
		CommonNameCheckEnabled: supertypes.NewBoolNull(),
	}

	tls.Enabled.SetPtr(data.SSLEnabled)
	tls.CommonNameCheckEnabled.SetPtr(data.CommonNameCheckEnabled)
	diags.Append(tls.DomainNames.Set(ctx, data.DomainNames)...)
	if diags.HasError() {
		return nil, true, diags
	}

	// prevent unexpected new value: .tls.ca_certificate_refs: was null, but now cty.ListValEmpty(cty.String).
	if len(data.CaCertificateRefs) != 0 {
		refs := []string{}
		for _, ca := range data.CaCertificateRefs {
			refs = append(refs, ca.ID)
		}

		diags.Append(tls.CaCertificateRefs.Set(ctx, refs)...)
		if diags.HasError() {
			return nil, true, diags
		}
	}

	diags.Append(stateRefreshed.TLS.Set(ctx, tls)...)
	if diags.HasError() {
		return nil, true, diags
	}

	// * Persistence
	persistence := &PoolModelPersistence{
		Type:  supertypes.NewStringNull(),
		Value: supertypes.NewStringNull(),
	}

	if data.PersistenceProfile != nil {
		persistence.Type.Set(string(data.PersistenceProfile.Type))
		persistence.Value.Set(data.PersistenceProfile.Value)
	}
	diags.Append(stateRefreshed.Persistence.Set(ctx, persistence)...)
	if diags.HasError() {
		return nil, true, diags
	}

	return stateRefreshed, true, nil
}
