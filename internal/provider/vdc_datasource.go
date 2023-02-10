package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &vdcsDataSource{}
	_ datasource.DataSourceWithConfigure = &vdcsDataSource{}
)

func NewVdcDataSource() datasource.DataSource {
	return &vdcDataSource{}
}

type vdcDataSource struct {
	client *CloudAvenueClient
}

type vdcDataSourceModel struct {
	ID                     types.String             `tfsdk:"id"`
	Name                   types.String             `tfsdk:"name"`
	Description            types.String             `tfsdk:"description"`
	VdcServiceClass        types.String             `tfsdk:"service_class"`
	VdcDisponibilityClass  types.String             `tfsdk:"disponibility_class"`
	VdcBillingModel        types.String             `tfsdk:"billing_model"`
	VcpuInMhz2             types.Float64            `tfsdk:"cpu_speed_in_mhz"`
	CPUAllocated           types.Float64            `tfsdk:"cpu_allocated"`
	MemoryAllocated        types.Float64            `tfsdk:"memory_allocated"`
	VdcStorageBillingModel types.String             `tfsdk:"storage_billing_model"`
	VdcStorageProfiles     []vdcStorageProfileModel `tfsdk:"storage_profile"`
	VdcGroup               types.String             `tfsdk:"vdc_group"`
}

func (d *vdcDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vdc"
}

func (d *vdcDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Provides a Cloud Avenue Organization VDC data source. An Organization VDC can be used to reference a VDC and use its data within other resources or data sources.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "ID is the Name of the VCD.",
			},
			"name": schema.StringAttribute{
				Required: true,
				MarkdownDescription: "The name of the org VDC. It must be unique in the organization.\n" +
					"The length must be between 2 and 27 characters.",
			},
			"description": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The description of the org VDC.",
			},
			"cpu_speed_in_mhz": schema.Float64Attribute{
				Computed: true,
				MarkdownDescription: "Specifies the clock frequency, in Mhz, for any virtual CPU that is allocated to a VM.\n" +
					"It must be at least 1200.",
			},
			"cpu_allocated": schema.Float64Attribute{
				Computed: true,
				MarkdownDescription: "CPU capacity in *MHz* that is committed to be available or used as a limit in PAYG mode.\n" +
					"It must be at least 5 * `cpu_speed_in_mhz` and at most 200 * `cpu_speed_in_mhz`.\n" +
					" *Note:* Reserved capacity is automatically set according to the service class.",
			},
			"memory_allocated": schema.Float64Attribute{
				Computed: true,
				MarkdownDescription: "Memory capacity in Gb that is committed to be available or used as a limit in PAYG mode.\n" +
					"It must be between 1 and 5000.",
			},
			"vdc_group": schema.StringAttribute{
				Computed: true,
				MarkdownDescription: "Name of an existing VDC group or a new one. This allows you to isolate your VDC.\n" +
					"VMs of VDCs which belong to the same VDC group can communicate together.",
			},
			"service_class": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The service class of the org VDC. It can be `ECO`, `STD`, `HP` or `VOIP`.",
			},
			"disponibility_class": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The disponibility class of the org VDC. It can be `ONE-ROOM`, `DUAL-ROOM` or `HA-DUAL-ROOM`.",
			},
			"billing_model": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Choose Billing model of compute resources. It can be `PAYG`, `DRAAS` or `RESERVED`.",
			},
			"storage_billing_model": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Choose Billing model of storage resources. It can be `PAYG` or `RESERVED`.",
			},
		},
		Blocks: map[string]schema.Block{
			"storage_profile": schema.ListNestedBlock{
				MarkdownDescription: "List of storage profiles for this VDC.",
				NestedObject: schema.NestedBlockObject{
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
							MarkdownDescription: "Set this storage profile as default for this VDC. Only one storage profile can be default per VDC.",
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

	client, ok := req.ProviderData.(*CloudAvenueClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *CloudAvenueClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
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

	vdcs, _, err := d.client.VDCApi.ApiCustomersV20VdcsVdcNameGet(d.client.auth, data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read vdc detail, got error: %s", err))
		return
	}

	var profiles []vdcStorageProfileModel
	for _, profile := range vdcs.Vdc.VdcStorageProfiles {
		p := vdcStorageProfileModel{
			Class:   types.StringValue(profile.Class),
			Limit:   types.Int64Value(int64(profile.Limit)),
			Default: types.BoolValue(profile.Default_),
		}
		profiles = append(profiles, p)
	}

	data = vdcDataSourceModel{
		ID:                     types.StringValue(vdcs.Vdc.Name),
		VdcGroup:               types.StringValue(vdcs.VdcGroup),
		Name:                   types.StringValue(vdcs.Vdc.Name),
		Description:            types.StringValue(vdcs.Vdc.Description),
		VdcServiceClass:        types.StringValue(vdcs.Vdc.VdcServiceClass),
		VdcDisponibilityClass:  types.StringValue(vdcs.Vdc.VdcDisponibilityClass),
		VdcBillingModel:        types.StringValue(vdcs.Vdc.VdcBillingModel),
		VcpuInMhz2:             types.Float64Value(vdcs.Vdc.VcpuInMhz2),
		CPUAllocated:           types.Float64Value(vdcs.Vdc.CpuAllocated),
		MemoryAllocated:        types.Float64Value(vdcs.Vdc.MemoryAllocated),
		VdcStorageBillingModel: types.StringValue(vdcs.Vdc.VdcStorageBillingModel),
		VdcStorageProfiles:     profiles,
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
