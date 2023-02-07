package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-validators/float64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	cloudavenue "github.com/orange-cloudavenue/cloudavenue-sdk-go"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &vdcResource{}
	_ resource.ResourceWithConfigure   = &vdcResource{}
	_ resource.ResourceWithImportState = &vdcResource{}
)

// NewVdcResource is a helper function to simplify the provider implementation.
func NewVdcResource() resource.Resource {
	return &vdcResource{}
}

// vdcResource is the resource implementation.
type vdcResource struct {
	client *CloudAvenueClient
}

type vdcResourceModel struct {
	ID                     types.String             `tfsdk:"id"`
	Name                   types.String             `tfsdk:"name"`
	Description            types.String             `tfsdk:"description"`
	VdcServiceClass        types.String             `tfsdk:"vdc_service_class"`
	VdcDisponibilityClass  types.String             `tfsdk:"vdc_disponibility_class"`
	VdcBillingModel        types.String             `tfsdk:"vdc_billing_model"`
	VcpuInMhz2             types.Float64            `tfsdk:"vcpu_in_mhz2"`
	CPUAllocated           types.Float64            `tfsdk:"cpu_allocated"`
	MemoryAllocated        types.Float64            `tfsdk:"memory_allocated"`
	VdcStorageBillingModel types.String             `tfsdk:"vdc_storage_billing_model"`
	VdcStorageProfiles     []vdcStorageProfileModel `tfsdk:"vdc_storage_profiles"`
	VdcGroup               types.String             `tfsdk:"vdc_group"`
}

type vdcStorageProfileModel struct {
	Class   types.String `tfsdk:"class"`
	Limit   types.Int64  `tfsdk:"limit"`
	Default types.Bool   `tfsdk:"default"`
}

// Metadata returns the resource type name.
func (r *vdcResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vdc"
}

