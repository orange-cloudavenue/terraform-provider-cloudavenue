// Package vapp provides a Terraform resource to manage vApps.
package vapp

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/org"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vapp"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vdc"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &vappResource{}
	_ resource.ResourceWithConfigure   = &vappResource{}
	_ resource.ResourceWithImportState = &vappResource{}
)

// NewVappResource is a helper function to simplify the provider implementation.
func NewVappResource() resource.Resource {
	return &vappResource{}
}

// vappResource is the resource implementation.
type vappResource struct {
	client *client.CloudAvenue
	vdc    vdc.VDC
	org    org.Org
}

type vappResourceModel struct {
	VAppName        types.String                  `tfsdk:"name"`
	VAppID          types.String                  `tfsdk:"id"`
	VDC             types.String                  `tfsdk:"vdc"`
	Description     types.String                  `tfsdk:"description"`
	Href            types.String                  `tfsdk:"href"`
	PowerON         types.Bool                    `tfsdk:"power_on"`
	GuestProperties map[types.String]types.String `tfsdk:"guest_properties"`
	StatusCode      types.Int64                   `tfsdk:"status_code"`
	StatusText      types.String                  `tfsdk:"status_text"`
	Lease           []vappLeaseModel              `tfsdk:"lease"`
}

// Metadata returns the resource type name.
func (r *vappResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName
}

// Schema defines the schema for the resource.
func (r *vappResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The Edge Gateway resource allows you to create and manage Edge Gateways in CloudAvenue.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID of the vApp.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "(ForceNew) Name of the vApp.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"vdc": vdc.Schema(),
			"description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Optional description of the vApp",
			},
			"power_on": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "A boolean value stating if this vApp should be powered on",
				// TODO default to false
			},
			"href": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "vApp Hyper Reference",
			},
			"guest_properties": schema.MapAttribute{
				Optional:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "Key/value settings for guest properties",
			},
			"status_code": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Shows the status code of the vApp",
			},
			"status_text": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Shows the status of the vApp",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"lease": schema.ListNestedBlock{
				MarkdownDescription: "Defines lease parameters for this vApp",
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
				},
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"runtime_lease_in_sec": schema.Int64Attribute{
							Optional:            true,
							Computed:            true,
							MarkdownDescription: "How long any of the VMs in the vApp can run before the vApp is automatically powered off or suspended. 0 means never expires. Max value is 3600",
							Validators: []validator.Int64{
								int64validator.Between(0, 3600),
							},
							PlanModifiers: []planmodifier.Int64{
								int64planmodifier.UseStateForUnknown(),
							},
						},
						"storage_lease_in_sec": schema.Int64Attribute{
							Optional:            true,
							Computed:            true,
							MarkdownDescription: "How long the vApp is available before being automatically deleted or marked as expired. 0 means never expires. Max value is 3600",
							Validators: []validator.Int64{
								int64validator.Between(0, 3600),
							},
							PlanModifiers: []planmodifier.Int64{
								int64planmodifier.UseStateForUnknown(),
							},
						},
					},
				},
			},
		},
	}
}

func (r *vappResource) Init(ctx context.Context, rm *vappResourceModel) (diags diag.Diagnostics) {
	r.org, diags = org.Init(r.client)
	if diags.HasError() {
		return
	}

	r.vdc, diags = vdc.Init(r.client, rm.VDC)

	return
}

func (r *vappResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.CloudAvenue)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.CloudAvenue, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

