package iam

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	superschema "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/schema"
)

type roleResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Rights      types.Set    `tfsdk:"rights"`
}

type roleDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	ReadOnly    types.Bool   `tfsdk:"read_only"`
	Rights      types.Set    `tfsdk:"rights"`
}

type roleSchemaOpts func(*roleSchemaParams)

type roleSchemaParams struct {
	resource   bool
	datasource bool
}

func withRoleDataSource() roleSchemaOpts {
	return func(params *roleSchemaParams) {
		params.datasource = true
	}
}

/*
roleSchema

This function is used to create the schema for the role resource and datasource.
Default is to create a resource schema. If you want to create a datasource schema
you must pass in the withDataSource() option.
*/
func roleSchema(opts ...roleSchemaOpts) superschema.Schema {
	params := &roleSchemaParams{}

	if len(opts) > 0 {
		for _, opt := range opts {
			opt(params)
		}
	} else {
		params.resource = true
	}

	_schema := superschema.Schema{}

	idAttribute := schema.StringAttribute{
		MarkdownDescription: "The ID is a unique identifier for the role.",
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	}

	nameAttribute := schema.StringAttribute{
		MarkdownDescription: "The name of the role.",
	}

	// Global schemas
	_schema.Attributes = map[string]schema.Attribute{
		"description": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "A description for the role",
		},
		"rights": schema.SetAttribute{
			Required:            true,
			MarkdownDescription: "A list of rights for the role",
			ElementType:         types.StringType,
		},
	}

	if params.resource {
		_schema.MarkdownDescription = "The role resource allows you to manage local role in Cloud Avenue."

		idAttribute.Computed = true
		_schema.Attributes["id"] = idAttribute

		nameAttribute.Required = true
		_schema.Attributes["name"] = nameAttribute
	}
	if params.datasource {
		_schema.MarkdownDescription = "The user data source allows you to read local role in Cloud Avenue."

		_schema = _schema.SetParam(superschema.Computed, "description", "rights")

		idAttribute.Optional = true
		idAttribute.Computed = true
		idAttribute.MarkdownDescription += " Required if `name` is not set."
		idAttribute.Validators = []validator.String{
			stringvalidator.ExactlyOneOf(path.MatchRoot("id"), path.MatchRoot("name")),
		}
		_schema.Attributes["id"] = idAttribute

		nameAttribute.Optional = true
		nameAttribute.Computed = true
		nameAttribute.MarkdownDescription += " Required if `id` is not set."
		nameAttribute.Validators = []validator.String{
			stringvalidator.ExactlyOneOf(path.MatchRoot("id"), path.MatchRoot("name")),
		}
		_schema.Attributes["name"] = nameAttribute

		_schema.Attributes["read_only"] = schema.BoolAttribute{
			Computed:            true,
			MarkdownDescription: "Indicates if the role is read only",
		}
	}
	return _schema
}
