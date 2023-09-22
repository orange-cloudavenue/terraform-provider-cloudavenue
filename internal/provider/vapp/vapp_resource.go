// Package vapp provides a Terraform resource to manage vApps.
package vapp

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/adminorg"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vapp"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vdc"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
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
	client   *client.CloudAvenue
	vdc      vdc.VDC
	adminorg adminorg.AdminOrg
	vapp     vapp.VAPP
}

// Metadata returns the resource type name.
func (r *vappResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName
}

// Schema defines the schema for the resource.
func (r *vappResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = vappSchema().GetResource(ctx)
}

func (r *vappResource) Init(ctx context.Context, rm *vappResourceModel) (diags diag.Diagnostics) {
	r.adminorg, diags = adminorg.Init(r.client)
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
func (r *vappResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer metrics.New("cloudavenue_vapp", r.client.GetOrgName(), metrics.Create)()

	// Retrieve values from plan
	var (
		plan  *vappResourceModel
		diags diag.Diagnostics
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

	// Create vApp
	r.vapp, diags = vapp.Create(r.vdc, plan.VAppName.ValueString(), plan.Description.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Wait for job to complete
	errRetry := retry.RetryContext(ctx, 90*time.Second, func() *retry.RetryError {
		currentStatus, errGetStatus := r.vapp.GetStatus()
		if errGetStatus != nil {
			retry.NonRetryableError(errGetStatus)
		}
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

	// Update vApp
	state := &vappResourceModel{
		Description: types.StringValue(r.vapp.GetDescription()),
	}
	resp.Diagnostics.Append(r.updateVapp(ctx, plan, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.VAppID = types.StringValue(r.vapp.GetID())
	plan.VDC = types.StringValue(r.vdc.GetName())

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *vappResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer metrics.New("cloudavenue_vapp", r.client.GetOrgName(), metrics.Read)()

	var (
		state *vappResourceModel
		diags diag.Diagnostics
	)

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

	r.vapp, diags = vapp.Init(r.client, r.vdc, state.VAppID, state.VAppName)
	if diags.Contains(vapp.DiagVAppNotFound) {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set data
	plan := &vappResourceModel{
		VAppID:          types.StringValue(r.vapp.GetID()),
		VAppName:        types.StringValue(r.vapp.GetName()),
		VDC:             types.StringValue(r.vdc.GetName()),
		Description:     utils.StringValueOrNull(r.vapp.GetDescription()),
		Lease:           types.ObjectNull(vappLeaseAttrTypes),
		GuestProperties: types.MapNull(types.StringType),
	}

	// Get guest properties
	guestProperties, diags := processGuestProperties(r.vapp)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	if len(guestProperties) > 0 {
		plan.GuestProperties, diags = types.MapValue(types.StringType, guestProperties)
		resp.Diagnostics.Append(diags...)
		if diags.HasError() {
			return
		}
	}

	// Get lease info
	leaseInfo, err := r.vapp.GetLease()
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving lease info", err.Error())
		return
	}

	if leaseInfo != nil {
		plan.Lease, diags = types.ObjectValueFrom(ctx, vappLeaseAttrTypes, vappLeaseModel{
			RuntimeLeaseInSec: types.Int64Value(int64(leaseInfo.DeploymentLeaseInSeconds)),
			StorageLeaseInSec: types.Int64Value(int64(leaseInfo.StorageLeaseInSeconds)),
		})
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	// Set refreshed plan
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *vappResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer metrics.New("cloudavenue_vapp", r.client.GetOrgName(), metrics.Update)()

	var (
		plan, state *vappResourceModel
		diags       diag.Diagnostics
	)

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
	r.vapp, diags = vapp.Init(r.client, r.vdc, plan.VAppID, plan.VAppName)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update vApp
	resp.Diagnostics.Append(r.updateVapp(ctx, plan, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *vappResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer metrics.New("cloudavenue_vapp", r.client.GetOrgName(), metrics.Delete)()

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
	defer metrics.New("cloudavenue_vapp", r.client.GetOrgName(), metrics.Import)()

	var (
		diags        diag.Diagnostics
		vAppIDOrName string
	)

	// Split req.ID with dot. ID format is EdgeGatewayIDOrName.StaticRouteNameOrID
	idParts := strings.Split(req.ID, ".")

	switch len(idParts) {
	case 1:
		vAppIDOrName = idParts[0]
		r.vdc, diags = vdc.Init(r.client, basetypes.NewStringNull())
		resp.Diagnostics.Append(diags...)
	case 2:
		vAppIDOrName = idParts[1]
		r.vdc, diags = vdc.Init(r.client, basetypes.NewStringValue(idParts[0]))
		resp.Diagnostics.Append(diags...)
	default:
		resp.Diagnostics.AddError("Invalid ID format", fmt.Sprintf("ID format is VDCName.VAppIDOrName or VAppIDOrName, got: %s", req.ID))
	}
	if resp.Diagnostics.HasError() {
		return
	}

	vapp, err := r.vdc.GetVAppByNameOrId(vAppIDOrName, true)
	if err != nil {
		if errors.Is(err, govcd.ErrorEntityNotFound) {
			resp.Diagnostics.AddError("vApp not found", err.Error())
			return
		}
		resp.Diagnostics.AddError("Error retrieving vApp", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), vapp.VApp.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), vapp.VApp.Name)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("vdc"), r.vdc.GetName())...)
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

// updateVapp make updates only on elements that must be updated.
func (r *vappResource) updateVapp(ctx context.Context, plan, state *vappResourceModel) (d diag.Diagnostics) {
	var runtimeLease, storageLease int

	// Get lease config
	if !plan.Lease.IsNull() {
		l := &vappLeaseModel{}
		d.Append(plan.Lease.As(ctx, l, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    false,
			UnhandledUnknownAsEmpty: false,
		})...)
		if d.HasError() {
			return
		}
		runtimeLease = int(l.RuntimeLeaseInSec.ValueInt64())
		storageLease = int(l.StorageLeaseInSec.ValueInt64())
	} else {
		runtimeLease = *r.adminorg.GetOrgVAppLeaseSettings().DeploymentLeaseSeconds
		storageLease = *r.adminorg.GetOrgVAppLeaseSettings().StorageLeaseSeconds
	}

	// Update lease if needed
	if runtimeLease != r.vapp.GetDeploymentLeaseInSeconds() ||
		storageLease != r.vapp.GetStorageLeaseInSeconds() {
		if err := r.vapp.RenewLease(runtimeLease, storageLease); err != nil {
			d.AddError("Error renewing lease", err.Error())
			return
		}
	}
	if err := r.vapp.RenewLease(runtimeLease, storageLease); err != nil {
		d.AddError("Error renewing lease", err.Error())
		return
	}

	// Update description if needed
	if !plan.Description.Equal(state.Description) {
		if err := r.vapp.UpdateDescription(plan.Description.ValueString()); err != nil {
			d.AddError("Error updating VApp description", err.Error())
			return
		}
	}

	// Update GuestProperties if needed
	if !reflect.DeepEqual(plan.GuestProperties, state.GuestProperties) {
		// Init guestProperties struct
		x := &govcdtypes.ProductSectionList{
			ProductSection: &govcdtypes.ProductSection{
				Info:     "Custom properties",
				Property: []*govcdtypes.Property{},
			},
		}

		// Extract values from plan
		if !plan.GuestProperties.IsNull() {
			guestProperties := make(map[string]string)
			d.Append(plan.GuestProperties.ElementsAs(ctx, &guestProperties, true)...)
			if d.HasError() {
				return
			}
			for k, v := range guestProperties {
				oneProp := &govcdtypes.Property{
					UserConfigurable: true,
					Type:             "string",
					Key:              k,
					Label:            k,
					Value:            &govcdtypes.Value{Value: v},
				}
				x.ProductSection.Property = append(x.ProductSection.Property, oneProp)
			}
		}
		if _, err := r.vapp.SetProductSectionList(x); err != nil {
			d.AddError("Error updating VApp guest properties", err.Error())
			return
		}
	}

	return nil
}
