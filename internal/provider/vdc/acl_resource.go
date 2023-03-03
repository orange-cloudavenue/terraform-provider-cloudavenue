// Package vdc provides a Terraform resource.
package vdc

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/acl"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vdc"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &aclResource{}
	_ resource.ResourceWithConfigure   = &aclResource{}
	_ resource.ResourceWithImportState = &aclResource{}
)

// NewACLResource is a helper function to simplify the provider implementation.
func NewACLResource() resource.Resource {
	return &aclResource{}
}

// aclResource is the resource implementation.
type aclResource struct {
	client *client.CloudAvenue
	vdc    vdc.VDC
}

type aclResourceModel struct {
	ID                  types.String `tfsdk:"id"`
	VDC                 types.String `tfsdk:"vdc"`
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
		MarkdownDescription: "Provides a Cloud Avenue vDC access control resource. This can be used to share vDC across users and/or groups.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID of the resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"vdc":                   vdc.Schema(),
			"everyone_access_level": acl.Schema(true)["everyone_access_level"],
			"shared_with":           acl.Schema(true)["shared_with"],
		},
	}
}

func (r *aclResource) Init(ctx context.Context, rm *aclResourceModel) (diags diag.Diagnostics) {
	r.vdc, diags = vdc.Init(r.client, rm.VDC)
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

	controlAccessParams, err := r.vdc.GetControlAccess(true)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving control access", err.Error())
		return
	}

	everyoneAccessLevel := ""
	if controlAccessParams.EveryoneAccessLevel != nil {
		everyoneAccessLevel = *controlAccessParams.EveryoneAccessLevel
	}

	plan := &aclResourceModel{
		ID:                  types.StringValue(r.vdc.GetID()),
		VDC:                 types.StringValue(r.vdc.GetName()),
		EveryoneAccessLevel: types.StringValue(everyoneAccessLevel),
	}

	if plan.EveryoneAccessLevel.ValueString() == "" {
		plan.EveryoneAccessLevel = types.StringNull()
	}

	plan.SharedWith = types.SetNull(types.ObjectType{AttrTypes: acl.SharedWithModelAttrTypes})
	if controlAccessParams.AccessSettings != nil {
		accessControlListSet, err := accessControlListToSharedSet(controlAccessParams.AccessSettings.AccessSetting)
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

	if _, err := r.vdc.DeleteControlAccess(true); err != nil {
		resp.Diagnostics.AddError("Error deleting control access", err.Error())
		return
	}
}

//go:generate go run github.com/FrangipaneTeam/tf-doc-extractor@latest -filename $GOFILE -example-dir ../../../examples -resource
func (r *aclResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("vdc"), req, resp)
}

func sharedSetToAccessControl(client *govcd.VCDClient, org *govcd.AdminOrg, input []acl.SharedWithModel) ([]*govcdtypes.AccessSetting, []*acl.SharedWithModel, error) {
	var output []*govcdtypes.AccessSetting
	var outputModel []*acl.SharedWithModel

	for _, item := range input {
		var subjectHref string
		var subjectType string
		var subjectName string
		var oModel *acl.SharedWithModel

		if !item.UserID.IsNull() && !item.UserID.IsUnknown() {
			userID := item.UserID.ValueString()
			user, err := org.GetUserById(userID, false)
			if err != nil {
				return nil, nil, fmt.Errorf("error retrieving user %s: %s", userID, err)
			}
			subjectHref = user.User.Href
			subjectType = user.User.Type
			subjectName = user.User.Name

			oModel = &acl.SharedWithModel{
				UserID:      types.StringValue("urn:vcloud:user:" + common.ExtractUUID(subjectHref)),
				SubjectName: types.StringValue(subjectName),
			}
		} else if !item.GroupID.IsNull() && !item.GroupID.IsUnknown() {
			groupID := item.GroupID.ValueString()
			group, err := org.GetGroupById(groupID, false)
			if err != nil {
				return nil, nil, fmt.Errorf("error retrieving group %s: %s", groupID, err)
			}
			subjectHref = group.Group.Href
			subjectType = group.Group.Type
			subjectName = group.Group.Name
			oModel = &acl.SharedWithModel{
				GroupID:     types.StringValue("urn:vcloud:group:" + common.ExtractUUID(subjectHref)),
				SubjectName: types.StringValue(subjectName),
			}
		}

		accessLevel := item.AccessLevel.ValueString()

		output = append(output, &govcdtypes.AccessSetting{
			Subject: &govcdtypes.LocalSubject{
				HREF: subjectHref,
				Name: subjectName,
				Type: subjectType,
			},
			ExternalSubject: nil,
			AccessLevel:     accessLevel,
		})
		oModel.AccessLevel = types.StringValue(accessLevel)
		outputModel = append(outputModel, oModel)
	}
	return output, outputModel, nil
}

