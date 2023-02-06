package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	_ "github.com/orange-cloudavenue/cloudavenue-sdk-go"
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
	ID       types.String `tfsdk:"id"`
	VdcGroup types.String `tfsdk:"vdc_group"`
	Name     types.String `tfsdk:"name"`
	Vdc      *vdcDetail   `tfsdk:"vdc"`
}

type vdcDetail struct {
	Name                   types.String         `tfsdk:"name"`
	Description            types.String         `tfsdk:"description"`
	VdcServiceClass        types.String         `tfsdk:"vdc_service_class"`
	VdcDisponibilityClass  types.String         `tfsdk:"vdc_disponibility_class"`
	VdcBillingModel        types.String         `tfsdk:"vdc_billing_model"`
	VcpuInMhz2             types.Float64        `tfsdk:"vcpu_in_mhz2"`
	CPUAllocated           types.Float64        `tfsdk:"cpu_allocated"`
	MemoryAllocated        types.Float64        `tfsdk:"memory_allocated"`
	VdcStorageBillingModel types.String         `tfsdk:"vdc_storage_billing_model"`
	VdcStorageProfiles     []vdcStorageProfiles `tfsdk:"vdc_storage_profiles"`
}

type vdcStorageProfiles struct {
	Class   types.String `tfsdk:"class"`
	Limit   types.Int64  `tfsdk:"limit"`
	Default types.Bool   `tfsdk:"default"`
}

func (d *vdcDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vdc"
}

func (d *vdcDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Show the vDC details.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "VDC name.",
				Required:            true,
			},
			"vdc_group": schema.StringAttribute{
				Computed: true,
			},
			"vdc": schema.SingleNestedAttribute{
				MarkdownDescription: "VDC details.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						MarkdownDescription: "VDC name.",
						Computed:            true,
					},
					"description": schema.StringAttribute{
						MarkdownDescription: "VDC UUID.",
						Computed:            true,
					},
					"vdc_service_class": schema.StringAttribute{
						MarkdownDescription: "VDC service class.",
						Computed:            true,
					},
					"vdc_disponibility_class": schema.StringAttribute{
						MarkdownDescription: "VDC disponibility class.",
						Computed:            true,
					},
					"vdc_billing_model": schema.StringAttribute{
						MarkdownDescription: "VDC billing model.",
						Computed:            true,
					},
					"vcpu_in_mhz2": schema.NumberAttribute{
						MarkdownDescription: "VDC CPU in Mhz2.",
						Computed:            true,
					},
					"cpu_allocated": schema.NumberAttribute{
						MarkdownDescription: "VDC CPU allocated.",
						Computed:            true,
					},
					"memory_allocated": schema.NumberAttribute{
						MarkdownDescription: "VDC memory allocated.",
						Computed:            true,
					},
					"vdc_storage_billing_model": schema.StringAttribute{
						MarkdownDescription: "VDC storage billing model.",
						Computed:            true,
					},
					"vdc_storage_profiles": schema.ListNestedAttribute{
						MarkdownDescription: "VDC storage profiles.",
						Computed:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"class": schema.StringAttribute{
									MarkdownDescription: "VDC storage profile class.",
									Computed:            true,
								},
								"limit": schema.NumberAttribute{
									MarkdownDescription: "VDC storage profile limit.",
									Computed:            true,
								},
								"default": schema.BoolAttribute{
									MarkdownDescription: "VDC storage profile default.",
									Computed:            true,
								},
							},
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

	var profiles []vdcStorageProfiles
	for _, profile := range vdcs.Vdc.VdcStorageProfiles {
		p := vdcStorageProfiles{
			Class:   types.StringValue(profile.Class),
			Limit:   types.Int64Value(int64(profile.Limit)),
			Default: types.BoolValue(profile.Default_),
		}
		profiles = append(profiles, p)
	}

	data = vdcDataSourceModel{
		VdcGroup: types.StringValue(vdcs.VdcGroup),
		Name:     types.StringValue(vdcs.Vdc.Name),
		Vdc: &vdcDetail{
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
		},
	}

	data.ID = types.StringValue("frangipane")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
