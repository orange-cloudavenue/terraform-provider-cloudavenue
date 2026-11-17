// SPDX-FileCopyrightText: Copyright (c) 2026 Orange
// SPDX-License-Identifier: MPL-2.0
//
// This software is distributed under the MPL-2.0 license.
// the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
// or see the "LICENSE" file for more details.

// SPDX-FileCopyrightText: Copyright (c) 2023-2026 Orange
// SPDX-License-Identifier: MPL-2.0
//
// This software is distributed under the MPL-2.0 license.
// the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
// or see the "LICENSE" file for more details.

// Package edgegw provides a Terraform resource.
package edgegw

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	edge "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/api/edgegateway/v1"
	vdcgroup "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/api/vdcgroup/v1"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	cav_types "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	cerrs "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/errors"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/mutex"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &ipSetResource{}
	_ resource.ResourceWithConfigure   = &ipSetResource{}
	_ resource.ResourceWithImportState = &ipSetResource{}
)

// NewIpSetResource is a helper function to simplify the provider implementation.
func NewIPSetResource() resource.Resource {
	return &ipSetResource{}
}

// ipSetResource is the resource implementation.
type ipSetResource struct {
	client         *client.CloudAvenue
	edgeClient     *edge.Client
	edgegw         *cav_types.ModelEdgeGateway
	vdcgroupClient *vdcgroup.Client
}

// Init Initializes the resource.
func (r *ipSetResource) Init(ctx context.Context, rm *IPSetModel) (diags diag.Diagnostics) {
	var err error
	var v2 cav.Client

	if v2, err = r.client.V2(); err == nil {
		r.edgeClient, err = edge.New(v2)
	}
	if err != nil {
		diags.AddError("Error initializing Edge Gateway client", err.Error())
		return diags
	}
	if r.vdcgroupClient, err = vdcgroup.New(v2); err == nil {
		diags.AddError("Error initializing Vdc Group client", err.Error())
		return diags
	}

	r.edgegw, err = r.edgeClient.GetEdgeGateway(ctx, cav_types.ParamsEdgeGateway{ID: rm.EdgeGatewayID.Get(), Name: rm.EdgeGatewayName.Get()})
	if err != nil {
		diags.AddError("Error retrieving Edge Gateway", err.Error())
		return diags
	}

	_, err = r.vdcgroupClient.GetVDCGroup(ctx, cav_types.ParamsGetVDCGroup{ID: r.edgegw.OwnerRef.ID, Name: r.edgegw.OwnerRef.Name})
	if err != nil {
		diags.AddError(
			"Edge Gateway belongs to a VDC Group",
			"Edge Gateway belongs to a VDC Group, please use the VDC Group resource instead `cloudavenue_vdcg_ip_set`.",
		)
	}

	return diags
}

// Metadata returns the resource type name.
func (r *ipSetResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_ip_set"
}

// Schema defines the schema for the resource.
func (r *ipSetResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ipSetSchema(ctx).GetResource(ctx)
}

