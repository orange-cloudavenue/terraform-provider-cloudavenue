package iam

import (
	fboolplanmodifier "github.com/FrangipaneTeam/terraform-plugin-framework-planmodifiers/boolplanmodifier"
	fint64planmodifier "github.com/FrangipaneTeam/terraform-plugin-framework-planmodifiers/int64planmodifier"
	fstringplanmodifier "github.com/FrangipaneTeam/terraform-plugin-framework-planmodifiers/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	superschema "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/schema"
)

type userResourceModel struct {
	// Base
	ID              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	RoleName        types.String `tfsdk:"role_name"`
	FullName        types.String `tfsdk:"full_name"`
	Email           types.String `tfsdk:"email"`
	Telephone       types.String `tfsdk:"telephone"`
	Enabled         types.Bool   `tfsdk:"enabled"`
	DeployedVMQuota types.Int64  `tfsdk:"deployed_vm_quota"`
	StoredVMQuota   types.Int64  `tfsdk:"stored_vm_quota"`

	// Specific
	Password      types.String `tfsdk:"password"`
	TakeOwnership types.Bool   `tfsdk:"take_ownership"`
}

type userDataSourceModel struct {
	// Base
	ID              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	RoleName        types.String `tfsdk:"role_name"`
	FullName        types.String `tfsdk:"full_name"`
	Email           types.String `tfsdk:"email"`
	Telephone       types.String `tfsdk:"telephone"`
	Enabled         types.Bool   `tfsdk:"enabled"`
	DeployedVMQuota types.Int64  `tfsdk:"deployed_vm_quota"`
	StoredVMQuota   types.Int64  `tfsdk:"stored_vm_quota"`

	// Specific
	ProviderType types.String `tfsdk:"provider_type"`
}

type userSchemaOpts func(*userSchemaParams)

type userSchemaParams struct {
	resource   bool
	datasource bool
}

func withDataSource() userSchemaOpts {
	return func(params *userSchemaParams) {
		params.datasource = true
	}
}

/*
userSchema

This function is used to create the schema for the user resource and datasource.
Default is to create a resource schema. If you want to create a datasource schema
you must pass in the withDataSource() option.
*/
func userSchema(opts ...userSchemaOpts) superschema.Schema {
	params := &userSchemaParams{}

	if len(opts) > 0 {
		for _, opt := range opts {
			opt(params)
		}
	} else {
		params.resource = true
	}

	_schema := superschema.Schema{}

	// Specific attributes
	idAttribute := schema.StringAttribute{
		MarkdownDescription: "The ID is a unique identifier for the user.",
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	}

	nameAttribute := schema.StringAttribute{
		MarkdownDescription: "The name of the user.",
	}

	roleNameAttribute := schema.StringAttribute{
		MarkdownDescription: "The role assigned to the user.",
	}

	enabledAttribute := schema.BoolAttribute{
		MarkdownDescription: "`true` if the user is enabled and can log in.",
	}

	// Global schemas
	_schema.Attributes = map[string]schema.Attribute{
		"full_name": schema.StringAttribute{
			MarkdownDescription: "The user's full name",
			Optional:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"email": schema.StringAttribute{
			MarkdownDescription: "The user's email address",
			Optional:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"telephone": schema.StringAttribute{
			MarkdownDescription: "The user's telephone number",
			Optional:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"deployed_vm_quota": schema.Int64Attribute{
			MarkdownDescription: "Quota of vApps that this user can deploy. A value of `0` specifies an unlimited quota.",
			Optional:            true,
			Computed:            true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
				fint64planmodifier.SetDefault(0),
			},
		},
		"stored_vm_quota": schema.Int64Attribute{
			MarkdownDescription: "Quota of vApps that this user can store. A value of `0` specifies an unlimited quota.",
			Optional:            true,
			Computed:            true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
				fint64planmodifier.SetDefault(0),
			},
		},
	}

	if params.resource {
		// * Settings schema
		_schema.MarkdownDescription = "The user resource allows you to manage local users in Cloud Avenue."

		// * Settings attributes
		// * ID
		idAttribute.Computed = true
		_schema.Attributes["id"] = idAttribute

		// * Name
		nameAttribute.Required = true
		nameAttribute.Description += " Only lowercase letters allowed."
		nameAttribute.PlanModifiers = []planmodifier.String{
			stringplanmodifier.RequiresReplace(),
			fstringplanmodifier.ToLower(),
		}
		_schema.Attributes["name"] = nameAttribute

		// * Role Name
		roleNameAttribute.Required = true
		roleNameAttribute.PlanModifiers = []planmodifier.String{
			stringplanmodifier.RequiresReplace(),
		}
		_schema.Attributes["role_name"] = roleNameAttribute

		// * Enabled
		enabledAttribute.MarkdownDescription += " Defaults to `true`."
		enabledAttribute.Optional = true
		enabledAttribute.Computed = true
		enabledAttribute.PlanModifiers = []planmodifier.Bool{
			fboolplanmodifier.SetDefault(true),
		}
		_schema.Attributes["enabled"] = enabledAttribute

		// * Password
		_schema.Attributes["password"] = schema.StringAttribute{
			MarkdownDescription: "The user's password. This value is never returned on read.",
			Required:            true,
			Sensitive:           true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(8),
			},
		}

		// * Take Ownership
		_schema.Attributes["take_ownership"] = schema.BoolAttribute{
			MarkdownDescription: "`true` if the user should take ownership of all vApps and media that are currently owned by the user that is being deleted.",
			Optional:            true,
		}
	}

	if params.datasource {
		// * Settings schema
		_schema.MarkdownDescription = "The user data source allows you to read users in Cloud Avenue."

		// * Settings attributes

		// * ID
		idAttribute.Optional = true
		idAttribute.Computed = true
		idAttribute.MarkdownDescription += " Required if `name` is not set."
		idAttribute.Validators = []validator.String{
			stringvalidator.ExactlyOneOf(path.MatchRoot("id"), path.MatchRoot("name")),
		}
		_schema.Attributes["id"] = idAttribute

		// * Name
		nameAttribute.Optional = true
		nameAttribute.Computed = true
		nameAttribute.MarkdownDescription += " Required if `id` is not set."
		nameAttribute.Validators = []validator.String{
			stringvalidator.ExactlyOneOf(path.MatchRoot("id"), path.MatchRoot("name")),
		}
		_schema.Attributes["name"] = nameAttribute

		// * Enabled
		enabledAttribute.Computed = true
		_schema.Attributes["enabled"] = enabledAttribute

		// * Role Name
		roleNameAttribute.Computed = true
		_schema.Attributes["role_name"] = roleNameAttribute

		// * Provider Type
		_schema.Attributes["provider_type"] = schema.StringAttribute{
			MarkdownDescription: "The type of provider used to authenticate the user.",
			Computed:            true,
		}

		// TODO Set to only compute (full_name, email, telephone, deployed_vm_quota, stored_vm_quota)
	}

	return _schema
}
