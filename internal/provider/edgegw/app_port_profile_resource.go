// Package network provides a Terraform resource.
package edgegw

import (
	"context"
	"fmt"
	"strings"

	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"

	superschema "github.com/FrangipaneTeam/terraform-plugin-framework-superschema"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/org"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/uuid"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &portProfilesResource{}
	_ resource.ResourceWithConfigure   = &portProfilesResource{}
	_ resource.ResourceWithImportState = &portProfilesResource{}
)

const (
	appPortProfileScope = "TENANT"
)

// NewPortProfilesResource is a helper function to simplify the provider implementation.
func NewPortProfilesResource() resource.Resource {
	return &portProfilesResource{}
}

// portProfilesResource is the resource implementation.
type portProfilesResource struct {
	client *client.CloudAvenue
	org    org.Org
}

type portProfilesResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	VDC         types.String `tfsdk:"vdc"`
	Description types.String `tfsdk:"description"`
	AppPorts    types.List   `tfsdk:"app_ports"`
}

type portProfilesResourceModelAppPorts []portProfilesResourceModelAppPort

type portProfilesResourceModelAppPort struct {
	Protocol types.String `tfsdk:"protocol"`
	Ports    types.Set    `tfsdk:"ports"`
}

func (rm *portProfilesResourceModel) AppPortsFromPlan(ctx context.Context) (appPorts portProfilesResourceModelAppPorts, diags diag.Diagnostics) {
	appPorts = make([]portProfilesResourceModelAppPort, 0)
	diags.Append(rm.AppPorts.ElementsAs(ctx, &appPorts, false)...)
	if diags.HasError() {
		return
	}

	return appPorts, diags
}

// * AppPort

func (rm *portProfilesResourceModelAppPort) PortsFromPlan(ctx context.Context) (ports []types.String, diags diag.Diagnostics) {
	if rm.Ports.IsNull() || rm.Ports.IsUnknown() {
		return
	}

	ports = make([]types.String, 0)
	diags.Append(rm.Ports.ElementsAs(ctx, &ports, false)...)
	if diags.HasError() {
		return
	}

	return ports, diags
}

// * AppPorts

// objectType.
func (p *portProfilesResourceModelAppPorts) ObjectType(ctx context.Context) types.ObjectType {
	return types.ObjectType{
		AttrTypes: p.AttrTypes(ctx),
	}
}

// attrTypes().
func (p *portProfilesResourceModelAppPorts) AttrTypes(_ context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"protocol": types.StringType,
		"ports":    types.SetType{ElemType: types.StringType},
	}
}

func (p *portProfilesResourceModelAppPorts) ToPlan(ctx context.Context) (basetypes.ListValue, diag.Diagnostics) {
	if p == nil {
		return types.ListNull(p.ObjectType(ctx)), nil
	}

	return types.ListValueFrom(ctx, p.ObjectType(ctx), p)
}

func (r *portProfilesResource) AppPortRead(ctx context.Context, portProfile *govcd.NsxtAppPortProfile) (appPorts portProfilesResourceModelAppPorts, diags diag.Diagnostics) {
	appPorts = make([]portProfilesResourceModelAppPort, 0)
	for _, appPort := range portProfile.NsxtAppPortProfile.ApplicationPorts {
		x := portProfilesResourceModelAppPort{
			Protocol: types.StringValue(appPort.Protocol),
			Ports:    types.SetNull(types.StringType),
		}

		p := make([]attr.Value, 0)
		for _, port := range appPort.DestinationPorts {
			p = append(p, types.StringValue(port))
		}

		ports, d := types.SetValue(types.StringType, p)
		if !d.HasError() {
			x.Ports = ports
		}

		appPorts = append(appPorts, x)
	}

	return
}

// Init Initializes the resource.
func (r *portProfilesResource) Init(ctx context.Context, rm *portProfilesResourceModel) (diags diag.Diagnostics) {
	r.org, diags = org.Init(r.client)

	return
}

// Metadata returns the resource type name.
func (r *portProfilesResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_app_port_profile"
}

// Schema defines the schema for the resource.
func (r *portProfilesResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = r.superSchema(ctx).GetResource(ctx)
}

