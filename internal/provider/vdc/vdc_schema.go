package vdc

import (
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"

	superschema "github.com/FrangipaneTeam/terraform-plugin-framework-superschema"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/v1/infrapi/rules"
)

const seeVDCRules = "See [Rules](https://registry.terraform.io/providers/orange-cloudavenue/cloudavenue/latest/docs/resources/vdc#rules) for more information."

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
					Required:            true,
					MarkdownDescription: seeVDCRules,
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
					MarkdownDescription: seeVDCRules,
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
				},
				DataSource: &schemaD.Int64Attribute{
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
						stringvalidator.OneOf(func() []string {
							var serviceClasses []string
							for _, sC := range rules.ALLServiceClasses {
								serviceClasses = append(serviceClasses, string(sC))
							}
							return serviceClasses
						}()...),
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
					Required:            true,
					MarkdownDescription: "The disponibility class available are different depending on the service class. " + seeVDCRules,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
					Validators: []validator.String{
						stringvalidator.OneOf(func() []string {
							var disponibilityClasses []string
							for _, dC := range rules.ALLDisponibilityClasses {
								disponibilityClasses = append(disponibilityClasses, string(dC))
							}
							return disponibilityClasses
						}()...),
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
					Required:            true,
					MarkdownDescription: "The billing model available are different depending on the service class. " + seeVDCRules,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
					Validators: []validator.String{
						stringvalidator.OneOf(func() []string {
							var billingModels []string
							for _, bM := range rules.ALLBillingModels {
								billingModels = append(billingModels, string(bM))
							}
							return billingModels
						}()...),
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
					Required:            true,
					MarkdownDescription: "The billing model available are different depending on the service class. " + seeVDCRules,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
					Validators: []validator.String{
						stringvalidator.OneOf(func() []string {
							var billingModels []string
							for _, bM := range rules.ALLStorageBillingModels {
								billingModels = append(billingModels, string(bM))
							}
							return billingModels
						}()...),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"storage_profiles": superschema.SuperSetNestedAttributeOf[vdcResourceModelVDCStorageProfile]{
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
							Required:            true,
							MarkdownDescription: "The storage class available are different depending on the service class. " + seeVDCRules,
							Validators: []validator.String{
								stringvalidator.OneOf(func() []string {
									var storageProfileClasses []string
									for _, sPC := range rules.ALLStorageProfilesClass {
										storageProfileClasses = append(storageProfileClasses, string(sPC))
									}
									return storageProfileClasses
								}()...),
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
