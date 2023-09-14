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

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/edgegw"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/mutex"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/org"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/uuid"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &ipSetResource{}
	_ resource.ResourceWithConfigure   = &ipSetResource{}
	_ resource.ResourceWithImportState = &ipSetResource{}
	// _ resource.ResourceWithModifyPlan     = &ipSetResource{}
	// _ resource.ResourceWithUpgradeState   = &ipSetResource{}
	// _ resource.ResourceWithValidateConfig = &ipSetResource{}.
)

// NewIpSetResource is a helper function to simplify the provider implementation.
func NewIPSetResource() resource.Resource {
	return &ipSetResource{}
}

// ipSetResource is the resource implementation.
type ipSetResource struct {
	client *client.CloudAvenue
	org    org.Org
	edgegw edgegw.EdgeGateway
}

// Init Initializes the resource.
func (r *ipSetResource) Init(ctx context.Context, rm *IPSetModel) (diags diag.Diagnostics) {
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
func (r *ipSetResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_ip_set"
}

// Schema defines the schema for the resource.
func (r *ipSetResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ipSetSchema(ctx).GetResource(ctx)
}

func (r *ipSetResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	vdcOrVDCGroup, err := r.edgegw.GetParent()
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving Edge Gateway parent", err.Error())
		return
	}

	var createdIPSet *govcd.NsxtFirewallGroup

	if vdcOrVDCGroup.IsVDCGroup() {
		mutex.GlobalMutex.KvLock(ctx, vdcOrVDCGroup.GetID())
		defer mutex.GlobalMutex.KvUnlock(ctx, vdcOrVDCGroup.GetID())
		ipSetConfig, d := plan.ToNsxtFirewallGroup(ctx, vdcOrVDCGroup.GetID())
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
		createdIPSet, err = vdcOrVDCGroup.SetIPSet(ipSetConfig)
	} else {
		mutex.GlobalMutex.KvLock(ctx, r.edgegw.GetID())
		defer mutex.GlobalMutex.KvUnlock(ctx, r.edgegw.GetID())
		ipSetConfig, d := plan.ToNsxtFirewallGroup(ctx, r.edgegw.GetID())
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
		createdIPSet, err = r.edgegw.SetIPSet(ipSetConfig)
	}
	if err != nil {
		resp.Diagnostics.AddError("Error creating IP Set", err.Error())
		return
	}

	plan.ID.Set(createdIPSet.NsxtFirewallGroup.ID)
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
	var (
		plan  = &IPSetModel{}
		state = &IPSetModel{}
	)

	defer metrics.New("cloudavenue_edgegateway_ip_set", r.client.GetOrgName(), metrics.Update)()

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

	var (
		ipSet       *govcd.NsxtFirewallGroup
		ipSetConfig *govcdtypes.NsxtFirewallGroup
		d           diag.Diagnostics
	)

	if vdcOrVDCGroup.IsVDCGroup() {
		mutex.GlobalMutex.KvLock(ctx, vdcOrVDCGroup.GetID())
		defer mutex.GlobalMutex.KvUnlock(ctx, vdcOrVDCGroup.GetID())
		ipSetConfig, d = plan.ToNsxtFirewallGroup(ctx, vdcOrVDCGroup.GetID())
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
		ipSet, err = vdcOrVDCGroup.GetIPSetByID(state.ID.Get())
	} else {
		mutex.GlobalMutex.KvLock(ctx, r.edgegw.GetID())
		defer mutex.GlobalMutex.KvUnlock(ctx, r.edgegw.GetID())
		ipSetConfig, d = plan.ToNsxtFirewallGroup(ctx, r.edgegw.GetID())
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
		ipSet, err = r.edgegw.GetIPSetByID(state.ID.Get())
	}
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving IP Set", err.Error())
		return
	}

	if _, err := ipSet.Update(ipSetConfig); err != nil {
		resp.Diagnostics.AddError("Error updating IP Set", err.Error())
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

	vdcOrVDCGroup, err := r.edgegw.GetParent()
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving Edge Gateway parent", err.Error())
		return
	}

	var ipSet *govcd.NsxtFirewallGroup

	if vdcOrVDCGroup.IsVDCGroup() {
		mutex.GlobalMutex.KvLock(ctx, vdcOrVDCGroup.GetID())
		defer mutex.GlobalMutex.KvUnlock(ctx, vdcOrVDCGroup.GetID())
		ipSet, err = vdcOrVDCGroup.GetIPSetByID(state.ID.Get())
	} else {
		mutex.GlobalMutex.KvLock(ctx, r.edgegw.GetID())
		defer mutex.GlobalMutex.KvUnlock(ctx, r.edgegw.GetID())
		ipSet, err = r.edgegw.GetIPSetByID(state.ID.Get())
	}
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving IP Set", err.Error())
		return
	}

	if err := ipSet.Delete(); err != nil {
		resp.Diagnostics.AddError("Error deleting IP Set", err.Error())
		return
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
		edgegwID, edgegwName string
		d                    diag.Diagnostics
		err                  error
		ipSet                *govcd.NsxtFirewallGroup
	)

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
		resp.Diagnostics.AddError("Failed to import Security Group.", err.Error())
		return
	}

	vdcOrVDCGroup, err := r.edgegw.GetParent()
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving Edge Gateway parent", err.Error())
		return
	}

	if vdcOrVDCGroup.IsVDCGroup() {
		ipSet, err = vdcOrVDCGroup.GetIPSetByName(idParts[1])
	} else {
		ipSet, err = r.edgegw.GetIPSetByName(idParts[1])
	}
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving IP Set", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), ipSet.NsxtFirewallGroup.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), ipSet.NsxtFirewallGroup.Name)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("edge_gateway_id"), ipSet.NsxtFirewallGroup.EdgeGatewayRef.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("edge_gateway_name"), ipSet.NsxtFirewallGroup.EdgeGatewayRef.Name)...)
}

