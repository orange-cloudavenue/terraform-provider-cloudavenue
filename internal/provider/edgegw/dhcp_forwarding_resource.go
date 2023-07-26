// Package edgegw provides a Terraform resource.
package edgegw

import (
	"context"
	"fmt"

	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/edgegw"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/mutex"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/org"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/uuid"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &dhcpForwardingResource{}
	_ resource.ResourceWithConfigure   = &dhcpForwardingResource{}
	_ resource.ResourceWithImportState = &dhcpForwardingResource{}
	_ resource.ResourceWithModifyPlan  = &dhcpForwardingResource{}
)

// NewDhcpForwardingResource is a helper function to simplify the provider implementation.
func NewDhcpForwardingResource() resource.Resource {
	return &dhcpForwardingResource{}
}

// dhcpForwardingResource is the resource implementation.
type dhcpForwardingResource struct {
	client *client.CloudAvenue
	org    org.Org
	edgegw edgegw.EdgeGateway
}

// Init Initializes the resource.
func (r *dhcpForwardingResource) Init(ctx context.Context, rm *DhcpForwardingModel) (diags diag.Diagnostics) {
	var err error

	r.org, diags = org.Init(r.client)
	if diags.HasError() {
		return
	}

	r.edgegw, err = r.org.GetEdgeGateway(edgegw.BaseEdgeGW{
		ID:   types.StringValue(rm.EdgeGatewayID.Get()),
		Name: types.StringValue(rm.EdgeGatewayName.Get()),
	})
	if err != nil {
		diags.AddError("Error retrieving Edge Gateway", err.Error())
		return
	}

	return
}

// Metadata returns the resource type name.
func (r *dhcpForwardingResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_dhcp_forwarding"
}

// Schema defines the schema for the resource.
func (r *dhcpForwardingResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = dhcpForwardingSchema(ctx).GetResource(ctx)
}

