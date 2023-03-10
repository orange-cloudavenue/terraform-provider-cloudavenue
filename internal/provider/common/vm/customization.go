package vm

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/vmware/go-vcloud-director/v2/govcd"

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

// ToAttrValue converts the Customization struct to a map of attr.Value.
func (c *Customization) ToAttrValue() map[string]attr.Value {
	return map[string]attr.Value{
		"enabled":                             c.Enabled,
		"force":                               c.Force,
		"change_sid":                          c.ChangeSid,
		"allow_local_admin_password":          c.AllowLocalAdminPassword,
		"must_change_password_on_first_login": c.MustChangePasswordOnFirstLogin,
		"admin_password":                      c.AdminPassword,
		"auto_generate_password":              c.AutoGeneratePassword,
		"number_of_auto_logons":               c.NumberOfAutoLogons,
		"join_domain":                         c.JoinDomain,
		"join_org_domain":                     c.JoinOrgDomain,
		"join_domain_name":                    c.JoinDomainName,
		"join_domain_user":                    c.JoinDomainUser,
		"join_domain_password":                c.JoinDomainPassword,
		"join_domain_account_ou":              c.JoinDomainAccountOU,
		"init_script":                         c.InitScript,
	}
}

// CustomizationAttrType.
func CustomizationAttrType() map[string]attr.Type {
	return map[string]attr.Type{
		"enabled":                             types.BoolType,
		"force":                               types.BoolType,
		"change_sid":                          types.BoolType,
		"allow_local_admin_password":          types.BoolType,
		"must_change_password_on_first_login": types.BoolType,
		"admin_password":                      types.StringType,
		"auto_generate_password":              types.BoolType,
		"number_of_auto_logons":               types.Int64Type,
		"join_domain":                         types.BoolType,
		"join_org_domain":                     types.BoolType,
		"join_domain_name":                    types.StringType,
		"join_domain_user":                    types.StringType,
		"join_domain_password":                types.StringType,
		"join_domain_account_ou":              types.StringType,
		"init_script":                         types.StringType,
	}
}

// CustomizationFromPlan converts the terraform plan to a Customization struct.
func CustomizationFromPlan(ctx context.Context, x types.Object) (*Customization, diag.Diagnostics) {
	if x.IsNull() || x.IsUnknown() {
		return &Customization{}, diag.Diagnostics{}
	}

	c := &Customization{}

	d := x.As(ctx, c, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})

	return c, d
}

// ToPlan converts the Customization struct to a terraform plan.
func (c *Customization) ToPlan() basetypes.ObjectValue {
	if c == nil {
		return types.ObjectNull(CustomizationAttrType())
	}

	return types.ObjectValueMust(CustomizationAttrType(), c.ToAttrValue())
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
				boolvalidator.ExactlyOneOf(path.MatchRoot("admin_password")),
			},
		},
		"admin_password": schema.StringAttribute{
			MarkdownDescription: "Manually specify admin password",
			Optional:            true,
			Computed:            true,
			Sensitive:           true,
			Validators: []validator.String{
				stringvalidator.ExactlyOneOf(path.MatchRoot("auto_generate_password")),
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

// CustomizationRead reads the customization block from the API response.
func CustomizationRead(vm *govcd.VM) (c Customization, err error) {
	customizationSection, err := vm.GetGuestCustomizationSection()
	if err != nil {
		return c, fmt.Errorf("unable to get guest customization section: %w", err)
	}

	c = Customization{
		Enabled:                        types.BoolValue(*customizationSection.Enabled),
		ChangeSid:                      types.BoolValue(*customizationSection.ChangeSid),
		AllowLocalAdminPassword:        types.BoolValue(*customizationSection.AdminPasswordEnabled),
		MustChangePasswordOnFirstLogin: types.BoolValue(*customizationSection.ResetPasswordRequired),
		AutoGeneratePassword:           types.BoolValue(*customizationSection.AdminPasswordAuto),
		AdminPassword:                  types.StringValue(customizationSection.AdminPassword),
		NumberOfAutoLogons:             types.Int64Value(int64(customizationSection.AdminAutoLogonCount)),
		JoinDomain:                     types.BoolValue(*customizationSection.JoinDomainEnabled),
		JoinOrgDomain:                  types.BoolValue(*customizationSection.UseOrgSettings),
		JoinDomainName:                 types.StringValue(customizationSection.DomainName),
		JoinDomainUser:                 types.StringValue(customizationSection.DomainUserName),
		JoinDomainPassword:             types.StringValue(customizationSection.DomainUserPassword),
		JoinDomainAccountOU:            types.StringValue(customizationSection.MachineObjectOU),
		InitScript:                     types.StringValue(customizationSection.CustomizationScript),
	}

	return c, nil
}
