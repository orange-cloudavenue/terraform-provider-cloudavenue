package vm

import (
	"fmt"

	fboolplanmodifier "github.com/FrangipaneTeam/terraform-plugin-framework-planmodifiers/boolplanmodifier"
	fstringvalidator "github.com/FrangipaneTeam/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	superschema "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/schema"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vdc"
)

type vmAffinityRuleResourceModel struct {
	ID       types.String `tfsdk:"id"`
	VDC      types.String `tfsdk:"vdc"`
	Name     types.String `tfsdk:"name"`
	Polarity types.String `tfsdk:"polarity"`
	Required types.Bool   `tfsdk:"required"`
	Enabled  types.Bool   `tfsdk:"enabled"`
	VMIDs    types.List   `tfsdk:"vm_ids"`
}

// type vmAffinityRuleDataSourceModel struct {
// 	ID       types.String `tfsdk:"id"`
// 	VDC      types.String `tfsdk:"vdc"`
// 	Name     types.String `tfsdk:"name"`
// 	Polarity types.String `tfsdk:"polarity"`
// 	Required types.Bool   `tfsdk:"required"`
// 	Enabled  types.Bool   `tfsdk:"enabled"`
// 	VMIDs    types.List   `tfsdk:"vm_ids"`
// 	RuleID   types.String `tfsdk:"rule_id"`
// }

type vmAffinityRuleSchemaOpts func(*vmAffinityRuleSchemaParams)

type vmAffinityRuleSchemaParams struct {
	resource   bool
	datasource bool
}

// func withDataSource() vmAffinityRuleSchemaOpts {
// 	return func(params *vmAffinityRuleSchemaParams) {
// 		params.datasource = true
// 	}
// }

/*
vmAffinityRuleSchema
This function is used to create the schema for the catalog resource and datasource.
Default is to create a resource schema.  If you want to create a datasource schema
you must pass in the withDataSource() option.
*/
func vmAffinityRuleSchema(opts ...vmAffinityRuleSchemaOpts) superschema.Schema {
	params := &vmAffinityRuleSchemaParams{}

	if len(opts) > 0 {
		for _, opt := range opts {
			opt(params)
		}
	} else {
		params.resource = true
	}

	_schema := superschema.Schema{}
	_schema.Attributes = map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "The ID of the affinity rule.",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"vdc": vdc.Schema(),
		"name": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "VM affinity rule name.",
		},
		"polarity": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: fmt.Sprintf("One of `%s`, `%s`", govcdtypes.PolarityAffinity, govcdtypes.PolarityAntiAffinity),
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
			Validators: []validator.String{
				stringvalidator.OneOf(govcdtypes.PolarityAffinity, govcdtypes.PolarityAntiAffinity),
			},
		},
		"required": schema.BoolAttribute{
			Optional: true,
			Computed: true,
			MarkdownDescription: "`True` if this affinity rule is required. When a rule is mandatory, " +
				"a host failover will not power on the VM if doing so would violate the rule. Default is `true`.",
			PlanModifiers: []planmodifier.Bool{
				fboolplanmodifier.SetDefault(true),
			},
		},
		"enabled": schema.BoolAttribute{
			Optional:            true,
			Computed:            true,
			MarkdownDescription: "`True` if this affinity rule is enabled. Default is `true`.",
			PlanModifiers: []planmodifier.Bool{
				fboolplanmodifier.SetDefault(true),
			},
		},
		"vm_ids": schema.ListAttribute{
			Required:            true,
			ElementType:         types.StringType,
			MarkdownDescription: "Set of VM IDs assigned to this rule.",
			Validators: []validator.List{
				listvalidator.SizeAtMost(2),
				listvalidator.ValueStringsAre(fstringvalidator.IsURN()),
			},
		},
	}

	if params.datasource {
		_schema.MarkdownDescription = "Provides a Cloud Avenue VM Affinity Rule. This can be used to read VM affinity and anti-affinity rules."
		// set computed for all attributes
		_schema = _schema.SetParam(superschema.Computed)

		_schema.Attributes["rule_id"] = schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "VM affinity rule ID. It's the preferred way of identifying a rule.",
			Validators: []validator.String{
				stringvalidator.ExactlyOneOf(path.MatchRoot("name"), path.MatchRoot("rule_id")),
			},
		}
	}

	if params.resource {
		_schema.MarkdownDescription = "Provides a Cloud Avenue VM Affinity Rule. This can be used to create, modify and delete VM affinity and anti-affinity rules."
	}

	return _schema
}