func (r *dhcpForwardingResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *dhcpForwardingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &DhcpForwardingModel{}

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

	// ! If Enabled is set to false, then DHCP Servers cannot be edited \0_o/
	if plan.DhcpServers.IsKnown() && !plan.Enabled.Get() {
		resp.Diagnostics.AddError("DHCP servers cannot be set", "DHCP servers can only be set when DHCP forwarding is enabled")
		return
	}

	resp.Diagnostics.Append(r.createOrUpdate(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.ID.Set(r.edgegw.GetID())
	state, _, d := r.read(ctx, plan)
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *dhcpForwardingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &DhcpForwardingModel{}

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
func (r *dhcpForwardingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var (
		plan  = &DhcpForwardingModel{}
		state = &DhcpForwardingModel{}
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

	// ! If Enabled is set to false, then DHCP Servers cannot be edited \0_o/
	if !plan.DhcpServers.Equal(state.DhcpServers) && !plan.Enabled.Get() {
		resp.Diagnostics.AddError("DHCP Servers cannot be edited", "DHCP servers can only be edited when DHCP forwarding is enabled")
		return
	}

	resp.Diagnostics.Append(r.createOrUpdate(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, _, d := r.read(ctx, plan)
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *dhcpForwardingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	state := &DhcpForwardingModel{}

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

	vdcOrVDCGroup, err := r.edgegw.GetParent()
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving Edge Gateway parent", err.Error())
		return
	}

	if vdcOrVDCGroup.IsVDCGroup() {
		mutex.GlobalMutex.KvLock(ctx, vdcOrVDCGroup.GetID())
		defer mutex.GlobalMutex.KvUnlock(ctx, vdcOrVDCGroup.GetID())
	} else {
		mutex.GlobalMutex.KvLock(ctx, r.edgegw.GetID())
		defer mutex.GlobalMutex.KvUnlock(ctx, r.edgegw.GetID())
	}

	// There is no "delete" for DHCP forwarding. It can only be updated to empty values (disabled)
	if _, err := r.edgegw.UpdateDhcpForwarder(&govcdtypes.NsxtEdgeGatewayDhcpForwarder{}); err != nil {
		resp.Diagnostics.AddError("Error deleting DHCP forwarding", err.Error())
		return
	}
}

func (r *dhcpForwardingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var (
		edgegwID, edgegwName string
		d                    diag.Diagnostics
		err                  error
	)

	r.org, d = org.Init(r.client)
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	if uuid.IsEdgeGateway(req.ID) {
		edgegwID = req.ID
	} else {
		edgegwName = req.ID
	}

	r.edgegw, err = r.org.GetEdgeGateway(edgegw.BaseEdgeGW{
		ID:   types.StringValue(edgegwID),
		Name: types.StringValue(edgegwName),
	})
	if err != nil {
		resp.Diagnostics.AddError("Failed to import DHCP Forwarding.", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), r.edgegw.GetID())...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("edge_gateway_id"), r.edgegw.GetID())...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("edge_gateway_name"), r.edgegw.GetName())...)
}

// ModifyPlan Check if DHCP servers can be edited.
func (r *dhcpForwardingResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	dPlan := &DhcpForwardingModel{}
	dState := &DhcpForwardingModel{}

	if d := req.Plan.Get(ctx, dPlan); d.HasError() {
		// return because plan is empty
		return
	}
	if d := req.State.Get(ctx, dState); d.HasError() {
		// return because state is empty
		return
	}

	// ! If Enabled is set to false, then DHCP Servers cannot be edited \0_o/
	if !dPlan.DhcpServers.Equal(dState.DhcpServers) && !dPlan.Enabled.Get() {
		resp.Diagnostics.AddError("DHCP Servers cannot be edited", "DHCP servers can only be edited when DHCP forwarding is enabled")
		return
	}
}

// * CustomFuncs

func (r *dhcpForwardingResource) read(ctx context.Context, planOrState *DhcpForwardingModel) (stateRefreshed *DhcpForwardingModel, found bool, diags diag.Diagnostics) {
	stateRefreshed = planOrState.Copy()

	dhcpForwardConfig, err := r.edgegw.GetDhcpForwarder()
	if err != nil {
		if govcd.ContainsNotFound(err) {
			return nil, false, nil
		}
		diags.AddError("Error retrieving NSX-T Edge Gateway DHCP forwarding", err.Error())
		return nil, true, diags
	}

	if !stateRefreshed.ID.IsKnown() {
		stateRefreshed.ID.Set(r.edgegw.GetID())
	}

	stateRefreshed.Enabled.Set(dhcpForwardConfig.Enabled)
	stateRefreshed.EdgeGatewayID.Set(r.edgegw.GetID())
	stateRefreshed.EdgeGatewayName.Set(r.edgegw.GetName())
	diags.Append(stateRefreshed.DhcpServers.Set(ctx, dhcpForwardConfig.DhcpServers)...)
	if diags.HasError() {
		return nil, true, diags
	}

	return stateRefreshed, true, nil
}

func (r *dhcpForwardingResource) createOrUpdate(ctx context.Context, plan *DhcpForwardingModel) (diags diag.Diagnostics) {
	vdcOrVDCGroup, err := r.edgegw.GetParent()
	if err != nil {
		diags.AddError("Error retrieving Edge Gateway parent", err.Error())
		return
	}

	if vdcOrVDCGroup.IsVDCGroup() {
		mutex.GlobalMutex.KvLock(ctx, vdcOrVDCGroup.GetID())
		defer mutex.GlobalMutex.KvUnlock(ctx, vdcOrVDCGroup.GetID())
	} else {
		mutex.GlobalMutex.KvLock(ctx, r.edgegw.GetID())
		defer mutex.GlobalMutex.KvUnlock(ctx, r.edgegw.GetID())
	}

	dhcpForwarderConfig, d := plan.ToNsxtEdgeGatewayDhcpForwarder(ctx)
	if d.HasError() {
		diags.Append(d...)
		return
	}

	if _, err := r.edgegw.UpdateDhcpForwarder(dhcpForwarderConfig); err != nil {
		diags.AddError("Error on change DHCP forwarding configuration", err.Error())
		return
	}

	return nil
}
