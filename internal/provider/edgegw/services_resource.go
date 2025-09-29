/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package edgegw

import (
	"context"
	"fmt"

	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	edgegateway "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/api/edgegateway/v1"
	sdktypes "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
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
	client *client.CloudAvenue
	// eClient is the Edge Gateway client from the SDK V2
	eClient *edgegateway.Client
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

	// Initialize SDK v2 EdgeGateway client
	eC, err := edgegateway.New(client.V2)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create Edge Gateway client, got error: %s", err),
		)
		return
	}
	r.eClient = eC
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

	// Build params
	params := r.paramsFromModel(plan)

	// Enable CloudAvenue services (sdkv2 already checks if enabled)
	cloudavenue.Lock(ctx)
	defer cloudavenue.Unlock(ctx)
	if err := r.eClient.EnableCloudavenueServices(ctx, params); err != nil {
		resp.Diagnostics.AddError("EdgeGateway Service Activation Error", fmt.Sprintf("Failed to enable network service: %s", err))
		return
	}

	// Use generic read function to refresh the state
	state, d := r.read(ctx, plan)
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

	// Refresh the state using SDK v2
	stateRefreshed, d := r.read(ctx, state)
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		// Remove from state when read fails (resource not available/disabled)
		resp.State.RemoveResource(ctx)
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

	// Disable CloudAvenue services via SDK v2 if enabled
	params := r.paramsFromModel(state)

	cloudavenue.Lock(ctx)
	defer cloudavenue.Unlock(ctx)

	if err := r.eClient.DisableCloudavenueServices(ctx, params); err != nil {
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
func (r *ServicesResource) read(ctx context.Context, planOrState *ServicesModel) (stateRefreshed *ServicesModel, diags diag.Diagnostics) {
	stateRefreshed = planOrState.Copy()

	// Read using SDK v2
	params := r.paramsFromModel(planOrState)
	svc, err := r.eClient.GetServices(ctx, params)
	if err != nil {
		diags.AddError("Client Error", fmt.Sprintf("Failed to retrieve edge gateway services: %s", err))
		return nil, diags
	}
	if svc == nil || svc.Services == nil {
		diags.AddError("Edge Gateway Service Disabled", "The edge gateway service is disabled")
		return nil, diags
	}

	// Top-level
	stateRefreshed.ID.Set(svc.Services.ID)

	stateRefreshed.EdgeGatewayName.Set(svc.Name)
	stateRefreshed.EdgeGatewayID.Set(svc.ID)
	stateRefreshed.Network.Set(svc.Services.Network)
	stateRefreshed.IPAddress.Set(svc.Services.IPAddress)

	// Catalog mapping
	catalogs := map[string]*ServicesModelCatalog{}
	for _, cat := range svc.Services.Services {
		catalogs[cat.Category] = &ServicesModelCatalog{
			Network:  supertypes.NewStringValueOrNull(cat.Network),
			Category: supertypes.NewStringValue(cat.Category),
			Services: supertypes.NewMapNestedObjectValueOfNull[ServicesModelCatalogService](ctx),
		}

		svMap := map[string]*ServicesModelCatalogService{}
		for _, s := range cat.Services {
			ports := make([]*ServicesModelCatalogServicePorts, len(s.Ports))
			for i, p := range s.Ports {
				ports[i] = &ServicesModelCatalogServicePorts{
					Port:     supertypes.NewInt32Null(),
					Protocol: supertypes.NewStringNull(),
				}
				ports[i].Port.SetInt(p.Port)
				ports[i].Protocol.Set(p.Protocol)
			}

			svMap[s.Name] = &ServicesModelCatalogService{
				Name:        supertypes.NewStringValue(s.Name),
				Description: supertypes.NewStringValue(s.Description),
				IPs:         supertypes.NewListValueOfSlice(ctx, s.IPs),
				FQDNs:       supertypes.NewListValueOfSlice(ctx, s.FQDNs),
				Ports:       supertypes.NewListNestedObjectValueOfSlice(ctx, ports),
			}
		}

		diags.Append(catalogs[cat.Category].Services.Set(ctx, svMap)...) // set services for category
		if diags.HasError() {
			return nil, diags
		}
	}

	diags.Append(stateRefreshed.Services.Set(ctx, catalogs)...) // set catalogs map

	return stateRefreshed, diags
}

// paramsFromModel builds SDK v2 params from model
func (r *ServicesResource) paramsFromModel(rm *ServicesModel) sdktypes.ParamsEdgeGateway {
	return sdktypes.ParamsEdgeGateway{
		ID:   rm.EdgeGatewayID.Get(),
		Name: rm.EdgeGatewayName.Get(),
	}
}