// Create creates the resource and sets the initial Terraform state.
func (r *vappResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) { //nolint: gocyclo
	// Retrieve values from plan
	var (
		plan *vappResourceModel
		err  error
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

	vapp, err := r.vdc.CreateRawVApp(plan.VAppName.ValueString(), plan.Description.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error creating vApp", err.Error())
		return
	}

	// Wait for job to complete
	errRetry := retry.RetryContext(ctx, 90*time.Second, func() *retry.RetryError {
		currentStatus, _ := vapp.GetStatus()
		tflog.Debug(ctx, fmt.Sprintf("Creating Vapp status: %s", currentStatus))
		if currentStatus == "UNRESOLVED" {
			return retry.RetryableError(fmt.Errorf("expected vapp status != UNRESOLVED"))
		}

		return nil
	})

	if errRetry != nil {
		resp.Diagnostics.AddError("Error waiting vapp to complete", errRetry.Error())
		return
	}

	if len(plan.GuestProperties) > 0 {
		x := plan.getGuestProperties()

		_, err = vapp.SetProductSectionList(x)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Adding guest properties",
				"Could not add guest properties, unexpected error: "+err.Error(),
			)
			return
		}
	} // end if !plan.GuestProperties.IsNull()

	var runtimeLease, storageLease int

	if len(plan.Lease) > 0 {
		runtimeLease = int(plan.Lease[0].RuntimeLeaseInSec.ValueInt64())
		storageLease = int(plan.Lease[0].StorageLeaseInSec.ValueInt64())
	} else {
		adminOrg, errGetAdminOrg := r.client.Vmware.GetAdminOrgById(r.org.GetID())
		if errGetAdminOrg != nil {
			resp.Diagnostics.AddError("Error retrieving Org", errGetAdminOrg.Error())
			return
		}

		if adminOrg.AdminOrg.OrgSettings == nil || adminOrg.AdminOrg.OrgSettings.OrgVAppLeaseSettings == nil {
			resp.Diagnostics.AddError("Error retrieving Org", "Org settings are not defined")
			return
		}

		runtimeLease = *adminOrg.AdminOrg.OrgSettings.OrgVAppLeaseSettings.DeploymentLeaseSeconds
		storageLease = *adminOrg.AdminOrg.OrgSettings.OrgVAppLeaseSettings.StorageLeaseSeconds
	}

	err = vapp.RenewLease(runtimeLease, storageLease)
	if err != nil {
		resp.Diagnostics.AddError("Error renewing lease", err.Error())
		return
	}

	if !plan.Description.IsNull() {
		err = vapp.UpdateDescription(plan.Description.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Error updating VApp description", err.Error())
			return
		}
	}

	if len(plan.GuestProperties) > 0 {
		x := plan.getGuestProperties()
		_, err = vapp.SetProductSectionList(x)
		if err != nil {
			resp.Diagnostics.AddError("Error updating VApp guest properties", err.Error())
			return
		}
	}

	// power_on
	if !plan.PowerON.IsNull() {
		if plan.PowerON.ValueBool() {
			task, errPowerOn := vapp.PowerOn()
			if errPowerOn != nil {
				resp.Diagnostics.AddError("Error powering on VApp", errPowerOn.Error())
				return
			}
			err = task.WaitTaskCompletion()
			if err != nil {
				resp.Diagnostics.AddError("Error powering on VApp", err.Error())
				return
			}
		} else {
			task, errUndeploy := vapp.Undeploy()
			if errUndeploy != nil {
				resp.Diagnostics.AddError("Error powering off VApp", errUndeploy.Error())
				return
			}
			err = task.WaitTaskCompletion()
			if err != nil {
				resp.Diagnostics.AddError("Error powering off VApp", err.Error())
				return
			}
		}
	}

	// Request vApp
	vappRefreshed, err := r.vdc.GetVAppByNameOrId(vapp.VApp.ID, true)
	if err != nil {
		if errors.Is(err, govcd.ErrorEntityNotFound) {
			resp.Diagnostics.AddError("vApp not found after creating", err.Error())
			return
		}
		resp.Diagnostics.AddError("Error retrieving vApp after creating", err.Error())
		return
	}

	guestProperties, err := vappRefreshed.GetProductSectionList()
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving guest properties", err.Error())
	}
	leaseInfo, err := vappRefreshed.GetLease()
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving lease info", err.Error())
		return
	}

	var statusText string

	statusText, err = vappRefreshed.GetStatus()
	if err != nil {
		statusText = vAppUnknownStatus
	}

	nPlan := &vappResourceModel{
		VAppID:      types.StringValue(vappRefreshed.VApp.ID),
		VAppName:    types.StringValue(vappRefreshed.VApp.Name),
		Description: types.StringValue(vappRefreshed.VApp.Description),
		StatusText:  types.StringValue(statusText),
		StatusCode:  types.Int64Value(int64(vappRefreshed.VApp.Status)),
		Href:        types.StringValue(vappRefreshed.VApp.HREF),

		VDC:     types.StringValue(r.vdc.GetName()),
		PowerON: plan.PowerON,
	}

	if plan.Lease != nil && len(plan.Lease) > 0 {
		nPlan.Lease = make([]vappLeaseModel, 1)
		nPlan.Lease = append(nPlan.Lease, vappLeaseModel{
			RuntimeLeaseInSec: types.Int64Value(int64(leaseInfo.DeploymentLeaseInSeconds)),
			StorageLeaseInSec: types.Int64Value(int64(leaseInfo.StorageLeaseInSeconds)),
		})
	}

	if guestProperties != nil && guestProperties.ProductSection != nil && guestProperties.ProductSection.Property != nil {
		for _, guestProperty := range guestProperties.ProductSection.Property {
			if guestProperty.Value != nil {
				nPlan.GuestProperties[types.StringValue(guestProperty.Key)] = types.StringValue(guestProperty.Value.Value)
			}
		}
	}

	tflog.Info(ctx, fmt.Sprintf("Vapp %s created", nPlan.VAppName.ValueString()))

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, &nPlan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *vappResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *vappResourceModel

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

	vappRefreshed, diags := vapp.Init(r.client, r.vdc, state.VAppID, state.VAppName)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	guestProperties, err := vappRefreshed.GetProductSectionList()
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving guest properties", err.Error())
	}
	leaseInfo, err := vappRefreshed.GetLease()
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving lease info", err.Error())
		return
	}

	var statusText string

	statusText, err = vappRefreshed.GetStatus()
	if err != nil {
		statusText = vAppUnknownStatus
	}

	plan := &vappResourceModel{
		VAppID:      types.StringValue(vappRefreshed.GetID()),
		VAppName:    types.StringValue(vappRefreshed.GetName()),
		Description: types.StringValue(vappRefreshed.GetDescription()),
		StatusText:  types.StringValue(statusText),
		StatusCode:  types.Int64Value(int64(vappRefreshed.GetStatusCode())),
		Href:        types.StringValue(vappRefreshed.GetHREF()),
		VDC:         types.StringValue(r.vdc.GetName()),

		PowerON: state.PowerON,
	}

	if state.Lease != nil && len(state.Lease) > 0 {
		plan.Lease = make([]vappLeaseModel, 1)
		plan.Lease = append(plan.Lease, vappLeaseModel{
			RuntimeLeaseInSec: types.Int64Value(int64(leaseInfo.DeploymentLeaseInSeconds)),
			StorageLeaseInSec: types.Int64Value(int64(leaseInfo.StorageLeaseInSeconds)),
		})
	}

	if guestProperties != nil && guestProperties.ProductSection != nil && guestProperties.ProductSection.Property != nil {
		for _, guestProperty := range guestProperties.ProductSection.Property {
			if guestProperty.Value != nil {
				plan.GuestProperties[types.StringValue(guestProperty.Key)] = types.StringValue(guestProperty.Value.Value)
			}
		}
	}

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *vappResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) { //nolint: gocyclo
	var plan, state *vappResourceModel

	// Get current state
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Init resource
	resp.Diagnostics.Append(r.Init(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Request vApp
	vapp, err := r.vdc.GetVAppByNameOrId(state.VAppID.ValueString(), true)
	if err != nil {
		if errors.Is(err, govcd.ErrorEntityNotFound) {
			resp.Diagnostics.AddError("vApp not found", err.Error())
			return
		}
		resp.Diagnostics.AddError("Error retrieving vApp", err.Error())
		return
	}

	if len(plan.GuestProperties) > 0 {
		x := plan.getGuestProperties()

		_, err = vapp.SetProductSectionList(x)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Adding guest properties",
				"Could not add guest properties, unexpected error: "+err.Error(),
			)
			return
		}
	} // end if !plan.GuestProperties.IsNull()

	var runtimeLease, storageLease int

	if len(plan.Lease) > 0 {
		runtimeLease = int(plan.Lease[0].RuntimeLeaseInSec.ValueInt64())
		storageLease = int(plan.Lease[0].StorageLeaseInSec.ValueInt64())
	} else {
		adminOrg, errGetAdminOrg := r.client.Vmware.GetAdminOrgById(r.org.GetID())
		if errGetAdminOrg != nil {
			resp.Diagnostics.AddError("Error retrieving Org", errGetAdminOrg.Error())
			return
		}

		if adminOrg.AdminOrg.OrgSettings == nil || adminOrg.AdminOrg.OrgSettings.OrgVAppLeaseSettings == nil {
			resp.Diagnostics.AddError("Error retrieving Org", "Org settings are not defined")
			return
		}

		runtimeLease = *adminOrg.AdminOrg.OrgSettings.OrgVAppLeaseSettings.DeploymentLeaseSeconds
		storageLease = *adminOrg.AdminOrg.OrgSettings.OrgVAppLeaseSettings.StorageLeaseSeconds
	}

	err = vapp.RenewLease(runtimeLease, storageLease)
	if err != nil {
		resp.Diagnostics.AddError("Error renewing lease", err.Error())
		return
	}

	if !plan.Description.IsNull() {
		err = vapp.UpdateDescription(plan.Description.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Error updating VApp description", err.Error())
			return
		}
	}

	if len(plan.GuestProperties) > 0 {
		x := plan.getGuestProperties()
		_, err = vapp.SetProductSectionList(x)
		if err != nil {
			resp.Diagnostics.AddError("Error updating VApp guest properties", err.Error())
			return
		}
	}

	// power_on
	if !plan.PowerON.IsNull() {
		if plan.PowerON.ValueBool() {
			task, errPowerOn := vapp.PowerOn()
			if errPowerOn != nil {
				resp.Diagnostics.AddError("Error powering on VApp", errPowerOn.Error())
				return
			}
			err = task.WaitTaskCompletion()
			if err != nil {
				resp.Diagnostics.AddError("Error powering on VApp", err.Error())
				return
			}
		} else {
			task, errUndeploy := vapp.Undeploy()
			if errUndeploy != nil {
				resp.Diagnostics.AddError("Error powering off VApp", errUndeploy.Error())
				return
			}
			err = task.WaitTaskCompletion()
			if err != nil {
				resp.Diagnostics.AddError("Error powering off VApp", err.Error())
				return
			}
		}
	}

	guestProperties, err := vapp.GetProductSectionList()
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving guest properties", err.Error())
	}
	leaseInfo, err := vapp.GetLease()
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving lease info", err.Error())
		return
	}

	// Request vApp
	vappRefreshed, err := r.vdc.GetVAppByNameOrId(state.VAppID.ValueString(), true)
	if err != nil {
		if errors.Is(err, govcd.ErrorEntityNotFound) {
			resp.Diagnostics.AddError("vApp not found after creating", err.Error())
			return
		}
		resp.Diagnostics.AddError("Error retrieving vApp after creating", err.Error())
		return
	}

	var statusText string

	statusText, err = vappRefreshed.GetStatus()
	if err != nil {
		statusText = vAppUnknownStatus
	}

	nPlan := &vappResourceModel{
		VAppID:      types.StringValue(vappRefreshed.VApp.ID),
		VAppName:    types.StringValue(vappRefreshed.VApp.Name),
		Description: types.StringValue(vappRefreshed.VApp.Description),
		StatusText:  types.StringValue(statusText),
		StatusCode:  types.Int64Value(int64(vappRefreshed.VApp.Status)),
		Href:        types.StringValue(vappRefreshed.VApp.HREF),

		VDC:     types.StringValue(r.vdc.GetName()),
		PowerON: plan.PowerON,
	}

	if plan.Lease != nil && len(plan.Lease) > 0 {
		nPlan.Lease = make([]vappLeaseModel, 1)
		nPlan.Lease = append(nPlan.Lease, vappLeaseModel{
			RuntimeLeaseInSec: types.Int64Value(int64(leaseInfo.DeploymentLeaseInSeconds)),
			StorageLeaseInSec: types.Int64Value(int64(leaseInfo.StorageLeaseInSeconds)),
		})
	}

	if guestProperties != nil && guestProperties.ProductSection != nil && guestProperties.ProductSection.Property != nil {
		for _, guestProperty := range guestProperties.ProductSection.Property {
			if guestProperty.Value != nil {
				nPlan.GuestProperties[types.StringValue(guestProperty.Key)] = types.StringValue(guestProperty.Value.Value)
			}
		}
	}

	tflog.Info(ctx, fmt.Sprintf("Vapp %s created", nPlan.VAppName.ValueString()))

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, &nPlan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *vappResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Get current state
	var state *vappResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Init resource
	resp.Diagnostics.Append(r.Init(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Request vApp
	vapp, err := r.vdc.GetVAppByNameOrId(state.VAppID.ValueString(), true)
	if err != nil {
		if errors.Is(err, govcd.ErrorEntityNotFound) {
			resp.Diagnostics.AddError("vApp not found", err.Error())
			return
		}
		resp.Diagnostics.AddError("Error retrieving vApp", err.Error())
		return
	}
	// to avoid network destroy issues - detach networks from vApp
	task, err := vapp.RemoveAllNetworks()
	if err != nil {
		resp.Diagnostics.AddError("Error delete VAPP", err.Error())
		return
	}
	err = task.WaitTaskCompletion()
	if err != nil {
		resp.Diagnostics.AddError("Error delete VAPP", err.Error())
		return
	}

	err = tryUndeploy(*vapp)
	if err != nil {
		resp.Diagnostics.AddError("Error delete VAPP", err.Error())
		return
	}

	task, err = vapp.Delete()
	if err != nil {
		resp.Diagnostics.AddError("Error delete VAPP", err.Error())
		return
	}

	err = task.WaitTaskCompletion()
	if err != nil {
		resp.Diagnostics.AddError("Error delete VAPP", err.Error())
		return
	}
}

func (r *vappResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

// getGuestProperties returns the guest properties of a vApp.
func (vapp *vappResourceModel) getGuestProperties() *govcdtypes.ProductSectionList {
	x := &govcdtypes.ProductSectionList{
		ProductSection: &govcdtypes.ProductSection{
			Info:     "Custom properties",
			Property: []*govcdtypes.Property{},
		},
	}

	for k, v := range vapp.GuestProperties {
		oneProp := &govcdtypes.Property{
			UserConfigurable: true,
			Type:             "string",
			Key:              k.String(),
			Label:            k.String(),
			Value:            &govcdtypes.Value{Value: v.String()},
		}
		x.ProductSection.Property = append(x.ProductSection.Property, oneProp)
	}

	return x
}

// tryUndeploy try to undeploy a vApp, but do not throw an error if the vApp is powered off.
// Very often the vApp is powered off at this point and Undeploy() would fail with error:
// "The requested operation could not be executed since vApp vApp_name is not running"
// So, if the error matches we just ignore it and the caller may fast forward to vapp.Delete().
func tryUndeploy(vapp govcd.VApp) error {
	task, err := vapp.Undeploy()
	reErr := regexp.MustCompile(`.*The requested operation could not be executed since vApp.*is not running.*`)
	if err != nil && reErr.MatchString(err.Error()) {
		// ignore - can't be undeployed
		return nil
	} else if err != nil {
		return fmt.Errorf("error undeploying vApp: %w", err)
	}

	err = task.WaitTaskCompletion()
	if err != nil {
		return fmt.Errorf("error undeploying vApp: %w", err)
	}
	return nil
}
