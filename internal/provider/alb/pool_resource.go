// Package alb provides a Terraform resource.
package alb

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/edgegw"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/org"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &poolResource{}
	_ resource.ResourceWithConfigure   = &poolResource{}
	_ resource.ResourceWithImportState = &poolResource{}
)

// NewALBPoolResource is a helper function to simplify the provider implementation.
func NewPoolResource() resource.Resource {
	return &poolResource{}
}

// albPoolResource is the resource implementation.
type poolResource struct {
	client  *client.CloudAvenue
	org     org.Org
	edgegw  edgegw.BaseEdgeGW
	albPool base
}

// Metadata returns the resource type name.
func (r *poolResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_pool"
}

// Schema defines the schema for the resource.
func (r *poolResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = albPoolSchema().GetResource(ctx)
}

func (r *poolResource) Init(ctx context.Context, rm *albPoolModel) (diags diag.Diagnostics) {
	r.albPool = base{
		name: rm.Name.ValueString(),
		id:   rm.ID.ValueString(),
	}

	r.edgegw = edgegw.BaseEdgeGW{
		ID:   rm.EdgeGatewayID,
		Name: rm.EdgeGatewayName,
	}

	r.org, diags = org.Init(r.client)
	return
}

func (r *poolResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *poolResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer metrics.New("cloudavenue_alb_pool", r.client.GetOrgName(), metrics.Create)()

	// Retrieve values from plan
	var (
		plan *albPoolModel
	)

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Init
	resp.Diagnostics.Append(r.Init(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Prepare config.
	albPoolConfig, err := r.getALBPoolConfig(ctx, plan)
	if err != nil {
		resp.Diagnostics.AddError("Unable to prepare ALB Pool Config", err.Error())
		return
	}

	// Lock EdgeGW
	edgeGW, err := r.org.GetEdgeGateway(r.edgegw)
	if err != nil {
		resp.Diagnostics.AddError("Unable to get Edge Gateway", err.Error())
		return
	}
	edgeGW.Lock(ctx)
	defer edgeGW.Unlock(ctx)

	// Create ALB Pool
	createdALBPool, err := r.client.Vmware.CreateNsxtAlbPool(albPoolConfig)
	if err != nil {
		resp.Diagnostics.AddError("Unable to create ALB Pool", err.Error())
		return
	}

	// Store ID
	plan.ID = utils.StringValueOrNull(createdALBPool.NsxtAlbPool.ID)

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *poolResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer metrics.New("cloudavenue_alb_pool", r.client.GetOrgName(), metrics.Read)()

	var (
		state *albPoolModel
		diags diag.Diagnostics
	)

	// Get current state
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(r.Init(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get albPool.
	albPool, err := r.GetALBPool()
	if err != nil {
		if govcd.ContainsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError("Unable to find ALB Pool", err.Error())
		return
	}

	// Set data
	plan := &albPoolModel{
		ID:                       utils.StringValueOrNull(albPool.NsxtAlbPool.ID),
		Name:                     state.Name,
		Description:              utils.StringValueOrNull(albPool.NsxtAlbPool.Description),
		EdgeGatewayID:            state.EdgeGatewayID,
		EdgeGatewayName:          state.EdgeGatewayName,
		Enabled:                  types.BoolValue(*albPool.NsxtAlbPool.Enabled),
		Algorithm:                utils.StringValueOrNull(albPool.NsxtAlbPool.Algorithm),
		DefaultPort:              types.Int64Value(int64(*albPool.NsxtAlbPool.DefaultPort)),
		GracefulTimeoutPeriod:    types.Int64Value(int64(*albPool.NsxtAlbPool.GracefulTimeoutPeriod)),
		PassiveMonitoringEnabled: types.BoolValue(*albPool.NsxtAlbPool.PassiveMonitoringEnabled),
		Members:                  types.SetNull(types.ObjectType{AttrTypes: memberAttrTypes}),
		HealthMonitors:           types.SetNull(types.StringType),
		PersistenceProfile:       types.ObjectNull(persistenceProfileAttrTypes),
	}

	// Set members
	if members := processMembers(albPool.NsxtAlbPool.Members); len(members) > 0 {
		plan.Members, diags = types.SetValueFrom(ctx, types.ObjectType{AttrTypes: memberAttrTypes}, members)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	// Set health monitors.
	healthMonitors := processHealthMonitors(albPool.NsxtAlbPool.HealthMonitors)

	if len(healthMonitors) > 0 {
		plan.HealthMonitors, diags = types.SetValueFrom(ctx, types.StringType, healthMonitors)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	// Set persistence profile
	p := processPersistenceProfile(albPool.NsxtAlbPool.PersistenceProfile)

	if p != (persistenceProfile{}) {
		plan.PersistenceProfile, diags = types.ObjectValueFrom(ctx, persistenceProfileAttrTypes, p)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *poolResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer metrics.New("cloudavenue_alb_pool", r.client.GetOrgName(), metrics.Update)()

	var plan *albPoolModel

	// Get current state
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Init
	resp.Diagnostics.Append(r.Init(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get albPool
	albPool, err := r.GetALBPool()
	if err != nil {
		resp.Diagnostics.AddError("Unable to find ALB Pool", err.Error())
		return
	}

	// Prepare config.
	albPoolConfig, err := r.getALBPoolConfig(ctx, plan)
	if err != nil {
		resp.Diagnostics.AddError("Unable to prepare ALB Pool Config", err.Error())
		return
	}

	// Lock EdgeGW
	edgeGW, err := r.org.GetEdgeGateway(r.edgegw)
	if err != nil {
		resp.Diagnostics.AddError("Unable to get Edge Gateway", err.Error())
		return
	}
	edgeGW.Lock(ctx)
	defer edgeGW.Unlock(ctx)

	// Update ALB Pool.
	_, err = albPool.Update(albPoolConfig)
	if err != nil {
		resp.Diagnostics.AddError("Unable to update ALB Pool", err.Error())
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *poolResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer metrics.New("cloudavenue_alb_pool", r.client.GetOrgName(), metrics.Delete)()

	var state *albPoolModel

	// Get current state
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Init
	resp.Diagnostics.Append(r.Init(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Lock EdgeGW
	edgeGW, err := r.org.GetEdgeGateway(r.edgegw)
	if err != nil {
		resp.Diagnostics.AddError("Unable to get Edge Gateway", err.Error())
		return
	}
	edgeGW.Lock(ctx)
	defer edgeGW.Unlock(ctx)

	// Get albPool
	albPool, err := r.GetALBPool()
	if err != nil {
		resp.Diagnostics.AddError("Unable to find ALB Pool", err.Error())
		return
	}

	err = albPool.Delete()
	if err != nil {
		resp.Diagnostics.AddError("Unable to delete ALB Pool", err.Error())
		return
	}
}

func (r *poolResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	defer metrics.New("cloudavenue_alb_pool", r.client.GetOrgName(), metrics.Import)()

	idParts := strings.Split(req.ID, ".")

	if len(idParts) != 2 {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: edge_gateway_name.alb_pool_name. Got: %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("edge_gateway_name"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), idParts[1])...)
}

// GetID returns the ID of the albPool.
func (r *poolResource) GetID() string {
	return r.albPool.id
}

// GetName returns the name of the albPool.
func (r *poolResource) GetName() string {
	return r.albPool.name
}

// GetALBPool returns the govcd.NsxtALBPool.
func (r *poolResource) GetALBPool() (*govcd.NsxtAlbPool, error) {
	if r.GetID() != "" {
		return r.client.Vmware.GetAlbPoolById(r.GetID())
	}

	nsxtEdge, err := r.org.GetEdgeGateway(r.edgegw)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve Edge gateway '%s'", r.edgegw.GetIDOrName())
	}
	return r.client.Vmware.GetAlbPoolByName(nsxtEdge.EdgeGateway.ID, r.GetName())
}

// getALBPoolConfig is the main function for getting *govcdtypes.NsxtAlbPool for API request. It nests multiple smaller
// functions for smaller types.
func (r *poolResource) getALBPoolConfig(ctx context.Context, d *albPoolModel) (*govcdtypes.NsxtAlbPool, error) {
	edge, err := r.org.GetEdgeGateway(r.edgegw)
	if err != nil {
		return nil, err
	}

	albPoolConfig := &govcdtypes.NsxtAlbPool{
		ID:          r.GetID(),
		Name:        d.Name.ValueString(),
		Description: d.Description.ValueString(),
		Enabled:     d.Enabled.ValueBoolPointer(),
		GatewayRef: govcdtypes.OpenApiReference{
			ID: edge.GetID(),
		},
		Algorithm:                d.Algorithm.ValueString(),
		DefaultPort:              utils.TakeIntPointer(int(d.DefaultPort.ValueInt64())),
		GracefulTimeoutPeriod:    utils.TakeIntPointer(int(d.GracefulTimeoutPeriod.ValueInt64())),
		PassiveMonitoringEnabled: d.PassiveMonitoringEnabled.ValueBoolPointer(),
	}

	poolMembers, err := r.getALBPoolMembersType(ctx, d)
	if err != nil {
		return nil, fmt.Errorf("error defining pool members: %w", err)
	}
	albPoolConfig.Members = poolMembers

	persistenceProfile, err := r.getALBPoolPersistenceProfileType(ctx, d)
	if err != nil && !errors.Is(err, ErrPersistenceProfileIsEmpty) {
		return nil, fmt.Errorf("error defining persistence profile: %w", err)
	}
	albPoolConfig.PersistenceProfile = persistenceProfile

	healthMonitors, err := r.getALBPoolHealthMonitorType(ctx, d)
	if err != nil {
		return nil, fmt.Errorf("error defining health monitors: %w", err)
	}
	albPoolConfig.HealthMonitors = healthMonitors

	return albPoolConfig, nil
}

// getALBPoolMembersType.
func (r *poolResource) getALBPoolMembersType(ctx context.Context, d *albPoolModel) ([]govcdtypes.NsxtAlbPoolMember, error) {
	var members []member
	diags := d.Members.ElementsAs(ctx, &members, true)
	if diags.HasError() {
		return nil, errors.New(diags[0].Detail())
	}
	memberSlice := make([]govcdtypes.NsxtAlbPoolMember, 0)
	for _, memberDefinition := range members {
		memberSlice = append(memberSlice, govcdtypes.NsxtAlbPoolMember{
			Enabled:   memberDefinition.Enabled.ValueBool(),
			IpAddress: memberDefinition.IPAddress.ValueString(),
			Ratio:     utils.TakeIntPointer(int(memberDefinition.Ratio.ValueInt64())),
			Port:      int(memberDefinition.Port.ValueInt64()),
		})
	}
	return memberSlice, nil
}

// getALBPoolPersistenceProfileType.
func (r *poolResource) getALBPoolPersistenceProfileType(ctx context.Context, d *albPoolModel) (*govcdtypes.NsxtAlbPoolPersistenceProfile, error) {
	if d.PersistenceProfile.IsNull() {
		return nil, ErrPersistenceProfileIsEmpty
	}

	p := &persistenceProfile{}
	if diags := d.PersistenceProfile.As(ctx, p, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	}); diags.HasError() {
		return nil, errors.New(diags[0].Detail())
	}

	return &govcdtypes.NsxtAlbPoolPersistenceProfile{
		Type:  p.Type.ValueString(),
		Value: p.Value.ValueString(),
	}, nil
}

// getALBPoolHealthMonitorType.
func (r *poolResource) getALBPoolHealthMonitorType(ctx context.Context, d *albPoolModel) (healthMonitors []govcdtypes.NsxtAlbPoolHealthMonitor, err error) {
	var healthMonitorsSlice []string

	if diags := d.HealthMonitors.ElementsAs(ctx, &healthMonitorsSlice, true); diags.HasError() {
		return nil, errors.New(diags[0].Detail())
	}

	for _, healthMonitor := range healthMonitorsSlice {
		healthMonitors = append(healthMonitors, govcdtypes.NsxtAlbPoolHealthMonitor{
			Type: healthMonitor,
		})
	}

	return
}
