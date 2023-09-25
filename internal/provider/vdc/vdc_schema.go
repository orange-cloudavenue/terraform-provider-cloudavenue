package vdc

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"

	superschema "github.com/FrangipaneTeam/terraform-plugin-framework-superschema"
)

/*
vdcSchema

This function is used to create the superschema for the vdc resource and datasource.
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
			"id": superschema.SuperStringAttribute{
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
			"name": superschema.SuperStringAttribute{
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
			"description": superschema.SuperStringAttribute{
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
			"cpu_speed_in_mhz": superschema.SuperInt64Attribute{
				Common: &schemaR.Int64Attribute{
					MarkdownDescription: "Specifies the clock frequency, in Mhz, for any virtual CPU that is allocated to a VM.",
				},
				Resource: &schemaR.Int64Attribute{
					Required: true,
					PlanModifiers: []planmodifier.Int64{
						int64planmodifier.RequiresReplaceIf(func(ctx context.Context, request planmodifier.Int64Request, resp *int64planmodifier.RequiresReplaceIfFuncResponse) {
							billingModel := new(types.String)
							resp.Diagnostics.Append(request.Plan.GetAttribute(ctx, path.Root("billing_model"), billingModel)...)
							if resp.Diagnostics.HasError() {
								return
							}
							if billingModel.ValueString() != "RESERVED" {
								resp.RequiresReplace = true
								resp.Diagnostics.AddAttributeWarning(path.Root("cpu_speed_in_mhz"), "Force replacement attributes", "You can change the cpu_speed_in_mhz attribute only if the billing_model is set to RESERVED.")
							}
						}, "", ""),
					},
					MarkdownDescription: "Force replacement attributes, however you can change the `cpu_speed_in_mhz` attribute only if the `billing_model` is set to **RESERVED**.",
					Validators: []validator.Int64{
						int64validator.AtLeast(1200),
					},
				},
				DataSource: &schemaD.Int64Attribute{
					Computed: true,
				},
			},

			"cpu_allocated": superschema.SuperInt64Attribute{
				Common: &schemaR.Int64Attribute{
					MarkdownDescription: "CPU capacity in *MHz* that is committed to be available or used as a limit in PAYG mode.",
				},
				Resource: &schemaR.Int64Attribute{
					MarkdownDescription: "\n\n -> Note: Reserved capacity is automatically set according to the service class.",
					Required:            true,
				},
				DataSource: &schemaD.Int64Attribute{
					Computed: true,
				},
			},
			"memory_allocated": superschema.SuperInt64Attribute{
				Common: &schemaR.Int64Attribute{
					MarkdownDescription: "Memory capacity in Gb that is committed to be available or used as a limit in PAYG mode.",
				},
				Resource: &schemaR.Int64Attribute{
					Required: true,
					Validators: []validator.Int64{
						int64validator.Between(1, 500),
					},
				},
				DataSource: &schemaD.Int64Attribute{
					Computed: true,
				},
			},
			"vdc_group": superschema.SuperStringAttribute{
				Deprecated: &superschema.Deprecated{
					DeprecationMessage:                "Remove the vdc_group attribute configuration as it replaced by the resource cloudavenue_vdc_group and the attribute will be removed in the version 0.12.0 of the provider.",
					ComputeMarkdownDeprecationMessage: true,
					Removed:                           true,
					FromAttributeName:                 "vdc_group",
					TargetRelease:                     "v0.12.0",
					LinkToMilestone:                   "https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/milestone/4",
					TargetResourceName:                "cloudavenue_vdc_group",
					LinkToResourceDoc:                 "https://registry.terraform.io/providers/orange-cloudavenue/cloudavenue/latest/docs/resources/vdc_group",
					LinkToIssue:                       "https://github.com/orange-cloudavenue/terraform-provider-cloudavenue/issues/448",
				},
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "vDC group name.",
				},
				Resource: &schemaR.StringAttribute{
					MarkdownDescription: "This can be an existing vDC group or a new one. This allows you to isolate your vDC.\n" +
						"VMs of vDCs which belong to the same vDC group can communicate together.",
					Optional: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplaceIfConfigured(),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"service_class": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The service class of the vDC.",
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
					Validators: []validator.String{
						stringvalidator.OneOf("ECO", "STD", "HP", "VOIP"),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"disponibility_class": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The disponibility class of the vDC.",
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
					Validators: []validator.String{
						stringvalidator.OneOf("ONE-ROOM", "DUAL-ROOM", "HA-DUAL-ROOM"),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"billing_model": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "Choose Billing model of compute resources.",
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
					Validators: []validator.String{
						stringvalidator.OneOf("PAYG", "DRAAS", "RESERVED"),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"storage_billing_model": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "Choose Billing model of storage resources.",
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
					Validators: []validator.String{
						stringvalidator.OneOf("PAYG", "RESERVED"),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"storage_profiles": superschema.SuperSetNestedAttribute{
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
					"class": superschema.SuperStringAttribute{
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
					"limit": superschema.SuperInt64Attribute{
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
					"default": superschema.SuperBoolAttribute{
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