// Schema defines the schema for the resource.
func (r *vdcResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "VDC resource allows you to create a org VDC in Cloud Avenue.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the org VDC. It must be unique in the organization. The length must be between 2 and 27 characters.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(2, 27),
				},
			},
			"description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The description of the org VDC.",
			},
			"vdc_group": schema.StringAttribute{
				Required: true,
				MarkdownDescription: "Name of an existing vDC group or a new one. This allows you to isolate your vDC." +
					"VMs of vDCs which belong to the same vDC group can communicate together.",
			},
			"vcpu_in_mhz2": schema.Float64Attribute{
				Required: true,
				MarkdownDescription: "Specifies the clock frequency, in Mhz, for any virtual CPU that is allocated to a VM." +
					"It must be at least 1200.",
				Validators: []validator.Float64{
					float64validator.AtLeast(1200),
				},
			},
			"cpu_allocated": schema.Float64Attribute{
				Required: true,
				MarkdownDescription: "Capacity that is committed to be available or used as a limit in PAYG mode." +
					"Unit for compute capacity allocated to this vdc is MHz. It must be beetwen 5 and 2500000.\n" +
					" *Note:* Reserved capacity is automatically set according to the service class.",
				Validators: []validator.Float64{
					float64validator.AtLeast(5),
					float64validator.AtMost(2500000),
				},
			},
			"memory_allocated": schema.Float64Attribute{
				Required: true,
				MarkdownDescription: "Memory capacity that is committed to be available or used as a limit in PAYG mode." +
					"Unit for memory capacity allocated to this vdc is Gb. It must be between 1 and 5000.",
				Validators: []validator.Float64{
					float64validator.AtLeast(1),
					float64validator.AtMost(5000),
				},
			},
			"vdc_service_class": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The service class of the org VDC. It can be `ECO`, `STD`, `HP` or `VOIP`.",
				Validators: []validator.String{
					stringvalidator.OneOf("ECO", "STD", "HP", "VOIP"),
				},
			},
			"vdc_disponibility_class": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The disponibility class of the org VDC. It can be `ONE-ROOM`, `DUAL-ROOM` or `HA-DUAL-ROOM`.",
				Validators: []validator.String{
					stringvalidator.OneOf("ONE-ROOM", "DUAL-ROOM", "HA-DUAL-ROOM"),
				},
			},
			"vdc_billing_model": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Choose Billing model of compute resources. It can be `PAYG`, `DRAAS` or `RESERVED`.",
				Validators: []validator.String{
					stringvalidator.OneOf("PAYG", "DRAAS", "RESERVED"),
				},
			},
			"vdc_storage_billing_model": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Choose Billing model of storage resources. It can be `PAYG` or `RESERVED`.",
				Validators: []validator.String{
					stringvalidator.OneOf("PAYG", "RESERVED"),
				},
			},
			"vdc_storage_profiles": schema.ListNestedAttribute{
				Required:            true,
				MarkdownDescription: "List of storage profiles for this VDC.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"class": schema.StringAttribute{
							Required: true,
							MarkdownDescription: "The storage class of the storage profile." +
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
							MarkdownDescription: "Set this storage profile as default for this VDC. Only one storage profile can be default per VDC.",
						},
					},
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

	client, ok := req.ProviderData.(*CloudAvenueClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *CloudAvenueClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

// Create creates the resource and sets the initial Terraform state.
func (r *vdcResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *vdcResourceModel

	// Read Terraform plan data into the model.
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Prepare the vdcObject to create a VDC.
	vdcObject := cloudavenue.CreateOrgVdcV2{
		VdcGroup: data.VdcGroup.ValueString(),
		Vdc: &cloudavenue.OrgVdcV2{
			Name:                   data.Name.ValueString(),
			Description:            data.Description.ValueString(),
			VdcServiceClass:        data.VdcServiceClass.ValueString(),
			VdcDisponibilityClass:  data.VdcDisponibilityClass.ValueString(),
			VdcBillingModel:        data.VdcBillingModel.ValueString(),
			VcpuInMhz2:             data.VcpuInMhz2.ValueFloat64(),
			CpuAllocated:           data.CPUAllocated.ValueFloat64(),
			MemoryAllocated:        data.MemoryAllocated.ValueFloat64(),
			VdcStorageBillingModel: data.VdcStorageBillingModel.ValueString(),
		},
	}

	// Iterate over the storage profiles and add them to the vdcObject.
	for _, storageProfile := range data.VdcStorageProfiles {
		vdcObject.Vdc.VdcStorageProfiles = append(vdcObject.Vdc.VdcStorageProfiles, cloudavenue.VdcStorageProfilesV2{
			Class:    storageProfile.Class.ValueString(),
			Limit:    int32(storageProfile.Limit.ValueInt64()),
			Default_: storageProfile.Default.ValueBool(),
		})
	}

	// Call API to create the resource and test for errors.
	_, _, err := r.client.VDCApi.ApiCustomersV20VdcsPost(r.client.auth, vdcObject)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating VDC",
			"Could not create VDC, unexpected error: "+err.Error(),
		)
		return
	}

	// Generate a unique ID for the resource.
	data.ID = utils.GenerateUUID(data.Name.ValueString())

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *vdcResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *vdcResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Call API to read the resource and test for errors.
	vdc, httpResp, err := r.client.VDCApi.ApiCustomersV20VdcsVdcNameGet(r.client.auth, data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading VDC",
			"Could not get VDC, unexpected error: "+err.Error(),
		)
		return
	}

	// If the resource does not exist, remove it from Terraform state.
	if httpResp.StatusCode == http.StatusNotFound {
		resp.State.RemoveResource(ctx)

		return
	}

	// Convert from the API data model to the Terraform data model
	// and refresh any attribute values.
	data = &vdcResourceModel{
		ID:                     utils.GenerateUUID(vdc.Vdc.Name),
		Name:                   types.StringValue(vdc.Vdc.Name),
		Description:            types.StringValue(vdc.Vdc.Description),
		VdcGroup:               types.StringValue(vdc.VdcGroup),
		VdcServiceClass:        types.StringValue(vdc.Vdc.VdcServiceClass),
		VdcDisponibilityClass:  types.StringValue(vdc.Vdc.VdcDisponibilityClass),
		VdcBillingModel:        types.StringValue(vdc.Vdc.VdcBillingModel),
		VcpuInMhz2:             types.Float64Value(vdc.Vdc.VcpuInMhz2),
		CPUAllocated:           types.Float64Value(vdc.Vdc.CpuAllocated),
		MemoryAllocated:        types.Float64Value(vdc.Vdc.MemoryAllocated),
		VdcStorageBillingModel: types.StringValue(vdc.Vdc.VdcStorageBillingModel),
		VdcStorageProfiles:     make([]vdcStorageProfileModel, len(vdc.Vdc.VdcStorageProfiles)),
	}

	for i, storageProfile := range vdc.Vdc.VdcStorageProfiles {
		data.VdcStorageProfiles[i] = vdcStorageProfileModel{
			Class:   types.StringValue(storageProfile.Class),
			Limit:   types.Int64Value(int64(storageProfile.Limit)),
			Default: types.BoolValue(storageProfile.Default_),
		}
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *vdcResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *vdcResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *vdcResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *vdcResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
