package vm

import (
	fstringvalidator "github.com/orange-cloudavenue/terraform-plugin-framework-validators/stringvalidator"

	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"

	superschema "github.com/FrangipaneTeam/terraform-plugin-framework-superschema"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vdc"
)

/*
vmAffinityRuleSchema
This function is used to create the schema for the catalog resource and datasource.
Default is to create a resource schema.  If you want to create a datasource schema
you must pass in the withDataSource() option.
*/
func vmAffinityRuleSchema() superschema.Schema {
	return superschema.Schema{
		Common: superschema.SchemaDetails{
			MarkdownDescription: "Provides a Cloud Avenue VM Affinity Rule.",
		},
		Resource: superschema.SchemaDetails{
			MarkdownDescription: "This can be used to create, modify and delete VM affinity and anti-affinity rules.",
		},
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: "This can be used to read VM affinity and anti-affinity rules.",
		},
		Attributes: map[string]superschema.Attribute{
			"id": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The ID of the affinity rule.",
				},
				Resource: &schemaR.StringAttribute{
					Computed: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Optional: true,
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot("name"), path.MatchRoot("id")),
					},
				},
			},
			"vdc": vdc.SuperSchema(),
			"name": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "VM affinity rule name.",
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
				},
				DataSource: &schemaD.StringAttribute{
					Optional: true,
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot("name"), path.MatchRoot("id")),
					},
				},
			},
			"polarity": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The polarity of the affinity rule.",
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
					Validators: []validator.String{
						stringvalidator.OneOf(govcdtypes.PolarityAffinity, govcdtypes.PolarityAntiAffinity),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"enabled": superschema.BoolAttribute{
				Resource: &schemaR.BoolAttribute{
					MarkdownDescription: "`True` if this affinity rule is enabled.",
					Optional:            true,
					Computed:            true,
					Default:             booldefault.StaticBool(true),
				},
				DataSource: &schemaD.BoolAttribute{
					MarkdownDescription: "`True` if this affinity rule is enabled.",
					Computed:            true,
				},
			},
			"required": superschema.BoolAttribute{
				Common: &schemaR.BoolAttribute{
					MarkdownDescription: "`True` if this affinity rule is required. When a rule is mandatory, " +
						"a host failover will not power on the VM if doing so would violate the rule.",
				},
				Resource: &schemaR.BoolAttribute{
					Optional: true,
					Computed: true,
					Default:  booldefault.StaticBool(true),
				},
				DataSource: &schemaD.BoolAttribute{
					Computed: true,
				},
			},
			"vm_ids": superschema.SetAttribute{
				Common: &schemaR.SetAttribute{
					MarkdownDescription: "List of VM IDs",
					ElementType:         types.StringType,
				},
				Resource: &schemaR.SetAttribute{
					MarkdownDescription: "to apply the affinity rule to.",
					Required:            true,
					Validators: []validator.Set{
						setvalidator.SizeAtLeast(2),
						setvalidator.ValueStringsAre(fstringvalidator.IsURN()),
					},
				},
				DataSource: &schemaD.SetAttribute{
					MarkdownDescription: "associated to the affinity rule.",
					Computed:            true,
				},
			},
		},
	}
}
