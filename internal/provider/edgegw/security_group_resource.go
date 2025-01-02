// Package edgegw provides a Terraform resource.
package edgegw

import (
	"context"
	"fmt"
	"strings"

	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/edgegw"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/mutex"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/org"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &securityGroupResource{}
	_ resource.ResourceWithConfigure   = &securityGroupResource{}
	_ resource.ResourceWithImportState = &securityGroupResource{}
)

// NewSecurityGroupResource is a helper function to simplify the provider implementation.
func NewSecurityGroupResource() resource.Resource {
	return &securityGroupResource{}
}

// securityGroupResource is the resource implementation.
type securityGroupResource struct {
	client *client.CloudAvenue
	org    org.Org
	edgegw edgegw.EdgeGateway
}

// Init Initializes the resource.
func (r *securityGroupResource) Init(ctx context.Context, rm *SecurityGroupModel) (diags diag.Diagnostics) {
	var err error

	r.org, diags = org.Init(r.client)
	if diags.HasError() {
		return
	}

	r.edgegw, err = r.org.GetEdgeGateway(edgegw.BaseEdgeGW{
		ID:   rm.EdgeGatewayID.StringValue,
		Name: rm.EdgeGatewayName.StringValue,
	})
	if err != nil {
		diags.AddError("Error retrieving Edge Gateway", err.Error())
		return
	}

	return
}

// Metadata returns the resource type name.
func (r *securityGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_security_group"
}

// Schema defines the schema for the resource.
func (r *securityGroupResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = securityGroupSchema(ctx).GetResource(ctx)
}

func (r *securityGroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *securityGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer metrics.New("cloudavenue_edgegateway_security_group", r.client.GetOrgName(), metrics.Create)()

	plan := &SecurityGroupModel{}

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

	mutex.GlobalMutex.KvLock(ctx, r.edgegw.GetID())
	defer mutex.GlobalMutex.KvUnlock(ctx, r.edgegw.GetID())

	values, d := plan.ToSDKSecurityGroupModel(ctx)
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	fwsg, err := r.edgegw.CreateFirewallSecurityGroup(values)
	if err != nil {
		resp.Diagnostics.AddError("Error creating security group", err.Error())
		return
	}

	// Set the ID
	plan.ID.Set(fwsg.ID)

	// Use generic read function to refresh the state
	state, found, d := r.read(ctx, plan)
	if !found {
		resp.Diagnostics.AddError("Security group not found", "The security group was not found after creation")
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
func (r *securityGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer metrics.New("cloudavenue_edgegateway_security_group", r.client.GetOrgName(), metrics.Read)()

	state := &SecurityGroupModel{}

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
func (r *securityGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer metrics.New("cloudavenue_edgegateway_security_group", r.client.GetOrgName(), metrics.Update)()

	var (
		plan  = &SecurityGroupModel{}
		state = &SecurityGroupModel{}
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

	mutex.GlobalMutex.KvLock(ctx, r.edgegw.GetID())
	defer mutex.GlobalMutex.KvUnlock(ctx, r.edgegw.GetID())

	fwsg, err := r.edgegw.GetFirewallSecurityGroup(state.ID.Get())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving security group", err.Error())
		return
	}

	values, d := plan.ToSDKSecurityGroupModel(ctx)
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	if err := fwsg.Update(values); err != nil {
		resp.Diagnostics.AddError("Error updating security group", err.Error())
		return
	}

	// Use generic read function to refresh the state
	stateRefreshed, found, d := r.read(ctx, plan)
	if !found {
		resp.Diagnostics.AddError("Security group not found", "The security group was not found after update")
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
func (r *securityGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer metrics.New("cloudavenue_edgegateway_security_group", r.client.GetOrgName(), metrics.Delete)()

	state := &SecurityGroupModel{}

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

	mutex.GlobalMutex.KvLock(ctx, r.edgegw.GetID())
	defer mutex.GlobalMutex.KvUnlock(ctx, r.edgegw.GetID())

	fwsg, err := r.edgegw.GetFirewallSecurityGroup(state.ID.Get())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving security group", err.Error())
		return
	}

	if err := fwsg.Delete(); err != nil {
		resp.Diagnostics.AddError("Error deleting security group", err.Error())
		return
	}
}

func (r *securityGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	defer metrics.New("cloudavenue_edgegateway_security_group", r.client.GetOrgName(), metrics.Import)()

	// id format is EdgeGatewayNameOrID.SecurityGroupNameOrID
	// * Import with custom logic
	idParts := strings.Split(req.ID, ".")

	if len(idParts) != 2 {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: EdgeGatewayNameOrID.SecurityGroupNameOrID Got: %q", req.ID),
		)
		return
	}
	edgeGatewayNameOrID, securityGroupNameOrID := idParts[0], idParts[1]

	x := &SecurityGroupModel{
		ID:              supertypes.NewStringNull(),
		Name:            supertypes.NewStringNull(),
		EdgeGatewayName: supertypes.NewStringNull(),
		EdgeGatewayID:   supertypes.NewStringNull(),
		Description:     supertypes.NewStringNull(),
		Members:         supertypes.NewSetValueOfNull[string](ctx),
	}

	if urn.IsEdgeGateway(edgeGatewayNameOrID) {
		x.EdgeGatewayID.Set(edgeGatewayNameOrID)
	} else {
		x.EdgeGatewayName.Set(edgeGatewayNameOrID)
	}

	if urn.IsSecurityGroup(securityGroupNameOrID) {
		x.ID.Set(securityGroupNameOrID)
	} else {
		x.Name.Set(securityGroupNameOrID)
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

func (r *securityGroupResource) read(ctx context.Context, planOrState *SecurityGroupModel) (stateRefreshed *SecurityGroupModel, found bool, diags diag.Diagnostics) {
	stateRefreshed = planOrState.Copy()

	/*
		Implement the resource read here
	*/

	idOrName := planOrState.Name.Get()
	if planOrState.ID.IsKnown() {
		idOrName = planOrState.ID.Get()
	}

	fwsg, err := r.edgegw.GetFirewallSecurityGroup(idOrName)
	if govcd.ContainsNotFound(err) {
		return nil, false, nil
	}
	if err != nil {
		diags.AddError("Error retrieving security group", err.Error())
		return nil, true, diags
	}

	stateRefreshed.ID.Set(fwsg.ID)
	stateRefreshed.Name.Set(fwsg.Name)
	stateRefreshed.Description.Set(fwsg.Description)
	stateRefreshed.EdgeGatewayID.Set(r.edgegw.GetID())
	stateRefreshed.EdgeGatewayName.Set(r.edgegw.GetName())

	if fwsg.Members != nil || len(fwsg.Members) > 0 {
		diags.Append(stateRefreshed.Members.Set(ctx, utils.OpenAPIReferenceToSliceID(fwsg.Members))...)
	} else {
		stateRefreshed.Members.SetNull(ctx)
	}

	return stateRefreshed, true, diags
}
