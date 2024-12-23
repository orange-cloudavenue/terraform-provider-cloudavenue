package vdcg

import (
	"context"
	"fmt"
	"strings"

	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
	v1 "github.com/orange-cloudavenue/cloudavenue-sdk-go/v1"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/mutex"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &NetworkIsolatedResource{}
	_ resource.ResourceWithConfigure   = &NetworkIsolatedResource{}
	_ resource.ResourceWithImportState = &NetworkIsolatedResource{}
)

// NewNetworkIsolatedResource is a helper function to simplify the provider implementation.
func NewNetworkIsolatedResource() resource.Resource {
	return &NetworkIsolatedResource{}
}

// NetworkIsolatedResource is the resource implementation.
type NetworkIsolatedResource struct {
	client *client.CloudAvenue
	vdcg   *v1.VDCGroup
}

// Init Initializes the resource.
func (r *NetworkIsolatedResource) Init(ctx context.Context, rm *networkIsolatedModel) (diags diag.Diagnostics) {
	var err error

	idOrName := rm.VDCGroupName.Get()
	if rm.VDCGroupID.IsKnown() {
		// Use the ID
		idOrName = rm.VDCGroupID.Get()
	}

	r.vdcg, err = r.client.CAVSDK.V1.VDC().GetVDCGroup(idOrName)
	if err != nil {
		diags.AddError("Error retrieving VDC Group", err.Error())
		return
	}
	return
}

// Metadata returns the resource type name.
func (r *NetworkIsolatedResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_network_isolated"
}

// Schema defines the schema for the resource.
func (r *NetworkIsolatedResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = networkIsolatedSchema(ctx).GetResource(ctx)
}

