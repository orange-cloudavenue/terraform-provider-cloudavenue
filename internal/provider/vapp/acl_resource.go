// Package vapp provides a Terraform resource.
package vapp

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/acl"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vapp"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vdc"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &aclResource{}
	_ resource.ResourceWithConfigure   = &aclResource{}
	_ resource.ResourceWithImportState = &aclResource{}
)

// NewaclResource is a helper function to simplify the provider implementation.
func NewACLResource() resource.Resource {
	return &aclResource{}
}

// aclResource is the resource implementation.
type aclResource struct {
	client *client.CloudAvenue
	vdc    vdc.VDC
	vapp   vapp.VApp
}

type aclResourceModel struct {
	ID                  types.String `tfsdk:"id"`
	VDC                 types.String `tfsdk:"vdc"`
	VAppID              types.String `tfsdk:"vapp_id"`
	VAppName            types.String `tfsdk:"vapp_name"`
	EveryoneAccessLevel types.String `tfsdk:"everyone_access_level"`
	SharedWith          types.Set    `tfsdk:"shared_with"`
}

// Metadata returns the resource type name.
func (r *aclResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_" + "acl"
}

// Schema defines the schema for the resource.
func (r *aclResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Provides a Cloud Avenue Access Control structure for a vApp. This can be used to create, update, and delete access control structures for a vApp.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the resource",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"vdc":                   vdc.Schema(),
			"vapp_id":               vapp.Schema()["vapp_id"],
			"vapp_name":             vapp.Schema()["vapp_name"],
			"everyone_access_level": acl.Schema(false)["everyone_access_level"],
			"shared_with":           acl.Schema(false)["shared_with"],
		},
	}
}

func (r *aclResource) Init(ctx context.Context, rm *aclResourceModel) (diags diag.Diagnostics) {
	r.vdc, diags = vdc.Init(r.client, rm.VDC)
	if diags.HasError() {
		return
	}

	r.vapp, diags = vapp.Init(r.client, r.vdc, rm.VAppID, rm.VAppName)
	if diags.HasError() {
		return
	}

	return
}

func (r *aclResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *aclResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var (
		plan *aclResourceModel
	)

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Init resource
	resp.Diagnostics.Append(r.Init(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create resource
	plan, diags := r.createOrUpdateACL(ctx, plan)
	if diags.HasError() {
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *aclResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *aclResourceModel

	// Get current state
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Init resource
	resp.Diagnostics.Append(r.Init(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Request acl
	accessControl, err := r.vapp.GetAccessControl(false)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving access control properties", err.Error())
		return
	}

	// SharedWithEveryone
	everyoneAccessLevel := ""
	if accessControl.EveryoneAccessLevel != nil {
		everyoneAccessLevel = *accessControl.EveryoneAccessLevel
	}

	plan := &aclResourceModel{
		ID:                  types.StringValue(r.vapp.GetID()),
		VDC:                 types.StringValue(r.vdc.GetName()),
		VAppID:              state.VAppID,
		VAppName:            state.VAppName,
		EveryoneAccessLevel: types.StringValue(everyoneAccessLevel),
	}

	if plan.EveryoneAccessLevel.ValueString() == "" {
		plan.EveryoneAccessLevel = types.StringNull()
	}

	plan.SharedWith = types.SetNull(types.ObjectType{AttrTypes: acl.SharedWithModelAttrTypes})
	if accessControl.AccessSettings != nil {
		accessControlListSet, err := acl.AccessControlListToSharedSet(accessControl.AccessSettings.AccessSetting)
		if err != nil {
			resp.Diagnostics.AddError("Error converting slice AccessSetting into set", err.Error())
			return
		}
		var diags diag.Diagnostics
		plan.SharedWith, diags = types.SetValueFrom(ctx, types.ObjectType{AttrTypes: acl.SharedWithModelAttrTypes}, accessControlListSet)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	// Remove resource if all rights are null
	if plan.EveryoneAccessLevel.IsNull() && plan.SharedWith.IsNull() {
		resp.State.RemoveResource(ctx)
		return
	}

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *aclResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan *aclResourceModel

	// Get current state
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Init resource
	resp.Diagnostics.Append(r.Init(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update resource
	plan, diags := r.createOrUpdateACL(ctx, plan)
	if diags.HasError() {
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *aclResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *aclResourceModel

	// Get current state
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Init resource
	resp.Diagnostics.Append(r.Init(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete vApp access control
	if err := r.vapp.RemoveAccessControl(false); err != nil {
		resp.Diagnostics.AddError("Error removing vApp", err.Error())
		return
	}
}

func (r *aclResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ".")

	if len(idParts) != 1 && len(idParts) != 2 {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: vdc.vapp_name vapp_name. Got: %q", req.ID),
		)
		return
	}

	if len(idParts) == 1 {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("vapp_name"), idParts[0])...)
	} else {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("vdc"), idParts[0])...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("vapp_name"), idParts[1])...)
	}
}

func (r *aclResource) createOrUpdateACL(ctx context.Context, plan *aclResourceModel) (*aclResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	var accessControl govcdtypes.ControlAccessParams

	everyoneAccessLevel := plan.EveryoneAccessLevel.ValueString()

	sharedList := []acl.SharedWithModel{}
	diags.Append(plan.SharedWith.ElementsAs(ctx, &sharedList, true)...)
	if diags.HasError() {
		return nil, diags
	}

	// Lock vApp
	diags.Append(r.vapp.LockParentVApp(ctx)...)
	if diags.HasError() {
		return nil, diags
	}
	defer diags.Append(r.vapp.UnlockParentVApp(ctx)...)

	var accessSettings []*govcdtypes.AccessSetting

	isSharedWithEveryone := !(plan.EveryoneAccessLevel.IsNull() || plan.EveryoneAccessLevel.IsUnknown())
	if isSharedWithEveryone {
		accessControl.IsSharedToEveryone = true
		accessControl.EveryoneAccessLevel = &everyoneAccessLevel
	} else {
		var sharedListOutput []*acl.SharedWithModel

		// Get admin Org
		adminOrg, err := r.client.Vmware.GetAdminOrgByNameOrId(r.client.GetOrg())
		if err != nil {
			diags.AddError("Error retrieving Org", err.Error())
			return nil, diags
		}

		accessSettings, sharedListOutput, err = acl.SharedSetToAccessControl(r.client.Vmware, adminOrg, sharedList)
		if err != nil {
			diags.AddError("Error when reading shared_with from schema.", err.Error())
			return nil, diags
		}

		plan.SharedWith = types.SetNull(types.ObjectType{AttrTypes: acl.SharedWithModelAttrTypes})
		if len(sharedListOutput) != 0 {
			accessControl.AccessSettings = &govcdtypes.AccessSettingList{
				AccessSetting: accessSettings,
			}
			plan.SharedWith, diags = types.SetValueFrom(ctx, types.ObjectType{AttrTypes: acl.SharedWithModelAttrTypes}, sharedListOutput)

			diags.Append(diags...)
			if diags.HasError() {
				return nil, diags
			}
		}
	}

	err := r.vapp.SetAccessControl(&accessControl, false)
	if err != nil {
		diags.AddError("Error setting access control for vApp", err.Error())
		return nil, diags
	}

	plan = &aclResourceModel{
		ID:                  types.StringValue(r.vapp.GetID()),
		VDC:                 types.StringValue(r.vdc.GetName()),
		VAppID:              plan.VAppID,
		VAppName:            plan.VAppName,
		EveryoneAccessLevel: plan.EveryoneAccessLevel,
		SharedWith:          plan.SharedWith,
	}

	return plan, diags
}
