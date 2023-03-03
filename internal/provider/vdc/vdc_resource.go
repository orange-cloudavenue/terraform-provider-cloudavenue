// Package vdc provides a resource to manage VDCs.
package vdc

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/float64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	sdkResource "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	apiclient "github.com/orange-cloudavenue/cloudavenue-sdk-go"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &vdcResource{}
	_ resource.ResourceWithConfigure   = &vdcResource{}
	_ resource.ResourceWithImportState = &vdcResource{}
)

// NewVDCResource is a helper function to simplify the provider implementation.
func NewVDCResource() resource.Resource {
	return &vdcResource{}
}

// vdcResource is the resource implementation.
type vdcResource struct {
	client *client.CloudAvenue
}

type vdcResourceModel struct {
	Timeouts               timeouts.Value           `tfsdk:"timeouts"`
	ID                     types.String             `tfsdk:"id"`
	Name                   types.String             `tfsdk:"name"`
	Description            types.String             `tfsdk:"description"`
	VDCServiceClass        types.String             `tfsdk:"service_class"`
	VDCDisponibilityClass  types.String             `tfsdk:"disponibility_class"`
	VDCBillingModel        types.String             `tfsdk:"billing_model"`
	VcpuInMhz2             types.Float64            `tfsdk:"cpu_speed_in_mhz"`
	CPUAllocated           types.Float64            `tfsdk:"cpu_allocated"`
	MemoryAllocated        types.Float64            `tfsdk:"memory_allocated"`
	VDCStorageBillingModel types.String             `tfsdk:"storage_billing_model"`
	VDCStorageProfiles     []vdcStorageProfileModel `tfsdk:"storage_profiles"`
	VDCGroup               types.String             `tfsdk:"vdc_group"`
}

// Metadata returns the resource type name.
func (r *vdcResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName
}

// Schema defines the schema for the resource.
func (r *vdcResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Provides a Cloud Avenue Organization vDC resource. This can be used to create, update and delete an Organization VDC.\n\n" +
			" -> Note: For more information about Organization vDC, please refer to the [Cloud Avenue documentation](https://wiki.cloudavenue.orange-business.com/w/index.php/Datacenter_virtuel).",
		Attributes: map[string]schema.Attribute{
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
				Read:   true,
				Delete: true,
				Update: true,
			}),
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID of the vDC.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
				MarkdownDescription: "(ForceNew) The name of the org vDC. It must be unique in the organization.\n" +
					"The length must be between 2 and 27 characters.\n" +
					helpers.ForceNewDescription,
				Validators: []validator.String{
					stringvalidator.LengthBetween(2, 27),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The description of the org vDC.",
			},
			"cpu_speed_in_mhz": schema.Float64Attribute{
				Required: true,
				MarkdownDescription: "Specifies the clock frequency, in Mhz, for any virtual CPU that is allocated to a VM.\n" +
					"It must be at least 1200.",
				Validators: []validator.Float64{
					float64validator.AtLeast(1200),
				},
			},
			"cpu_allocated": schema.Float64Attribute{
				Required: true,
				MarkdownDescription: "CPU capacity in *MHz* that is committed to be available or used as a limit in PAYG mode.\n" +
					"It must be at least 5 * `cpu_speed_in_mhz`.\n\n" +
					" -> Note: Reserved capacity is automatically set according to the service class.",
			},
			"memory_allocated": schema.Float64Attribute{
				Required: true,
				MarkdownDescription: "Memory capacity in Gb that is committed to be available or used as a limit in PAYG mode.\n" +
					"It must be between 1 and 5000.",
				Validators: []validator.Float64{
					float64validator.AtLeast(1),
					float64validator.AtMost(5000),
				},
			},
			"vdc_group": schema.StringAttribute{
				Required: true,
				MarkdownDescription: "(ForceNew) Name of an existing vDC group or a new one. This allows you to isolate your vDC.\n" +
					"VMs of vDCs which belong to the same vDC group can communicate together.\n" +
					helpers.ForceNewDescription,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"service_class": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The service class of the org vDC. It can be `ECO`, `STD`, `HP` or `VOIP`.",
				Validators: []validator.String{
					stringvalidator.OneOf("ECO", "STD", "HP", "VOIP"),
				},
			},
			"disponibility_class": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The disponibility class of the org vDC. It can be `ONE-ROOM`, `DUAL-ROOM` or `HA-DUAL-ROOM`.",
				Validators: []validator.String{
					stringvalidator.OneOf("ONE-ROOM", "DUAL-ROOM", "HA-DUAL-ROOM"),
				},
			},
			"billing_model": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Choose Billing model of compute resources. It can be `PAYG`, `DRAAS` or `RESERVED`.",
				Validators: []validator.String{
					stringvalidator.OneOf("PAYG", "DRAAS", "RESERVED"),
				},
			},
			"storage_billing_model": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Choose Billing model of storage resources. It can be `PAYG` or `RESERVED`.",
				Validators: []validator.String{
					stringvalidator.OneOf("PAYG", "RESERVED"),
				},
			},
			"storage_profiles": schema.SetNestedAttribute{
				Required:            true,
				MarkdownDescription: "List of storage profiles for this vDC.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"class": schema.StringAttribute{
							Required: true,
							MarkdownDescription: "The storage class of the storage profile.\n" +
								"It can be `silver`, `silver_r1`, `silver_r2`, `gold`, `gold_r1`, `gold_r2`, `gold_hm`, `platinum3k`, `platinum3k_r1`, `platinum3k_r2`, `platinum3k_hm`, `platinum7k`, `platinum7k_r1`, `platinum7k_r2`, `platinum7k_hm`.",
							Validators: []validator.String{
								stringvalidator.OneOf("silver", "silver_r1", "silver_r2", "gold", "gold_r1", "gold_r2", "gold_hm", "platinum3k", "platinum3k_r1", "platinum3k_r2", "platinum3k_hm", "platinum7k", "platinum7k_r1", "platinum7k_r2", "platinum7k_hm"),
							},
						},
						"limit": schema.Int64Attribute{
							Required:            true,
							MarkdownDescription: "Max number of units allocated for this storage profile. In Gb. It must be between 500 and 10000.",
							Validators: []validator.Int64{
								int64validator.AtLeast(500),
								int64validator.AtMost(10000),
							},
						},
						"default": schema.BoolAttribute{
							Required:            true,
							MarkdownDescription: "Set this storage profile as default for this vDC. Only one storage profile can be default per vDC.",
						},
					},
				},
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
				},
			},
		},
	}
}

