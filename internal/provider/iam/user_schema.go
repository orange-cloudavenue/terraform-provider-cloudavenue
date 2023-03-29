package iam

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"

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

/*
userSchema

This function is used to create the schema for the user resource and datasource.
*/
func userSchema() superschema.Schema {
	return superschema.Schema{
		Common: superschema.SchemaDetails{
			MarkdownDescription: "The user",
		},
		Resource: superschema.SchemaDetails{
			MarkdownDescription: " resource allows you to manage local users in Cloud Avenue.",
		},
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: " data source allows you to read users in Cloud Avenue.",
		},
		Attributes: map[string]superschema.Attribute{
			"id": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The ID of the user.",
				},
				Resource: &schemaR.StringAttribute{
					Computed: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
				DataSource: &schemaD.StringAttribute{
					MarkdownDescription: " Required if `name` is not set.",
					Optional:            true,
					Computed:            true,
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot("name"), path.MatchRoot("id")),
					},
				},
			},
			"name": superschema.StringAttribute{
				Resource: &schemaR.StringAttribute{
					MarkdownDescription: "(ForceNew) The name of the user.",
					Required:            true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
				DataSource: &schemaD.StringAttribute{
					MarkdownDescription: "The name of the user. Required if `id` is not set.",
					Optional:            true,
					Computed:            true,
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot("name"), path.MatchRoot("id")),
					},
				},
			},
			"role_name": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The role assigned to the user.",
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"full_name": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The user's full name.",
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"email": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The user's email address.",
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"telephone": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The user's telephone number.",
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"enabled": superschema.BoolAttribute{
				Common: &schemaR.BoolAttribute{
					MarkdownDescription: "`true` if the user is enabled and can log in. (Default to `true`)",
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
			"deployed_vm_quota": superschema.Int64Attribute{
				Common: &schemaR.Int64Attribute{
					MarkdownDescription: "Quota of vApps that this user can deploy. A value of `0` specifies an unlimited quota. (Default to `0`)",
				},
				Resource: &schemaR.Int64Attribute{
					Optional: true,
					Computed: true,
					Default:  int64default.StaticInt64(0),
				},
				DataSource: &schemaD.Int64Attribute{
					Computed: true,
				},
			},
			"stored_vm_quota": superschema.Int64Attribute{
				Common: &schemaR.Int64Attribute{
					MarkdownDescription: "Quota of vApps that this user can store. A value of `0` specifies an unlimited quota. (Default to `0`)",
				},
				Resource: &schemaR.Int64Attribute{
					Optional: true,
					Computed: true,
					Default:  int64default.StaticInt64(0),
				},
				DataSource: &schemaD.Int64Attribute{
					Computed: true,
				},
			},
			"password": superschema.StringAttribute{
				Resource: &schemaR.StringAttribute{
					MarkdownDescription: "The user's password. This value is never returned on read.",
					Required:            true,
					Sensitive:           true,
				},
			},
			"take_ownership": superschema.BoolAttribute{
				Resource: &schemaR.BoolAttribute{
					MarkdownDescription: "`true` if the user should take ownership of all vApps and media that are currently owned by the user that is being deleted. (Default to `true`)",
					Optional:            true,
					Computed:            true,
					Default:             booldefault.StaticBool(true),
				},
			},
			"provider_type": superschema.StringAttribute{
				DataSource: &schemaD.StringAttribute{
					MarkdownDescription: "Identity provider type for this this user. One of: `INTEGRATED`, `SAML`, `OAUTH`.",
					Computed:            true,
				},
			},
		},
	}
}