func (r *portProfilesResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *portProfilesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &portProfilesResourceModel{}

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

	appPorts, d := plan.AppPortsFromPlan(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	/*
		Implement the resource creation logic here.
	*/

	appPortProfileConfig := &govcdtypes.NsxtAppPortProfile{
		Name:            plan.Name.ValueString(),
		Description:     plan.Description.ValueString(),
		Scope:           appPortProfileScope,
		ContextEntityId: plan.VDC.ValueString(),
		OrgRef:          &govcdtypes.OpenApiReference{ID: r.org.GetID()},
	}

	applicationPorts := make([]govcdtypes.NsxtAppPortProfilePort, len(appPorts))
	for index, singlePort := range appPorts {
		ports, d := singlePort.PortsFromPlan(ctx)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}

		applicationPorts[index] = govcdtypes.NsxtAppPortProfilePort{
			Protocol:         singlePort.Protocol.ValueString(),
			DestinationPorts: utils.SliceTypesStringToSliceString(ports),
		}
	}
	appPortProfileConfig.ApplicationPorts = applicationPorts

	createdAppPortProfile, err := r.org.CreateNsxtAppPortProfile(appPortProfileConfig)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating NSX-T App Port Profile",
			fmt.Sprintf("Error creating NSX-T App Port Profile: %s", err),
		)
		return
	}

	appPortsState, d := r.AppPortRead(ctx, createdAppPortProfile)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	state := *plan
	state.ID = types.StringValue(createdAppPortProfile.NsxtAppPortProfile.ID)
	state.AppPorts, d = appPortsState.ToPlan(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *portProfilesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &portProfilesResourceModel{}

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

	portProfiles, err := r.org.GetNsxtAppPortProfileByName(state.Name.ValueString(), appPortProfileScope)
	if err != nil {
		if govcd.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading NSX-T App Port Profile",
			fmt.Sprintf("Error reading NSX-T App Port Profile: %s", err),
		)
		return
	}

	appPortsState, d := r.AppPortRead(ctx, portProfiles)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan := &portProfilesResourceModel{
		ID:          types.StringValue(portProfiles.NsxtAppPortProfile.ID),
		Name:        types.StringValue(portProfiles.NsxtAppPortProfile.Name),
		Description: utils.StringValueOrNull(portProfiles.NsxtAppPortProfile.Description),
		VDC:         state.VDC,
	}
	plan.AppPorts, d = appPortsState.ToPlan(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *portProfilesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var (
		plan  = &portProfilesResourceModel{}
		state = &portProfilesResourceModel{}
	)

	// Get current state
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

	appPorts, d := plan.AppPortsFromPlan(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	/*
		Implement the resource update here
	*/

	portProfiles, err := r.org.GetNsxtAppPortProfileByName(state.Name.ValueString(), appPortProfileScope)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading NSX-T App Port Profile",
			fmt.Sprintf("Error reading NSX-T App Port Profile: %s", err),
		)
		return
	}

	newPortProfiles := portProfiles

	newPortProfiles.NsxtAppPortProfile.Name = plan.Name.ValueString()
	newPortProfiles.NsxtAppPortProfile.Description = plan.Description.ValueString()
	newPortProfiles.NsxtAppPortProfile.ContextEntityId = plan.VDC.ValueString()

	applicationPorts := make([]govcdtypes.NsxtAppPortProfilePort, len(appPorts))
	for index, singlePort := range appPorts {
		ports, d := singlePort.PortsFromPlan(ctx)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}

		applicationPorts[index] = govcdtypes.NsxtAppPortProfilePort{
			Protocol:         singlePort.Protocol.ValueString(),
			DestinationPorts: utils.SliceTypesStringToSliceString(ports),
		}
	}
	newPortProfiles.NsxtAppPortProfile.ApplicationPorts = applicationPorts

	newUpdated, err := portProfiles.Update(newPortProfiles.NsxtAppPortProfile)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating NSX-T App Port Profile",
			fmt.Sprintf("Error updating NSX-T App Port Profile: %s", err),
		)
		return
	}

	appPortsStateUpdated, d := r.AppPortRead(ctx, newUpdated)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan = &portProfilesResourceModel{
		ID:          types.StringValue(newUpdated.NsxtAppPortProfile.ID),
		Name:        types.StringValue(newUpdated.NsxtAppPortProfile.Name),
		Description: utils.StringValueOrNull(newUpdated.NsxtAppPortProfile.Description),
		VDC:         state.VDC,
	}
	plan.AppPorts, d = appPortsStateUpdated.ToPlan(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *portProfilesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	state := &portProfilesResourceModel{}

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

	portProfiles, err := r.org.GetNsxtAppPortProfileByName(state.Name.ValueString(), appPortProfileScope)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading NSX-T App Port Profile",
			fmt.Sprintf("Error reading NSX-T App Port Profile: %s", err),
		)
		return
	}

	if err = portProfiles.Delete(); err != nil {
		resp.Diagnostics.AddError(
			"Error deleting NSX-T App Port Profile",
			fmt.Sprintf("Error deleting NSX-T App Port Profile: %s", err),
		)
		return
	}
}

func (r *portProfilesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ".")

	if len(idParts) != 2 {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: VDCOrVDCGroupID.appPortProfileName. Got: %q", req.ID),
		)
		return
	}

	vdcID := uuid.Normalize(uuid.VDC, idParts[0]).String()

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("vdc"), vdcID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), idParts[1])...)
}

func (r *portProfilesResource) superSchema(_ context.Context) superschema.Schema {
	return superschema.Schema{
		Resource: superschema.SchemaDetails{
			MarkdownDescription: "Provides a NSX-T App Port Profile resource",
		},
		Attributes: map[string]superschema.Attribute{
			"id": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The ID of the VM.",
					Computed:            true,
				},
				Resource: &schemaR.StringAttribute{
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
			},
			"name": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "Application Port Profile name.",
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
				},
			},
			"description": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "Application Port Profile description.",
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
			},
			"vdc": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "ID of VDC or VDC Group",
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
						stringplanmodifier.UseStateForUnknown(),
					},
				},
			},
			"app_ports": superschema.ListNestedAttribute{
				Common: &schemaR.ListNestedAttribute{
					MarkdownDescription: "List of application ports.",
					Required:            true,
				},
				Attributes: map[string]superschema.Attribute{
					"ports": superschema.SetAttribute{
						Common: &schemaR.SetAttribute{
							MarkdownDescription: "Set of ports or ranges.",
							Computed:            true,
						},
						Resource: &schemaR.SetAttribute{
							Optional:    true,
							ElementType: types.StringType,
							PlanModifiers: []planmodifier.Set{
								setplanmodifier.UseStateForUnknown(),
							},
						},
					},
					"protocol": superschema.StringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "Protocol.",
						},
						Resource: &schemaR.StringAttribute{
							Required: true,
							Validators: []validator.String{
								stringvalidator.OneOf("ICMPv4", "ICMPv6", "TCP", "UDP"),
							},
						},
					},
				},
			},
		},
	}
}