// Configure configures the resource.
func (r *vdcResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *vdcResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan *vdcResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create() is passed a default timeout to use if no value
	// has been supplied in the Terraform configuration.
	createTimeout, errTO := plan.Timeouts.Create(ctx, 8*time.Minute)
	if errTO != nil {
		resp.Diagnostics.AddError(
			"Error creating timeout",
			"Could not create timeout, unexpected error",
		)
		return
	}

	ctxTO, cancel := context.WithTimeout(ctx, createTimeout)
	defer cancel()

	auth, errCtx := helpers.GetAuthContextWithTO(r.client.Auth, ctxTO)
	if errCtx != nil {
		resp.Diagnostics.AddError(
			"Error creating context",
			"Could not create context, context value token is not a string ?",
		)
		return
	}

	// Prepare the body to create a VDC.
	body := apiclient.CreateOrgVdcV2{
		VdcGroup: plan.VDCGroup.ValueString(),
		Vdc: &apiclient.OrgVdcV2{
			Name:                   plan.Name.ValueString(),
			Description:            plan.Description.ValueString(),
			VdcServiceClass:        plan.VDCServiceClass.ValueString(),
			VdcDisponibilityClass:  plan.VDCDisponibilityClass.ValueString(),
			VdcBillingModel:        plan.VDCBillingModel.ValueString(),
			VcpuInMhz2:             plan.VcpuInMhz2.ValueFloat64(),
			CpuAllocated:           plan.CPUAllocated.ValueFloat64(),
			MemoryAllocated:        plan.MemoryAllocated.ValueFloat64(),
			VdcStorageBillingModel: plan.VDCStorageBillingModel.ValueString(),
		},
	}

	// Iterate over the storage profiles and add them to the body.
	for _, storageProfile := range plan.VDCStorageProfiles {
		body.Vdc.VdcStorageProfiles = append(body.Vdc.VdcStorageProfiles, apiclient.VdcStorageProfilesV2{
			Class:    storageProfile.Class.ValueString(),
			Limit:    int32(storageProfile.Limit.ValueInt64()),
			Default_: storageProfile.Default.ValueBool(),
		})
	}

	var err error
	var job apiclient.Jobcreated
	var httpR *http.Response

	// Call API to create the resource and test for errors.
	job, httpR, err = r.client.APIClient.VDCApi.CreateOrgVdc(auth, body)
	if apiErr := helpers.CheckAPIError(err, httpR); apiErr != nil {
		defer httpR.Body.Close()
		resp.Diagnostics.Append(apiErr.GetTerraformDiagnostic())
		if resp.Diagnostics.HasError() {
			return
		}
	}

	// Wait for job to complete
	createStateConf := &sdkResource.StateChangeConf{
		Delay: 10 * time.Second,
		Refresh: func() (interface{}, string, error) {
			jobStatus, errGetJob := helpers.GetJobStatus(auth, r.client, job.JobId)
			if errGetJob != nil {
				return nil, "", err
			}
			return jobStatus, jobStatus.String(), nil
		},
		MinTimeout: 5 * time.Second,
		Timeout:    5 * time.Minute,
		Pending:    helpers.JobStatePending(),
		Target:     helpers.JobStateDone(),
	}

	_, err = createStateConf.WaitForStateContext(ctxTO)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating VDC",
			"Could not create vdc, unexpected error: "+err.Error(),
		)
		return
	}

	// Get vDC UUID by parsing vDCs list and set URN ID
	var ID string
	vdcs, httpR, err := r.client.APIClient.VDCApi.GetOrgVdcs(auth)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read vdcs detail, got error: %s", err))
		return
	}
	defer httpR.Body.Close()

	for _, v := range vdcs {
		if plan.Name.ValueString() == v.VdcName {
			ID = common.NormalizeID("urn:vcloud:vdc:", v.VdcUuid)
			break
		}
	}
	plan.ID = types.StringValue(ID)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "VDC created")

	// Save plan into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *vdcResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state *vdcResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Call API to read the resource and test for errors.
	readTimeout, errTO := state.Timeouts.Read(ctx, 8*time.Minute)
	if errTO != nil {
		resp.Diagnostics.AddError(
			"Error creating timeout",
			"Could not create timeout, unexpected error",
		)
		return
	}

	ctxTO, cancel := context.WithTimeout(ctx, readTimeout)
	defer cancel()

	auth, errCtx := helpers.GetAuthContextWithTO(r.client.Auth, ctxTO)
	if errCtx != nil {
		resp.Diagnostics.AddError(
			"Error creating context",
			"Could not create context, context value token is not a string ?",
		)
		return
	}

	// Get vDC info
	vdc, httpR, err := r.client.APIClient.VDCApi.GetOrgVdcByName(auth, state.Name.ValueString())
	if apiErr := helpers.CheckAPIError(err, httpR); apiErr != nil {
		resp.Diagnostics.Append(apiErr.GetTerraformDiagnostic())
		if resp.Diagnostics.HasError() {
			return
		}

		resp.State.RemoveResource(ctx)

		return
	}
	defer httpR.Body.Close()

	// Get vDC UUID by parsing vDCs list and set URN ID
	var ID string
	vdcs, httpR, err := r.client.APIClient.VDCApi.GetOrgVdcs(auth)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read vdcs detail, got error: %s", err))
		return
	}
	defer httpR.Body.Close()

	for _, v := range vdcs {
		if state.Name.ValueString() == v.VdcName {
			ID = common.NormalizeID("urn:vcloud:vdc:", v.VdcUuid)
			break
		}
	}

	// Get storageProfile
	var profiles []vdcStorageProfileModel
	for _, profile := range vdc.Vdc.VdcStorageProfiles {
		p := vdcStorageProfileModel{
			Class:   types.StringValue(profile.Class),
			Limit:   types.Int64Value(int64(profile.Limit)),
			Default: types.BoolValue(profile.Default_),
		}
		profiles = append(profiles, p)
	}

	// Convert from the API data model to the Terraform data model
	// and refresh any attribute values.
	state = &vdcResourceModel{
		Timeouts:               state.Timeouts,
		ID:                     types.StringValue(ID),
		Name:                   types.StringValue(vdc.Vdc.Name),
		Description:            types.StringValue(vdc.Vdc.Description),
		VDCGroup:               types.StringValue(vdc.VdcGroup),
		VDCServiceClass:        types.StringValue(vdc.Vdc.VdcServiceClass),
		VDCDisponibilityClass:  types.StringValue(vdc.Vdc.VdcDisponibilityClass),
		VDCBillingModel:        types.StringValue(vdc.Vdc.VdcBillingModel),
		VcpuInMhz2:             types.Float64Value(vdc.Vdc.VcpuInMhz2),
		CPUAllocated:           types.Float64Value(vdc.Vdc.CpuAllocated),
		MemoryAllocated:        types.Float64Value(vdc.Vdc.MemoryAllocated),
		VDCStorageBillingModel: types.StringValue(vdc.Vdc.VdcStorageBillingModel),
		VDCStorageProfiles:     profiles,
	}

	// Save updated state into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *vdcResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan *vdcResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update() is passed a default timeout to use if no value
	// has been supplied in the Terraform configuration.
	updateTimeout, errTO := plan.Timeouts.Update(ctx, 8*time.Minute)
	if errTO != nil {
		resp.Diagnostics.AddError(
			"Error creating timeout",
			"Could not create timeout, unexpected error",
		)
		return
	}

	ctxTO, cancel := context.WithTimeout(ctx, updateTimeout)
	defer cancel()

	auth, errCtx := helpers.GetAuthContextWithTO(r.client.Auth, ctxTO)
	if errCtx != nil {
		resp.Diagnostics.AddError(
			"Error creating context",
			"Could not create context, context value token is not a string ?",
		)
		return
	}

	// Convert from Terraform data model into API data model
	body := apiclient.UpdateOrgVdcV2{
		VdcGroup: plan.VDCGroup.ValueString(),
		Vdc: &apiclient.OrgVdcV2{
			Name:                   plan.Name.ValueString(),
			Description:            plan.Description.ValueString(),
			VdcServiceClass:        plan.VDCServiceClass.ValueString(),
			VdcDisponibilityClass:  plan.VDCDisponibilityClass.ValueString(),
			VdcBillingModel:        plan.VDCBillingModel.ValueString(),
			VcpuInMhz2:             plan.VcpuInMhz2.ValueFloat64(),
			CpuAllocated:           plan.CPUAllocated.ValueFloat64(),
			MemoryAllocated:        plan.MemoryAllocated.ValueFloat64(),
			VdcStorageBillingModel: plan.VDCStorageBillingModel.ValueString(),
		},
	}

	// Iterate over the storage profiles and add them to the body.
	for _, storageProfile := range plan.VDCStorageProfiles {
		body.Vdc.VdcStorageProfiles = append(body.Vdc.VdcStorageProfiles, apiclient.VdcStorageProfilesV2{
			Class:    storageProfile.Class.ValueString(),
			Limit:    int32(storageProfile.Limit.ValueInt64()),
			Default_: storageProfile.Default.ValueBool(),
		})
	}

	var err error
	var job apiclient.Jobcreated
	var httpR *http.Response

	// Call API to update the resource and test for errors.
	job, httpR, err = r.client.APIClient.VDCApi.UpdateOrgVdc(auth, body, body.Vdc.Name)
	if apiErr := helpers.CheckAPIError(err, httpR); apiErr != nil {
		defer httpR.Body.Close()
		resp.Diagnostics.Append(apiErr.GetTerraformDiagnostic())
		if resp.Diagnostics.HasError() {
			return
		}
	}

	// Wait for job to complete
	updateStateConf := &sdkResource.StateChangeConf{
		Delay: 10 * time.Second,
		Refresh: func() (interface{}, string, error) {
			jobStatus, errGetJob := helpers.GetJobStatus(auth, r.client, job.JobId)
			if errGetJob != nil {
				return nil, "", err
			}
			return jobStatus, jobStatus.String(), nil
		},
		MinTimeout: 5 * time.Second,
		Timeout:    5 * time.Minute,
		Pending:    helpers.JobStatePending(),
		Target:     helpers.JobStateDone(),
	}

	_, err = updateStateConf.WaitForStateContext(ctxTO)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating VDC",
			"Could not update vdc, unexpected error: "+err.Error(),
		)
		return
	}

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "VDC updated")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *vdcResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Get current state
	var state *vdcResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete() is passed a default timeout to use if no value
	// has been supplied in the Terraform configuration.
	deleteTimeout, errTO := state.Timeouts.Delete(ctx, 8*time.Minute)
	if errTO != nil {
		resp.Diagnostics.AddError(
			"Error creating timeout",
			"Could not create timeout, unexpected error",
		)
		return
	}

	ctxTO, cancel := context.WithTimeout(ctx, deleteTimeout)
	defer cancel()

	auth, errCtx := helpers.GetAuthContextWithTO(r.client.Auth, ctxTO)
	if errCtx != nil {
		resp.Diagnostics.AddError(
			"Error creating context",
			"Could not create context, context value token is not a string ?",
		)
		return
	}

	// Delete the VDC
	job, httpR, err := r.client.APIClient.VDCApi.DeleteOrgVdc(auth, state.Name.ValueString())
	if apiErr := helpers.CheckAPIError(err, httpR); apiErr != nil {
		defer httpR.Body.Close()
		resp.Diagnostics.Append(apiErr.GetTerraformDiagnostic())
		if resp.Diagnostics.HasError() {
			return
		}
	}

	// Wait for job to complete
	deleteStateConf := &sdkResource.StateChangeConf{
		Delay: 10 * time.Second,
		Refresh: func() (interface{}, string, error) {
			jobStatus, errGetJob := helpers.GetJobStatus(auth, r.client, job.JobId)
			if errGetJob != nil {
				return nil, "", errGetJob
			}

			return jobStatus, jobStatus.String(), nil
		},
		MinTimeout: 5 * time.Second,
		Timeout:    5 * time.Minute,
		Pending:    helpers.JobStatePending(),
		Target:     helpers.JobStateDone(),
	}

	_, err = deleteStateConf.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting vdc",
			"Could not delete vdc, unexpected error: "+err.Error(),
		)
		return
	}

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "VDC deleted")
}

func (r *vdcResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
