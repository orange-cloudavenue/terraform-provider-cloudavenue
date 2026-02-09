/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

// Package vapp provides a Terraform resource.
package vapp

import (
	"context"
	"errors"
	"fmt"
	"strings"

	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vapp"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vdc"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &isolatedNetworkResource{}
	_ resource.ResourceWithConfigure   = &isolatedNetworkResource{}
	_ resource.ResourceWithImportState = &isolatedNetworkResource{}
)

// NewIsolatedNetworkResource is a helper function to simplify the provider implementation.
func NewIsolatedNetworkResource() resource.Resource {
	return &isolatedNetworkResource{}
}

// isolatedNetworkResource is the resource implementation.
type isolatedNetworkResource struct {
	client *client.CloudAvenue
	vdc    vdc.VDC
	vapp   vapp.VAPP
}

// Metadata returns the resource type name.
func (r *isolatedNetworkResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_" + "isolated_network"
}

// Schema defines the schema for the resource.
func (r *isolatedNetworkResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = isolatedNetworkSchema().GetResource(ctx)
}

func (r *isolatedNetworkResource) Init(_ context.Context, rm *isolatedNetworkModel) (diags diag.Diagnostics) {
	r.vdc, diags = vdc.Init(r.client, rm.VDC.StringValue)
	if diags.HasError() {
		return diags
	}

	r.vapp, diags = vapp.Init(r.client, r.vdc, rm.VAppID.StringValue, rm.VAppName.StringValue)

	return diags
}

func (r *isolatedNetworkResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *isolatedNetworkResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer metrics.New("cloudavenue_vapp_isolated_network", r.client.GetOrgName(), metrics.Create)()

	plan := &isolatedNetworkModel{}

	// Retrieve values from plan
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Init resource
	resp.Diagnostics.Append(r.Init(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Lock vApp
	resp.Diagnostics.Append(r.vapp.LockVAPP(ctx)...)
	if resp.Diagnostics.HasError() {
		return
	}
	defer r.vapp.UnlockVAPP(ctx)

	vappNetworkSettings, d := r.buildVappNetworkObject(ctx, plan)
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	// Create network
	_, err := r.vapp.CreateVappNetwork(vappNetworkSettings, nil)
	if err != nil {
		resp.Diagnostics.AddError("Error creating VApp isolated network", err.Error())
		return
	}

	// Use generic read function to refresh the state
	stateRefreshed, _, d := r.read(ctx, plan)
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, stateRefreshed)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *isolatedNetworkResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer metrics.New("cloudavenue_vapp_isolated_network", r.client.GetOrgName(), metrics.Read)()

	state := &isolatedNetworkModel{}

	// Get current state
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Init resource
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
func (r *isolatedNetworkResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer metrics.New("cloudavenue_vapp_isolated_network", r.client.GetOrgName(), metrics.Update)()

	var (
		plan  = &isolatedNetworkModel{}
		state = &isolatedNetworkModel{}
	)

	// Get current plan and state
	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Init resource
	resp.Diagnostics.Append(r.Init(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Lock vApp
	resp.Diagnostics.Append(r.vapp.LockVAPP(ctx)...)
	if resp.Diagnostics.HasError() {
		return
	}
	defer r.vapp.UnlockVAPP(ctx)

	vappNetworkSettings, d := r.buildVappNetworkObject(ctx, plan)
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	// Update network
	_, err := r.vapp.UpdateNetwork(vappNetworkSettings, nil)
	if err != nil {
		resp.Diagnostics.AddError("Error updating vApp network", err.Error())
		return
	}

	// Use generic read function to refresh the state
	stateRefreshed, _, d := r.read(ctx, plan)
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, stateRefreshed)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *isolatedNetworkResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer metrics.New("cloudavenue_vapp_isolated_network", r.client.GetOrgName(), metrics.Delete)()

	state := &isolatedNetworkModel{}

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

	// Lock vApp
	resp.Diagnostics.Append(r.vapp.LockVAPP(ctx)...)
	if resp.Diagnostics.HasError() {
		return
	}
	defer r.vapp.UnlockVAPP(ctx)

	_, err := r.vapp.RemoveNetwork(state.ID.String())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting vApp network", err.Error())
		return
	}
}

func (r *isolatedNetworkResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	defer metrics.New("cloudavenue_vapp_isolated_network", r.client.GetOrgName(), metrics.Import)()

	idParts := strings.Split(req.ID, ".")

	if len(idParts) != 3 && len(idParts) != 2 {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: vdc.vapp_name.network_name or vapp_name.network_name. Got: %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("vapp_name"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), idParts[1])...)

	if len(idParts) == 3 {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("vdc"), idParts[0])...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("vapp_name"), idParts[1])...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), idParts[2])...)
	}
}