func (r *ipSetResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *ipSetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer metrics.New("cloudavenue_edgegateway_ip_set", r.client.GetOrgName(), metrics.Create)()

	plan := &IPSetModel{}

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

	var createdIPSet *cav_types.ModelGetFirewallGroup

	_, err := r.vdcgroupClient.GetVDCGroup(ctx, cav_types.ParamsGetVDCGroup{ID: r.edgegw.OwnerRef.ID, Name: r.edgegw.OwnerRef.Name})
	if err != nil {
		mutex.GlobalMutex.KvLock(ctx, r.edgegw.OwnerRef.ID)
		defer mutex.GlobalMutex.KvUnlock(ctx, r.edgegw.OwnerRef.ID)

		ipAddrs := make([]string, 0)
		resp.Diagnostics.Append(plan.IPAddresses.Get(ctx, &ipAddrs, false)...)

		createdIPSet, err = r.vdcgroupClient.CreateIPSet(ctx, cav_types.ParamsCreateIPSet{Name: plan.Name.Get(), Description: plan.Description.Get(), VDCGroupID: r.edgegw.OwnerRef.ID, VDCGroupName: r.edgegw.OwnerRef.Name, IPAddresses: ipAddrs})
	} else {
		mutex.GlobalMutex.KvLock(ctx, r.edgegw.ID)
		defer mutex.GlobalMutex.KvUnlock(ctx, r.edgegw.ID)
		ipAddrs := make([]string, 0)
		resp.Diagnostics.Append(plan.IPAddresses.Get(ctx, &ipAddrs, false)...)

		createdIPSet, err = r.edgeClient.CreateIPSet(ctx, cav_types.ParamsCreateEdgeGatewayIPSet{Name: plan.Name.Get(), Description: plan.Description.Get(), EdgeGatewayID: r.edgegw.ID, EdgeGatewayName: r.edgegw.Name, IPAddresses: ipAddrs})
	}
	if err != nil {
		resp.Diagnostics.AddError("Error creating IP Set", err.Error())
		return
	}

	plan.ID.Set(createdIPSet.ID)
	stateRefreshed, _, d := r.read(ctx, plan)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, stateRefreshed)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *ipSetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer metrics.New("cloudavenue_edgegateway_ip_set", r.client.GetOrgName(), metrics.Read)()

	state := &IPSetModel{}

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
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, stateRefreshed)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *ipSetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer metrics.New("cloudavenue_edgegateway_ip_set", r.client.GetOrgName(), metrics.Update)()
	var (
		plan  = &IPSetModel{}
		state = &IPSetModel{}
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

	var (
		d           diag.Diagnostics
		vdcGroupErr error
	)

	ipAddrs := make([]string, 0)
	resp.Diagnostics.Append(plan.IPAddresses.Get(ctx, &ipAddrs, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if _, vdcGroupErr = r.vdcgroupClient.GetVDCGroup(ctx, cav_types.ParamsGetVDCGroup{ID: r.edgegw.OwnerRef.ID, Name: r.edgegw.OwnerRef.Name}); vdcGroupErr != nil {
		mutex.GlobalMutex.KvLock(ctx, r.edgegw.OwnerRef.ID)
		defer mutex.GlobalMutex.KvUnlock(ctx, r.edgegw.OwnerRef.ID)
		_, err := r.vdcgroupClient.GetIPSet(ctx, cav_types.ParamsGetIPSet{ID: state.ID.Get(), Name: state.Name.Get(), VDCGroupID: r.edgegw.OwnerRef.ID, VDCGroupName: r.edgegw.OwnerRef.Name})
		if err != nil {
			resp.Diagnostics.AddError("Error retrieving existing IP Set", err.Error())
			return
		}
		_, err = r.vdcgroupClient.UpdateIPSet(ctx, cav_types.ParamsUpdateIPSet{ID: state.ID.Get(), Name: plan.Name.Get(), Description: plan.Description.Get(), VDCGroupID: r.edgegw.OwnerRef.ID, VDCGroupName: r.edgegw.OwnerRef.Name, IPAddresses: ipAddrs})
		if err != nil {
			resp.Diagnostics.AddError("Error updating IP Set", err.Error())
			return
		}
	} else {
		mutex.GlobalMutex.KvLock(ctx, r.edgegw.ID)
		defer mutex.GlobalMutex.KvUnlock(ctx, r.edgegw.ID)
		_, err := r.edgeClient.GetIPSet(ctx, cav_types.ParamsGetEdgeGatewayIPSet{ID: state.ID.Get(), Name: state.Name.Get(), EdgeGatewayID: r.edgegw.ID, EdgeGatewayName: r.edgegw.Name})
		if err != nil {
			resp.Diagnostics.AddError("Error retrieving existing IP Set", err.Error())
			return
		}
		_, err = r.edgeClient.UpdateIPSet(ctx, cav_types.ParamsUpdateEdgeGatewayIPSet{ID: state.ID.Get(), Name: plan.Name.Get(), Description: plan.Description.Get(), EdgeGatewayID: r.edgegw.ID, EdgeGatewayName: r.edgegw.Name, IPAddresses: ipAddrs})
		if err != nil {
			resp.Diagnostics.AddError("Error updating IP Set", err.Error())
			return
		}
	}

	stateRefreshed, _, d := r.read(ctx, plan)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, stateRefreshed)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *ipSetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer metrics.New("cloudavenue_edgegateway_ip_set", r.client.GetOrgName(), metrics.Delete)()

	state := &IPSetModel{}

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

	var vdcGroupErr error

	if _, vdcGroupErr = r.vdcgroupClient.GetVDCGroup(ctx, cav_types.ParamsGetVDCGroup{ID: r.edgegw.OwnerRef.ID, Name: r.edgegw.OwnerRef.Name}); vdcGroupErr != nil {
		mutex.GlobalMutex.KvLock(ctx, r.edgegw.OwnerRef.ID)
		defer mutex.GlobalMutex.KvUnlock(ctx, r.edgegw.OwnerRef.ID)
		_, err := r.vdcgroupClient.GetIPSet(ctx, cav_types.ParamsGetIPSet{ID: state.ID.Get(), VDCGroupID: r.edgegw.OwnerRef.ID, VDCGroupName: r.edgegw.OwnerRef.Name})
		if err != nil {
			resp.Diagnostics.AddError("Error retrieving IP Set", err.Error())
			return
		}
		err = r.vdcgroupClient.DeleteIPSet(ctx, cav_types.ParamsDeleteIPSet{ID: state.ID.Get(), VDCGroupID: r.edgegw.OwnerRef.ID, VDCGroupName: r.edgegw.OwnerRef.Name})
		if err != nil {
			resp.Diagnostics.AddError("Error deleting IP Set", err.Error())
			return
		}
	} else {
		mutex.GlobalMutex.KvLock(ctx, r.edgegw.ID)
		defer mutex.GlobalMutex.KvUnlock(ctx, r.edgegw.ID)
		_, err := r.edgeClient.GetIPSet(ctx, cav_types.ParamsGetEdgeGatewayIPSet{ID: state.ID.Get(), EdgeGatewayID: r.edgegw.ID, EdgeGatewayName: r.edgegw.Name})
		if err != nil {
			resp.Diagnostics.AddError("Error retrieving IP Set", err.Error())
			return
		}
		err = r.edgeClient.DeleteIPSet(ctx, cav_types.ParamsDeleteEdgeGatewayIPSet{ID: state.ID.Get(), EdgeGatewayID: r.edgegw.ID, EdgeGatewayName: r.edgegw.Name})
		if err != nil {
			resp.Diagnostics.AddError("Error deleting IP Set", err.Error())
			return
		}
	}
}

func (r *ipSetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	defer metrics.New("cloudavenue_edgegateway_ip_set", r.client.GetOrgName(), metrics.Import)()

	// id format is edgeGatewayIDOrName.ipSetName
	idParts := strings.Split(req.ID, ".")

	if len(idParts) != 2 {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: edgeGatewayIDOrName.ipSetName. Got: %q", req.ID),
		)
		return
	}

	var (
		err   error
		ipSet *cav_types.ModelGetFirewallGroup
	)

	state := &IPSetModel{}
	state.EdgeGatewayID.Set(idParts[0])
	state.Name.Set(idParts[1])

	// Init the resource
	resp.Diagnostics.Append(r.Init(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if _, vdcGroupErr := r.vdcgroupClient.GetVDCGroup(ctx, cav_types.ParamsGetVDCGroup{ID: r.edgegw.OwnerRef.ID, Name: r.edgegw.OwnerRef.Name}); vdcGroupErr != nil {
		ipSet, err = r.vdcgroupClient.GetIPSet(ctx, cav_types.ParamsGetIPSet{Name: state.Name.Get(), VDCGroupID: r.edgegw.OwnerRef.ID, VDCGroupName: r.edgegw.OwnerRef.Name})
	} else {
		ipSet, err = r.edgeClient.GetIPSet(ctx, cav_types.ParamsGetEdgeGatewayIPSet{Name: state.Name.Get(), EdgeGatewayID: r.edgegw.ID, EdgeGatewayName: r.edgegw.Name})
	}
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving IP Set", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), ipSet.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), ipSet.Name)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(edgeGatewayID), r.edgegw.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(edgeGatewayName), r.edgegw.Name)...)
}

func (r *ipSetResource) read(ctx context.Context, planOrState *IPSetModel) (stateRefreshed *IPSetModel, found bool, diags diag.Diagnostics) {
	stateRefreshed = planOrState.Copy()

	var ipSetConfig *cav_types.ModelGetFirewallGroup

	var err error
	var vdcGroupErr error
	if _, vdcGroupErr = r.vdcgroupClient.GetVDCGroup(ctx, cav_types.ParamsGetVDCGroup{ID: r.edgegw.OwnerRef.ID, Name: r.edgegw.OwnerRef.Name}); vdcGroupErr != nil {
		ipSetConfig, err = r.vdcgroupClient.GetIPSet(ctx, cav_types.ParamsGetIPSet{ID: stateRefreshed.ID.Get(), Name: stateRefreshed.Name.Get(), VDCGroupID: r.edgegw.OwnerRef.ID, VDCGroupName: r.edgegw.OwnerRef.Name})
	} else {
		ipSetConfig, err = r.edgeClient.GetIPSet(ctx, cav_types.ParamsGetEdgeGatewayIPSet{ID: stateRefreshed.ID.Get(), Name: stateRefreshed.Name.Get(), EdgeGatewayID: r.edgegw.ID, EdgeGatewayName: r.edgegw.Name})
	}
	if err != nil {
		if cerrs.IsNotFound(err) {
			return stateRefreshed, false, nil
		}
		cerrs.AddError(&diags, cerrs.ActionRead, "IP Set", err)
		return stateRefreshed, found, diags
	}

	stateRefreshed.ID.Set(ipSetConfig.ID)
	stateRefreshed.Name.Set(ipSetConfig.Name)
	stateRefreshed.Description.Set(ipSetConfig.Description)
	stateRefreshed.EdgeGatewayID.Set(r.edgegw.ID)
	stateRefreshed.EdgeGatewayName.Set(r.edgegw.Name)
	diags.Append(stateRefreshed.IPAddresses.Set(ctx, ipSetConfig.IPAddresses)...)

	return stateRefreshed, true, diags
}
