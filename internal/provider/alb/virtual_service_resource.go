package alb

import (
	"context"
	"fmt"
	"strings"

	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"

	v1 "github.com/orange-cloudavenue/cloudavenue-sdk-go/v1"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/edgegw"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/mutex"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/org"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/uuid"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &VirtualServiceResource{}
	_ resource.ResourceWithConfigure   = &VirtualServiceResource{}
	_ resource.ResourceWithImportState = &VirtualServiceResource{}
)

// NewVirtualServiceResource is a helper function to simplify the provider implementation.
func NewVirtualServiceResource() resource.Resource {
	return &VirtualServiceResource{}
}

// VirtualServiceResource is the resource implementation.
type VirtualServiceResource struct {
	client *client.CloudAvenue
	edgegw edgegw.EdgeGateway
	org    org.Org
}

// Init Initializes the resource.
func (r *VirtualServiceResource) Init(ctx context.Context, rm *VirtualServiceModel) (diags diag.Diagnostics) {
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
func (r *VirtualServiceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_virtual_service"
}

// Schema defines the schema for the resource.
func (r *VirtualServiceResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = virtualServiceSchema(ctx).GetResource(ctx)
}

func (r *VirtualServiceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *VirtualServiceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer metrics.New("cloudavenue_alb_virtual_service", r.client.GetOrgName(), metrics.Create)()

	plan := &VirtualServiceModel{}

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
	// Convert the plan to NSXT ALB Virtual Service
	albConfig, diags := plan.toALBVirtualService(ctx, r)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Lock object EGW or VDC Group
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

	// Create the ALB Virtual Service
	albVS, err := r.edgegw.CreateALBVirtualService(albConfig.VirtualService)
	if err != nil {
		resp.Diagnostics.AddError("Error creating ALB Virtual Service", err.Error())
		return
	}

	// Set the ID of the created resource
	plan.ID.Set(albVS.VirtualService.ID)

	// Use generic read function to refresh the state
	state, found, d := r.read(ctx, plan)

	if !found {
		resp.State.RemoveResource(ctx)
		resp.Diagnostics.AddError("ALB Virtual Service not found", "ALB Virtual Service not found after creation")
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
func (r *VirtualServiceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer metrics.New("cloudavenue_alb_virtual_service", r.client.GetOrgName(), metrics.Read)()

	state := &VirtualServiceModel{}
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

	// Lock object EGW or VDC Group
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

	// Refresh the state
	stateRefreshed, found, d := r.read(ctx, state)
	if !found {
		resp.State.RemoveResource(ctx)
		resp.Diagnostics.AddError("ALB Virtual Service not found", "ALB Virtual Service not found after refresh")
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
func (r *VirtualServiceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer metrics.New("cloudavenue_alb_virtual_service", r.client.GetOrgName(), metrics.Update)()

	var (
		plan  = &VirtualServiceModel{}
		state = &VirtualServiceModel{}
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
		albVS *v1.EdgeGatewayALBVirtualService
		err   error
	)

	// Lock object EGW or VDC Group
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

	// Get the current ALB Virtual Service from API
	if state.ID.IsKnown() {
		albVS, err = r.edgegw.GetALBVirtualService(state.ID.Get())
	} else {
		albVS, err = r.edgegw.GetALBVirtualService(state.Name.Get())
	}
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving ALB Virtual Service", err.Error())
		return
	}

	// Get the ALB Virtual Service from plan
	newALBConfig, diags := plan.toALBVirtualService(ctx, r)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	newALBConfig.VirtualService.ID = albVS.VirtualService.ID

	// Update the ALB Virtual Service
	_, err = albVS.Update(newALBConfig.VirtualService)
	if err != nil {
		resp.Diagnostics.AddError("Error updating ALB Virtual Service", err.Error())
		return
	}

	// Use generic read function to refresh the state
	stateRefreshed, found, d := r.read(ctx, plan)
	if !found {
		resp.State.RemoveResource(ctx)
		resp.Diagnostics.AddError("ALB Virtual Service not found", "ALB Virtual Service not found after creation")
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
func (r *VirtualServiceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer metrics.New("cloudavenue_alb_virtual_service", r.client.GetOrgName(), metrics.Delete)()

	state := &VirtualServiceModel{}
	diags := diag.Diagnostics{}

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

	// Lock object EGW or VDC Group
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

	// Get the current ALB Virtual Service from API
	var albVS *v1.EdgeGatewayALBVirtualService
	if state.ID.IsKnown() {
		albVS, err = r.edgegw.GetALBVirtualService(state.ID.Get())
	} else {
		albVS, err = r.edgegw.GetALBVirtualService(state.Name.Get())
	}
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving ALB Virtual Service", err.Error())
		return
	}

	// Delete the ALB Virtual Service
	err = albVS.Delete()
	if err != nil {
		diags.AddError("Error deleting ALB Virtual Service", err.Error())
		return
	}
}

func (r *VirtualServiceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	defer metrics.New("cloudavenue_alb_virtual_service", r.client.GetOrgName(), metrics.Import)()

	// * Import with custom logic
	idParts := strings.Split(req.ID, ".")

	if len(idParts) != 2 {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: edge_gateway_NameOrID.alb_virtual_service_Name Got: %q", req.ID),
		)
		return
	}

	// Get EdgeGW is ID or Name
	var edgegwID, edgegwName string
	if uuid.IsEdgeGateway(idParts[0]) {
		edgegwID = idParts[0]
	} else {
		edgegwName = idParts[0]
	}

	x := &VirtualServiceModel{
		Name:            supertypes.NewStringNull(),
		EdgeGatewayID:   supertypes.NewStringNull(),
		EdgeGatewayName: supertypes.NewStringNull(),
	}

	x.Name.Set(idParts[1])
	x.EdgeGatewayID.Set(edgegwID)
	x.EdgeGatewayName.Set(edgegwName)

	// Init the resource
	resp.Diagnostics.Append(r.Init(ctx, x)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the ALB Virtual Service from API
	stateRefreshed, found, d := r.read(ctx, x)
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

// * CustomFuncs

// read is a generic read function that can be used by the resource Create, Read and Update functions.
func (r *VirtualServiceResource) read(ctx context.Context, planOrState *VirtualServiceModel) (stateRefreshed *VirtualServiceModel, found bool, diags diag.Diagnostics) {
	stateRefreshed = planOrState.Copy()
	var albVS *v1.EdgeGatewayALBVirtualService
	var err error

	// Get the current ALB Virtual Service from API
	if stateRefreshed.ID.IsKnown() {
		albVS, err = r.edgegw.GetALBVirtualService(stateRefreshed.ID.Get())
	} else {
		albVS, err = r.edgegw.GetALBVirtualService(stateRefreshed.Name.Get())
	}
	if err != nil {
		if govcd.IsNotFound(err) {
			return stateRefreshed, false, nil
		}
		diags.AddError("Error retrieving ALB Virtual Service", err.Error())
		return stateRefreshed, true, diags
	}

	// Populate the state with the data from the API
	stateRefreshed.ID.Set(albVS.VirtualService.ID)
	stateRefreshed.Name.Set(albVS.VirtualService.Name)
	stateRefreshed.Description.Set(albVS.VirtualService.Description)
	stateRefreshed.Enabled.SetPtr(albVS.VirtualService.Enabled)
	stateRefreshed.VirtualIP.Set(albVS.VirtualService.VirtualIPAddress)
	stateRefreshed.EdgeGatewayID.Set(r.edgegw.GetID())
	stateRefreshed.EdgeGatewayName.Set(r.edgegw.GetName())
	stateRefreshed.PoolID.Set(albVS.VirtualService.LoadBalancerPoolRef.ID)
	stateRefreshed.PoolName.Set(albVS.VirtualService.LoadBalancerPoolRef.Name)
	stateRefreshed.ServiceEngineGroupName.Set(albVS.VirtualService.ServiceEngineGroupRef.Name)
	stateRefreshed.ServiceType.Set(albVS.VirtualService.ApplicationProfile.Type)
	if albVS.VirtualService.CertificateRef != nil {
		stateRefreshed.CertificateID.Set(albVS.VirtualService.CertificateRef.ID)
	}

	// Populate Service Ports
	x := make([]*VirtualServiceModelServicePort, 0)
	for _, svcPort := range albVS.VirtualService.ServicePorts {
		y := &VirtualServiceModelServicePort{
			PortStart: supertypes.NewInt64Null(),
			PortEnd:   supertypes.NewInt64Null(),
			PortSSL:   supertypes.NewBoolNull(),
			PortType:  supertypes.NewStringNull(),
		}
		y.PortStart.SetIntPtr(svcPort.PortStart)
		y.PortEnd.SetIntPtr(svcPort.PortEnd)
		y.PortSSL.SetPtr(svcPort.SslEnabled)
		y.PortType.Set(svcPort.TcpUdpProfile.Type)

		x = append(x, y)
	}
	d := stateRefreshed.ServicePorts.Set(ctx, x)
	if d.HasError() {
		diags.Append(d...)
		return stateRefreshed, true, diags
	}

	return stateRefreshed, true, nil
}