// * CustomFuncs

// read is a generic read function that can be used by the resource Create, Read and Update functions.
func (r *isolatedNetworkResource) read(ctx context.Context, planOrState *isolatedNetworkModel) (stateRefreshed *isolatedNetworkModel, found bool, diags diag.Diagnostics) {
	stateRefreshed = planOrState.Copy()

	net, err := r.findNetwork(planOrState.Name.Get())
	if err != nil {
		diags.AddError("Error finding network", err.Error())
		return nil, false, diags
	}

	// Get UUID.
	networkID, err := govcd.GetUuidFromHref(net.Link.HREF, false)
	if err != nil {
		diags.AddError("Error creating vApp network uuid", err.Error())
		return stateRefreshed, found, diags
	}

	planOrState.ID.Set(urn.Normalize(urn.Network, networkID).String())
	planOrState.Name.Set(net.NetworkName)
	planOrState.VDC.Set(r.vdc.GetName())
	planOrState.VAppName.Set(r.vapp.GetName())
	planOrState.VAppID.Set(r.vapp.GetID())
	planOrState.Description.Set(net.Description)
	planOrState.GuestVLANAllowed.SetPtr(net.Configuration.GuestVlanAllowed)
	planOrState.RetainIPMacEnabled.SetPtr(net.Configuration.RetainNetInfoAcrossDeployments)

	if len(net.Configuration.IPScopes.IPScope) > 0 {
		planOrState.Netmask.Set(net.Configuration.IPScopes.IPScope[0].Netmask)
		planOrState.Gateway.Set(net.Configuration.IPScopes.IPScope[0].Gateway)
		planOrState.DNS1.Set(net.Configuration.IPScopes.IPScope[0].DNS1)
		planOrState.DNS2.Set(net.Configuration.IPScopes.IPScope[0].DNS2)
		planOrState.DNSSuffix.Set(net.Configuration.IPScopes.IPScope[0].DNSSuffix)
	}

	if net.Configuration.IPScopes.IPScope[0].IPRanges != nil {
		ipPool := make([]*isolatedNetworkModelStaticIPPool, 0)
		for _, ipRange := range net.Configuration.IPScopes.IPScope[0].IPRanges.IPRange {
			ipPool = append(ipPool, &isolatedNetworkModelStaticIPPool{
				StartAddress: supertypes.NewStringValue(ipRange.StartAddress),
				EndAddress:   supertypes.NewStringValue(ipRange.EndAddress),
			})
		}
		diags.Append(planOrState.StaticIPPool.Set(ctx, ipPool)...)
		if diags.HasError() {
			return stateRefreshed, found, diags
		}
	} else {
		planOrState.StaticIPPool.SetNull(ctx)
	}

	return planOrState, true, diags
}

// find network in network list.
func (r *isolatedNetworkResource) findNetwork(networkName string) (network *govcdtypes.VAppNetworkConfiguration, err error) {
	vAppNetworkConfig, err := r.vapp.GetNetworkConfig()
	if err != nil {
		return nil, err
	}

	for _, networkConfig := range vAppNetworkConfig.NetworkConfig {
		if networkConfig.NetworkName == networkName {
			return &networkConfig, nil
		}
	}

	return nil, errors.New("network not found")
}

// build vapp isolated network object.
func (r *isolatedNetworkResource) buildVappNetworkObject(ctx context.Context, plan *isolatedNetworkModel) (vappNetworkSettings *govcd.VappNetworkSettings, diags diag.Diagnostics) {
	staticIPPools, d := plan.StaticIPPool.Get(ctx)
	if d.HasError() {
		diags.Append(d...)
		return vappNetworkSettings, diags
	}

	staticIPRanges := make([]*govcdtypes.IPRange, 0)
	for _, staticIPPool := range staticIPPools {
		staticIPRanges = append(staticIPRanges, &govcdtypes.IPRange{
			StartAddress: staticIPPool.StartAddress.Get(),
			EndAddress:   staticIPPool.EndAddress.Get(),
		})
	}

	return &govcd.VappNetworkSettings{
		Name:               plan.Name.Get(),
		Description:        plan.Description.Get(),
		Gateway:            plan.Gateway.Get(),
		NetMask:            plan.Netmask.Get(),
		DNS1:               plan.DNS1.Get(),
		DNS2:               plan.DNS2.Get(),
		DNSSuffix:          plan.DNSSuffix.Get(),
		StaticIPRanges:     staticIPRanges,
		RetainIpMacEnabled: plan.RetainIPMacEnabled.GetPtr(),
		GuestVLANAllowed:   plan.GuestVLANAllowed.GetPtr(),
	}, diags
}
