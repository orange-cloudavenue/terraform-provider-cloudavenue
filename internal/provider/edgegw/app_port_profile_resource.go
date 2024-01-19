package edgegw

import (
	"context"
	"fmt"
	"strings"

	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/edgegw"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/org"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vdc"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/uuid"
)

// NewAppPortProfileResource is a helper function to simplify the provider implementation.
func NewAppPortProfileResource() resource.Resource {
	return &appPortProfileResource{}
}

// appPortProfileResource is the resource implementation.
type appPortProfileResource struct {
	client *client.CloudAvenue
	org    org.Org
}

// Init Initializes the resource.
func (r *appPortProfileResource) Init(ctx context.Context, rm *AppPortProfileModel) (diags diag.Diagnostics) {
	r.org, diags = org.Init(r.client)
	return
}

// Metadata returns the resource type name.
func (r *appPortProfileResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_app_port_profile"
}

// Schema defines the schema for the resource.
func (r *appPortProfileResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = appPortProfilesSchema(ctx).GetResource(ctx)
}

func (r *appPortProfileResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *appPortProfileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer metrics.New("cloudavenue_edgegateway_app_port_profile", r.client.GetOrgName(), metrics.Create)()

	plan := &AppPortProfileModel{}

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

	var vdcID string

	// Retrieve VDC from edge gateway (VDC attribute is deprecated)
	if !plan.EdgeGatewayID.IsKnown() && !plan.EdgeGatewayName.IsKnown() {
		// TODO - Deprecated - Remove in version v0.19.0
		// Get VDC by the deprecated attribute
		vdcData, d := vdc.Init(r.client, plan.VDC.StringValue)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}

		vdcID = vdcData.GetID()
	} else {
		// Get VDC by the edge gateway attribute
		edgegw, err := r.org.GetEdgeGateway(edgegw.BaseEdgeGW{
			ID:   types.StringValue(plan.EdgeGatewayID.Get()),
			Name: types.StringValue(plan.EdgeGatewayName.Get()),
		})
		if err != nil {
			resp.Diagnostics.AddError("Error retrieving Edge Gateway", err.Error())
			return
		}

		vdcData, err := edgegw.GetParent()
		if err != nil {
			resp.Diagnostics.AddError("Error retrieving Edge Gateway parent", err.Error())
			return
		}
		vdcID = vdcData.GetID()
	}

	appPorts, d := plan.toNsxtAppPortProfilePorts(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	appPortProfile, err := r.org.CreateNsxtAppPortProfile(&govcdtypes.NsxtAppPortProfile{
		Name:             plan.Name.Get(),
		Description:      plan.Description.Get(),
		ContextEntityId:  vdcID,
		ApplicationPorts: appPorts,
		OrgRef:           &govcdtypes.OpenApiReference{ID: r.org.GetID()},
		Scope:            govcdtypes.ApplicationPortProfileScopeTenant,
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating App Port Profile", err.Error())
		return
	}

	plan.ID.Set(appPortProfile.NsxtAppPortProfile.ID)
	state, found, d := r.read(ctx, plan)
	if !found {
		resp.State.RemoveResource(ctx)
		resp.Diagnostics.AddError("App Port Profile not found", "App Port Profile not found after creation")
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
func (r *appPortProfileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer metrics.New("cloudavenue_edgegateway_app_port_profile", r.client.GetOrgName(), metrics.Read)()

	state := &AppPortProfileModel{}

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
func (r *appPortProfileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer metrics.New("cloudavenue_edgegateway_app_port_profile", r.client.GetOrgName(), metrics.Update)()

	var (
		plan  = &AppPortProfileModel{}
		state = &AppPortProfileModel{}
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

	appPortProfile, err := r.org.GetNsxtAppPortProfileById(state.ID.Get())
	if err != nil {
		resp.Diagnostics.AddError("Error reading App Port Profile", err.Error())
		return
	}

	appPorts, d := plan.toNsxtAppPortProfilePorts(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	appPortProfile.NsxtAppPortProfile.ApplicationPorts = appPorts
	if _, err := appPortProfile.Update(appPortProfile.NsxtAppPortProfile); err != nil {
		resp.Diagnostics.AddError("Error updating App Port Profile", err.Error())
		return
	}

	stateRefreshed, _, d := r.read(ctx, state)
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, stateRefreshed)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *appPortProfileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer metrics.New("cloudavenue_edgegateway_app_port_profile", r.client.GetOrgName(), metrics.Delete)()

	state := &AppPortProfileModel{}

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

	appPortProfile, err := r.org.GetNsxtAppPortProfileById(state.ID.Get())
	if err != nil {
		if govcd.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading App Port Profile", err.Error())
		return
	}

	if err := appPortProfile.Delete(); err != nil {
		resp.Diagnostics.AddError("Error deleting App Port Profile", err.Error())
		return
	}
}

func (r *appPortProfileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	defer metrics.New("cloudavenue_edgegateway_app_port_profile", r.client.GetOrgName(), metrics.Import)()

	var d diag.Diagnostics

	r.org, d = org.Init(r.client)
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	// split req.ID into edge gateway ID and app port profile ID/name
	split := strings.Split(req.ID, ".")
	if len(split) != 2 {
		resp.Diagnostics.AddError("Invalid import ID", "Import ID must be in the format <edge_gateway_id_or_name>.<app_port_profile_id_or_name>")
		return
	}
	edgeIDOrName, appPortProfileIDOrName := split[0], split[1]

	x := &AppPortProfileModel{
		ID:              supertypes.NewStringNull(),
		Name:            supertypes.NewStringNull(),
		EdgeGatewayID:   supertypes.NewStringNull(),
		EdgeGatewayName: supertypes.NewStringNull(),
	}

	if uuid.IsEdgeGateway(edgeIDOrName) {
		x.EdgeGatewayID.Set(edgeIDOrName)
	} else {
		x.EdgeGatewayName.Set(edgeIDOrName)
	}

	if uuid.IsAppPortProfile(appPortProfileIDOrName) {
		x.ID.Set(appPortProfileIDOrName)
	} else {
		x.Name.Set(appPortProfileIDOrName)
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

func (r *appPortProfileResource) read(ctx context.Context, planOrState *AppPortProfileModel) (stateRefreshed *AppPortProfileModel, found bool, diags diag.Diagnostics) {
	stateRefreshed = planOrState.Copy()

	var (
		appPortProfile *govcd.NsxtAppPortProfile
		err            error
	)

	if planOrState.ID.IsKnown() {
		appPortProfile, err = r.org.GetNsxtAppPortProfileById(stateRefreshed.ID.Get())
	} else {
		appPortProfile, err = r.org.GetNsxtAppPortProfileByName(stateRefreshed.Name.Get(), govcdtypes.ApplicationPortProfileScopeTenant)
	}
	if err != nil {
		if govcd.IsNotFound(err) {
			return nil, false, nil
		}
		diags.AddError("Error reading App Port Profile", err.Error())
		return
	}

	appPorts := make([]*AppPortProfileModelAppPort, len(appPortProfile.NsxtAppPortProfile.ApplicationPorts))
	for index, singlePort := range appPortProfile.NsxtAppPortProfile.ApplicationPorts {
		ap := &AppPortProfileModelAppPort{
			Protocol: supertypes.NewStringNull(),
			Ports:    supertypes.NewSetValueOfNull[string](ctx),
		}

		ap.Protocol.Set(singlePort.Protocol)
		// DestinationPorts is optional
		if len(singlePort.DestinationPorts) > 0 {
			diags.Append(ap.Ports.Set(ctx, singlePort.DestinationPorts)...)
			if diags.HasError() {
				return
			}
		}
		appPorts[index] = ap
	}

	stateRefreshed.ID.Set(appPortProfile.NsxtAppPortProfile.ID)
	stateRefreshed.Name.Set(appPortProfile.NsxtAppPortProfile.Name)
	stateRefreshed.Description.Set(appPortProfile.NsxtAppPortProfile.Description)
	stateRefreshed.AppPorts.Set(ctx, appPorts)

	return stateRefreshed, true, nil
}
