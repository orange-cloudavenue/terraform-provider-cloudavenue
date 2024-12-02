// Package edgegw provides a Terraform resource.
package edgegw

import (
	"context"
	"fmt"
	"strings"

	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"

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
func (r *securityGroupResource) Init(ctx context.Context, rm *securityGroupModel) (diags diag.Diagnostics) {
	var err error

	r.org, diags = org.Init(r.client)
	if diags.HasError() {
		return
	}

	r.edgegw, err = r.org.GetEdgeGateway(edgegw.BaseEdgeGW{
		ID:   rm.EdgeGatewayID,
		Name: rm.EdgeGatewayName,
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

	plan := &securityGroupModel{}

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

	mutex.GlobalMutex.KvLock(ctx, vdcOrVDCGroup.GetID())
	defer mutex.GlobalMutex.KvUnlock(ctx, vdcOrVDCGroup.GetID())

	securityGroup, d := r.securityGroupToNsxtFirewallGroup(ctx, plan)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	newSecGroup, err := r.edgegw.CreateNsxtFirewallGroup(securityGroup)
	if err != nil {
		resp.Diagnostics.AddError("Error creating Security Group", err.Error())
		return
	}

	state, d := r.read(ctx, newSecGroup)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *securityGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer metrics.New("cloudavenue_edgegateway_security_group", r.client.GetOrgName(), metrics.Read)()

	state := &securityGroupModel{}

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

	secGroup, err := r.getSecurityGroup(ctx, state)
	if err != nil {
		if govcd.ContainsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error retrieving Security Group", err.Error())
		return
	}

	plan, d := r.read(ctx, secGroup)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *securityGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer metrics.New("cloudavenue_edgegateway_security_group", r.client.GetOrgName(), metrics.Update)()

	var (
		plan  = &securityGroupModel{}
		state = &securityGroupModel{}
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

	mutex.GlobalMutex.KvLock(ctx, vdcOrVDCGroup.GetID())
	defer mutex.GlobalMutex.KvUnlock(ctx, vdcOrVDCGroup.GetID())

	securityGroup, d := r.securityGroupToNsxtFirewallGroup(ctx, plan)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get current Security Group
	secGroup, err := r.getSecurityGroup(ctx, state)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving Security Group", err.Error())
		return
	}

	// Set actually Security Group ID in new object
	securityGroup.ID = secGroup.NsxtFirewallGroup.ID

	// Update Security Group
	secGroupUpdated, err := secGroup.Update(securityGroup)
	if err != nil {
		resp.Diagnostics.AddError("Error updating Security Group", err.Error())
		return
	}

	// Read updated Security Group
	plan, d = r.read(ctx, secGroupUpdated)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *securityGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer metrics.New("cloudavenue_edgegateway_security_group", r.client.GetOrgName(), metrics.Delete)()

	state := &securityGroupModel{}

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

	mutex.GlobalMutex.KvLock(ctx, vdcOrVDCGroup.GetID())
	defer mutex.GlobalMutex.KvUnlock(ctx, vdcOrVDCGroup.GetID())

	// Get current Security Group
	secGroup, err := r.getSecurityGroup(ctx, state)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving Security Group", err.Error())
		return
	}

	if err := secGroup.Delete(); err != nil {
		resp.Diagnostics.AddError("Error deleting Security Group", err.Error())
		return
	}
}

func (r *securityGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	defer metrics.New("cloudavenue_edgegateway_security_group", r.client.GetOrgName(), metrics.Import)()

	// id format is edgeGatewayIDOrName.securityGroupIDOrName
	idParts := strings.Split(req.ID, ".")

	if len(idParts) != 2 {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: edgeGatewayIDOrName.securityGroupIDOrName. Got: %q", req.ID),
		)
		return
	}

	var (
		id, name, edgegwID, edgegwName string
		d                              diag.Diagnostics
		err                            error
	)

	r.org, d = org.Init(r.client)
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	if urn.IsEdgeGateway(idParts[0]) {
		edgegwID = idParts[0]
	} else {
		edgegwName = idParts[0]
	}

	if urn.IsSecurityGroup(idParts[1]) {
		id = idParts[1]
	} else {
		name = idParts[1]
	}

	r.edgegw, err = r.org.GetEdgeGateway(edgegw.BaseEdgeGW{
		ID:   types.StringValue(edgegwID),
		Name: types.StringValue(edgegwName),
	})
	if err != nil {
		resp.Diagnostics.AddError("Failed to import Security Group.", err.Error())
		return
	}

	securityGroup, err := r.getSecurityGroup(ctx, &securityGroupModel{
		ID:   utils.StringValueOrNull(id),
		Name: utils.StringValueOrNull(name),
	})
	if err != nil {
		resp.Diagnostics.AddError("Failed to import Security Group.", err.Error())
		return
	}

	// ID
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), securityGroup.NsxtFirewallGroup.ID)...)
	// Name
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), securityGroup.NsxtFirewallGroup.Name)...)
	// edge_gateway_id
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("edge_gateway_id"), r.edgegw.GetID())...)
	// edge_gateway_name
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("edge_gateway_name"), r.edgegw.GetName())...)
}

