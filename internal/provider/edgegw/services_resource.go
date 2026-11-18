// SPDX-FileCopyrightText: Copyright (c) 2026 Orange
// SPDX-License-Identifier: MPL-2.0
//
// This software is distributed under the MPL-2.0 license.
// the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
// or see the "LICENSE" file for more details.

// SPDX-FileCopyrightText: Copyright (c) 2025-2026 Orange
// SPDX-License-Identifier: MPL-2.0
//
// This software is distributed under the MPL-2.0 license.
// the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
// or see the "LICENSE" file for more details.

package edgegw

import (
	"context"
	"fmt"

	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	edge "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/api/edgegateway/v1"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/cloudavenue"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &ServicesResource{}
	_ resource.ResourceWithConfigure   = &ServicesResource{}
	_ resource.ResourceWithImportState = &ServicesResource{}
)

// NewServicesResource is a helper function to simplify the provider implementation.
func NewServicesResource() resource.Resource {
	return &ServicesResource{}
}

// ServicesResource is the resource implementation.
type ServicesResource struct {
	client     *client.CloudAvenue
	edgeClient *edge.Client
	edgegw     *types.ModelEdgeGateway
}

// Init Initializes the resource.
func (r *ServicesResource) Init(ctx context.Context, rm *ServicesModel) (diags diag.Diagnostics) {
	var err error
	var v2 cav.Client

	if v2, err = r.client.V2(); err == nil {
		r.edgeClient, err = edge.New(v2)
	}
	if err != nil {
		diags.AddError("Error initializing EdgeGateway client", err.Error())
		return diags
	}

	r.edgegw, err = r.edgeClient.GetEdgeGateway(ctx, types.ParamsEdgeGateway{
		ID:   rm.EdgeGatewayID.Get(),
		Name: rm.EdgeGatewayName.Get(),
	})
	if err != nil {
		diags.AddError("Error retrieving Edge Gateway", err.Error())
		return diags
	}

	return diags
}

// Metadata returns the resource type name.
func (r *ServicesResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_services"
}

// Schema defines the schema for the resource.
func (r *ServicesResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = servicesSchema(ctx).GetResource(ctx)
}

func (r *ServicesResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *ServicesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer metrics.New("cloudavenue_edgegateway_services", r.client.GetOrgName(), metrics.Create)()

	plan := &ServicesModel{}

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

	// * Allow the service to already be enabled as some edges may have the parameter enabled by default.
	networkServices, err := r.edgeClient.GetServices(ctx, types.ParamsEdgeGateway{ID: r.edgegw.ID, Name: r.edgegw.Name})
	if err != nil {
		resp.Diagnostics.AddError("Edge Gateway Services Retrieval Error", fmt.Sprintf("Failed to retrieve services of edge gateway: %s", err))
		return
	}
	if networkServices.Services != nil {
		cloudavenue.Lock(ctx)
		defer cloudavenue.Unlock(ctx)

		if err := r.edgeClient.EnableCloudavenueServices(ctx, types.ParamsEdgeGateway{ID: r.edgegw.ID, Name: r.edgegw.Name}); err != nil {
			resp.Diagnostics.AddError("EdgeGateway Service Activation Error", fmt.Sprintf("Failed to enable network service: %s", err))
			return
		}
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
func (r *ServicesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer metrics.New("cloudavenue_edgegateway_services", r.client.GetOrgName(), metrics.Read)()

	state := &ServicesModel{}

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
func (r *ServicesResource) Update(_ context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer metrics.New("cloudavenue_edgegateway_services", r.client.GetOrgName(), metrics.Update)()

	// No update available for this resource
	resp.Diagnostics.AddError("Resource Update Error", "The resource does not support updates")
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *ServicesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer metrics.New("cloudavenue_edgegateway_services", r.client.GetOrgName(), metrics.Delete)()

	state := &ServicesModel{}

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

	cloudavenue.Lock(ctx)
	defer cloudavenue.Unlock(ctx)

	if err := r.edgeClient.DisableCloudavenueServices(ctx, types.ParamsEdgeGateway{ID: r.edgegw.ID, Name: r.edgegw.Name}); err != nil {
		resp.Diagnostics.AddError("EdgeGateway Service Disablement Error", fmt.Sprintf("Failed to disable network service: %s", err))
		return
	}
}

func (r *ServicesResource) ImportState(_ context.Context, _ resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	defer metrics.New("cloudavenue_edgegateway_services", r.client.GetOrgName(), metrics.Import)()

	// This resource does not support import. Create a new resource instead of importing.
	resp.Diagnostics.AddError("Resource Import Error", "The resource does not support import. Create a new resource instead of importing.")
}

// * CustomFuncs

// read is a generic read function that can be used by the resource Create, Read and Update functions.
func (r *ServicesResource) read(ctx context.Context, planOrState *ServicesModel) (stateRefreshed *ServicesModel, found bool, diags diag.Diagnostics) {
	stateRefreshed = planOrState.Copy()

	/*
		Implement the resource read here
	*/

	networkServices, err := r.edgeClient.GetServices(ctx, types.ParamsEdgeGateway{ID: r.edgegw.ID, Name: r.edgegw.Name})
	if err != nil {
		diags.AddError("Edge Gateway Services Retrieval Error", fmt.Sprintf("Failed to retrieve services of edge gateway: %s", err))
		return nil, false, diags
	}
	if networkServices.Services == nil {
		diags.AddError("Edge Gateway Service Disabled", "The edge gateway service is disabled")
		return nil, false, diags
	}

	stateRefreshed.ID.Set(r.edgegw.ID)
	stateRefreshed.EdgeGatewayName.Set(r.edgegw.Name)
	stateRefreshed.EdgeGatewayID.Set(r.edgegw.ID)
	stateRefreshed.Network.Set(networkServices.Services.Network)
	stateRefreshed.IPAddress.Set(networkServices.Services.IPAddress)
	svcs := map[string]*ServicesModelServices{}

	for _, svc := range networkServices.Services.Services {
		if _, ok := svcs[svc.Category]; !ok {
			svcs[svc.Category] = &ServicesModelServices{
				Network:  supertypes.NewStringValueOrNull(svc.Network),
				Services: supertypes.NewMapNestedObjectValueOfNull[ServicesModelService](ctx),
			}
		}

		svs := map[string]*ServicesModelService{}

		for _, s := range svc.Services {
			if _, ok := svs[s.Name]; !ok {
				ports := make([]*ServicesModelServicePorts, len(s.Ports))
				for i, p := range s.Ports {
					ports[i] = &ServicesModelServicePorts{
						Port:     supertypes.NewInt32Null(),
						Protocol: supertypes.NewStringNull(),
					}

					ports[i].Port.SetInt(p.Port)
					ports[i].Protocol.Set(p.Protocol)
				}

				svs[s.Name] = &ServicesModelService{
					Name:        supertypes.NewStringValue(s.Name),
					Description: supertypes.NewStringValue(s.Description),
					IPs:         supertypes.NewListValueOfSlice(ctx, s.IPs),
					FQDNs:       supertypes.NewListValueOfSlice(ctx, s.FQDNs),
					Ports:       supertypes.NewListNestedObjectValueOfSlice(ctx, ports),
				}
			}
		}

		diags.Append(svcs[svc.Category].Services.Set(ctx, svs)...)
		if diags.HasError() {
			return nil, true, diags
		}
	}

	diags.Append(stateRefreshed.Services.Set(ctx, svcs)...)

	return stateRefreshed, true, diags
}