func accessControlListToSharedSet(input []*govcdtypes.AccessSetting) ([]acl.SharedWithModel, error) {
	var output []acl.SharedWithModel

	for _, item := range input {
		o := acl.SharedWithModel{}

		switch item.Subject.Type {
		case govcdtypes.MimeAdminUser:
			o.UserID = types.StringValue("urn:vcloud:user:" + common.ExtractUUID(item.Subject.HREF))
		case govcdtypes.MimeAdminGroup:
			o.GroupID = types.StringValue("urn:vcloud:group:" + common.ExtractUUID(item.Subject.HREF))
		default:
			return nil, fmt.Errorf("unhandled type '%s' for item %s", item.Subject.Type, item.Subject.Name)
		}
		o.AccessLevel = types.StringValue(item.AccessLevel)
		o.SubjectName = types.StringValue(item.Subject.Name)

		output = append(output, o)
	}
	return output, nil
}

func (r *aclResource) createOrUpdateACL(ctx context.Context, plan *aclResourceModel) (*aclResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	everyoneAccessLevel := plan.EveryoneAccessLevel.ValueString()

	sharedList := []acl.SharedWithModel{}
	diags = plan.SharedWith.ElementsAs(ctx, &sharedList, true)
	if diags.HasError() {
		return nil, diags
	}

	adminOrg, err := r.client.Vmware.GetAdminOrgByNameOrId(r.client.GetOrg())
	if err != nil {
		diags.AddError("Error retrieving Org", err.Error())
		return nil, diags
	}

	var accessSettings []*govcdtypes.AccessSetting
	sharedListOutput := []*acl.SharedWithModel{}

	isSharedWithEveryone := !(plan.EveryoneAccessLevel.IsNull() || plan.EveryoneAccessLevel.IsUnknown())
	if !isSharedWithEveryone {
		everyoneAccessLevel = ""
		accessSettings, sharedListOutput, err = sharedSetToAccessControl(r.client.Vmware, adminOrg, sharedList)
		if err != nil {
			diags.AddError("Error when reading shared_with from schema.", err.Error())
			return nil, diags
		}
	}

	_, err = r.vdc.SetControlAccess(isSharedWithEveryone, everyoneAccessLevel, accessSettings, true)
	if err != nil {
		diags.AddError("Error setting control access", err.Error())
		return nil, diags
	}

	plan = &aclResourceModel{
		ID:                  types.StringValue(r.vdc.GetID()),
		VDC:                 types.StringValue(r.vdc.GetName()),
		EveryoneAccessLevel: plan.EveryoneAccessLevel,
		SharedWith:          plan.SharedWith,
	}

	plan.SharedWith = types.SetNull(types.ObjectType{AttrTypes: acl.SharedWithModelAttrTypes})
	if len(sharedListOutput) != 0 {
		plan.SharedWith, diags = types.SetValueFrom(ctx, types.ObjectType{AttrTypes: acl.SharedWithModelAttrTypes}, sharedListOutput)

		diags.Append(diags...)
		if diags.HasError() {
			return nil, diags
		}
	}
	return plan, diags
}
