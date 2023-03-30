package vdc

import (
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework-validators/float64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"

	superschema "github.com/FrangipaneTeam/terraform-plugin-framework-superschema"
)

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

/*
roleSchema

This function is used to create the schema for the role resource and datasource.
Default is to create a resource schema. If you want to create a datasource schema
you must pass in the withDataSource() option.
*/
func vdcSchema() superschema.Schema {
	return superschema.Schema{
		Common: superschema.SchemaDetails{
			MarkdownDescription: "Provides a Cloud Avenue vDC (Virtual Data Center) ",
		},
		Resource: superschema.SchemaDetails{
			MarkdownDescription: "resource. This can be used to create, update and delete vDC.\n\n",
		},
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: "data source. This can be used to reference a vDC and use its data within other resources or data sources.",
		},
		Attributes: map[string]superschema.Attribute{
			"timeouts": &superschema.TimeoutAttribute{
				Resource: &superschema.ResourceTimeoutAttribute{
					Create: true,
					Read:   true,
					Delete: true,
					Update: true,
				},
			},
			"id": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The ID of the vDC.",
					Computed:            true,
				},
				Resource: &schemaR.StringAttribute{
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
			},
			"name": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The name of the vDC.",
					Required:            true,
					Validators: []validator.String{
						stringvalidator.LengthBetween(2, 27),
					},
				},
				Resource: &schemaR.StringAttribute{
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
			},
			"description": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "A description of the vDC.",
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"cpu_speed_in_mhz": superschema.Float64Attribute{
				Common: &schemaR.Float64Attribute{
					MarkdownDescription: "Specifies the clock frequency, in Mhz, for any virtual CPU that is allocated to a VM.",
				},
				Resource: &schemaR.Float64Attribute{
					Required: true,
					Validators: []validator.Float64{
						float64validator.AtLeast(1200),
					},
				},
				DataSource: &schemaD.Float64Attribute{
					Computed: true,
				},
			},
			"cpu_allocated": superschema.Float64Attribute{
				Common: &schemaR.Float64Attribute{
					MarkdownDescription: "CPU capacity in *MHz* that is committed to be available or used as a limit in PAYG mode.",
				},
				Resource: &schemaR.Float64Attribute{
					MarkdownDescription: "\n\n -> Note: Reserved capacity is automatically set according to the service class.",
					Required:            true,
				},
				DataSource: &schemaD.Float64Attribute{
					Computed: true,
				},
			},
			"memory_allocated": superschema.Float64Attribute{
				Common: &schemaR.Float64Attribute{
					MarkdownDescription: "Memory capacity in Gb that is committed to be available or used as a limit in PAYG mode.",
				},
				Resource: &schemaR.Float64Attribute{
					Required: true,
					Validators: []validator.Float64{
						float64validator.Between(1, 500),
					},
				},
				DataSource: &schemaD.Float64Attribute{
					Computed: true,
				},
			},
			"vdc_group": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "vDC group name.",
				},
				Resource: &schemaR.StringAttribute{
					MarkdownDescription: "This can be an existing vDC group or a new one. This allows you to isolate your vDC.\n" +
						"VMs of vDCs which belong to the same vDC group can communicate together.",
					Required: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"service_class": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The service class of the vDC.",
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
					Validators: []validator.String{
						stringvalidator.OneOf("ECO", "STD", "HP", "VOIP"),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"disponibility_class": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The disponibility class of the vDC.",
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
					Validators: []validator.String{
						stringvalidator.OneOf("ONE-ROOM", "DUAL-ROOM", "HA-DUAL-ROOM"),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"billing_model": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "Choose Billing model of compute resources.",
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
					Validators: []validator.String{
						stringvalidator.OneOf("PAYG", "DRAAS", "RESERVED"),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"storage_billing_model": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "Choose Billing model of storage resources.",
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
					Validators: []validator.String{
						stringvalidator.OneOf("PAYG", "RESERVED"),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"storage_profiles": superschema.SetNestedAttribute{
				Common: &schemaR.SetNestedAttribute{
					MarkdownDescription: "List of storage profiles for this vDC.",
				},
				Resource: &schemaR.SetNestedAttribute{
					Required: true,
					Validators: []validator.Set{
						setvalidator.SizeAtLeast(1),
					},
				},
				DataSource: &schemaD.SetNestedAttribute{
					Computed: true,
				},
				Attributes: map[string]superschema.Attribute{
					"class": superschema.StringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The storage class of the storage profile.",
						},
						Resource: &schemaR.StringAttribute{
							Required: true,
							Validators: []validator.String{
								stringvalidator.OneOf("silver", "silver_r1", "silver_r2", "gold", "gold_r1", "gold_r2", "gold_hm", "platinum3k", "platinum3k_r1", "platinum3k_r2", "platinum3k_hm", "platinum7k", "platinum7k_r1", "platinum7k_r2", "platinum7k_hm"),
							},
						},
						DataSource: &schemaD.StringAttribute{
							Computed: true,
						},
					},
					"limit": superschema.Int64Attribute{
						Common: &schemaR.Int64Attribute{
							MarkdownDescription: "Max number in *Gb* of units allocated for this storage profile.",
						},
						Resource: &schemaR.Int64Attribute{
							Required: true,
							Validators: []validator.Int64{
								int64validator.Between(500, 10000),
							},
						},
						DataSource: &schemaD.Int64Attribute{
							Computed: true,
						},
					},
					"default": superschema.BoolAttribute{
						Common: &schemaR.BoolAttribute{
							MarkdownDescription: "Set this storage profile as default for this vDC. Only one storage profile can be default per vDC.",
						},
						Resource: &schemaR.BoolAttribute{
							Required: true,
						},
						DataSource: &schemaD.BoolAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}
