/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

// Package network provides a Terraform resource.
package network

import (
	"context"
	"fmt"
	"strings"

	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/mutex"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/org"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &dhcpBindingResource{}
	_ resource.ResourceWithConfigure   = &dhcpBindingResource{}
	_ resource.ResourceWithImportState = &dhcpBindingResource{}
)

// NewDhcpBindingResource is a helper function to simplify the provider implementation.
func NewDhcpBindingResource() resource.Resource {
	return &dhcpBindingResource{}
}

// dhcpBindingResource is the resource implementation.
type dhcpBindingResource struct {
	client *client.CloudAvenue
	org    org.Org
}

// Init Initializes the resource.
func (r *dhcpBindingResource) Init(_ context.Context, _ *DHCPBindingModel) (diags diag.Diagnostics) {
	r.org, diags = org.Init(r.client)

	return diags
}

// Metadata returns the resource type name.
func (r *dhcpBindingResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_dhcp_binding"
}

// Schema defines the schema for the resource.
func (r *dhcpBindingResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = dhcpBindingSchema(ctx).GetResource(ctx)
}

func (r *dhcpBindingResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *dhcpBindingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer metrics.New("cloudavenue_network_dhcp_binding", r.client.GetOrgName(), metrics.Create)()

	plan := &DHCPBindingModel{}

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

	mutex.GlobalMutex.KvLock(ctx, plan.OrgNetworkID.Get())
	defer mutex.GlobalMutex.KvUnlock(ctx, plan.OrgNetworkID.Get())

	orgNetwork, err := r.org.GetOpenApiOrgVdcNetworkById(plan.OrgNetworkID.Get())
	if err != nil {
		resp.Diagnostics.AddError("Failed to get network DHCP", err.Error())
		return
	}

	// Check if the DHCP is enabled
	if !orgNetwork.IsDhcpEnabled() {
		resp.Diagnostics.AddError("DHCP is not enabled on the network", "Please use 'cloudavenue_network_dhcp' resource to enable it")
		return
	}

	dhcpBindingConfig, d := plan.ToNetworkDhcpBindingType(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	createdDhcpBinding, err := orgNetwork.CreateOpenApiOrgVdcNetworkDhcpBinding(dhcpBindingConfig)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create DHCP binding", err.Error())
		return
	}

	plan.ID.Set(createdDhcpBinding.OpenApiOrgVdcNetworkDhcpBinding.ID)
	state, _, d := r.read(ctx, plan)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *dhcpBindingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer metrics.New("cloudavenue_network_dhcp_binding", r.client.GetOrgName(), metrics.Read)()

	state := &DHCPBindingModel{}

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

	plan, found, d := r.read(ctx, state)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !found {
		resp.State.RemoveResource(ctx)
		return
	}

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *dhcpBindingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer metrics.New("cloudavenue_network_dhcp_binding", r.client.GetOrgName(), metrics.Update)()

	var (
		plan  = &DHCPBindingModel{}
		state = &DHCPBindingModel{}
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

	mutex.GlobalMutex.KvLock(ctx, plan.OrgNetworkID.Get())
	defer mutex.GlobalMutex.KvUnlock(ctx, plan.OrgNetworkID.Get())

	orgNetwork, err := r.org.GetOpenApiOrgVdcNetworkById(plan.OrgNetworkID.Get())
	if err != nil {
		resp.Diagnostics.AddError("Failed to get network DHCP", err.Error())
		return
	}

	dhcpBinding, err := orgNetwork.GetOpenApiOrgVdcNetworkDhcpBindingById(state.ID.Get())
	if err != nil {
		resp.Diagnostics.AddError("Failed to get DHCP binding", err.Error())
		return
	}

	dhcpBindingConfig, d := plan.ToNetworkDhcpBindingType(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	if _, err := dhcpBinding.Update(dhcpBindingConfig); err != nil {
		resp.Diagnostics.AddError("Failed to update DHCP binding", err.Error())
		return
	}

	stateUpdated, _, d := r.read(ctx, plan)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, stateUpdated)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *dhcpBindingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer metrics.New("cloudavenue_network_dhcp_binding", r.client.GetOrgName(), metrics.Delete)()

	state := &DHCPBindingModel{}

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

	mutex.GlobalMutex.KvLock(ctx, state.OrgNetworkID.Get())
	defer mutex.GlobalMutex.KvUnlock(ctx, state.OrgNetworkID.Get())

	orgNetwork, err := r.org.GetOpenApiOrgVdcNetworkById(state.OrgNetworkID.Get())
	if err != nil {
		resp.Diagnostics.AddError("Failed to get network DHCP", err.Error())
		return
	}

	dhcpBinding, err := orgNetwork.GetOpenApiOrgVdcNetworkDhcpBindingById(state.ID.Get())
	if err != nil {
		resp.Diagnostics.AddError("Failed to get DHCP binding", err.Error())
		return
	}

	if err := dhcpBinding.Delete(); err != nil {
		resp.Diagnostics.AddError("Failed to delete DHCP binding", err.Error())
		return
	}
}

func (r *dhcpBindingResource) read(ctx context.Context, planOrState *DHCPBindingModel) (state *DHCPBindingModel, found bool, diags diag.Diagnostics) {
	refreshed := planOrState.Copy()

	orgNetwork, err := r.org.GetOpenApiOrgVdcNetworkById(refreshed.OrgNetworkID.Get())
	if err != nil {
		diags.AddError("Failed to get network DHCP", err.Error())
		return nil, true, diags
	}

	var dhcpBinding *govcd.OpenApiOrgVdcNetworkDhcpBinding

	if refreshed.ID.IsKnown() {
		dhcpBinding, err = orgNetwork.GetOpenApiOrgVdcNetworkDhcpBindingById(refreshed.ID.Get())
	} else {
		dhcpBinding, err = orgNetwork.GetOpenApiOrgVdcNetworkDhcpBindingByName(refreshed.Name.Get())
	}
	if err != nil {
		if govcd.ContainsNotFound(err) {
			return nil, false, nil
		}
		diags.AddError("Failed to get DHCP binding", err.Error())
		return nil, true, diags
	}

	refreshed.ID.Set(dhcpBinding.OpenApiOrgVdcNetworkDhcpBinding.ID)
	refreshed.Name.Set(dhcpBinding.OpenApiOrgVdcNetworkDhcpBinding.Name)
	refreshed.Description.Set(dhcpBinding.OpenApiOrgVdcNetworkDhcpBinding.Description)
	refreshed.IPAddress.Set(dhcpBinding.OpenApiOrgVdcNetworkDhcpBinding.IpAddress)
	refreshed.MacAddress.Set(dhcpBinding.OpenApiOrgVdcNetworkDhcpBinding.MacAddress)
	refreshed.LeaseTime.SetIntPtr(dhcpBinding.OpenApiOrgVdcNetworkDhcpBinding.LeaseTime)

	diags.Append(refreshed.DNSServers.Set(ctx, dhcpBinding.OpenApiOrgVdcNetworkDhcpBinding.DnsServers)...)
	if diags.HasError() {
		return nil, true, diags
	}

	if dhcpBinding.OpenApiOrgVdcNetworkDhcpBinding.DhcpV4BindingConfig != nil {
		x := DHCPBindingModelDhcpV4Config{}
		x.GatewayAddress.Set(dhcpBinding.OpenApiOrgVdcNetworkDhcpBinding.DhcpV4BindingConfig.GatewayIPAddress)
		x.Hostname.Set(dhcpBinding.OpenApiOrgVdcNetworkDhcpBinding.DhcpV4BindingConfig.HostName)

		diags.Append(refreshed.DhcpV4Config.Set(ctx, x)...)
		if diags.HasError() {
			return nil, true, diags
		}
	}

	return refreshed, true, diags
}

// ImportState imports a resource from orgNetworkID.DhcpBindingName.
func (r *dhcpBindingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	defer metrics.New("cloudavenue_network_dhcp_binding", r.client.GetOrgName(), metrics.Import)()

	// Init the resource
	resp.Diagnostics.Append(r.Init(ctx, nil)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get URI from import ID
	resourceURI := strings.Split(req.ID, ".")
	if len(resourceURI) != 2 {
		resp.Diagnostics.AddError("Invalid import ID format.", "The import ID should be in the format orgNetworkID.DhcpBindingName")
		return
	}
	orgNetworkID, bindingName := resourceURI[0], resourceURI[1]

	orgNetwork, err := r.org.GetOpenApiOrgVdcNetworkById(orgNetworkID)
	if err != nil {
		resp.Diagnostics.AddError("Failed to get org network", err.Error())
		return
	}

	dhcpBinding, err := orgNetwork.GetOpenApiOrgVdcNetworkDhcpBindingByName(bindingName)
	if err != nil {
		resp.Diagnostics.AddError("Failed to get DHCP binding", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), dhcpBinding.OpenApiOrgVdcNetworkDhcpBinding.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("org_network_id"), orgNetworkID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), dhcpBinding.OpenApiOrgVdcNetworkDhcpBinding.Name)...)
}
