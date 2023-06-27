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

type securityGroupModel struct {
	ID                  types.String `tfsdk:"id"`
	EdgeGatewayID       types.String `tfsdk:"edge_gateway_id"`
	EdgeGatewayName     types.String `tfsdk:"edge_gateway_name"`
	Name                types.String `tfsdk:"name"`
	Description         types.String `tfsdk:"description"`
	MemberOrgNetworkIDs types.Set    `tfsdk:"member_org_network_ids"`
}

type securityGroupModelMemberOrgNetworkIDs []string

// func securityGroupToNsxtFirewallGroup.
func (r *securityGroupResource) securityGroupToNsxtFirewallGroup(ctx context.Context, rm *securityGroupModel) (securityGroup *govcdtypes.NsxtFirewallGroup, diags diag.Diagnostics) {
	parentEdgeGW, err := r.edgegw.GetParent()
	if err != nil {
		diags.AddError("Unable to get parent edge gateway", fmt.Sprintf("Unable to get parent edge gateway: %s", err.Error()))
		return
	}

	var ownerID string

	if parentEdgeGW.IsVDCGroup() {
		ownerID = parentEdgeGW.GetID()
	} else {
		ownerID = r.edgegw.GetID()
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

// MemberOrgNetworkIDsFromPlan returns the member_org_network_ids from the plan.
func (rm *securityGroupModel) MemberOrgNetworkIDsFromPlan(ctx context.Context) (securityGroupModelMemberOrgNetworkIDs, diag.Diagnostics) {
	ids := securityGroupModelMemberOrgNetworkIDs{}
	return ids, rm.MemberOrgNetworkIDs.ElementsAs(ctx, &ids, false)
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

	secGroup, err := state.GetSecurityGroup(ctx, r)
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
	secGroup, err := state.GetSecurityGroup(ctx, r)
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
	secGroup, err := state.GetSecurityGroup(ctx, r)
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
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (rm *securityGroupModel) GetSecurityGroup(ctx context.Context, r *securityGroupResource) (*govcd.NsxtFirewallGroup, error) {
	parentEdgeGW, err := r.edgegw.GetParent()
	if err != nil {
		return nil, err
	}

	if parentEdgeGW.IsVDCGroup() {
		return parentEdgeGW.GetNsxtFirewallGroupByID(rm.ID.ValueString())
	}

	return r.edgegw.GetNsxtFirewallGroupById(rm.ID.ValueString())
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