func (r *securityGroupResource) read(ctx context.Context, securityGroup *govcd.NsxtFirewallGroup) (plan *securityGroupModel, diags diag.Diagnostics) {
	if securityGroup == nil || securityGroup.NsxtFirewallGroup == nil {
		diags.AddError("Error retrieving Security Group", "Security Group not found")
		return nil, diags
	}

	return &securityGroupModel{
		ID:                  types.StringValue(securityGroup.NsxtFirewallGroup.ID),
		Name:                types.StringValue(securityGroup.NsxtFirewallGroup.Name),
		EdgeGatewayID:       types.StringValue(r.edgegw.GetID()),
		EdgeGatewayName:     types.StringValue(r.edgegw.GetName()),
		Description:         utils.StringValueOrNull(securityGroup.NsxtFirewallGroup.Description),
		MemberOrgNetworkIDs: utils.OpenAPIReferenceToSliceID(securityGroup.NsxtFirewallGroup.Members).ToTerraformTypesStringSet(ctx),
	}, nil
}

// GetSecurityGroup retrieves the Security Group from the API.
func (r *securityGroupResource) getSecurityGroup(_ context.Context, rm *securityGroupModel) (*govcd.NsxtFirewallGroup, error) {
	parentEdgeGW, err := r.edgegw.GetParent()
	if err != nil {
		return nil, err
	}

	if parentEdgeGW.IsVDCGroup() {
		return parentEdgeGW.GetSecurityGroupByNameOrID(rm.GetIDOrName().ValueString())
	}

	if err := r.edgegw.Refresh(); err != nil {
		return nil, err
	}
	return r.edgegw.GetSecurityGroupByNameOrID(rm.GetIDOrName().ValueString())
}

// func securityGroupToNsxtFirewallGroup.
func (r *securityGroupResource) securityGroupToNsxtFirewallGroup(ctx context.Context, rm *securityGroupModel) (securityGroup *govcdtypes.NsxtFirewallGroup, diags diag.Diagnostics) {
	parentEdgeGW, err := r.edgegw.GetParent()
	if err != nil {
		diags.AddError("Unable to get parent edge gateway", fmt.Sprintf("Unable to get parent edge gateway: %s", err.Error()))
		return
	}

	ownerID := r.edgegw.GetID()

	if parentEdgeGW.IsVDCGroup() {
		ownerID = parentEdgeGW.GetID()
	}

	members, d := rm.MemberOrgNetworkIDsFromPlan(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return
	}

	memberReferences := make([]govcdtypes.OpenApiReference, len(members))
	for index, member := range members {
		memberReferences[index].ID = member
	}

	return &govcdtypes.NsxtFirewallGroup{
		Name:        rm.Name.ValueString(),
		Description: rm.Description.ValueString(),
		OwnerRef:    &govcdtypes.OpenApiReference{ID: ownerID},
		Type:        govcdtypes.FirewallGroupTypeSecurityGroup,
		Members:     memberReferences,
	}, diags
}
