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

// Package publicip provides a Terraform resource to manage public IPs.
package publicip

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/cloudavenue"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/edgegw"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/org"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &publicIPResource{}
	_ resource.ResourceWithConfigure   = &publicIPResource{}
	_ resource.ResourceWithImportState = &publicIPResource{}
)

// NewPublicIPResource returns a new resource implementing the public_ip resource.
func NewPublicIPResource() resource.Resource {
	return &publicIPResource{}
}

// publicIPResource is the resource implementation.
type publicIPResource struct {
	client *client.CloudAvenue
	org    org.Org
	edgegw edgegw.EdgeGateway
}

// Metadata returns the resource type name.
func (r *publicIPResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName
}

// Init.
func (r *publicIPResource) Init(_ context.Context, rm *publicIPResourceModel) (diags diag.Diagnostics) {
	var err error

	r.org, diags = org.Init(r.client)
	if diags.HasError() {
		return diags
	}

	// EdgeGatewayID and EdgeGatewayName are Null if ImportState
	if rm.EdgeGatewayID.IsKnown() || rm.EdgeGatewayName.IsKnown() {
		r.edgegw, err = r.org.GetEdgeGateway(edgegw.BaseEdgeGW{
			ID:   rm.EdgeGatewayID.StringValue,
			Name: rm.EdgeGatewayName.StringValue,
		})
		if err != nil {
			diags.AddError("Error retrieving Edge Gateway", err.Error())
			return diags
		}
	}

	return diags
}

// Schema defines the schema for the resource.
func (r *publicIPResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = publicIPSchema().GetResource(ctx)
}

func (r *publicIPResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *publicIPResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer metrics.New("cloudavenue_publicip", r.client.GetOrgName(), metrics.Create)()

	// Retrieve values from plan
	plan := &publicIPResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(r.Init(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create() is passed a default timeout to use if no value
	// has been supplied in the Terraform configuration.
	createTimeout, errTO := plan.Timeouts.Create(ctx, 5*time.Minute)
	if errTO != nil {
		resp.Diagnostics.AddError(
			"Error creating timeout",
			"Could not create timeout, unexpected error",
		)
		return
	}

	cloudavenue.Lock(ctx)
	defer cloudavenue.Unlock(ctx)

	// * Create the public IP
	job, err := r.client.CAVSDK.V1.PublicIP.New(r.edgegw.GetID())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Public IP",
			fmt.Sprintf("Could not create Public IP, unexpected error: %s", err),
		)
		return
	}

	// * Wait for job to complete
	if err := job.Wait(3, int(createTimeout.Seconds())); err != nil {
		resp.Diagnostics.AddError(
			"Error creating Public IP",
			fmt.Sprintf("Could not create Public IP, unexpected error: %s", err),
		)
		return
	}

	ipRefreshed, err := r.client.CAVSDK.V1.PublicIP.GetIPByJob(job)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting Public IP",
			fmt.Sprintf("Could not get Public IP, unexpected error: %s", err),
		)
		return
	}

	plan.ID.Set(ipRefreshed.UplinkIP)
	plan.PublicIP.Set(ipRefreshed.UplinkIP)
	state, found, d := r.read(ctx, plan)
	if !found {
		resp.State.RemoveResource(ctx)
		resp.Diagnostics.AddError("Error reading Public IP", "Public IP not found after creation")
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
func (r *publicIPResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer metrics.New("cloudavenue_publicip", r.client.GetOrgName(), metrics.Read)()

	state := &publicIPResourceModel{}

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
func (r *publicIPResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	plan := &publicIPResourceModel{}
	state := &publicIPResourceModel{}

	// Get current plan
	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get current state
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	/*
		Implement the resource read here
	*/

	stateRefreshed, found, d := r.read(ctx, state)
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, stateRefreshed)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *publicIPResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer metrics.New("cloudavenue_publicip", r.client.GetOrgName(), metrics.Delete)()

	state := &publicIPResourceModel{}

	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(r.Init(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete() is passed a default timeout to use if no value
	// has been supplied in the Terraform configuration.
	deleteTimeout, errTO := state.Timeouts.Delete(ctx, 8*time.Minute)
	if errTO != nil {
		resp.Diagnostics.AddError(
			"Error creating timeout",
			"Could not create timeout, unexpected error",
		)
		return
	}

	ip, err := r.client.CAVSDK.V1.PublicIP.GetIP(state.ID.Get())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting Public IP",
			fmt.Sprintf("Could not get Public IP, unexpected error: %s", err),
		)
		return
	}

	cloudavenue.Lock(ctx)
	defer cloudavenue.Unlock(ctx)

	// * Delete the public IP
	job, err := ip.Delete()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Public IP",
			fmt.Sprintf("Could not delete Public IP, unexpected error: %s", err),
		)
		return
	}

	// * Wait for job to complete
	if err := job.Wait(3, int(deleteTimeout.Seconds())); err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Public IP",
			fmt.Sprintf("Could not delete Public IP, unexpected error: %s", err),
		)
		return
	}
}

func (r *publicIPResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	defer metrics.New("cloudavenue_publicip", r.client.GetOrgName(), metrics.Import)()

	// Slipt the ID into EdgeGatewayIDOrName and PublicIP
	// ID format: EdgeGatewayIDOrName.PublicIP
	idSplit := strings.Split(req.ID, ".")
	if len(idSplit) != 5 {
		resp.Diagnostics.AddError(
			"Error importing Public IP",
			fmt.Sprintf("Could not import Public IP, unexpected ID format: %s", req.ID),
		)
		return
	}

	edgeGwNameOrID, publicIP := idSplit[0], fmt.Sprintf("%s.%s.%s.%s", idSplit[1], idSplit[2], idSplit[3], idSplit[4])
	if urn.IsEdgeGateway(edgeGwNameOrID) {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("edge_gateway_id"), edgeGwNameOrID)...)
	} else {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("edge_gateway_name"), edgeGwNameOrID)...)
	}

	ip, err := r.client.CAVSDK.V1.PublicIP.GetIP(publicIP)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting Public IP",
			fmt.Sprintf("Could not get Public IP, unexpected error: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), ip.UplinkIP)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("public_ip"), ip.UplinkIP)...)
}

// * CustomFuncs

func (r *publicIPResource) read(_ context.Context, planOrState *publicIPResourceModel) (stateRefreshed *publicIPResourceModel, found bool, diags diag.Diagnostics) {
	stateRefreshed = planOrState.Copy()

	pubIP, err := r.client.CAVSDK.V1.PublicIP.GetIP(planOrState.ID.Get())
	if err != nil {
		if errors.Is(err, fmt.Errorf("not found")) {
			return stateRefreshed, false, nil
		}
		diags.AddError("Error getting Public IP", err.Error())
		return stateRefreshed, true, diags
	}

	stateRefreshed.ID.Set(pubIP.UplinkIP)
	stateRefreshed.PublicIP.Set(pubIP.UplinkIP)
	stateRefreshed.EdgeGatewayID.Set(r.edgegw.EdgeGateway.ID)
	stateRefreshed.EdgeGatewayName.Set(r.edgegw.EdgeGateway.Name)

	return stateRefreshed, true, nil
}
