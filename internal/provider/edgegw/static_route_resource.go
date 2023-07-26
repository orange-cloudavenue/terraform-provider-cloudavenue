// Package edgegw provides a Terraform resource.
package edgegw

import (
	"context"
	"fmt"
	"strings"

	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/edgegw"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/mutex"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/org"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/uuid"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &staticRouteResource{}
	_ resource.ResourceWithConfigure   = &staticRouteResource{}
	_ resource.ResourceWithImportState = &staticRouteResource{}
)

// NewStaticRouteResource is a helper function to simplify the provider implementation.
func NewStaticRouteResource() resource.Resource {
	return &staticRouteResource{}
}

// staticRouteResource is the resource implementation.
type staticRouteResource struct {
	client *client.CloudAvenue
	org    org.Org
	edgegw edgegw.EdgeGateway
}

// Init Initializes the resource.
func (r *staticRouteResource) Init(ctx context.Context, rm *StaticRouteModel) (diags diag.Diagnostics) {
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
func (r *staticRouteResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_static_route"
}

// Schema defines the schema for the resource.
func (r *staticRouteResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = staticRouteSchema(ctx).GetResource(ctx)
}

func (r *staticRouteResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *staticRouteResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &StaticRouteModel{}

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

	stateRouteConfig, d := plan.ToNsxtEdgeGatewayStaticRoute(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	createdStaticRoute, err := r.edgegw.CreateStaticRoute(stateRouteConfig)
	if err != nil {
		resp.Diagnostics.AddError("Error creating static route", err.Error())
		return
	}

	plan.ID.Set(createdStaticRoute.NsxtEdgeGatewayStaticRoute.ID)
	state, _, d := r.read(ctx, plan)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *staticRouteResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &StaticRouteModel{}

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
func (r *staticRouteResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var (
		plan  = &StaticRouteModel{}
		state = &StaticRouteModel{}
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

	staticRoute, err := r.edgegw.GetStaticRouteById(plan.ID.Get())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving static route", err.Error())
		return
	}

	staticRouteConfig, d := plan.ToNsxtEdgeGatewayStaticRoute(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	staticRouteConfig.ID = staticRoute.NsxtEdgeGatewayStaticRoute.ID
	staticRouteConfig.Version = staticRoute.NsxtEdgeGatewayStaticRoute.Version

	if _, err := staticRoute.Update(staticRouteConfig); err != nil {
		resp.Diagnostics.AddError("Error updating static route", err.Error())
		return
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
func (r *staticRouteResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	state := &StaticRouteModel{}

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

	staticRoute, err := r.edgegw.GetStaticRouteById(state.ID.Get())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving static route", err.Error())
		return
	}

	if err := staticRoute.Delete(); err != nil {
		resp.Diagnostics.AddError("Error deleting static route", err.Error())
		return
	}
}

func (r *staticRouteResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var (
		edgegwID, edgegwName string
		d                    diag.Diagnostics
		err                  error
		staticRoute          *govcd.NsxtEdgeGatewayStaticRoute
	)

	// Split req.ID with dot. ID format is EdgeGatewayIDOrName.StaticRouteNameOrID
	idParts := strings.Split(req.ID, ".")

	if len(idParts) != 2 {
		resp.Diagnostics.AddError("Invalid ID format", "ID format is EdgeGatewayIDOrName.StaticRouteNameOrID")
		return
	}

	r.org, d = org.Init(r.client)
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	if uuid.IsEdgeGateway(idParts[0]) {
		edgegwID = idParts[0]
	} else {
		edgegwName = idParts[0]
	}

	r.edgegw, err = r.org.GetEdgeGateway(edgegw.BaseEdgeGW{
		ID:   types.StringValue(edgegwID),
		Name: types.StringValue(edgegwName),
	})
	if err != nil {
		resp.Diagnostics.AddError("Failed to import DHCP Forwarding.", err.Error())
		return
	}

	// Static Route ID is not a URN
	if uuid.IsUUIDV4(idParts[1]) {
		staticRoute, err = r.edgegw.GetStaticRouteById(idParts[1])
	} else {
		staticRoute, err = r.edgegw.GetStaticRouteByName(idParts[1])
	}
	if err != nil {
		resp.Diagnostics.AddError("Failed to Get DHCP Forwarding.", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), staticRoute.NsxtEdgeGatewayStaticRoute.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), staticRoute.NsxtEdgeGatewayStaticRoute.Name)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("edge_gateway_id"), r.edgegw.GetID())...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("edge_gateway_name"), r.edgegw.GetName())...)
}

// * CustomFuncs

func (r *staticRouteResource) read(ctx context.Context, planOrState *StaticRouteModel) (stateRefreshed *StaticRouteModel, found bool, diags diag.Diagnostics) {
	stateRefreshed = planOrState.Copy()

	var (
		staticRoute *govcd.NsxtEdgeGatewayStaticRoute
		err         error
	)

	if planOrState.ID.IsKnown() {
		staticRoute, err = r.edgegw.GetStaticRouteById(planOrState.ID.Get())
	} else {
		staticRoute, err = r.edgegw.GetStaticRouteByName(planOrState.Name.Get())
	}
	if err != nil {
		if govcd.ContainsNotFound(err) {
			return nil, false, nil
		}
		diags.AddError("Error retrieving static route", err.Error())
		return nil, true, diags
	}

	stateRefreshed.ID.Set(staticRoute.NsxtEdgeGatewayStaticRoute.ID)
	stateRefreshed.Name.Set(staticRoute.NsxtEdgeGatewayStaticRoute.Name)
	stateRefreshed.Description = utils.SuperStringValueOrNull(staticRoute.NsxtEdgeGatewayStaticRoute.Description)
	stateRefreshed.NetworkCidr.Set(staticRoute.NsxtEdgeGatewayStaticRoute.NetworkCidr)
	stateRefreshed.EdgeGatewayID.Set(r.edgegw.GetID())
	stateRefreshed.EdgeGatewayName.Set(r.edgegw.GetName())

	nHs := make(StaticRouteModelNextHops, 0)
	for _, nextHop := range staticRoute.NsxtEdgeGatewayStaticRoute.NextHops {
		nH := StaticRouteModelNextHop{}
		nH.AdminDistance.Set(int64(nextHop.AdminDistance))
		nH.IPAddress.Set(nextHop.IPAddress)
		nHs = append(nHs, nH)
	}
	stateRefreshed.NextHops.Set(ctx, nHs)

	return stateRefreshed, true, nil
}
