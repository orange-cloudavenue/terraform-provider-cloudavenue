// Package vdc provides a resource to manage VDCs.
package vdc

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common"
)

var (
	_ datasource.DataSource              = &vdcsDataSource{}
	_ datasource.DataSourceWithConfigure = &vdcsDataSource{}
)

// NewVDCDataSource returns a new resource implementing the vdcs data source.
func NewVDCDataSource() datasource.DataSource {
	return &vdcDataSource{}
}

type vdcDataSource struct {
	client *client.CloudAvenue
}

type vdcDataSourceModel struct {
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

func (d *vdcDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName
}

func (d *vdcDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Provides a Cloud Avenue Organization vDC data source. An Organization VDC can be used to reference a vDC and use its data within other resources or data sources.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID of the vDC.",
			},
			"name": schema.StringAttribute{
				Required: true,
				MarkdownDescription: "The name of the org vDC. It must be unique in the organization.\n" +
					"The length must be between 2 and 27 characters.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(2, 27),
				},
			},
			"description": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The description of the org vDC.",
			},
			"cpu_speed_in_mhz": schema.Float64Attribute{
				Computed: true,
				MarkdownDescription: "Specifies the clock frequency, in Mhz, for any virtual CPU that is allocated to a VM.\n" +
					"It must be at least 1200.",
			},
			"cpu_allocated": schema.Float64Attribute{
				Computed: true,
				MarkdownDescription: "CPU capacity in *MHz* that is committed to be available or used as a limit in PAYG mode.\n" +
					"It must be at least 5 * `cpu_speed_in_mhz` and at most 200 * `cpu_speed_in_mhz`.\n\n" +
					" -> Note: Reserved capacity is automatically set according to the service class.",
			},
			"memory_allocated": schema.Float64Attribute{
				Computed: true,
				MarkdownDescription: "Memory capacity in Gb that is committed to be available or used as a limit in PAYG mode.\n" +
					"It must be between 1 and 5000.",
			},
			"vdc_group": schema.StringAttribute{
				Computed: true,
				MarkdownDescription: "Name of an existing vDC group or a new one. This allows you to isolate your VDC.\n" +
					"VMs of vDCs which belong to the same vDC group can communicate together.",
			},
			"service_class": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The service class of the org vDC. It can be `ECO`, `STD`, `HP` or `VOIP`.",
			},
			"disponibility_class": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The disponibility class of the org vDC. It can be `ONE-ROOM`, `DUAL-ROOM` or `HA-DUAL-ROOM`.",
			},
			"billing_model": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Choose Billing model of compute resources. It can be `PAYG`, `DRAAS` or `RESERVED`.",
			},
			"storage_billing_model": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Choose Billing model of storage resources. It can be `PAYG` or `RESERVED`.",
			},
			"storage_profiles": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "List of storage profiles for this vDC.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"class": schema.StringAttribute{
							Computed: true,
							MarkdownDescription: "The storage class of the storage profile.\n" +
								"It can be `silver`, `silver_r1`, `silver_r2`, `gold`, `gold_r1`, `gold_r2`, `gold_hm`, `platinum3k`, `platinum3k_r1`, `platinum3k_r2`, `platinum3k_hm`, `platinum7k`, `platinum7k_r1`, `platinum7k_r2`, `platinum7k_hm`.",
						},
						"limit": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "Max number of units allocated for this storage profile. In Gb. It must be between 500 and 10000.",
						},
						"default": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Set this storage profile as default for this vDC. Only one storage profile can be default per vDC.",
						},
					},
				},
			},
		},
	}
}

func (d *vdcDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = client
}

func (d *vdcDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data vdcDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get vDC info
	vdc, httpR, err := d.client.APIClient.VDCApi.GetOrgVdcByName(d.client.Auth, data.Name.ValueString())
	if apiErr := helpers.CheckAPIError(err, httpR); apiErr != nil {
		resp.Diagnostics.Append(apiErr.GetTerraformDiagnostic())
		if resp.Diagnostics.HasError() {
			return
		}
		return
	}
	defer httpR.Body.Close()

	// Get vDC UUID by parsing vDCs list and set URN ID
	var ID string
	vdcs, httpR, err := d.client.APIClient.VDCApi.GetOrgVdcs(d.client.Auth)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read vdcs detail, got error: %s", err))
		return
	}
	defer httpR.Body.Close()

	for _, v := range vdcs {
		if data.Name.ValueString() == v.VdcName {
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

	data = vdcDataSourceModel{
		ID:                     types.StringValue(ID),
		VDCGroup:               types.StringValue(vdc.VdcGroup),
		Name:                   types.StringValue(vdc.Vdc.Name),
		Description:            types.StringValue(vdc.Vdc.Description),
		VDCServiceClass:        types.StringValue(vdc.Vdc.VdcServiceClass),
		VDCDisponibilityClass:  types.StringValue(vdc.Vdc.VdcDisponibilityClass),
		VDCBillingModel:        types.StringValue(vdc.Vdc.VdcBillingModel),
		VcpuInMhz2:             types.Float64Value(vdc.Vdc.VcpuInMhz2),
		CPUAllocated:           types.Float64Value(vdc.Vdc.CpuAllocated),
		MemoryAllocated:        types.Float64Value(vdc.Vdc.MemoryAllocated),
		VDCStorageBillingModel: types.StringValue(vdc.Vdc.VdcStorageBillingModel),
		VDCStorageProfiles:     profiles,
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