func (r *ipSetResource) read(ctx context.Context, planOrState *IPSetModel) (stateRefreshed *IPSetModel, found bool, diags diag.Diagnostics) {
	stateRefreshed = planOrState.Copy()

	vdcOrVDCGroup, err := r.edgegw.GetParent()
	if err != nil {
		diags.AddError("Error retrieving Edge Gateway parent", err.Error())
		return stateRefreshed, true, nil
	}

	nameOrID := stateRefreshed.ID.Get()
	if !stateRefreshed.ID.IsKnown() {
		nameOrID = stateRefreshed.Name.Get()
	}

	var ipSetConfig *govcd.NsxtFirewallGroup

	if vdcOrVDCGroup.IsVDCGroup() {
		ipSetConfig, err = vdcOrVDCGroup.GetIPSetByNameOrID(nameOrID)
	} else {
		ipSetConfig, err = r.edgegw.GetIPSetByNameOrID(nameOrID)
	}
	if err != nil {
		if govcd.IsNotFound(err) {
			return stateRefreshed, false, nil
		}
		diags.AddError("Error retrieving IP Set", err.Error())
		return
	}

	stateRefreshed.ID.Set(ipSetConfig.NsxtFirewallGroup.ID)
	stateRefreshed.Name.Set(ipSetConfig.NsxtFirewallGroup.Name)
	stateRefreshed.Description.Set(ipSetConfig.NsxtFirewallGroup.Description)
	stateRefreshed.EdgeGatewayID.Set(ipSetConfig.NsxtFirewallGroup.EdgeGatewayRef.ID)
	stateRefreshed.EdgeGatewayName.Set(ipSetConfig.NsxtFirewallGroup.EdgeGatewayRef.Name)
	diags.Append(stateRefreshed.IPAddresses.Set(ctx, ipSetConfig.NsxtFirewallGroup.IpAddresses)...)

	return stateRefreshed, true, diags
}
