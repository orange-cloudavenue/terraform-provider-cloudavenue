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
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	sdkResource "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers"
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
}

type vappResourceModel struct {
	ID              types.String                  `tfsdk:"id"`
	VappName        types.String                  `tfsdk:"vapp_name"`
	VappID          types.String                  `tfsdk:"vapp_id"`
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
	resp.TypeName = req.ProviderTypeName + "_vapp"
}

// Schema defines the schema for the resource.
func (r *vappResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The Edge Gateway resource allows you to create and manage Edge Gateways in CloudAvenue.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				CustomType:          types.StringType,
				Computed:            true,
				MarkdownDescription: "The ID is a `vapp_id`.",
			},
			"vapp_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "A name for the vApp, unique within the VDC. Required if `vapp_id` is not set.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"vapp_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID of vApp",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"vdc": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The name of VDC to use, optional if defined at provider level",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
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

	// If VDC is not defined at data source level, use the one defined at provider level
	if plan.VDC.IsNull() || plan.VDC.IsUnknown() {
		if r.client.DefaultVDCExist() {
			plan.VDC = types.StringValue(r.client.GetDefaultVDC())
		} else {
			resp.Diagnostics.AddError("Missing VDC", "VDC is required when not defined at provider level")
			return
		}
	}

	org, vdc, err := r.client.GetOrgAndVDC(r.client.GetOrg(), plan.VDC.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving VDC", err.Error())
		return
	}

	vapp, err := vdc.CreateRawVApp(plan.VappName.ValueString(), plan.Description.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error creating vApp", err.Error())
		return
	}

	// Wait for job to complete
	createStateConf := &sdkResource.StateChangeConf{
		Delay: 5 * time.Second,
		Refresh: func() (interface{}, string, error) {
			currentStatus, _ := vapp.GetStatus()
			tflog.Debug(ctx, fmt.Sprintf("Creating Vapp status: %s", currentStatus))
			if currentStatus == "UNRESOLVED" {
				return nil, helpers.PENDING.String(), nil
			}
			return helpers.DONE.String(), helpers.DONE.String(), nil
		},
		MinTimeout: 5 * time.Second,
		Timeout:    90 * time.Second,
		Pending:    []string{helpers.PENDING.String()},
		Target:     []string{helpers.DONE.String()},
	}

	// Wait vapp status is not UNRESOLVED
	_, err = createStateConf.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating VDC",
			"Could not create vdc, unexpected error: "+err.Error(),
		)
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
		adminOrg, errGetAdminOrg := r.client.Vmware.GetAdminOrgById(org.Org.ID)
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
	vappRefreshed, err := vdc.GetVAppByNameOrId(vapp.VApp.ID, true)
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
		ID:          types.StringValue(vappRefreshed.VApp.ID),
		VappName:    types.StringValue(vappRefreshed.VApp.Name),
		VappID:      types.StringValue(vappRefreshed.VApp.ID),
		Description: types.StringValue(vappRefreshed.VApp.Description),
		StatusText:  types.StringValue(statusText),
		StatusCode:  types.Int64Value(int64(vappRefreshed.VApp.Status)),
		Href:        types.StringValue(vappRefreshed.VApp.HREF),

		VDC:     types.StringValue(vdc.Vdc.Name),
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

	tflog.Info(ctx, fmt.Sprintf("Vapp %s created", nPlan.VappName.ValueString()))

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

	// If VDC is not defined at data source level, use the one defined at provider level
	if state.VDC.IsNull() || state.VDC.IsUnknown() {
		if r.client.DefaultVDCExist() {
			state.VDC = types.StringValue(r.client.GetDefaultVDC())
		} else {
			resp.Diagnostics.AddError("Missing VDC", "VDC is required when not defined at provider level")
			return
		}
	}

	_, vdc, err := r.client.GetOrgAndVDC(r.client.GetOrg(), state.VDC.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving VDC", err.Error())
		return
	}

	var vappNameID string
	if state.VappID.IsNull() || state.VappID.IsUnknown() {
		vappNameID = state.VappName.ValueString()
	} else {
		vappNameID = state.VappID.ValueString()
	}

	tflog.Info(ctx, fmt.Sprintf("Reading vApp %s", vappNameID))
	// Request vApp
	vappRefreshed, err := vdc.GetVAppByNameOrId(vappNameID, true)
	if err != nil {
		if errors.Is(err, govcd.ErrorEntityNotFound) {
			resp.Diagnostics.AddError("vApp not found", err.Error())
			return
		}
		resp.Diagnostics.AddError("Error retrieving vApp", err.Error())
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
		ID:          types.StringValue(vappRefreshed.VApp.ID),
		VappName:    types.StringValue(vappRefreshed.VApp.Name),
		VappID:      types.StringValue(vappRefreshed.VApp.ID),
		Description: types.StringValue(vappRefreshed.VApp.Description),
		StatusText:  types.StringValue(statusText),
		StatusCode:  types.Int64Value(int64(vappRefreshed.VApp.Status)),
		Href:        types.StringValue(vappRefreshed.VApp.HREF),
		VDC:         types.StringValue(vdc.Vdc.Name),

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

	org, vdc, err := r.client.GetOrgAndVDC(r.client.GetOrg(), plan.VDC.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving VDC", err.Error())
		return
	}

	// Request vApp
	vapp, err := vdc.GetVAppByNameOrId(state.ID.String(), true)
	if err != nil {
		if errors.Is(err, govcd.ErrorEntityNotFound) {
			resp.Diagnostics.AddError("vApp not found after creating", err.Error())
			return
		}
		resp.Diagnostics.AddError("Error retrieving vApp after creating", err.Error())
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
		adminOrg, errGetAdminOrg := r.client.Vmware.GetAdminOrgById(org.Org.ID)
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
	vappRefreshed, err := vdc.GetVAppByNameOrId(state.ID.String(), true)
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
		ID:          types.StringValue(vappRefreshed.VApp.ID),
		VappName:    types.StringValue(vappRefreshed.VApp.Name),
		VappID:      types.StringValue(vappRefreshed.VApp.ID),
		Description: types.StringValue(vappRefreshed.VApp.Description),
		StatusText:  types.StringValue(statusText),
		StatusCode:  types.Int64Value(int64(vappRefreshed.VApp.Status)),
		Href:        types.StringValue(vappRefreshed.VApp.HREF),

		VDC:     types.StringValue(vdc.Vdc.Name),
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

	tflog.Info(ctx, fmt.Sprintf("Vapp %s created", nPlan.VappName.ValueString()))

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

	_, vdc, err := r.client.GetOrgAndVDC(r.client.GetOrg(), state.VDC.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving VDC", err.Error())
		return
	}

	// Request vApp
	vapp, err := vdc.GetVAppByNameOrId(state.ID.String(), true)
	if err != nil {
		if errors.Is(err, govcd.ErrorEntityNotFound) {
			resp.Diagnostics.AddError("vApp not found after creating", err.Error())
			return
		}
		resp.Diagnostics.AddError("Error retrieving vApp after creating", err.Error())
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
	resource.ImportStatePassthroughID(ctx, path.Root("vapp_id"), req, resp)
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