func (r *NetworkIsolatedResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *NetworkIsolatedResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer metrics.New("cloudavenue_vdcg_network_isolated", r.client.GetOrgName(), metrics.Create)()

	plan := &networkIsolatedModel{}

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
	values, d := plan.ToSDK(ctx)
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	mutex.GlobalMutex.KvLock(ctx, r.vdcg.GetID())
	defer mutex.GlobalMutex.KvUnlock(ctx, r.vdcg.GetID())

	networkIsolated, err := r.vdcg.CreateNetworkIsolated(values)
	if err != nil {
		resp.Diagnostics.AddError("Error creating isolated network", err.Error())
		return
	}

	plan.ID.Set(networkIsolated.ID)

	// Use generic read function to refresh the state
	state, found, d := r.read(ctx, plan)
	if !found {
		resp.State.RemoveResource(ctx)
		resp.Diagnostics.AddError("Resource not found after creation", "The resource was not found after creation.")
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
func (r *NetworkIsolatedResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer metrics.New("cloudavenue_vdcg_network_isolated", r.client.GetOrgName(), metrics.Read)()

	state := &networkIsolatedModel{}

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
func (r *NetworkIsolatedResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer metrics.New("cloudavenue_vdcg_network_isolated", r.client.GetOrgName(), metrics.Update)()

	var (
		plan  = &networkIsolatedModel{}
		state = &networkIsolatedModel{}
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

	values, d := plan.ToSDK(ctx)
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	mutex.GlobalMutex.KvLock(ctx, r.vdcg.GetID())
	defer mutex.GlobalMutex.KvUnlock(ctx, r.vdcg.GetID())

	net, err := r.vdcg.GetNetworkIsolated(state.ID.Get())
	if err != nil {
		resp.Diagnostics.AddError("Error getting isolated network", err.Error())
		return
	}

	values.ID = state.ID.Get()

	// Update the network
	if err := net.Update(values); err != nil {
		resp.Diagnostics.AddError("Error updating isolated network", err.Error())
		return
	}

	// Use generic read function to refresh the state
	stateRefreshed, found, d := r.read(ctx, state)
	if !found {
		resp.State.RemoveResource(ctx)
		resp.Diagnostics.AddError("Resource not found after update", "The resource was not found after update. Please refresh the state.")
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
func (r *NetworkIsolatedResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer metrics.New("cloudavenue_vdcg_network_isolated", r.client.GetOrgName(), metrics.Delete)()

	state := &networkIsolatedModel{}

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

	mutex.GlobalMutex.KvLock(ctx, r.vdcg.GetID())
	defer mutex.GlobalMutex.KvUnlock(ctx, r.vdcg.GetID())

	net, err := r.vdcg.GetNetworkIsolated(state.ID.Get())
	if err != nil {
		resp.Diagnostics.AddError("Error getting isolated network", err.Error())
		return
	}

	// Delete the network
	if err := net.Delete(); err != nil {
		resp.Diagnostics.AddError("Error deleting isolated network", err.Error())
		return
	}
}

func (r *NetworkIsolatedResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	defer metrics.New("cloudavenue_vdcg_network_isolated", r.client.GetOrgName(), metrics.Import)()

	// * Import with custom logic
	idParts := strings.Split(req.ID, ".")

	if len(idParts) != 2 {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: vdcGroupNameOrID.networkNameOrID Got: %q", req.ID),
		)
		return
	}

	x := &networkIsolatedModel{
		ID:           supertypes.NewStringNull(),
		Name:         supertypes.NewStringNull(),
		VDCGroupName: supertypes.NewStringNull(),
		VDCGroupID:   supertypes.NewStringNull(),
	}

	if urn.IsVDCGroup(idParts[0]) {
		x.VDCGroupID.Set(idParts[0])
	} else {
		x.VDCGroupName.Set(idParts[0])
	}

	if urn.IsNetwork(idParts[1]) {
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
func (r *NetworkIsolatedResource) read(ctx context.Context, planOrState *networkIsolatedModel) (stateRefreshed *networkIsolatedModel, found bool, diags diag.Diagnostics) {
	stateRefreshed = planOrState.Copy()

	/*
		Implement the resource read here
	*/

	var (
		net *v1.VDCNetworkIsolated
		err error
	)

	if urn.IsNetwork(planOrState.ID.Get()) {
		net, err = r.vdcg.GetNetworkIsolated(planOrState.ID.Get())
	} else {
		net, err = r.vdcg.GetNetworkIsolated(planOrState.Name.Get())
	}
	if err != nil {
		if govcd.ContainsNotFound(err) {
			return nil, false, nil
		}
		diags.AddError("Error getting isolated network", err.Error())
		return
	}

	// Populate the state with the network data
	stateRefreshed.ID.Set(net.ID)
	stateRefreshed.Name.Set(net.Name)
	stateRefreshed.Description.Set(net.Description)
	stateRefreshed.VDCGroupName.Set(r.vdcg.GetName())
	stateRefreshed.VDCGroupID.Set(r.vdcg.GetID())
	stateRefreshed.Gateway.Set(net.Subnet.Gateway)
	stateRefreshed.PrefixLength.SetInt(net.Subnet.PrefixLength)
	stateRefreshed.DNS1.Set(net.Subnet.DNSServer1)
	stateRefreshed.DNS2.Set(net.Subnet.DNSServer2)
	stateRefreshed.DNSSuffix.Set(net.Subnet.DNSSuffix)
	stateRefreshed.GuestVLANAllowed.SetPtr(net.GuestVLANTaggingAllowed)

	x := []*networkIsolatedModelStaticIPPool{}
	for _, ipRange := range net.Subnet.IPRanges {
		n := &networkIsolatedModelStaticIPPool{
			StartAddress: supertypes.NewStringNull(),
			EndAddress:   supertypes.NewStringNull(),
		}
		n.StartAddress.Set(ipRange.StartAddress)
		n.EndAddress.Set(ipRange.EndAddress)
		x = append(x, n)
	}

	diags.Append(stateRefreshed.StaticIPPool.Set(ctx, x)...)
	if diags.HasError() {
		return nil, true, diags
	}

	return stateRefreshed, true, nil
}
