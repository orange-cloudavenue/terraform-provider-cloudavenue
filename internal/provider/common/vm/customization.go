package vm

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/boolpm"
)

type Customization struct {
	Enabled                        types.Bool   `tfsdk:"enabled"`
	Force                          types.Bool   `tfsdk:"force"`
	ChangeSid                      types.Bool   `tfsdk:"change_sid"`
	AllowLocalAdminPassword        types.Bool   `tfsdk:"allow_local_admin_password"`
	MustChangePasswordOnFirstLogin types.Bool   `tfsdk:"must_change_password_on_first_login"`
	AdminPassword                  types.String `tfsdk:"admin_password"`
	AutoGeneratePassword           types.Bool   `tfsdk:"auto_generate_password"`
	NumberOfAutoLogons             types.Int64  `tfsdk:"number_of_auto_logons"`
	JoinDomain                     types.Bool   `tfsdk:"join_domain"`
	JoinOrgDomain                  types.Bool   `tfsdk:"join_org_domain"`
	JoinDomainName                 types.String `tfsdk:"join_domain_name"`
	JoinDomainUser                 types.String `tfsdk:"join_domain_user"`
	JoinDomainPassword             types.String `tfsdk:"join_domain_password"`
	JoinDomainAccountOU            types.String `tfsdk:"join_domain_account_ou"`
	InitScript                     types.String `tfsdk:"init_script"`
}

// CustomizationSchema returns the schema for the customization block.
func CustomizationSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"force": schema.BoolAttribute{
			MarkdownDescription: "`true` value will cause the VM to reboot on every `apply` operation",
			Optional:            true,
			PlanModifiers: []planmodifier.Bool{
				boolpm.SetDefault(false),
			},
		},
		"enabled": schema.BoolAttribute{
			MarkdownDescription: "`true` value will enable guest customization. It may occur on first boot or when `force` is used",
			Optional:            true,
			Computed:            true,
		},
		"change_sid": schema.BoolAttribute{
			MarkdownDescription: "`true` value will change the SID of the guest OS. Applicable only for Windows VMs.",
			Optional:            true,
			Computed:            true,
		},
		"allow_local_admin_password": schema.BoolAttribute{
			MarkdownDescription: "Allow local administrator password",
			Optional:            true,
			Computed:            true,
		},
		"must_change_password_on_first_login": schema.BoolAttribute{
			MarkdownDescription: "Require Administrator to change password on first login",
			Optional:            true,
			Computed:            true,
		},
		"auto_generate_password": schema.BoolAttribute{
			MarkdownDescription: "Auto generate password",
			Optional:            true,
			Computed:            true,
			Validators: []validator.Bool{
				boolvalidator.ExactlyOneOf(path.MatchRoot("auto_generate_password"), path.MatchRoot("admin_password")),
			},
		},
		"admin_password": schema.StringAttribute{
			MarkdownDescription: "Manually specify admin password",
			Optional:            true,
			Computed:            true,
			Sensitive:           true,
			Validators: []validator.String{
				stringvalidator.ExactlyOneOf(path.MatchRoot("auto_generate_password"), path.MatchRoot("admin_password")),
			},
		},
		"number_of_auto_logons": schema.Int64Attribute{
			MarkdownDescription: "Number of times to log on automatically. `0` is disabled.",
			Optional:            true,
			Computed:            true,
			Validators: []validator.Int64{
				int64validator.AtLeast(0),
			},
		},
		"join_domain": schema.BoolAttribute{
			MarkdownDescription: "Enable this VM to join a domain",
			Optional:            true,
			Computed:            true,
		},
		"join_org_domain": schema.BoolAttribute{
			MarkdownDescription: "Use organization's domain for joining",
			Optional:            true,
			Computed:            true,
		},
		"join_domain_name": schema.StringAttribute{
			MarkdownDescription: "Custom domain name for join",
			Optional:            true,
			Computed:            true,
		},
		"join_domain_user": schema.StringAttribute{
			MarkdownDescription: "Username for custom domain name join",
			Optional:            true,
			Computed:            true,
		},
		"join_domain_password": schema.StringAttribute{
			MarkdownDescription: "Password for custom domain name join",
			Optional:            true,
			Computed:            true,
			Sensitive:           true,
		},
		"join_domain_account_ou": schema.StringAttribute{
			MarkdownDescription: "Organizational Unit (OU) for custom domain name join",
			Optional:            true,
			Computed:            true,
		},
		"init_script": schema.StringAttribute{
			MarkdownDescription: "Script to run on initial boot or with `customization.force=true` set",
			Optional:            true,
			Computed:            true,
		},
	}
}
